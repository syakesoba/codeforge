//go:build ignore

package main

import (
	"errors"
	"fmt"
)

type User struct {
	ID   int
	Name string
}

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	FindByID(id int) (*User, error)
}

type InMemoryUserRepository struct {
	users map[int]*User
}

func NewInMemoryUserRepository(users map[int]*User) *InMemoryUserRepository {
	return &InMemoryUserRepository{users: users}
}

func (r *InMemoryUserRepository) FindByID(id int) (*User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("find user %d: %w", id, ErrUserNotFound)
	}
	return user, nil
}

func UserNames(repo UserRepository, ids []int) ([]string, error) {
	names := make([]string, 0, len(ids))
	for _, id := range ids {
		user, err := repo.FindByID(id)
		if err != nil {
			return nil, err
		}
		names = append(names, user.Name)
	}
	return names, nil
}

func main() {}
