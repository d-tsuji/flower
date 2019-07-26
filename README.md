# flower
this is task flow executor engine.

flowerはタスクをflowのように実行するツールです。
予めマスタに定義されたタスク定義に従ってタスクを実行します。
本ツール自体はRest Clientのみの役割を担い、業務ロジックはRest Serverで実行されることを想定しています。

## 特徴

- タスクの順序制御/流量制御が可能
  - プログラマブルにタスクを登録することができます。
  - JP1などの複雑になりがちなタスク定義をマスタでコントロールすることができます。
  - 同時実行数を制御することでシステムの負荷をコントロールできます。
  
- パラメータの引き継ぎが可能
  - 先行タスクが完了したときのレスポンスボディのパラメータを後続のリクエストボディのパラメータとして用いることができます。

## デモ

サンプルタスクとして、マスタ(ms_task)にタスクを登録してあります。
タスクを実行してみましょう。

サンプルでは以下のようなDAGの4つのタスクを順次実行します。

![sample](https://user-images.githubusercontent.com/24369487/60518099-cd90af80-9d1b-11e9-8068-44e5296ec495.PNG)

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

app    | 2019-07-26T07:53:18.507Z       INFO    app/watcher.go:12       Task Watching.
app    | 2019-07-26T07:53:38.510Z       INFO    app/watcher.go:12       Task Watching.
app    | 2019-07-26T07:53:38.511Z       INFO    app/watcher.go:20       Executable task found. Put channel. %v
app    | 2019-07-26T07:53:38.520Z       INFO    app/runner.go:14        Task starting.
app    | 2019-07-26T07:53:38.567Z       INFO    mock/jobRegister.go:45  Request received!
app    | 2019-07-26T07:53:40.569Z       INFO    app/runner.go:38
app    | 2019-07-26T07:53:40.580Z       INFO    app/runner.go:48        Task finished : %v
app    | 2019-07-26T07:53:40.581Z       INFO    app/watcher.go:12       Task Watching.
app    | 2019-07-26T07:53:40.583Z       INFO    app/watcher.go:20       Executable task found. Put channel. %v
app    | 2019-07-26T07:53:40.590Z       INFO    app/runner.go:14        Task starting.
app    | 2019-07-26T07:53:42.601Z       INFO    app/runner.go:38        {"parameters":{"executeHostname":"localhost","fromYM":"201801","toYM": "201905"}}
app    |
app    | 2019-07-26T07:53:42.620Z       INFO    app/runner.go:48        Task finished : %v
app    | 2019-07-26T07:53:42.620Z       INFO    app/watcher.go:12       Task Watching.
app    | 2019-07-26T07:53:42.622Z       INFO    app/watcher.go:20       Executable task found. Put channel. %v
app    | 2019-07-26T07:53:42.636Z       INFO    app/runner.go:14        Task starting.
app    | 2019-07-26T07:53:52.643Z       INFO    mock/jobRegister.go:54  Heavy Process start.
app    | 2019-07-26T07:53:52.644Z       INFO    app/runner.go:38
app    | 2019-07-26T07:53:52.664Z       INFO    app/runner.go:48        Task finished : %v
app    | 2019-07-26T07:53:52.664Z       INFO    app/watcher.go:12       Task Watching.
app    | 2019-07-26T07:53:52.666Z       INFO    app/watcher.go:20       Executable task found. Put channel. %v
app    | 2019-07-26T07:53:52.674Z       INFO    app/runner.go:14        Task starting.
app    | 2019-07-26T07:53:52.684Z       INFO    mock/jobRegister.go:45  Request received!
app    | 2019-07-26T07:53:54.686Z       INFO    app/runner.go:38
app    | 2019-07-26T07:53:54.698Z       INFO    app/runner.go:48        Task finished : %v
app    | 2019-07-26T07:53:54.698Z       INFO    app/watcher.go:12       Task Watching.
app    | 2019-07-26T07:53:58.521Z       INFO    app/watcher.go:12       Task Watching.

```
