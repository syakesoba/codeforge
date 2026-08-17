package main

import "errors"

// Book は蔵書データです。
type Book struct {
	ID     int
	Title  string
	OnLoan bool
}

// ErrBookNotFound / ErrBookAlreadyOnLoan は業務ルール上のエラーです。
var (
	ErrBookNotFound      = errors.New("book not found")
	ErrBookAlreadyOnLoan = errors.New("book already on loan")
)

// BookRepository は蔵書データの永続化を抽象化します（Lesson 2のリポジトリパターン）。
type BookRepository interface {
	FindByID(id int) (*Book, error)
	Save(b *Book) error
}

// InMemoryBookRepository は実装済みです。
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

// Notifier は通知手段を抽象化します（Lesson 1のDI）。
type Notifier interface {
	Notify(message string) error
}

// LibraryService は「本を借りる」という業務ルールを持つサービス層（Lesson 3）です。
// BookRepository（データアクセス）とNotifier（通知）の2つの依存を、
// どちらも具体的な実装ではなくインターフェースとして受け取ります。
type LibraryService struct {
	repo     BookRepository
	notifier Notifier
}

// NewLibraryService はコンストラクタです（実装済み）。
func NewLibraryService(repo BookRepository, notifier Notifier) *LibraryService {
	return &LibraryService{repo: repo, notifier: notifier}
}

// Borrow は本を貸し出します。
func (s *LibraryService) Borrow(bookID int) error {
	// TODO:
	// 1. book, err := s.repo.FindByID(bookID) 。エラーならそのまま返す
	// 2. book.OnLoan が true なら ErrBookAlreadyOnLoan を返す
	// 3. book.OnLoan = true にする
	// 4. s.repo.Save(book) を呼び、エラーならそのまま返す
	// 5. s.notifier.Notify(fmt.Sprintf("「%s」を貸し出しました", book.Title)) を呼び、
	//    エラーならそのまま返す
	// 6. 成功したら nil を返す
	return nil
}

func main() {}
