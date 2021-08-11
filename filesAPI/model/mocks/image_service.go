package mocks

import (
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/model"
	"github.com/h2non/bimg"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
)

// MockImageService is a mock type for model.ImageService
type MockImageService struct {
	mock.Mock
}

func (m *MockImageService) UploadImage(
	imageFileHeader *multipart.FileHeader,
) (*model.ImageInfo, error) {
	ret := m.Called(imageFileHeader)

	// first value passed to "Return"
	var r0 *model.ImageInfo
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.ImageInfo)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockImageService) OptimizeImage(path string, options bimg.Options) error {
	ret := m.Called(path, options)

	var r1 error

	if ret.Get(0) != nil {
		r1 = ret.Get(1).(error)
	}

	return r1
}