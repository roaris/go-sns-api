FROM golang:alpine

RUN apk update && apk add git

RUN mkdir /app

WORKDIR /app

COPY ./app .

RUN go get -u github.com/cosmtrek/air
CMD ["air", "-c", ".air.toml"]
