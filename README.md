# dajaredetector

自分がフォローしているアカウントの投稿にダジャレが含まれていたら指摘する mastodon bot です。

## 機能

- [x] Streaming API を用いてホームタイムラインを監視（[go-mastodon](https://github.com/mattn/go-mastodon)を利用）
- [x] 投稿本文にダジャレが含まれているか評価（[dajarep](https://github.com/kurehajime/dajarep)を利用）
- [x] ダジャレが含まれていたらリプライを送る
- [ ] フォローしてきたアカウントをフォローする
- [ ] 特定条件でフォロー解除

## 使い方

1. `.env.example` を参考に `.env` ファイルを作成（アクセストークンが必要です）
1. `go run dajaredetector.go`  
   たぶん `go build` とかで実行ファイルを作って `./dajaredetector` で良いと思います。
