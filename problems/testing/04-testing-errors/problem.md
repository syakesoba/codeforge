# Lesson 4: エラーケースをテストする

## 今回学ぶこと

- 戻り値が `(値, error)` の形をした関数のテストの書き方
- `errors.Is` で「特定の種類のエラーが返ってきたか」を確認する方法
- 「値が正しいこと」と「エラーが正しいこと」を1つのテーブルテストで両方扱う方法

### errors.Is で期待するエラーかどうかを確認する

Goでは `errors.New` で作った変数（センチネルエラー）を使い、「どんな種類のエラーが起きたか」を区別するのが一般的です。

```go
var ErrNotFound = errors.New("not found")

func findUser(id int) (*User, error) {
	if id == 0 {
		return nil, ErrNotFound
	}
	// ...
}
```

テストでは、`err == ErrNotFound` のように直接比較することもできますが、`errors.Is(err, ErrNotFound)` を使うのが推奨されます。標準ライブラリの関数（`strconv.Atoi` など）が返すエラーは `%w` でラップされていることがあり、直接比較では見つけられない場合でも `errors.Is` はラップの中まで辿って正しく判定してくれるためです。

```go
func TestFindUser(t *testing.T) {
	_, err := findUser(0)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("findUser(0) error = %v, want ErrNotFound", err)
	}
}
```

### 値のケースとエラーのケースを1つのテーブルにまとめる

「成功したときの値」と「失敗したときのエラー」を両方扱いたい場合、テーブルの各ケースに「期待するエラー（無ければ nil）」というフィールドを持たせると、1つのループでどちらのケースも処理できます。

```go
type divCase struct {
	name      string
	a, b      int
	want      int
	wantErrIs error // nilなら成功を期待
}

func runDivCases(t *testing.T, cases []divCase) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := divide(c.a, c.b)

			if c.wantErrIs != nil {
				if !errors.Is(err, c.wantErrIs) {
					t.Errorf("divide error = %v, want %v", err, c.wantErrIs)
				}
				return // エラーケースでは値の比較はしない
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != c.want {
				t.Errorf("divide(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
			}
		})
	}
}
```

`wantErrIs` が設定されているケースでは値の比較をスキップして `return` している点に注目してください。エラーが返ったときの `got` の値（多くの場合ゼロ値）を比較しても意味が無いためです。

## 演習

以下は**すでに実装済み**です（変更不要です）。

```go
var ErrNotPositive = errors.New("value must be positive")

// parsePositiveInt は文字列を正の整数としてパースする。
// パース自体に失敗すればstrconv.Atoiのエラーを、
// パースできても0以下ならErrNotPositiveを返す。
func parsePositiveInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if n <= 0 {
		return 0, ErrNotPositive
	}
	return n, nil
}
```

`runParsePositiveIntCases` を実装してください。`wantErrIs` が `nil` でなければ `errors.Is(err, c.wantErrIs)` で確認し、`nil` なら `err` が無いこと・`got` が `want` と一致することを確認してください。

## ヒント

- `errors.Is` を使うには `"errors"` パッケージのimportが必要です。書いてから「インポートを自動修正」で追加できます
- エラーケースでは、値（`got`）の比較は行わず `return` してサブテストを終えましょう
- `strconv.Atoi` が返すパースエラーは `strconv.ErrSyntax` を包んでいるため、`errors.Is(err, strconv.ErrSyntax)` で判定できます（非公開テストで使っています）
