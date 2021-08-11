package main

import (
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/datasrc"
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/model"
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/repository"
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/service"
	"github.com/h2non/bimg"
	"log"
	"os"
)

func main() {
	log.Println("Image processing...")
	ds, err := datasrc.InitDS()

	if err != nil {
		log.Fatalf("Unable to initialize data sources: %v\n", err)
	}

	amqpService := service.NewAmqpService(ds.RabbitMQConnection)

	queueConfig := model.QueueConfig{
		Name:             os.Getenv("AMQP_QUEUE_NAME"),
		Durable:          true,
		DeleteWhenUnused: false,
		Exclusive:        false,
		NoWait:           false,
		Arguments:        nil,
	}
	consumeConfig := model.ConsumeConfig{
		Consumer:      "",
		AutoAck:       false,
		Exclusive:     false,
		NoLocal:       false,
		NoWait:        false,
		Arguments:     nil,
		FuncToPerform: func(uid string) error {
			options := bimg.Options{Width: 240, Height: 160}
			folderName := os.Getenv("IMAGE_FOLDER")
			imageRepository := repository.NewDiskImageRepository(folderName)
			imageService := service.NewImageService(&service.ImageConfig{
				ImageRepository: imageRepository,
			})
			return imageService.OptimizeImage(uid, options)
		},
	}
	err = amqpService.Consume(queueConfig, consumeConfig)
	if err != nil {
		log.Println("Consume error")
	}
	if err = ds.Close(); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: %v\n", err)
	}
}
