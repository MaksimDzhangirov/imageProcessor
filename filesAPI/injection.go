package main

import (
	"fmt"
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/datasrc"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/handler"
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/repository"
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/service"
)

func inject(d *datasrc.DataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	folderName := os.Getenv("IMAGE_FOLDER")
	imageRepository := repository.NewDiskImageRepository(folderName)

	/*
	 * service layer
	 */
	imageService := service.NewImageService(&service.ImageConfig{
		ImageRepository: imageRepository,
	})

	amqpService := service.NewAmqpService(d.RabbitMQConnection)

	// initialize gin.Engine
	router := gin.Default()
	router.Static("/images", fmt.Sprintf("./%s", folderName))

	maxBodyBytes := os.Getenv("MAX_BODY_BYTES")
	mbb, err := strconv.ParseInt(maxBodyBytes, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse MAX_BODY_BYTES as int: %w", err)
	}

	handler.NewHandler(&handler.Config{
		R:            router,
		ImageService: imageService,
		AmqpService:  amqpService,
		MaxBodyBytes: mbb,
	})

	return router, nil
}
