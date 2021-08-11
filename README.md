# Project Overview

## Summary

The goal is to create a system for uploading and processing files. It should 
consist of two microservices, one responsible for file uploading and retrieval
(Files API), and another responsible for file processing (Processing API).
Each microservice should reside in a separate main package.

## Files API

The Files API microservice should expose an HTTP endpoint for uploading files.
Once a file is uploaded, it should be saved to the file system, and its ID 
should be sent to the Processing API via a RabbitMQ queue.

## Processing API

The Processing API microservice should provide functionality for processing 
and optimizing images. It should accept file ID from a RabbitMQ queue, reduce 
the image size and overwrite the file in the file system.

## How to run

```shell
docker-compose up
```

http://file.upload.local:8080/ - main page with form for image upload. 
http://file.upload.local:8080/image (POST request) - url for image saving on disk

Stored images can be accessed from browser using url
http://file.upload.local:8080/images/{image_name}

You can config some parameters using environment variables (file `.env.dev`).

RabbitMQ UI Url Address - `http://rabbit.client.local:15672`. Use `guest` as 
username and password.

## How to run tests

```shell
docker-compose run files_api sh -c "cd handler && go test -v ./..."
```