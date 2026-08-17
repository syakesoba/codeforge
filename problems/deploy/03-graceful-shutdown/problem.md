# Lesson 3: グレースフルシャットダウンを実装する

## 今回学ぶこと

- `docker stop` がコンテナに何を送っているのか
- 処理中のリクエストを中断させずに終了する「グレースフルシャットダウン」
- `context.WithTimeout` と `Shutdown(ctx)` パターン

### docker stop は何をしているのか

`docker stop` を実行すると、Dockerはまずコンテナ内のプロセスに **SIGTERM**（「そろそろ終了してください」という合図）を送ります。プロセスが一定時間（デフォルト10秒）以内に自分から終了しなければ、Dockerは強制的に **SIGKILL** で終了させます。

もしアプリがSIGTERMを無視していたら、常にSIGKILLによる強制終了になり、その瞬間に処理中だったリクエストは中断され、クライアントにはエラーが返ってしまいます。デプロイのたびに（新しいバージョンに切り替えるたびに）これが起きると、ユーザー体験を損ないます。

### グレースフルシャットダウン

そこで、SIGTERMを受け取ったら「新しいリクエストの受付は止めつつ、処理中のリクエストは最後まで終わらせてから終了する」という挙動を実装します。これを**グレースフルシャットダウン（graceful shutdown）**と呼びます。

Goの `*http.Server` には、まさにこのための `Shutdown(ctx context.Context) error` メソッドが用意されています。

```go
srv := &http.Server{Addr: ":8080"}

// (SIGTERMを受け取ったら…)
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
	log.Printf("shutdown error: %v", err)
}
```

`Shutdown` は新しい接続の受付を即座に停止し、処理中のリクエストが終わるのを待ってから返ります。ただし、`ctx` に指定した時間を過ぎても処理中のリクエストが終わらない場合は、そこで強制的に打ち切り、`ctx` のエラー（`context.DeadlineExceeded`）を返します。「無期限に待ち続けてコンテナがいつまでも終了しない」事態を防ぐための安全弁です。

### なぜタイムアウトの管理を自分で書かないのか

`Shutdown` はすでに `ctx` のキャンセル・タイムアウトを内部で正しく扱う設計になっています。そのため、呼び出し側が「タイムアウトしたかどうか」を別途タイマーなどで判定する必要はなく、`context.WithTimeout` で期限付きのcontextを渡し、その戻り値をそのまま返すだけで十分です。

## 演習

`GracefulShutdown(shutdowner Shutdowner, timeout time.Duration) error` を実装してください。

- `context.WithTimeout(context.Background(), timeout)` で期限付きのcontextを作る
- `defer cancel()` を必ず呼ぶ（contextリソースの解放漏れを防ぐため）
- `shutdowner.Shutdown(ctx)` を呼び、その戻り値をそのまま返す

## ヒント

- `Shutdowner` は `*http.Server` が実際に満たしているインターフェースです。テストでは、この演習のようにモック実装を使って「タイムアウトした場合」の挙動も検証できます
- `context.WithTimeout` は `(context.Context, context.CancelFunc)` のペアを返します。`cancel` を呼び忘れると、たとえ処理が早く終わってもcontextの内部リソースが解放されないままになるため、必ず `defer cancel()` してください
