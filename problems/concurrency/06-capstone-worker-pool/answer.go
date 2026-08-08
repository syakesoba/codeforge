//go:build ignore

package main

import "sync"

// processJobs は jobs を workerCount 個のワーカーゴルーチンで並行処理し、
// 各ジョブに process を適用した結果を返します（結果の順序は保証されません）。
func processJobs(jobs []int, workerCount int, process func(int) int) []int {
	jobsCh := make(chan int)
	resultsCh := make(chan int)

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

	go func() {
		for _, j := range jobs {
			jobsCh <- j
		}
		close(jobsCh)
	}()

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var results []int
	for r := range resultsCh {
		results = append(results, r)
	}
	return results
}

func main() {}
