# go-sns-api

[![golangci-lint](https://github.com/roaris/go-sns-api/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/roaris/go-sns-api/actions/workflows/golangci-lint.yml)
[![test](https://github.com/roaris/go-sns-api/actions/workflows/test.yml/badge.svg)](https://github.com/roaris/go-sns-api/actions/workflows/test.yml)

## 動作環境
```
$ docker -v
Docker version 20.10.12, build e91ed57
```

## 検証環境(docker-compose)
```
git clone https://github.com/roaris/go-sns-api
cd go-sns-api
docker-compose build
docker-compose up
```

APIサーバー: [localhost:8080](http://localhost:8080)

Swagger: [localhost:8081](http://localhost:8081)

ハンドラのテスト: `docker-compose exec app go test ./handlers`

検証が終わったら、`docker-compose down`

## 本番環境
デプロイ先: https://go-sns-api.herokuapp.com/
```
$ curl https://go-sns-api.herokuapp.com/ping
pong
```

## エンドポイント一覧
| パス | HTTPメソッド | 概要
| :-- | :-- | :--
| /api/v1/users | POST | ユーザー作成
| /api/v1/auth | POST | JWTトークンを返す
| /api/v1/users/me | GET | ログイン中のユーザー取得
| /api/v1/users/me | PATCH | ログイン中のユーザー情報更新
| /api/v1/users/me/followees | POST | ユーザーのフォロー
| /api/v1/users/me/followees/:id | DELETE | ユーザーのアンフォロー
| /api/v1/users/:id/followees | GET | ユーザーのフォロイー一覧
| /api/v1/users/:id/followers | GET | ユーザーのフォロワー一覧
| /api/v1/posts/:id | GET | 投稿取得
| /api/v1/posts | GET | タイムライン取得
| /api/v1/posts | POST | 投稿作成
| /api/v1/posts/:id | PATCH | 投稿編集
| /api/v1/posts/:id | DELETE | 投稿削除
| /api/v1/posts/:id/likes | POST | いいねする
| /api/v1/posts/:id/likes | DELETE | いいね取り消し
