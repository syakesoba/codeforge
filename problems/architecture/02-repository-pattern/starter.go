package main

import "errors"

// User はユーザーデータです。
type User struct {
	ID   int
	Name string
}

// ErrUserNotFound はユーザーが見つからなかったことを表すセンチネルエラーです。
var ErrUserNotFound = errors.New("user not found")

// UserRepository はユーザーデータの永続化を抽象化するインターフェースです。
// 「インメモリのmap」「SQLデータベース」「外部API」など、データの実体が
// 何であるかをこのインターフェースの向こう側に隠すことで、業務ロジック側は
// 実装の詳細を気にせずに済むようになります（リポジトリパターン）。
type UserRepository interface {
	FindByID(id int) (*User, error)
}

// InMemoryUserRepository は UserRepository のインメモリ実装です。
// テストや開発初期段階で、本物のDBを用意せずに動作確認するのに便利です。
type InMemoryUserRepository struct {
	users map[int]*User
}

// NewInMemoryUserRepository はコンストラクタです（実装済み）。
func NewInMemoryUserRepository(users map[int]*User) *InMemoryUserRepository {
	return &InMemoryUserRepository{users: users}
}

// FindByID は UserRepository インターフェースの実装です。
func (r *InMemoryUserRepository) FindByID(id int) (*User, error) {
	// TODO:
	// r.users[id] をカンマOKイディオムで検索する
	// 見つかれば user, nil を返す
	// 見つからなければ nil, fmt.Errorf("find user %d: %w", id, ErrUserNotFound) を返す
	return nil, nil
}

// UserNames は、複数のidに対応するユーザー名をまとめて取得します。
// 引数の型を「具体的な InMemoryUserRepository」ではなく「UserRepositoryインターフェース」
// にすることで、将来SQL実装などに差し替えても、この関数は一切変更不要になります。
func UserNames(repo UserRepository, ids []int) ([]string, error) {
	// TODO:
	// names := make([]string, 0, len(ids)) で用意する
	// idsをループし、repo.FindByID(id) を呼ぶ
	// エラーがあれば nil, err をそのまま返す
	// 無ければ user.Name を names に追加していく
	// 最後に names, nil を返す
	return nil, nil
}

func main() {}
