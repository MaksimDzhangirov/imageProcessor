package model

import (
	"github.com/google/uuid"
)

// ImageInfo contains image information
type ImageInfo struct {
	UID       uuid.UUID `json:"uid"`
	ImageType string    `json:"image_type"`
	ImagePath string    `json:"image_path"`
}
