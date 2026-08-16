# Lesson 4: 関数オプションパターンで柔軟な初期化をする

## 今回学ぶこと

- 設定項目が多い構造体を初期化する際の困りごと
- **関数オプションパターン（Functional Options Pattern）** という、Goでよく使われる解決策
- 可変長引数 `...Option` の使い方

### 設定項目が多いと、コンストラクタが破綻する

`Server` のように複数の設定項目（`Host`, `Port`, `Timeout`, …）を持つ構造体を、素朴にコンストラクタで初期化しようとすると、こうなりがちです。

```go
func NewServer(host string, port int, timeout int) *Server { ... }

// 呼び出し側は「全部の値」を毎回指定しなければならない
srv := NewServer("localhost", 8080, 30)
```

これだと、「ポート番号だけ変えたい、他はデフォルトのままでいい」というよくある要望に応えづらく、項目が増えるたびに引数の並び順を覚えるのも大変になります（設定項目が数十個になるライブラリのコードでよく問題になります）。

### 関数オプションパターン

Goでは、「構造体を変更する関数」を可変長引数として受け取る、というパターンがよく使われます。

```go
type Option func(*Server)

func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

func NewServer(opts ...Option) *Server {
	s := &Server{Host: "localhost", Port: 8080, Timeout: 30} // まずデフォルト値
	for _, opt := range opts {
		opt(s) // 渡されたoptionだけ順番に適用
	}
	return s
}
```

呼び出し側は、変更したい項目だけを指定できるようになります。

```go
srv1 := NewServer()                              // 全部デフォルト
srv2 := NewServer(WithPort(9000))                // ポートだけ変更
srv3 := NewServer(WithPort(9000), WithTimeout(5)) // 複数指定も自然にできる
```

`WithPort(9000)` は、それ自体は `Server` を直接変更するのではなく、「`Server` を受け取って `Port` を書き換える関数」を**返す**点がポイントです。`NewServer` の中でその関数が実際に呼ばれて初めて `Server` が変更されます。

> この設計により、将来 `WithHost(host string)` のような新しいオプションを追加しても、既存の `NewServer()` や `NewServer(WithPort(9000))` の呼び出しコードは一切壊れません。これは「設定項目が多い＆将来も増えていく」種類の構造体（HTTPクライアント、DB接続、外部ライブラリの設定など）で特に威力を発揮します。

## 演習

`WithPort`、`WithTimeout`、`NewServer` を実装してください。

**`WithPort(port int) Option`**
- `s.Port = port` を行う `Option`（`func(*Server)`）を返す

**`WithTimeout(timeout int) Option`**
- `s.Timeout = timeout` を行う `Option` を返す

**`NewServer(opts ...Option) *Server`**
- `Host: "localhost"`, `Port: 8080`, `Timeout: 30` を初期値とする `*Server` を作る
- `opts` を先頭から順番に1つずつ適用する
- 適用し終わった `*Server` を返す

## ヒント

- `WithPort` の中身は `return func(s *Server) { s.Port = port }` の1行です。クロージャが `port` を覚えている点がポイントです
- `NewServer` は `for _, opt := range opts { opt(s) }` で全ての `Option` を順番に呼び出します
