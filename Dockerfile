FROM golang:alpine

RUN apk update && apk add git

# ワーキングディレクトリの設定
WORKDIR /app

# ホストのファイルをコンテナの/appにコピーする
COPY * ./

# airを使ってホットリロード
RUN go get -u github.com/cosmtrek/air
CMD ["air", "-c", ".air.toml"]
