package main

// HashPassword は平文パスワードをbcryptでハッシュ化して返します。
func HashPassword(password string) (string, error) {
	// TODO:
	// 1. bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) を呼ぶ
	// 2. エラーなら "", err を返す
	// 3. 成功したら string(hashed), nil を返す
	return "", nil
}

// CheckPassword は保存済みハッシュと平文パスワードを照合します。
// パスワードが正しければ true、間違っていれば false を返します。
func CheckPassword(hashedPassword, password string) bool {
	// TODO: bcrypt.CompareHashAndPassword を使って照合する
	return false
}

func main() {}
