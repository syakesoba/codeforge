# Lesson 6（道場）: ワーカープールを作ろう

## 今回学ぶこと

このレッスンは「道場」、つまりCourse「ゴルーチン・並行処理」のまとめ問題です。Lesson 1〜4で学んだ内容を組み合わせて、**ワーカープール**という実務でもよく使われるパターンを実装します。

- Lesson 1: ゴルーチンと `sync.WaitGroup`
- Lesson 2: `channel` での値のやり取り
- Lesson 4: `select`（今回は直接使いませんが、考え方はここでも活きています）

## 解説: ワーカープールとは

大量のジョブ（仕事）を処理したいとき、ジョブの数だけゴルーチンを起動すると、数が多すぎてリソースを使いすぎることがあります。そこで、**あらかじめ決まった数のワーカー（作業員ゴルーチン）**を用意し、ジョブを1つのチャネル（ジョブキュー）に流し込んで、ワーカーたちがそこから順番に取り出して処理する、という設計がよく使われます。

```
jobs channel:    [1] [2] [3] [4] [5] ...
                   ↓   ↓   ↓
              worker1 worker2 worker3  ← 決まった数のワーカーが並行に処理
                   ↓   ↓   ↓
results channel: [1²] [4²] [9²] ...
```

### 実装の流れ

1. **jobsチャネル**を作り、そこにジョブを流し込む
2. **workerCount個のワーカーゴルーチン**を起動する。各ワーカーは `for j := range jobsCh` でジョブを受け取り続け、処理結果を **resultsチャネル**に送る
3. 全ワーカーが処理を終えたら resultsチャネルを閉じる必要があるが、そのタイミングは「全ワーカーの完了」なので、ここで `sync.WaitGroup` が活躍する
4. 呼び出し側は resultsチャネルを `range` で受け取り、スライスに集める

```go
func processJobs(jobs []int, workerCount int, process func(int) int) []int {
	jobsCh := make(chan int)
	resultsCh := make(chan int)

	// (1) ワーカーを起動する
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobsCh {
				resultsCh <- process(j)
			}
		}()
	}

	// (2) ジョブを流し込むゴルーチン
	go func() {
		for _, j := range jobs {
			jobsCh <- j
		}
		close(jobsCh) // 全部流し終わったのでjobsChを閉じる
	}()

	// (3) 全ワーカーの完了を待ってからresultsChを閉じるゴルーチン
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// (4) 結果を集める
	var results []int
	for r := range resultsCh {
		results = append(results, r)
	}
	return results
}
```

ポイントは、`jobsCh` を閉じることで各ワーカーの `for range jobsCh` が自然に終了し、`wg.Wait()` が完了し、`resultsCh` が閉じられ、呼び出し側の `for range resultsCh` も自然に終了する——という「閉じる連鎖」がきれいに設計されていることです。

`process` の実行順序やワーカー間の処理順序は保証されないため、`results` の**順序も保証されません**（採点では、順序を無視して中身の集合が一致するかを確認します）。

## 演習

上記の `processJobs` 関数を実装してください。シグネチャは `func processJobs(jobs []int, workerCount int, process func(int) int) []int` です。

## ヒント

- `jobsCh` を閉じるのは「ジョブを流し込むゴルーチン」の役目、`resultsCh` を閉じるのは「`wg.Wait()` を待つゴルーチン」の役目です。役割を混同しないよう注意してください
- `workerCount` が `jobs` の数より多くても少なくても正しく動く必要があります
