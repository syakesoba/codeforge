//go:build ignore

package main

// gradeScore はわざと境界条件にバグを仕込んだ実装です
// （90点ちょうどが "A" ではなく "B" になってしまう）。
// ユーザーが書いたテストがこのバグを検出できるかを確認するために使います。
func gradeScore(score int) string {
	switch {
	case score > 90:
		return "A"
	case score >= 70:
		return "B"
	case score >= 50:
		return "C"
	default:
		return "F"
	}
}
