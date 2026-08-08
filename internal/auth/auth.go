// Package auth はサイト自体のユーザー認証（サインアップ・ログイン・セッション）を提供します。
//
// セッションはランダムなトークンをHttpOnly Cookieで発行し、サーバー側の
// sessions テーブルと突き合わせる方式です。JWTはCourse「認証・JWT」の
// 学習コンテンツとして扱うため、サイト自体の認証には使っていません。
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/syakesoba/codeforge/internal/store"
)

// SessionDuration はセッションの有効期間です。
const SessionDuration = 30 * 24 * time.Hour

// SessionCookieName はセッショントークンを保存するCookie名です。
const SessionCookieName = "codeforge_session"

// MinPasswordLength はパスワードの最低文字数です。
const MinPasswordLength = 8

var (
	// ErrInvalidCredentials はメールアドレスまたはパスワードが誤っていることを表します。
	// 「ユーザーが存在しない」と「パスワードが違う」を区別すると、
	// 攻撃者に登録済みメールアドレスを推測されるため、意図的に同じエラーにしています。
	ErrInvalidCredentials = errors.New("invalid email or password")

	// ErrInvalidEmail はメールアドレスの形式が不正であることを表します。
	ErrInvalidEmail = errors.New("invalid email format")

	// ErrPasswordTooShort はパスワードが短すぎることを表します。
	ErrPasswordTooShort = fmt.Errorf("password must be at least %d characters", MinPasswordLength)

	// ErrEmailTaken は既に登録済みのメールアドレスであることを表します。
	ErrEmailTaken = store.ErrEmailTaken

	// ErrNoSession は有効なセッションが存在しないことを表します。
	ErrNoSession = errors.New("no valid session")
)

// Service は認証機能を提供します。
type Service struct {
	store *store.Store
}

// New は Service を作成します。
func New(s *store.Store) *Service {
	return &Service{store: s}
}

// Session はログイン成功時に発行されるセッション情報です。
type Session struct {
	Token     string
	ExpiresAt time.Time
	User      *store.User
}

// SignUp は新規ユーザーを登録し、そのままログイン状態にします。
func (s *Service) SignUp(email, password string) (*Session, error) {
	email = normalizeEmail(email)

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrInvalidEmail
	}
	if len(password) < MinPasswordLength {
		return nil, ErrPasswordTooShort
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.store.CreateUser(email, string(hashed))
	if err != nil {
		return nil, err // ErrEmailTaken を含む
	}

	return s.newSession(user)
}

// LogIn はメールアドレスとパスワードを検証し、セッションを発行します。
func (s *Service) LogIn(email, password string) (*Session, error) {
	email = normalizeEmail(email)

	user, err := s.store.UserByEmail(email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// タイミング攻撃を避けるため、ユーザーが存在しない場合も
			// ハッシュ計算と同程度の時間をかけてから失敗させる。
			bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}

	return s.newSession(user)
}

// LogOut はセッションを破棄します。
func (s *Service) LogOut(token string) error {
	if token == "" {
		return nil
	}
	return s.store.DeleteSession(token)
}

// UserByToken はセッショントークンからユーザーを取得します。
func (s *Service) UserByToken(token string) (*store.User, error) {
	if token == "" {
		return nil, ErrNoSession
	}
	user, err := s.store.UserBySessionToken(token)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrNoSession
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) newSession(user *store.User) (*Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().UTC().Add(SessionDuration)
	if err := s.store.CreateSession(token, user.ID, expiresAt); err != nil {
		return nil, err
	}

	return &Session{Token: token, ExpiresAt: expiresAt, User: user}, nil
}

func generateToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// dummyHash は存在しないユーザーに対しても bcrypt の計算時間をかけるための
// ダミーのハッシュです（パスワード "dummy" の bcrypt ハッシュ）。
var dummyHash = []byte("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy")
