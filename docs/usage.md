# 使い方 Usage
## Remote - Containers
VS Code の拡張機能 Remote - Containers(識別子: ms-vscode-remote.remote-containers) を使って開発する  
コンテナ内で VS Code を起動し go 言語のための VS Code の拡張機能 Go(識別子: golang.go) を使う  
[golang.go](https://marketplace.visualstudio.com/items?itemName=golang.Go)  

golang.go は 様々なモジュールをインストールする必要がある  
特に gopls が go のコード補完してくれる便利なライブラリで使いたい  
ただ そのモジュールを Docker image に含めないようにするため  
コンテナを起動した後の コンテナ内 VS Code で変更を加える  
.devcontainer/devcontainer.json で gopls だけ install しているが VSCode でコンテナに入った時は認識してくれないから  
install が終わったら gopls だけ VSCode の通知から install にすればよい  

コンテナ内の git では日本語が使えないため コミットするときは ローカルの git bash 等を使う  
```bash
$ docker-compose build
$ docker-compose up -d
# VS Code より Remote - Containers で接続する
# VS Code の通知(golang.go)より gopls だけ install をする

# 基本は VS Code 内のターミナルで良いが ローカルの PowerShell からアクセスしたくなった場合
$ docker-compose exec gcp-training bash
# 初めて開発を始めるときは go mod を生成する
/go/src/github.com/ozaki-physics/gcp-training# go mod init $REPOSITORY
# 終えるとき
# VS Code より Remote - Containers で接続をやめる
$ docker-compose down
```

ちなみに image の時点で go build は済んでおり  
image から直接 run または docker-compose.yml の command をコメントアウトで確認できる  
```bash
$ docker container run --rm -d -p 8080:8080 --name go_gcp_training go1.16:gcp_training
$ docker container stop go_gcp_training
```

### 外部モジュールのバージョンアップ
例として github.com/gin-gonic/gin をバージョンアップする  
1. コンテナにアタッチする
2. モジュールのバージョンアップして go.mod を更新する
3. コンテナを削除してもバージョンアップが反映されるように docker image を作り直す

```bash
$ docker-compose up -d
$ docker-compose exec gcp_training bash

# モジュールのバージョンアップ
/go/src/github.com/ozaki-physics/gcp-training# go install -v github.com/gin-gonic/gin
# 不要モジュールの削除
/go/src/github.com/ozaki-physics/gcp-training# go mod tidy -v

$ docker-compose down
$ docker image rm go1.16:gcp_training
# docker image の作り直し
$ docker-compose build
```

## gcloud docker
[GCP Container Repository](https://console.cloud.google.com/gcr/images/google.com:cloudsdktool/GLOBAL/cloud-sdk?authuser=9)をもとに pull する image を決める  
```bash
$ docker pull gcr.io/google.com/cloudsdktool/cloud-sdk:372.0.0
# カレントディレクトリの内容を全部マウントする
$ docker run -it --rm -v "$(pwd):/gcp-training" gcr.io/google.com/cloudsdktool/cloud-sdk:372.0.0 bash
$ cd gcp-training
$ gcloud init
You must log in to continue. Would you like to log in (Y/n)? Yでエンター
# URL が表示されるから ブラウザに入力して Google アカウントを選択
# ワンタイムキーみたいなのが表示されるから console に入力
# GCP のプロジェクトを選ぶ
$ gcloud app deploy
# リージョンは us-west1-a
```
