package main

// gradeScore はテストの点数（0〜100）を A/B/C/F の評価に変換します（実装済み・変更不要）。
func gradeScore(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 70:
		return "B"
	case score >= 50:
		return "C"
	default:
		return "F"
	}
}
