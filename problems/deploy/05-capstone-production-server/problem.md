# Lesson 5（道場）: 本番向けのHTTPサーバーを組み立てる

## 今回学ぶこと

- これまでのレッスンで学んだ要素を組み合わせて、本番運用を意識したサーバーの骨組みを作る
- 設定（Lesson 1）に応じてルーティングを切り替える設計
- （まとめ）Goアプリを実際にDocker化するときの流れ

### これまで学んだことの総仕上げ

このレッスンは「道場」＝総仕上げ問題であると同時に、「デプロイ・Docker化」コース全体、そしてCodeForgeのGo応用編全コースのまとめでもあります。

- 環境変数による設定（Lesson 1）: `Config.Debug` によって、登録するルートを変える
- ヘルスチェック（Lesson 2）: `HealthHandler` をそのまま組み込む
- 依存性の注入（実践的なアーキテクチャコース）: `HealthChecker` や `routes` を外から受け取ることで、`NewApp` 自体はテストしやすいまま保たれる

```go
func NewApp(cfg Config, check HealthChecker, routes map[string]http.HandlerFunc) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", HealthHandler(check))
	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}
	if cfg.Debug {
		mux.HandleFunc("/debug", DebugHandler)
	}
	return mux
}
```

`cfg.Debug` が `false`（本番のデフォルト）であれば `/debug` は一切登録されず、そのパスへのアクセスは404になります。デバッグ用の情報（内部状態、環境変数の値など）を本番で誤って公開してしまう事故を、コードの構造そのもので防いでいます。

## 演習

`Config`、`HealthHandler`、`DebugHandler` は実装済みです。`NewApp` を実装してください。

**`NewApp(cfg Config, check HealthChecker, routes map[string]http.HandlerFunc) *http.ServeMux`**
1. `http.NewServeMux()` で空のmuxを作る
2. `"/healthz"` に `HealthHandler(check)` を登録する
3. `routes` の各エントリを、そのパスとハンドラーでmuxに登録する
4. `cfg.Debug` が `true` の場合のみ、`"/debug"` に `DebugHandler` を登録する
5. 組み立てたmuxを返す

## ヒント

- `mux.HandleFunc(path, handler)` で登録します。`routes` は `map[string]http.HandlerFunc` なので、for-rangeでキーと値を両方取り出せます
- Debugルートの登録は、他のルート登録がすべて終わった後の `if` 文1つで十分です

---

## コースのまとめ: Goアプリを実際にDocker化する

ここまでのレッスンで、Dockerでの運用を意識したGoコードの書き方（環境変数設定・ヘルスチェック・グレースフルシャットダウン・構造化ロギング）を学んできました。最後に、実際にこれらを組み込んだGoアプリをコンテナ化する際の典型的な `Dockerfile` を紹介します。

```dockerfile
# --- ビルド用ステージ ---
FROM golang:1.22 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# CGO無効・静的リンクにすることで、後段の軽量イメージでも単体で動くバイナリになる
RUN CGO_ENABLED=0 GOOS=linux go build -o /app ./cmd/server

# --- 実行用ステージ ---
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app /app
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s CMD ["/app", "-healthcheck"]
ENTRYPOINT ["/app"]
```

- **マルチステージビルド**: 1つ目のステージ（`builder`）でコンパイルし、2つ目のステージでは**コンパイル済みのバイナリだけ**をコピーします。Goコンパイラや依存パッケージのソースコードといったビルド時にしか使わないものを最終イメージに含めないため、イメージサイズが数百MB〜数GB→数十MB程度まで小さくなります
- **`HEALTHCHECK`**: Lesson 2で学んだヘルスチェックエンドポイントを、Docker自身が定期的に呼び出す設定です
- **環境変数**: `docker run -e PORT=3000 -e DEBUG=true myapp` のように、Lesson 1で学んだ `LoadConfig` が読み取る環境変数を起動時に注入します
- **`docker stop`**: Lesson 3で学んだグレースフルシャットダウンが、SIGTERMを受け取った時点で効いてきます
- **ログ収集**: `ENTRYPOINT` で起動したプロセスの標準出力は自動的に `docker logs` で見られるので、Lesson 4のJSON構造化ログがそのまま活用できます

これで、CodeForgeのGo応用編で扱った内容は一通り完了です。お疲れさまでした。
