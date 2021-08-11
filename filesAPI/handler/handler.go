package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/model"
)

// Handler struct holds required services for handler to function
type Handler struct {
	ImageService model.ImageService
	AmqpService  model.AmqpService
	MaxBodyBytes int64
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R            *gin.Engine
	ImageService model.ImageService
	AmqpService  model.AmqpService
	MaxBodyBytes int64
}

// NewHandler initializes the handler with required injected services along with http routes
// Does not return as it deals directly with a reference to the gin Engine
func NewHandler(c *Config) {
	h := &Handler{
		ImageService: c.ImageService,
		AmqpService:  c.AmqpService,
		MaxBodyBytes: c.MaxBodyBytes,
	}

	c.R.GET("/", h.MainPage)

	c.R.POST("/image", h.Image)
}
