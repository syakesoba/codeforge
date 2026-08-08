//go:build ignore

package main

// firstResult は ch1, ch2 のうち、先に値が届いた方を返します。
func firstResult(ch1, ch2 <-chan int) int {
	select {
	case v := <-ch1:
		return v
	case v := <-ch2:
		return v
	}
}

func main() {}
