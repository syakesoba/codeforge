# Lesson 1: テーブル駆動テストを書く

## 今回学ぶこと

- Goで標準的な「テーブル駆動テスト（table-driven test）」の書き方
- `t.Run` でサブテストとしてケースを実行する方法
- `t.Errorf` でテストを失敗させ、原因がわかるメッセージを残す方法

### テーブル駆動テストとは

同じ関数を「入力」と「期待する出力」を変えながら何度も検証したいとき、ケースごとに`if`文と`t.Errorf`を並べて書くと同じコードの繰り返しになってしまいます。

Goでは、ケースを**構造体のスライス（テーブル）**として定義し、それを1つのループで回すのが定番の書き方です。

```go
type addCase struct {
	name string
	a, b int
	want int
}

func TestAdd(t *testing.T) {
	cases := []addCase{
		{"positive numbers", 2, 3, 5},
		{"with zero", 0, 5, 5},
		{"negative numbers", -2, -3, -5},
	}

	for _, c := range cases {
		got := add(c.a, c.b)
		if got != c.want {
			t.Errorf("add(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}
```

ケースを追加したいときは、テーブルに1行足すだけで済みます。テストのロジック（比較して失敗させる部分）を毎回書き直す必要がありません。

### t.Run でサブテストにする

上のコードにはひとつ弱点があります。途中のケースが失敗しても、ループはそのまま最後まで進んでしまい、**どのケースで失敗したのか**はメッセージを読むまでわかりません。

`t.Run(名前, func(t *testing.T) { ... })` を使うと、ケースごとに**名前付きの独立したサブテスト**として実行できます。

```go
for _, c := range cases {
	c := c // Go 1.22以降は不要だが、書いても害はない
	t.Run(c.name, func(t *testing.T) {
		got := add(c.a, c.b)
		if got != c.want {
			t.Errorf("add(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	})
}
```

こうすると `go test -v` の出力に `--- PASS: TestAdd/positive_numbers` のようにケース単位で結果が表示され、どのケースが失敗したか一目でわかります。CodeForgeの「実行結果ログ」にもこの形式で表示されます。

### t.Fatal と t.Errorf の違い

- `t.Errorf`: 失敗を記録するが、**そのテスト関数の実行は続行する**（複数の問題を一度に報告できる）
- `t.Fatalf`: 失敗を記録して、**その場でテスト関数を打ち切る**（後続のコードが失敗の原因で panic するのを防ぎたいときに使う）

サブテストの中で「値を比較して失敗を記録するだけ」の場合は `t.Errorf` を使うのが一般的です。

## 演習

以下の `gradeScore` 関数は、0〜100点のテストの点数を `"A"`/`"B"`/`"C"`/`"F"` の評価に変換する関数で、**すでに実装済み**です（変更不要です）。

```go
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
```

`runGradeCases` 関数を実装してください。`cases` の各要素について `gradeScore` を呼び出し、`t.Run` でサブテストとして実行し、結果が `want` と一致しなければ `t.Errorf` で失敗させてください。

## ヒント

- ループの中で `t.Run(c.name, func(t *testing.T) { ... })` を呼び出します
- サブテストの中で `gradeScore(c.score)` の結果を `c.want` と比較します
- 一致しない場合は `t.Errorf("gradeScore(%d) = %q, want %q", c.score, got, c.want)` のように、実際の値・期待値の両方をメッセージに含めてください
