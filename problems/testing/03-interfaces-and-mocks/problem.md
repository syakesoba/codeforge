# Lesson 3: インターフェースでモックする

## 今回学ぶこと

- 「現在時刻」「外部API」など、テストのたびに結果が変わってしまう依存をインターフェースで抽象化する方法
- テスト用の「フェイク（モック）」実装を書く方法
- 依存性注入（Dependency Injection）の基本的な考え方

### なぜテストしにくいコードがあるのか

以下のような関数は、そのままではテストしづらいコードの典型例です。

```go
func isBusinessHours() bool {
	h := time.Now().Hour() // 実行するたびに結果が変わってしまう
	return h >= 9 && h < 18
}
```

`time.Now()` を直接呼んでいるため、テストを実行するたびに結果が変わります。「9時に実行したときの動作」をテストで固定的に再現することができません。

### インターフェースで抽象化する

「現在時刻を取得する」という操作をインターフェースとして切り出し、関数の引数として受け取るようにすると、本番では本物の時計を、テストでは偽物の時計を渡せるようになります。

```go
// Clock は現在時刻を返すインターフェース。
type Clock interface {
	Now() time.Time
}

// realClock は本番で使う実装。
type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// isBusinessHours は、渡された clock の時刻を使って判定する。
func isBusinessHours(clock Clock) bool {
	h := clock.Now().Hour()
	return h >= 9 && h < 18
}
```

このように「実際に使うもの」を直接呼ばず、インターフェース越しに受け取ることを**依存性注入（Dependency Injection）**と呼びます。

### テスト用のフェイクを書く

テストでは、`Clock` インターフェースを実装した「フェイク」を用意し、好きな時刻を固定して渡します。

```go
// fakeClock はテスト用の Clock 実装。固定した時刻を返す。
type fakeClock struct {
	now time.Time
}

func (f fakeClock) Now() time.Time {
	return f.now
}

func TestIsBusinessHours(t *testing.T) {
	clock := fakeClock{now: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)} // 10時に固定
	if !isBusinessHours(clock) {
		t.Errorf("10時は営業時間内のはずが、falseが返った")
	}
}
```

`fakeClock` は本物の `time.Time` を保持しているだけの、ごく小さな構造体です。モックというと大げさに聞こえますが、実体はこの程度のシンプルなもので十分なことがほとんどです。

## 演習

以下は**すでに実装済み**です（変更不要です）。

```go
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// isExpired は、clock.Now() が expiresAt より後（after）であれば true を返す。
func isExpired(clock Clock, expiresAt time.Time) bool {
	return clock.Now().After(expiresAt)
}
```

以下の2つを実装してください。

1. `fakeClock.Now()`: 保持している `f.now` をそのまま返す
2. `runIsExpiredCases`: `cases` の各要素について、その `now` を持つ `fakeClock` を作り、`isExpired(clock, c.expiresAt)` を呼び出して `want` と比較する（`t.Run` でサブテスト化し、不一致なら `t.Errorf`）

## ヒント

- `fakeClock.Now()` は1行で書けます: `return f.now`
- `clock := fakeClock{now: c.now}` で偽物の時計を作ってから `isExpired(clock, c.expiresAt)` を呼び出してください
- ループ・サブテストの組み立て方はLesson 1・2と同じです
