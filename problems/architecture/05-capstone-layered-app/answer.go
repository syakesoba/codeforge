//go:build ignore

package main

import (
	"errors"
	"fmt"
)

type Book struct {
	ID     int
	Title  string
	OnLoan bool
}

var (
	ErrBookNotFound      = errors.New("book not found")
	ErrBookAlreadyOnLoan = errors.New("book already on loan")
)

type BookRepository interface {
	FindByID(id int) (*Book, error)
	Save(b *Book) error
}

type InMemoryBookRepository struct {
	books map[int]*Book
}

func NewInMemoryBookRepository(books map[int]*Book) *InMemoryBookRepository {
	return &InMemoryBookRepository{books: books}
}

func (r *InMemoryBookRepository) FindByID(id int) (*Book, error) {
	b, ok := r.books[id]
	if !ok {
		return nil, ErrBookNotFound
	}
	return b, nil
}

func (r *InMemoryBookRepository) Save(b *Book) error {
	r.books[b.ID] = b
	return nil
}

type Notifier interface {
	Notify(message string) error
}

type LibraryService struct {
	repo     BookRepository
	notifier Notifier
}

func NewLibraryService(repo BookRepository, notifier Notifier) *LibraryService {
	return &LibraryService{repo: repo, notifier: notifier}
}

func (s *LibraryService) Borrow(bookID int) error {
	book, err := s.repo.FindByID(bookID)
	if err != nil {
		return err
	}
	if book.OnLoan {
		return ErrBookAlreadyOnLoan
	}
	book.OnLoan = true
	if err := s.repo.Save(book); err != nil {
		return err
	}
	return s.notifier.Notify(fmt.Sprintf("「%s」を貸し出しました", book.Title))
}

func main() {}
