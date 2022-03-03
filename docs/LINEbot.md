# LINE bot
## `go install` する
LINE official SDKs for the Messaging API. の [Go 言語 SDK](https://github.com/line/line-bot-sdk-go/blob/master/examples/echo_bot/server.go) を参考にする  

`go install github.com/line/line-bot-sdk-go/v7/linebot`  

>no required module provides package github.com/line/line-bot-sdk-go/v7/linebot; to add it:
>        go get github.com/line/line-bot-sdk-go/v7/linebot

って怒られたので `go get -u` でインストールする  
`go get: added github.com/line/line-bot-sdk-go/v7 v7.13.0` ってなった  

一応 `go mod tidy` を実行しておく  


1. LINE Developers でログイン  
2. Providers を選択  
3. Channel を選択  
  - Basic settings の Channel secret を取得  
  - Messaging API の Channel access token を取得  
  Secret Manager に追加  
4. Webhook URL を設定  
  `https://プロジェクトID.uw.r.appspot.com/callback` など  
5. Use webhook を オンにする  


`/callback-go`, `/callback-python` のようにちゃんとハンドリングすれば 1個の GAE で 2個の LINE bot が動かせた!  
