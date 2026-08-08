package main

// processJobs は jobs を workerCount 個のワーカーゴルーチンで並行処理し、
// 各ジョブに process を適用した結果を返します（結果の順序は保証されません）。
func processJobs(jobs []int, workerCount int, process func(int) int) []int {
	// TODO:
	// 1. jobsCh, resultsCh の2つのchannelを作る
	// 2. workerCount個のワーカーゴルーチンを起動する
	//    （sync.WaitGroupで完了を管理し、各ワーカーは for j := range jobsCh で処理する）
	// 3. jobsを流し込むゴルーチンを起動し、流し終わったら jobsCh を close する
	// 4. 全ワーカーの完了(wg.Wait())を待ってから resultsCh を close するゴルーチンを起動する
	// 5. resultsCh を range で受け取り、スライスに集めて返す
	return nil
}

func main() {}
