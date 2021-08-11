package handler

import (
	"encoding/json"
	"errors"
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/model"
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/model/fixture"
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/model/mocks"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestImage(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	mockImageService := new(mocks.MockImageService)
	mockAmqpService := new(mocks.MockAmqpService)

	NewHandler(&Config{
		R:            router,
		ImageService: mockImageService,
		AmqpService:  mockAmqpService,
		MaxBodyBytes: 10 * 1024 * 1024,
	})

	t.Run("Success", func(t *testing.T) {
		rr := httptest.NewRecorder()

		imageURL := "/images/1234.png"

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()

		setImageArgs := mock.Arguments{
			mock.AnythingOfType("*multipart.FileHeader"),
		}

		imageInfo := model.ImageInfo{
			UID:       uuid.UUID{},
			ImageType: "png",
			ImagePath: imageURL,
		}

		mockImageService.On("UploadImage", setImageArgs...).Return(&imageInfo, nil)

		setAmqpArgs := mock.Arguments{
			mock.AnythingOfType("model.QueueConfig"),
			mock.AnythingOfType("model.PublishConfig"),
		}
		mockAmqpService.On("Send", setAmqpArgs...).Return(nil)

		request, _ := http.NewRequest(http.MethodPost, "/image", multipartImageFixture.MultipartBody)
		request.Header.Set("Content-Type", multipartImageFixture.ContentType)

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"imageUrl": imageURL,
			"message":  "success",
		})

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockImageService.AssertCalled(t, "UploadImage", setImageArgs...)
		mockAmqpService.AssertCalled(t, "Send", setAmqpArgs...)
	})

	t.Run("Disallowed mimetype", func(t *testing.T) {
		rr := httptest.NewRecorder()

		multipartImageFixture := fixture.NewMultipartImage("image.txt", "mage/svg+xml")
		defer multipartImageFixture.Close()

		request, _ := http.NewRequest(http.MethodPost, "/image", multipartImageFixture.MultipartBody)
		request.Header.Set("Content-Type", "multipart/form-data")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockImageService.AssertNotCalled(t, "UploadImage")
	})

	t.Run("No image file provided", func(t *testing.T) {
		rr := httptest.NewRecorder()

		request, _ := http.NewRequest(http.MethodPost, "/image", nil)
		request.Header.Set("Content-Type", "multipart/form-data")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockImageService.AssertNotCalled(t, "UploadImage")
	})

	t.Run("Error from UploadImage", func(t *testing.T) {

		router := gin.Default()

		mockImageService := new(mocks.MockImageService)

		NewHandler(&Config{
			R:            router,
			ImageService: mockImageService,
			MaxBodyBytes: 10 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()

		setImageArgs := mock.Arguments{
			mock.AnythingOfType("*multipart.FileHeader"),
		}

		mockError := errors.New("internal server error")

		mockImageService.On("UploadImage", setImageArgs...).Return(nil, mockError)

		request, _ := http.NewRequest(http.MethodPost, "/image", multipartImageFixture.MultipartBody)
		request.Header.Set("Content-Type", multipartImageFixture.ContentType)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		mockImageService.AssertCalled(t, "UploadImage", setImageArgs...)
	})
}
