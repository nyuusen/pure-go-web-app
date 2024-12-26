# pure-go-web-app

## これは何？

[伸び悩んでいる3年目Webエンジニアのための、Python Webアプリケーション自作入門](https://zenn.dev/bigen1925/books/introduction-to-web-application-with-python)をGoで実装してみたリポジトリです。

一応、POSTリクエストを受け取る所まで実装しております。  
時間の関係で、クエリストリングとテンプレートエンジンとCookieとリファクタリングはスキップしちゃいました。  
コードが鬼汚いのはご愛嬌ということで...。

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

