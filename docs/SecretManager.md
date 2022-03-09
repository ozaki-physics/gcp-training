# Secret Manager
## `go install` する
[Secret Manager client libraries](https://cloud.google.com/secret-manager/docs/reference/libraries#create-service-account-console) を参考にする  

`go install cloud.google.com/go/secretmanager/apiv1`  

>no required module provides package cloud.google.com/go/secretmanager/apiv1; to add it:
>        go get cloud.google.com/go/secretmanager/apiv1

って怒られたので `go get` でインストールする  
`go get: added cloud.google.com/go/secretmanager v1.2.0` ってなった  

`go install google.golang.org/genproto/googleapis/cloud/secretmanager/v1`  
を実行したら何も表示されず終わった  
一応 `go get google.golang.org/genproto/googleapis/cloud/secretmanager/v1` しておく  
少しダウンロードされた  

一応 `go mod tidy` を実行しておく  

GAE に Secret Manager の read 権限渡すの忘れて 1時間ぐらい悩んだ  
