# Lesson 2: リポジトリパターンでデータアクセスを抽象化する

## 今回学ぶこと

- リポジトリパターン: データの永続化をインターフェースの向こうに隠す設計
- 「引数の型をインターフェースにする」ことで得られる柔軟性
- インメモリ実装を使った、DB不要なテスト・開発

### なぜデータアクセスを抽象化するのか

業務ロジックの中に `sql.DB` への直接のクエリが散らばっていると、次のような問題が起きます。

- ロジックをテストするたびに、本物のDBを用意する必要がある
- DBの種類（MySQL→PostgreSQLなど）を変えると、ロジックのコードまで書き換わってしまう

そこで、「データをどう取得するか」という**手段**を、インターフェースとして切り出します。

```go
type UserRepository interface {
	FindByID(id int) (*User, error)
}
```

業務ロジック側はこのインターフェースの向こう側が何なのか（インメモリのmapなのか、SQLデータベースなのか）を一切知る必要がありません。これを**リポジトリパターン**と呼びます。

### インメモリ実装

開発の初期段階やテストでは、本物のDBの代わりに単純なmapを使った実装で十分なことがよくあります。

```go
type InMemoryUserRepository struct {
	users map[int]*User
}

func (r *InMemoryUserRepository) FindByID(id int) (*User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("find user %d: %w", id, ErrUserNotFound)
	}
	return user, nil
}
```

`InMemoryUserRepository` は `UserRepository` インターフェースの `FindByID(id int) (*User, error)` というメソッドを持っているので、自動的に `UserRepository` を満たします（Goでは「このインターフェースを実装します」という宣言は不要で、メソッドセットが一致すれば自動的にそのインターフェースとして扱えます）。

### 引数の型をインターフェースにする

リポジトリパターンの効果を最大限に活かすには、それを使う側の関数も「具体的な実装」ではなく「インターフェース」を受け取るようにします。

```go
// 悪い例: 具体的な実装に直接依存してしまっている
func UserNames(repo *InMemoryUserRepository, ids []int) ([]string, error) { ... }

// 良い例: インターフェースに依存している
func UserNames(repo UserRepository, ids []int) ([]string, error) { ... }
```

良い例のように書いておけば、`InMemoryUserRepository` でも、将来作るかもしれない `SQLUserRepository` でも、テスト用の偽物実装でも、`FindByID(id int) (*User, error)` というメソッドさえ持っていれば `UserNames` にそのまま渡せます。**`UserNames` 自身は一切変更する必要がありません。**

## 演習

`InMemoryUserRepository.FindByID` と `UserNames` を実装してください。

**`(r *InMemoryUserRepository) FindByID(id int) (*User, error)`**
- `r.users` から `id` を検索する
- 見つかれば、そのユーザーと `nil` を返す
- 見つからなければ、`nil` と `fmt.Errorf("find user %d: %w", id, ErrUserNotFound)` を返す

**`UserNames(repo UserRepository, ids []int) ([]string, error)`**
- `ids` の各要素について `repo.FindByID(id)` を呼ぶ
- エラーが発生したら、その時点で `nil, err` を返す
- エラーがなければ、名前を集めたスライスと `nil` を返す

## ヒント

- `FindByID` の検索は `user, ok := r.users[id]` の形（カンマOKイディオム）で行います
- `UserNames` は `names := make([]string, 0, len(ids))` で結果を溜めるスライスを用意してから、`ids` をfor-rangeでループします
- 引数の型が `UserRepository`（インターフェース）であって `*InMemoryUserRepository`（具体型）でないことに注目してください。これにより `UserNames` はどんなリポジトリ実装に対しても動作します
