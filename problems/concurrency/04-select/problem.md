# Lesson 4: selectで複数のチャネルを扱う

## 今回学ぶこと

- `select` 文で複数のチャネルを同時に待つ方法
- 「先に届いた方を採用する」というよくあるパターン

### selectとは

`select` 文は、複数のチャネル操作のうち**最初に準備ができたもの**を実行する構文です。`switch` に似ていますが、各 `case` がチャネルの送受信になります。

```go
select {
case v := <-ch1:
	fmt.Println("ch1から受信:", v)
case v := <-ch2:
	fmt.Println("ch2から受信:", v)
}
```

`ch1` と `ch2` の両方に値が届く可能性がある状況で、`select` はどちらか先に準備できた方の `case` を実行します。両方ともまだ準備できていなければ、どちらかが準備できるまでブロックします。

### 「先に届いた方を採用する」パターン

複数の処理を並行に走らせて、一番早く終わったものだけを使いたい、というのはよくある要求です。次の例は、2つのゴルーチンに同じ計算をさせて、先に終わった方の結果を採用します。

```go
func firstEvenSquare(nums []int) int {
	ch := make(chan int)

	for _, n := range nums {
		go func() {
			if n%2 == 0 {
				ch <- n * n
			}
		}()
	}

	return <-ch // 最初にchannelへ届いた値を受け取る
}
```

このレッスンの演習では、あらかじめ用意された2つの受信専用チャネル（`<-chan int`）を受け取り、`select` でどちらか先に値が届いた方を返す、というシンプルな形で `select` の基本を確認します。

## 演習

`firstResult` 関数を実装してください。2つの受信専用チャネル `ch1`, `ch2` を受け取り、`select` を使って**どちらか先に値が届いた方**を返してください。

## ヒント

- 関数のシグネチャはすでに `func firstResult(ch1, ch2 <-chan int) int` と定義されています
- `select { case v := <-ch1: return v; case v := <-ch2: return v }` の形になります
