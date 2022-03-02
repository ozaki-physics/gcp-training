# gcp-training
Google Cloud Platform で様々なことをやってみる練習リポジトリ

## 環境 Requirement
- Docker

## 使い方 Usage
[詳細](./docs/usage.md)
1. VS Code の拡張機能 Remote - Containers で開発する  
gopls だけ VSCode の通知から install する  
2. 開発が終わったら Remote - Containers を閉じる  
3. gcloud が使えるコンテナを開発したコードをマウントしながら起動して アタッチする  
4. 開発したコードをデプロイする  

## 試したこと
- Secret Manager から値を取得する
