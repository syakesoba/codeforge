package main

// buildCubeMap は nums の各要素をキーとして、その3乗の値をmapに格納して返します。
// 複数のゴルーチンから安全に書き込めるよう、sync.Mutexで保護してください。
func buildCubeMap(nums []int) map[int]int {
	// TODO:
	// 1. 結果を格納する map[int]int を用意する
	// 2. sync.Mutex を1つ用意する
	// 3. 各要素についてゴルーチンを起動し、3乗の計算はロックの外、
	//    mapへの書き込みは mu.Lock() / mu.Unlock() で挟んで保護する
	// 4. sync.WaitGroup で全ゴルーチンの完了を待つ
	return nil
}

func main() {}
