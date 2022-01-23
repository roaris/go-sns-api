FROM golang:alpine

RUN apk update && apk add git

RUN mkdir /app

WORKDIR /app

COPY /app /app
