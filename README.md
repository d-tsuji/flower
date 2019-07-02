# flower
this is task flow executor engine.

flowerはタスクをflowのように実行するツールです。
予めマスタに定義されたタスク定義に従ってタスクを実行します。
本ツール自体はRest Clientのみの役割を担い、業務ロジックはRest Serverで実行されることを想定しています。

## 特徴

- タスクの順序制御が可能
  - プログラマブルにタスクを登録することができます。
  - JP1などの複雑になりがちなタスク定義をマスタでコントロールすることができます。

- タスクの流量制御が可能
  - 同時実行数を制御することでシステムの負荷をコントロールできます。
  
- パラメータの引き継ぎが可能
  - 先行タスクが完了したときのレスポンスボディのパラメータを後続のリクエストボディのパラメータとして用いることができます。

## デモ

サンプルタスクとして、マスタ(ms_task)にタスクを登録してあります。
タスクを実行してみましょう。

GitHubからソースを取得
```
$ git clone https://github.com/d-tsuji/flower.git
```

ディレクトリに移動
```
$ cd flower
```

docker-composeでアプリケーションを起動
```
$ docker-compose up -d
```

httpクライアントでタスクを登録
```
$ curl http://localhost:8021/register?taskId=hello
```

コンテナのログを確認
```
$ docker-compose logs

app    | 2019/07/02 13:31:58 main() watching...
app    | 2019/07/02 13:31:58 Task Watching...
app    | 2019/07/02 13:32:08 main() watching...
app    | 2019/07/02 13:32:08 Task Watching...
app    | 2019/07/02 13:32:08 Executable task found. Put channel. {ff88e3d7-f746-4493-bb64-1573d2eb05c4 hello 1 }
app    | 2019/07/02 13:32:08 Task starting... : {ff88e3d7-f746-4493-bb64-1573d2eb05c4 hello 1 }
app    | 2019/07/02 13:32:08 Request received!
app    | 2019/07/02 13:32:10 Hello world!
app    |
app    | 2019/07/02 13:32:10 Task finished : {ff88e3d7-f746-4493-bb64-1573d2eb05c4 hello 1 }
app    | 2019/07/02 13:32:10 Task Watching...
app    | 2019/07/02 13:32:10 Executable task found. Put channel. {ff88e3d7-f746-4493-bb64-1573d2eb05c4 hello 2 Hello world!
app    | }
app    | 2019/07/02 13:32:10 Task starting... : {ff88e3d7-f746-4493-bb64-1573d2eb05c4 hello 2 Hello world!
app    | }
app    | 2019/07/02 13:32:12 [Method] POST
app    | 2019/07/02 13:32:12 [Header] User-Agent: Go-http-client/1.1
app    | 2019/07/02 13:32:12 [Header] Content-Length: 81
app    | 2019/07/02 13:32:12 [Header] Accept-Encoding: gzip
app    | 2019/07/02 13:32:12 {"parameters":{"executeHostname":"localhost","fromYM":"201801","toYM": "201905"}}
app    |
app    | 2019/07/02 13:32:12 Task finished : {ff88e3d7-f746-4493-bb64-1573d2eb05c4 hello 2 Hello world!
app    | }
app    | 2019/07/02 13:32:12 Task Watching...
app    | 2019/07/02 13:32:12 Executable task found. Put channel. {ff88e3d7-f746-4493-bb64-1573d2eb05c4 hello 3 {"parameters":{"executeHostname":"localhost","fromYM":"201801","toYM": "201905"}}
app    | }
app    | 2019/07/02 13:32:12 Task starting... : {ff88e3d7-f746-4493-bb64-1573d2eb05c4 hello 3 {"parameters":{"executeHostname":"localhost","fromYM":"201801","toYM": "201905"}}
app    | }
app    | 2019/07/02 13:32:18 main() watching...
app    | 2019/07/02 13:32:18 Task Watching...
app    | 2019/07/02 13:32:22 Heavy Process start.
app    | 2019/07/02 13:32:22 Heavy Process finish.
app    |
app    | 2019/07/02 13:32:22 Task finished : {ff88e3d7-f746-4493-bb64-1573d2eb05c4 hello 3 {"parameters":{"executeHostname":"localhost","fromYM":"201801","toYM": "201905"}}
app    | }
app    | 2019/07/02 13:32:22 Task Watching...
app    | 2019/07/02 13:32:22 Executable task found. Put channel. {ff88e3d7-f746-4493-bb64-1573d2eb05c4 hello 4 Heavy Process finish.
app    | }
app    | 2019/07/02 13:32:22 Task starting... : {ff88e3d7-f746-4493-bb64-1573d2eb05c4 hello 4 Heavy Process finish.
app    | }
app    | 2019/07/02 13:32:22 Request received!
app    | 2019/07/02 13:32:24 Hello world!
app    |
app    | 2019/07/02 13:32:24 Task finished : {ff88e3d7-f746-4493-bb64-1573d2eb05c4 hello 4 Heavy Process finish.
app    | }
app    | 2019/07/02 13:32:24 Task Watching...
app    | 2019/07/02 13:32:28 main() watching...
app    | 2019/07/02 13:32:28 Task Watching...

```
