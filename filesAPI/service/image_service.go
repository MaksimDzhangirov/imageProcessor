package service

import (
	"errors"
	"log"
	"mime/multipart"
	"path/filepath"

	"github.com/h2non/bimg"

	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/model"
)

type imageService struct {
	ImageRepository model.ImageRepository
}

type ImageConfig struct {
	ImageRepository model.ImageRepository
}

// NewUserService is a factory function for
// initializing a ImageService with its repository layer dependencies
func NewImageService(c *ImageConfig) model.ImageService {
	return &imageService{
		ImageRepository: c.ImageRepository,
	}
}

func (s *imageService) UploadImage(
	imageFileHeader *multipart.FileHeader,
) (*model.ImageInfo, error) {
	imageFile, err := imageFileHeader.Open()
	if err != nil {
		log.Printf("Failed to open image file: %v\n", err)
		return nil, errors.New("failed to open image file")
	}

	extension := filepath.Ext(imageFileHeader.Filename)
	// Upload user's image to ImageRepository
	imageInfo, err := s.ImageRepository.Save(extension, imageFile)

	if err != nil {
		log.Printf("Unable to upload image: %v\n", err)
		return nil, err
	}

	return imageInfo, nil
}

func (s *imageService) OptimizeImage(path string, options bimg.Options) error {
	buffer, err := bimg.Read(path)
	if err != nil {
		return err
	}

	newImage, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		return err
	}

	return bimg.Write(path, newImage)
}
