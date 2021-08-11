package model

import (
	"github.com/h2non/bimg"
	"mime/multipart"
)

// ImageService defines methods the handler layer expects
// any service it interacts with to implement
type ImageService interface {
	UploadImage(imageFileHeader *multipart.FileHeader) (*ImageInfo, error)
	OptimizeImage(path string, options bimg.Options) error
}

// AmqpService is an interface to interact with RabbitMQ
type AmqpService interface {
	Send(queueConfig QueueConfig, publishConfig PublishConfig) error
	Consume(queueConfig QueueConfig, consumeConfig ConsumeConfig) error
}

// ImageRepository is an interface to store images
type ImageRepository interface {
	// Save saves a new image to the store
	Save(imageType string, imageData multipart.File) (*ImageInfo, error)
}
