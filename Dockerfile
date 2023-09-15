FROM golang:1.21.1-alpine3.18

RUN apk update && apk add git

# ワーキングディレクトリの設定
WORKDIR /app

# ホストのファイルをコンテナの/appにコピーする
COPY * ./

# airを使ってホットリロード
RUN go install github.com/cosmtrek/air@v1.45.0
CMD ["air", "-c", ".air.toml"]
