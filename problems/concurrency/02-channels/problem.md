# Lesson 2: channelでやり取りする

## 今回学ぶこと

- `channel` の基本（`make`、送信 `<-`、受信 `<-`）
- `close` でチャネルを閉じる意味
- `range` でチャネルから値を受け取り続ける方法

### channelとは

`channel` は、ゴルーチン同士が値をやり取りするための「パイプ」です。`make(chan 型)` で作成し、`<-` 演算子で送受信します。

```go
ch := make(chan int)

go func() {
	ch <- 42 // 送信
}()

v := <-ch // 受信（値が送られてくるまでブロックする）
fmt.Println(v) // 42
```

Lesson 1で使った `sync.WaitGroup` は「完了を待つだけ」でしたが、channelは「完了を待ちながら値も受け取れる」という点が違います。

### closeとrange

送信側がすべて送り終わったら `close(ch)` でチャネルを閉じます。受信側は `for v := range ch` と書くことで、チャネルが閉じられるまで値を受け取り続けられます。

```go
func collectDoubled(nums []int) []int {
	ch := make(chan int)

	go func() {
		defer close(ch) // 送信が終わったら必ず閉じる
		for _, n := range nums {
			ch <- n * 2
		}
	}()

	var result []int
	for v := range ch {
		result = append(result, v)
	}
	return result
}
```

このコードでは、1つのゴルーチンが順番に値を送信しているため、`result` の順序は `nums` と同じ順序になります。

`close` を忘れると、受信側の `for range` が永遠にブロックし続けてしまう（プログラムがフリーズする）ので注意してください。

## 演習

`collectDoubled` と同じ考え方で、`collectSquares` 関数を実装してください。`nums` の各要素を**2乗**した値をchannel経由で送り、呼び出し側で `range` を使って集めてください。

## ヒント

- ゴルーチンの中で `defer close(ch)` を書くのを忘れないでください
- 受信側は `for v := range ch { result = append(result, v) }` の形になります
