# BigQuery
## `go install` する
[クイックスタート: クライアント ライブラリの使用](https://cloud.google.com/bigquery/docs/quickstarts/quickstart-client-libraries?hl=ja#client-libraries-install-go) を参考にする  

`go get cloud.google.com/go/bigquery`  
一応 `go mod tidy` を実行しておく  

1回では なんか依存関係が解決できなかったのか `cloud.google.com/go/bigquery` が使うことができなかったから  
再度 `go get cloud.google.com/go/bigquery` と `go mod tidy` を実行したら 使えるようになった  

BigQuery の sql は  FORM が ``` `テーブル名` ``` とバッククオートで囲まないといけない  
だから go で書くと 以下のようになる  
```go
`
SELECT
  name,
  count
FROM
  ` + "`babynames.names_2014`" + `
WHERE
  gender = 'M'
ORDER BY
  count DESC
LIMIT
  5
`
```
