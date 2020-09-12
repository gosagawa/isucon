
# 概要
isucon用の秘伝のMakefileとそれをDockerで確認できる環境です。
Makefileの動作確認や、各種ツールの利用確認に役立ててください。

# 事前準備

## slackcatの使用準備
slackcat用のトークンを準備します。slackcatの利用が初めてであれば、http://slackcat.chat/ の右上、Add to Slackからxoxpで始まるトークンを取得しておきます。

## 検証用Dockerの起動
以下を実行し、コンテナにsshログインできる事を確認してください。

```
docker-compose build
docker-compose up -d
make ssh
```
# 使用方法

以下、一通りの手順を実行後にMakefileで定義された内容を実施してください。

## make setup
一番最初に実行するツール類の設定です。

以下のアプリをインストールします。
- percona-toolkit
- kataribe
- myprofiler
- slackcat

slackcatインストール時には以下を聞かれるので入力してください。
```
nickname for team:t (複数チームで使わないのであればなんでも良いが何か入れる必要あり)
token issued:xoxp-XXXXXXX（事前準備で準備したslackcatのトークン）
```

その後.slackcatファイルを変更して、default_channnelに投稿したいチャンネルを設定する。

## pprofの設定
このレポジトリ内では追加ずみですが、実際の競技ではpprofの設定を追加します。

- importに追加

```
       _ "net/http/pprof"
```

- main関数に追加

```
       go func() {
               log.Println(http.ListenAndServe("localhost:6060", nil))
       }()
```

## TODOの処理

Makeファイル内でTODOとなっている部分は個別設定が必要であるはずなので、適宜設定を行ってください。このレポジトリを試すだけであれば設定不要です。

# 謝辞

このレポジトリ、およびMakefileを作るにあたり以下の記事とMakefileを参考にさせていただいています。

[ISUCON9予選1日目で最高スコアを出しました](https://to-hutohu.com/2019/09/09/isucon9-qual/)

[Makefile](https://github.com/tohutohu/isucon9/blob/master/Makefile)
