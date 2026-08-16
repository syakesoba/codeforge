//go:build ignore

package main

import "errors"

type Account struct {
	ID      int
	Balance int
}

var ErrInsufficientBalance = errors.New("insufficient balance")

type AccountRepository interface {
	FindByID(id int) (*Account, error)
	Save(acc *Account) error
}

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

type TransferService struct {
	repo AccountRepository
}

func NewTransferService(repo AccountRepository) *TransferService {
	return &TransferService{repo: repo}
}

func (s *TransferService) Transfer(fromID, toID, amount int) error {
	from, err := s.repo.FindByID(fromID)
	if err != nil {
		return err
	}
	to, err := s.repo.FindByID(toID)
	if err != nil {
		return err
	}
	if from.Balance < amount {
		return ErrInsufficientBalance
	}
	from.Balance -= amount
	to.Balance += amount
	if err := s.repo.Save(from); err != nil {
		return err
	}
	if err := s.repo.Save(to); err != nil {
		return err
	}
	return nil
}

func main() {}
