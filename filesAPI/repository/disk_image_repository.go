package repository

import (
	"fmt"
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/model"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
)

// DiskImageRepository stores image on disk
type DiskImageRepository struct {	
	imageFolder string	
}

// NewDiskImageRepository returns a new DiskImageRepository
func NewDiskImageRepository(imageFolder string) *DiskImageRepository {
	return &DiskImageRepository{
		imageFolder: imageFolder,		
	}
}

// Save saves a new image to the store
func (store *DiskImageRepository) Save(
	imageType string,
	imageFile multipart.File,
) (*model.ImageInfo, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("cannot generate image id: %w", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageID, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot create image file: %w", err)
	}

	_, err = io.Copy(file, imageFile)
	if err != nil {
		return nil, fmt.Errorf("cannot write image to file: %w", err)
	}

	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("file.Close: %v", err)
	}

	imageInfo := &model.ImageInfo{
		UID:       imageID,
		ImageType: imageType,
		ImagePath: imagePath,
	}

	return imageInfo, nil
}