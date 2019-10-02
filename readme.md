
# 概要
isucon用Makefile

# 事前準備

## slackcatの使用準備
slackcat用のトークンを準備します。slackcatの利用が初めてであれば、http://slackcat.chat/の右上、Add to Slackからxoxpで始まるトークンを取得しておきます。

## 検証用Dockerの起動
以下でdockerに入れればOK

```
docker-compose build
docker-compose up -d
make ssh
```
# 使用方法

## make setup
一番最初に実行するツール類の設定です。

以下のアプリをインストールします。
- percona-toolkit
- kataribe
- myprofiler
- slackcat

slackcatインストール時には以下を聞かれるので入力してください。
nickname for team:t (複数チームで使わないのであればなんでも良いの)
token issued:xoxp-XXXXXXX（事前準備で準備したslackcatのトークン）

その後.slackcatファイルを変更して、default_channnelに投稿したいチャンネルを設定する。

## pprofの設定

