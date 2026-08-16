# Lesson 2: サブテストを使いこなす

## 今回学ぶこと

- 入力値からサブテストの名前を組み立てる方法
- サブテストの命名で意識すべきポイント
- `t.Parallel()` でサブテストを並列実行する方法（紹介のみ）

### ケースに name フィールドが無いとき

Lesson 1では、テーブルの各ケースに `name` フィールドを用意していました。しかし、入力値がそのまま良い名前になる場合は、わざわざ `name` を書かなくても `fmt.Sprintf` で組み立てられます。

```go
type divCase struct {
	a, b int
	want int
}

func TestDivide(t *testing.T) {
	cases := []divCase{
		{10, 2, 5},
		{9, 3, 3},
	}

	for _, c := range cases {
		name := fmt.Sprintf("%d/%d", c.a, c.b)
		t.Run(name, func(t *testing.T) {
			got := divide(c.a, c.b)
			if got != c.want {
				t.Errorf("divide(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
			}
		})
	}
}
```

`go test -v` を実行すると `--- PASS: TestDivide/10/2` のように表示され、どの入力の組み合わせで実行されたかが一目でわかります。

### サブテスト名に使えない文字

サブテスト名にスペースが含まれていると、`go test -run` で指定するときにアンダースコア `_` に置き換わって表示されます（実行結果には影響しません）。迷ったら、`/` や `=` を区切りに使うと入力値が読み取りやすくなります。

```go
name := fmt.Sprintf("month=%d/leap=%v", c.month, c.isLeapYear)
// 表示例: month=2/leap=true
```

### t.Parallel() で並列実行する（紹介）

サブテスト同士に依存関係が無く、それぞれが独立して実行できる場合は、サブテストの先頭で `t.Parallel()` を呼ぶことで、複数のサブテストを並列に実行できます。

```go
t.Run(name, func(t *testing.T) {
	t.Parallel() // このサブテストは他の並列サブテストと同時に実行される
	got := divide(c.a, c.b)
	if got != c.want {
		t.Errorf("divide(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
	}
})
```

ケース数が多く、1件ずつの実行に時間がかかるテスト（ネットワークI/Oを伴うテストなど）で特に効果があります。今回の演習では必須ではありませんが、覚えておくと便利です。

## 演習

以下の `daysInMonth` 関数は、月（1〜12）とうるう年かどうかを受け取り、その月の日数を返す関数で、**すでに実装済み**です（変更不要です）。

```go
func daysInMonth(month int, isLeapYear bool) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if isLeapYear {
			return 29
		}
		return 28
	default:
		return 0
	}
}
```

`runDaysInMonthCases` 関数を実装してください。`daysCase` に `name` フィールドは**ありません**。`cases` の各要素について、`fmt.Sprintf` で入力値からサブテスト名を組み立て、`t.Run` でサブテストとして実行し、`daysInMonth` の結果が `want` と一致しなければ `t.Errorf` で失敗させてください。

## ヒント

- サブテスト名の例: `fmt.Sprintf("month=%d/leap=%v", c.month, c.isLeapYear)`
- ループの中身の骨格はLesson 1とほぼ同じです。「名前をどう作るか」だけが違います
- `daysInMonth(c.month, c.isLeapYear)` の結果を `c.want` と比較し、一致しなければ `t.Errorf("daysInMonth(%d, %v) = %d, want %d", c.month, c.isLeapYear, got, c.want)` のように実際の値・期待値の両方をメッセージに含めてください
