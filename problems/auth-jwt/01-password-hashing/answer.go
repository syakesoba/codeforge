//go:build ignore

package main

import "golang.org/x/crypto/bcrypt"

// HashPassword は平文パスワードをbcryptでハッシュ化して返します。
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPassword は保存済みハッシュと平文パスワードを照合します。
// パスワードが正しければ true、間違っていれば false を返します。
func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func main() {}
