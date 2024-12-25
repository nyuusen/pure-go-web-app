# pure-go-web-app

[伸び悩んでいる3年目Webエンジニアのための、Python Webアプリケーション自作入門](https://zenn.dev/bigen1925/books/introduction-to-web-application-with-python)をGoで実装してみるリポジトリです。

## (ローカル)server起動

```
go run ./server/main.go
```

※listen状態で接続を待ち受け続ける

## Docker起動(apache)


```
docker compose up -d
```

## client起動

```
go run ./client/main.go -port=8080
```

