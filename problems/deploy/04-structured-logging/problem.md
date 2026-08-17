# Lesson 4: 構造化ロギングでコンテナ環境に対応する

## 今回学ぶこと

- コンテナ環境では「ログファイル」ではなく「標準出力」にログを出す理由
- 人間が読みやすいログと、機械が解析しやすいログの違い
- `log/slog` パッケージによる構造化ロギング（JSON形式）

### コンテナのログは標準出力に出す

コンテナは「使い捨て」が前提です。コンテナが再起動・再作成されると、コンテナ内部に書いたログファイルは失われてしまいます。そのため、コンテナ化されたアプリは**標準出力（stdout）にログを出す**のが基本です。`docker logs` コマンドや、Kubernetesのログ収集の仕組みは、コンテナの標準出力を自動的に収集してくれます。

### なぜ「文章のログ」では困るのか

```
2024/01/15 10:23:41 user 42 placed order 1001 for 3 items
```

人間が読む分にはこれで十分ですが、大量のログを集約・検索・集計するツール（CloudWatch Logs、Datadog、Elasticsearchなど）にとっては、この形式は解析しづらいものです。「status=500のログだけ抽出する」ためには、文章から数値を正規表現で抜き出す必要が出てきます。

### JSON形式で構造化する

そこで、ログの各項目を最初からキーと値のペア（構造化データ）として出力します。

```json
{"time":"2024-01-15T10:23:41Z","level":"INFO","msg":"request handled","method":"POST","path":"/api/orders","status":201}
```

これなら、ログ収集ツール側で `status` フィールドだけを条件に検索したり、`method` ごとに集計したりするのが簡単になります。

### log/slog パッケージ

Go標準ライブラリの `log/slog`（Go 1.21〜）を使うと、この形式のログを簡単に出力できます。

```go
handler := slog.NewJSONHandler(os.Stdout, nil)
logger := slog.New(handler)

logger.Info("request handled", "method", "POST", "path", "/api/orders", "status", 201)
```

`Info` の第1引数はメッセージ文字列、それ以降は `key1, value1, key2, value2, ...` という形でキーと値のペアを交互に渡します（**属性(attribute)** と呼びます）。これにより、上記のようなJSON1行が出力されます。

`slog.NewJSONHandler` は第1引数に「出力先」（`io.Writer`）を取ります。本番では `os.Stdout` を渡しますが、テストでは `bytes.Buffer` のような別の `io.Writer` を渡せば、実際に出力された内容を検証できます。

## 演習

`NewRequestLogger` と `LogRequest` を実装してください。

**`NewRequestLogger(w io.Writer) *slog.Logger`**
- `slog.NewJSONHandler(w, nil)` でJSON形式のハンドラーを作る
- `slog.New(handler)` でLoggerを作って返す

**`LogRequest(logger *slog.Logger, method, path string, status int)`**
- `logger.Info(...)` を呼び、メッセージは `"request handled"`
- 属性として `"method"`, `"path"`, `"status"` というキーで、それぞれ引数の値を渡す

## ヒント

- `NewRequestLogger` は2行（ハンドラー作成 → Logger作成）、あるいは1行にまとめても書けます
- `LogRequest` は `logger.Info("request handled", "method", method, "path", path, "status", status)` の1行です
