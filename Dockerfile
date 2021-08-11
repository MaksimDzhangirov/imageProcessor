FROM golang:alpine as builder

WORKDIR /go/src/app

RUN apk update && apk add curl gcc pkgconfig build-base glib-dev expat-dev libpng-dev

RUN curl -s https://raw.githubusercontent.com/h2non/bimg/master/preinstall.sh | sh -

ENV GO111MODULE=on

COPY filesAPI /go/src/filesAPI
COPY processingAPI .

RUN go mod download

RUN go build -o ./run .

FROM alpine:latest
WORKDIR /root/

#Copy executable from builder
COPY --from=builder /go/src/app/run .

CMD ["./run"]