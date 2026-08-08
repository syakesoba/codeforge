# Lesson 1: パスワードをハッシュ化する

## 今回学ぶこと

- パスワードを平文で保存してはいけない理由
- `bcrypt` によるハッシュ化と検証
- なぜ「比較」ではなく専用の関数を使うのか

### パスワードは絶対に平文で保存しない

もしDBにパスワードをそのまま（平文で）保存していると、DBが漏洩した瞬間に全ユーザーのパスワードが流出します。多くの人が複数のサービスで同じパスワードを使い回しているため、被害は自分のサービスだけに留まりません。

そこで、パスワードは**ハッシュ化**して保存します。ハッシュ化は「元に戻せない変換」なので、漏洩してもパスワード自体は分かりません。

### bcryptを使う

ハッシュ化には、パスワード保存用に設計された **bcrypt** を使います。Goでは `golang.org/x/crypto/bcrypt` パッケージが標準的です。

```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
```

- `GenerateFromPassword` は `[]byte` を受け取り `[]byte` を返すので、文字列との変換が必要です
- 第2引数の `cost` は計算の重さです。値を大きくするほど総当たり攻撃に強くなる代わりに時間がかかります。通常は `bcrypt.DefaultCost` を使います

> **なぜSHA-256ではなくbcryptなのか**: SHA-256などの汎用ハッシュ関数は「高速であること」が長所ですが、パスワード保存では逆に短所になります。攻撃者が1秒間に何億回も試行できてしまうからです。bcryptは意図的に低速に設計されており、さらに後述の「ソルト」を自動で扱ってくれます。

### 検証は CompareHashAndPassword で行う

ここが重要です。bcryptは同じパスワードでも**毎回違うハッシュ値**を生成します。内部でランダムな「ソルト」を混ぜているためです。

```go
h1, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
h2, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
// h1 と h2 は異なる文字列になる！
```

そのため、「入力されたパスワードをハッシュ化して、保存済みハッシュと文字列比較する」というやり方は**動きません**。検証には専用の関数を使います。

```go
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
```

`CompareHashAndPassword` は、保存済みハッシュに埋め込まれたソルトを取り出して同じ手順で計算し、一致するかを判定してくれます。パスワードが正しければ `nil`、間違っていれば `bcrypt.ErrMismatchedHashAndPassword` を返します。

## 演習

`HashPassword` と `CheckPassword` の2つの関数を実装してください。

**`HashPassword(password string) (string, error)`**
- `bcrypt.GenerateFromPassword` でハッシュ化し、文字列として返す
- コストは `bcrypt.DefaultCost` を使う

**`CheckPassword(hashedPassword, password string) bool`**
- `bcrypt.CompareHashAndPassword` で検証し、パスワードが正しければ `true`、間違っていれば `false` を返す

## ヒント

- `bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)` の戻り値は `([]byte, error)` です
- `CheckPassword` は `return bcrypt.CompareHashAndPassword(...) == nil` の1行で書けます
