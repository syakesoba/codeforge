# Lesson 3: sync.Mutexで排他制御する

## 今回学ぶこと

- 複数のゴルーチンが同じデータへ同時に書き込むと何が起きるか
- `sync.Mutex` で「同時に1つのゴルーチンだけがアクセスできる区間」を作る方法

## 解説

### 同時書き込みの問題

Lesson 1では、各ゴルーチンが配列の**別々のインデックス**に書き込んでいたので安全でした。しかし、複数のゴルーチンが**同じ変数**（例えば同じmapや同じ合計値）に同時に書き込もうとすると、データが壊れます。

Goのmapは特にわかりやすく、保護なしで複数のゴルーチンから同時に書き込むと、次のような**ランタイムエラーでプログラムごと異常終了**します。

```
fatal error: concurrent map writes
```

これは実行するたびに起きるとは限らない厄介なバグです（タイミング次第で偶然動いてしまうこともある）。

```go
func buildSquareMapUnsafe(nums []int) map[int]int {
	result := make(map[int]int)

	var wg sync.WaitGroup
	for _, n := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result[n] = n * n // 複数ゴルーチンが同時に書き込む可能性がある → 危険
		}()
	}
	wg.Wait()

	return result
}
```

### sync.Mutexで守る

`sync.Mutex` は「鍵」のようなものです。`mu.Lock()` した区間は、他のゴルーチンが同時に `mu.Lock()` しようとしてもブロックされ、`mu.Unlock()` されるまで待たされます。これにより「同時に1つのゴルーチンだけが実行できる区間（クリティカルセクション）」を作れます。

```go
func buildSquareMap(nums []int) map[int]int {
	result := make(map[int]int)
	var mu sync.Mutex

	var wg sync.WaitGroup
	for _, n := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			square := n * n

			mu.Lock()
			result[n] = square // ここは同時に1ゴルーチンしか実行できない
			mu.Unlock()
		}()
	}
	wg.Wait()

	return result
}
```

計算そのもの（`n * n`）はロックの外で行い、mapへの書き込みだけをロックで挟むのがポイントです。ロックする範囲は必要最小限にするのが定石です。

## 演習

`buildCubeMap` 関数を実装してください。`nums` の各要素をキーとして、その**3乗**の値を共有のmapに格納します。複数のゴルーチンから安全に書き込めるよう、`sync.Mutex` で保護してください。

## ヒント

- `mu.Lock()` と `mu.Unlock()` は必ずペアで呼び出してください（`defer` は使えません。ロック区間をループの中の一部分だけに限定したいので、Lock直後にdeferすると関数全体が終わるまでUnlockされず意味がなくなってしまいます）
- ロックするのは `result[n] = cube` の行だけで十分です
