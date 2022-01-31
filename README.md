# go-sns-api

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

検証が終わったら、`docker-compose down`

## エンドポイント一覧
| パス | HTTPメソッド | 概要 | トークン保護
| :-- | :-- | :-- | :--
| /api/v1/users | POST | ユーザー作成 | NO
| /api/v1/auth | POST | JWTトークンを返す | NO
| /api/v1/posts/:id | GET | 投稿取得 | NO
| /api/v1/posts | POST | 投稿作成 | YES
| /api/v1/posts/:id | PATCH | 投稿編集 | YES
| /api/v1/posts/:id | DELETE | 投稿削除 | YES
