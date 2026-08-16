# Lesson 2: ヘルスチェックエンドポイントを実装する

## 今回学ぶこと

- コンテナオーケストレーション（Docker、Kubernetes）が「アプリが生きているか」をどう知るのか
- ヘルスチェック用エンドポイントの役割
- `HealthChecker` のようなチェック内容を差し替え可能にする設計

### なぜヘルスチェックが必要なのか

コンテナで動くアプリは、プロセスが起動していても「正常に動作しているとは限りません」。たとえば、アプリ自体は起動していてもDBへの接続が切れていて、リクエストを処理できない状態かもしれません。

Dockerの `HEALTHCHECK` 命令や、Kubernetesの liveness/readiness probe は、定期的にアプリの特定のエンドポイント（多くの場合 `/healthz` のようなパス）にHTTPリクエストを送り、そのレスポンスで正常性を判断します。

- 200 OKが返れば「正常」とみなし、トラフィックを送り続ける
- エラー系のステータス（503など）が返れば「異常」とみなし、そのコンテナへのトラフィックを止めたり、再起動したりする

つまり、**ヘルスチェック用エンドポイントを実装すること自体が、コンテナ運用の前提**になります。

### チェック内容を差し替え可能にする

「正常かどうか」の判定基準は、アプリによって異なります（DBに接続できるか、外部APIが呼べるか、など）。判定ロジックをハンドラーに直接書いてしまうと、テストのたびに本物のDBが必要になってしまいます。

そこで、判定処理そのものを関数の型として抽象化します。

```go
type HealthChecker func() error

func HealthHandler(check HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := check(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
```

`HealthHandler` は `HealthChecker` を受け取って `http.HandlerFunc` を**返す**関数です。実際のチェック内容（DB接続確認など）は、呼び出し側が `HealthChecker` として渡すので、`HealthHandler` 自身はチェックの中身を一切知りません。テストでは、常に成功/失敗を返すだけの単純な関数を渡せば、本物のDBなしにハンドラーの振る舞いを検証できます。

## 演習

`HealthHandler(check HealthChecker) http.HandlerFunc` を実装してください。

- 返す `http.HandlerFunc` の中で `check()` を呼ぶ
- エラーが `nil` なら、ステータスコード `200`（`http.StatusOK`）と本文 `"ok"` を書き込む
- エラーが返れば、ステータスコード `503`（`http.StatusServiceUnavailable`）と、そのエラーメッセージ（`err.Error()`）を本文として書き込む

## ヒント

- `http.HandlerFunc` は `func(http.ResponseWriter, *http.Request)` と同じ形なので、その形のクロージャを `return func(w http.ResponseWriter, r *http.Request) { ... }` として返します
- ステータスコードは `w.WriteHeader(...)` で、本文は `w.Write([]byte(...))` で書き込みます。**`WriteHeader` は `Write` より先に呼ぶ**必要があります
