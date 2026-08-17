package main

// ValidationError は入力検証に失敗したことを表すエラー型です。
// どのフィールドが、なぜ不正だったのかを構造化して持てるのが
// 文字列だけの errors.New との違いです。
type ValidationError struct {
	Field   string
	Message string
}

// Error は error インターフェースを満たすために実装します。
// error インターフェースは「Error() string を持つ型」であるだけなので、
// このメソッドを実装すれば ValidationError はそのまま error として扱えます。
func (e *ValidationError) Error() string {
	// TODO: fmt.Sprintf("%s: %s", e.Field, e.Message) を返す
	return ""
}

// ValidateAge は age が 0〜150 の範囲かどうかを検証します。
// 範囲外なら *ValidationError を、問題なければ nil を返してください。
func ValidateAge(age int) error {
	// TODO:
	// age < 0 または age > 150 なら
	//   &ValidationError{Field: "age", Message: "0から150の範囲で指定してください"} を返す
	// それ以外なら nil を返す
	return nil
}

func main() {}
