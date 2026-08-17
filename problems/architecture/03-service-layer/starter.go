package main

import "errors"

// Account は口座の残高情報です。
type Account struct {
	ID      int
	Balance int
}

// ErrInsufficientBalance は残高不足を表すセンチネルエラーです。
var ErrInsufficientBalance = errors.New("insufficient balance")

// AccountRepository は口座データの永続化を抽象化します（Lesson 2で学んだリポジトリパターン）。
type AccountRepository interface {
	FindByID(id int) (*Account, error)
	Save(acc *Account) error
}

// InMemoryAccountRepository は AccountRepository のインメモリ実装です（実装済み）。
type InMemoryAccountRepository struct {
	accounts map[int]*Account
}

func NewInMemoryAccountRepository(accounts map[int]*Account) *InMemoryAccountRepository {
	return &InMemoryAccountRepository{accounts: accounts}
}

func (r *InMemoryAccountRepository) FindByID(id int) (*Account, error) {
	acc, ok := r.accounts[id]
	if !ok {
		return nil, errors.New("account not found")
	}
	return acc, nil
}

func (r *InMemoryAccountRepository) Save(acc *Account) error {
	r.accounts[acc.ID] = acc
	return nil
}

// TransferService は「口座間の送金」という業務ルールを持つサービス層です。
// データの出し入れ（AccountRepository）と、金額チェックなどの業務ルールを
// 分離しておくことで、ルール部分だけを単体でテストしたり、リポジトリの
// 実装（インメモリ→SQL等）を差し替えたりしやすくなります。
type TransferService struct {
	repo AccountRepository
}

// NewTransferService はコンストラクタです（実装済み）。
func NewTransferService(repo AccountRepository) *TransferService {
	return &TransferService{repo: repo}
}

// Transfer は fromID の口座から toID の口座へ amount を送金します。
func (s *TransferService) Transfer(fromID, toID, amount int) error {
	// TODO:
	// 1. from, err := s.repo.FindByID(fromID) 。エラーならそのまま返す
	// 2. to, err := s.repo.FindByID(toID) 。エラーならそのまま返す
	// 3. from.Balance が amount より小さいなら ErrInsufficientBalance を返す
	// 4. from.Balance -= amount
	// 5. to.Balance += amount
	// 6. s.repo.Save(from) を呼び、エラーならそのまま返す
	// 7. s.repo.Save(to) を呼び、エラーならそのまま返す
	// 8. 成功したら nil を返す
	return nil
}

func main() {}
