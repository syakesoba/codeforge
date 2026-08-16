package main

// ValidateForm はフォーム入力を検証し、問題があればすべてまとめて返します。
// 見つかった問題ぶんの error を errors.Join でまとめてください。
// 問題が1つも無ければ nil を返します（errors.Join は要素が全てnilか、
// スライスが空のときnilを返すので、自然に実現できます）。
func ValidateForm(name string, age int, email string) error {
	// TODO:
	// 1. var errs []error を用意する
	// 2. name が "" なら errs に errors.New("name is required") を追加する
	// 3. age が 0未満または150より大きいなら errs に errors.New("age out of range") を追加する
	// 4. email に "@" が含まれていなければ errs に errors.New("invalid email") を追加する
	//    （strings.Contains(email, "@") で判定できます）
	// 5. return errors.Join(errs...) で返す
	return nil
}

func main() {}
