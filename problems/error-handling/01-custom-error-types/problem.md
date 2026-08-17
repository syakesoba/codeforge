# Lesson 1: カスタムエラー型を定義する

## 今回学ぶこと

- Goの `error` インターフェースの正体
- `errors.New` だけでは表現できない「構造化されたエラー情報」
- 独自のエラー型を定義する方法

### error インターフェースの正体

Goの `error` は、実はとてもシンプルなインターフェースです。

```go
type error interface {
	Error() string
}
```

つまり「`Error() string` というメソッドを持つ型」であれば、何でも `error` として扱えます。`errors.New("...")` が返しているのも、内部的にはこのインターフェースを満たす小さな構造体にすぎません。

### errors.New だけでは足りない場面

```go
if age < 0 || age > 150 {
	return errors.New("age: 0から150の範囲で指定してください")
}
```

これでも動きますが、呼び出し側が「どのフィールドが」「どんな理由で」不正だったのかをプログラム的に取り出すことができません。文字列を正規表現でパースするような無理をしない限り、`err.Error()` の内容を人間が読むことしかできないのです。

### 独自のエラー型を定義する

そこで、構造体に `Error() string` メソッドを実装して、独自のエラー型を作ります。

```go
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

こうすると、`ValidationError` は `error` インターフェースを満たすので `error` として返せる一方、呼び出し側は必要に応じて型アサーションで元の構造体を取り出し、`Field` や他のフィールドに直接アクセスできます。

```go
func ValidateAge(age int) error {
	if age < 0 || age > 150 {
		return &ValidationError{Field: "age", Message: "0から150の範囲で指定してください"}
	}
	return nil
}

// 呼び出し側
err := ValidateAge(-1)
if ve, ok := err.(*ValidationError); ok {
	fmt.Println("不正なフィールド:", ve.Field)
}
```

> **ポインタ型で実装する理由**: `Error()` を `*ValidationError`（ポインタレシーバ）に実装し、返り値も `&ValidationError{...}`（ポインタ）にするのが一般的です。値渡しだと構造体がコピーされ、後述の `errors.As` での比較や、大きな構造体でのコピーコストの面で不利になるためです。

## 演習

`ValidationError` 型と、それを使う `ValidateAge` 関数を実装してください。

**`(e *ValidationError) Error() string`**
- `fmt.Sprintf("%s: %s", e.Field, e.Message)` の形式で文字列を返す

**`ValidateAge(age int) error`**
- `age` が `0`未満または`150`より大きい場合、`&ValidationError{Field: "age", Message: "0から150の範囲で指定してください"}` を返す
- それ以外の場合は `nil` を返す

## ヒント

- `Error()` メソッドの中身は `fmt.Sprintf` の1行だけです。`fmt` パッケージのimportが必要になるので、エディタの「インポートを自動修正」機能を使ってください
- `ValidateAge` はまず条件をチェックし、条件に当てはまれば `&ValidationError{...}` を、そうでなければ `nil` を返す、というシンプルな分岐です
