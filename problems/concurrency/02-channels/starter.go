package main

// collectSquares は nums の各要素を2乗した値をchannel経由で送信し、
// 呼び出し側でrangeを使って集めた結果を返します。
func collectSquares(nums []int) []int {
	// TODO:
	// 1. make(chan int) でchannelを作る
	// 2. ゴルーチンの中で nums を順番に2乗してchannelへ送信し、送り終わったら close する
	// 3. for range でchannelから値を受け取り、スライスに集めて返す
	return nil
}

func main() {}
