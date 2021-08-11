package handler

import (
	"fmt"
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/model"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Image handler
func (h *Handler) Image(c *gin.Context) {

	// limit overly large request bodies
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.MaxBodyBytes)

	imageFileHeader, err := c.FormFile("imageFile")

	// check for error before checking for non-nil header
	if err != nil {
		// should be a validation error
		log.Printf("Unable parse multipart/form-data: %+v", err)

		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": fmt.Sprintf("Max request body size is %v bytes\n", h.MaxBodyBytes),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to parse multipart/form-data",
		})
		return
	}

	if imageFileHeader == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Must include an imageFile",
		})
		return
	}

	mimeType := imageFileHeader.Header.Get("Content-Type")

	// Validate image mime-type is allowable
	if valid := isAllowedImageType(mimeType); !valid {
		log.Println("Image is not an allowable mime-type")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ImageFile must be 'image/jpeg', 'image/jpg' or 'image/png'",
		})
		return
	}

	imageInfo, err := h.ImageService.UploadImage(imageFileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}

	queueConfig := model.QueueConfig{
		Name:             os.Getenv("AMQP_QUEUE_NAME"),
		Durable:          true,
		DeleteWhenUnused: false,
		Exclusive:        false,
		NoWait:           false,
		Arguments:        nil,
	}
	publishConfig := model.PublishConfig{
		Exchange:        "",
		Mandatory:       false,
		Immediate:       false,
		StringToPublish: imageInfo.ImagePath,
	}
	err = h.AmqpService.Send(queueConfig, publishConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"imageUrl": imageInfo.ImagePath,
		"message":  "success",
	})
}
