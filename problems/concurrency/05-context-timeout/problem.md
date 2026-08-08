# Lesson 5: contextでタイムアウト制御する

## 今回学ぶこと

- `context.WithTimeout` で「これ以上待たない」時間を設定する方法
- `ctx.Done()` と `select` を組み合わせたタイムアウト処理
- `ctx.Err()` でタイムアウトの理由を取得する方法

### なぜタイムアウトが必要か

外部のAPIを呼んだり、重い処理を待ったりするとき、相手が固まってしまうと自分のプログラムまで無限に待ち続けてしまいます。`context` パッケージは、「一定時間で諦める」「呼び出し元がキャンセルしたら止める」といった制御を統一的に扱う仕組みです。

### context.WithTimeout

`context.WithTimeout(親ctx, 時間)` は、指定した時間が経過すると自動的に「完了」扱いになる `context.Context` を返します。完了したかどうかは `ctx.Done()` が返すchannelで判定でき、Lesson 4で学んだ `select` と組み合わせて使います。

```go
func runWithTimeout(ctx context.Context, work func() int) (int, error) {
	resultCh := make(chan int, 1)
	go func() {
		resultCh <- work()
	}()

	select {
	case result := <-resultCh:
		return result, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}
```

- `work` を別ゴルーチンで実行し、結果を `resultCh` に送る
- `select` で「`work` の結果が届く」のと「`ctx` がタイムアウトする」のを同時に待つ
- タイムアウトが先に来たら、`ctx.Err()`（`context.DeadlineExceeded` など）を返す

呼び出す側は次のように `context.WithTimeout` を使います。`cancel` 関数は、タイムアウトを待たずに正常終了した場合でもリソースを解放するために、`defer cancel()` で必ず呼び出してください。

```go
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

result, err := runWithTimeout(ctx, someSlowFunc)
```

> **注意**: `runWithTimeout` の中でタイムアウトしても、`work` を実行しているゴルーチン自体は止まりません（Goには他のゴルーチンを強制終了する仕組みがないため）。実務では `work` 側もコンテキストを受け取って自発的に中断する作りにするのが望ましいですが、このレッスンでは `select` によるタイムアウト検知の基本に絞って学びます。

## 演習

`runWithTimeout` 関数はすでに実装されています。あなたが実装するのは、それを使う `sumWithDeadline` 関数です。

`sumWithDeadline(nums []int, timeout time.Duration) (int, error)` は、`nums` の合計を計算する処理を `timeout` 時間以内に完了させたい、という関数です。

1. `context.WithTimeout(context.Background(), timeout)` でタイムアウト付きコンテキストを作る（`defer cancel()` を忘れずに）
2. `runWithTimeout(ctx, work)` を呼び出す。`work` は `nums` の合計を計算して返す関数（クロージャ）にする
3. `runWithTimeout` の戻り値をそのまま返す

## ヒント

- `work := func() int { sum := 0; for _, n := range nums { sum += n }; return sum }` のようにクロージャを作れます
- 合計の計算自体は一瞬で終わるので、このテストではタイムアウトは発生しない想定です（タイムアウトが実際に発動するケースは `runWithTimeout` 自体のテストで別途確認します）
