// Package store はユーザー・セッション・学習進捗をSQLiteに永続化します。
package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// ErrNotFound は該当するレコードが存在しないことを表します。
var ErrNotFound = errors.New("not found")

// ErrEmailTaken は既に登録済みのメールアドレスであることを表します。
var ErrEmailTaken = errors.New("email already registered")

// Store はアプリケーションのデータアクセスをまとめた型です。
type Store struct {
	db *sql.DB
}

// User は登録済みユーザーです。
type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

// Open は指定パスのSQLiteファイルを開き、必要なテーブルを作成します。
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// SQLiteは書き込みを直列化するため、接続数を1に絞って
	// "database is locked" を避ける。
	db.SetMaxOpenConns(1)

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

// Close はデータベース接続を閉じます。
func (s *Store) Close() error {
	return s.db.Close()
}

// Ping はデータベースへの疎通を確認します。ヘルスチェックエンドポイントから使います。
func (s *Store) Ping() error {
	return s.db.Ping()
}

func (s *Store) migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			email         TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at    DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			token      TEXT PRIMARY KEY,
			user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			expires_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS progress (
			user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			problem_id TEXT NOT NULL,
			passed_at  DATETIME NOT NULL,
			PRIMARY KEY (user_id, problem_id)
		)`,
		`CREATE TABLE IF NOT EXISTS code_drafts (
			user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			problem_id TEXT NOT NULL,
			code       TEXT NOT NULL,
			updated_at DATETIME NOT NULL,
			PRIMARY KEY (user_id, problem_id)
		)`,
	}
	for _, stmt := range stmts {
		if _, err := s.db.Exec(stmt); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	return nil
}

// CreateUser は新しいユーザーを登録します。
// 既に同じメールアドレスが登録済みの場合は ErrEmailTaken を返します。
func (s *Store) CreateUser(email, passwordHash string) (*User, error) {
	var exists int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ?`, email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if exists > 0 {
		return nil, ErrEmailTaken
	}

	result, err := s.db.Exec(
		`INSERT INTO users (email, password_hash, created_at) VALUES (?, ?, ?)`,
		email, passwordHash, time.Now().UTC(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user id: %w", err)
	}

	return &User{ID: id, Email: email, PasswordHash: passwordHash}, nil
}

// UserByEmail はメールアドレスでユーザーを取得します。
func (s *Store) UserByEmail(email string) (*User, error) {
	var u User
	err := s.db.QueryRow(
		`SELECT id, email, password_hash FROM users WHERE email = ?`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &u, nil
}

// CreateSession はセッションを作成します。
func (s *Store) CreateSession(token string, userID int64, expiresAt time.Time) error {
	_, err := s.db.Exec(
		`INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)`,
		token, userID, expiresAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

// UserBySessionToken は有効なセッショントークンからユーザーを取得します。
// トークンが存在しない、または期限切れの場合は ErrNotFound を返します。
func (s *Store) UserBySessionToken(token string) (*User, error) {
	var u User
	var expiresAt time.Time
	err := s.db.QueryRow(`
		SELECT u.id, u.email, u.password_hash, s.expires_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token = ?
	`, token).Scan(&u.ID, &u.Email, &u.PasswordHash, &expiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if time.Now().UTC().After(expiresAt) {
		// 期限切れセッションはこの機会に掃除しておく
		_ = s.DeleteSession(token)
		return nil, ErrNotFound
	}

	return &u, nil
}

// DeleteSession はセッションを削除します（ログアウト）。
func (s *Store) DeleteSession(token string) error {
	if _, err := s.db.Exec(`DELETE FROM sessions WHERE token = ?`, token); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// MarkProgress は問題の合格を記録します。既に記録済みの場合は何もしません。
func (s *Store) MarkProgress(userID int64, problemID string) error {
	_, err := s.db.Exec(`
		INSERT INTO progress (user_id, problem_id, passed_at) VALUES (?, ?, ?)
		ON CONFLICT (user_id, problem_id) DO NOTHING
	`, userID, problemID, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to mark progress: %w", err)
	}
	return nil
}

// ProgressEntry は合格済み問題1件分の記録です。
type ProgressEntry struct {
	ProblemID string
	PassedAt  time.Time
}

// CompletedProgress はユーザーが合格済みの問題を、合格日時付きで返します。
func (s *Store) CompletedProgress(userID int64) ([]ProgressEntry, error) {
	rows, err := s.db.Query(`SELECT problem_id, passed_at FROM progress WHERE user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list progress: %w", err)
	}
	defer rows.Close()

	entries := []ProgressEntry{}
	for rows.Next() {
		var e ProgressEntry
		if err := rows.Scan(&e.ProblemID, &e.PassedAt); err != nil {
			return nil, fmt.Errorf("failed to scan progress: %w", err)
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate progress: %w", err)
	}
	return entries, nil
}

// SaveDraft はレッスンの編集中コードを保存します（自動保存用）。
func (s *Store) SaveDraft(userID int64, problemID, code string) error {
	_, err := s.db.Exec(`
		INSERT INTO code_drafts (user_id, problem_id, code, updated_at) VALUES (?, ?, ?, ?)
		ON CONFLICT (user_id, problem_id) DO UPDATE SET code = excluded.code, updated_at = excluded.updated_at
	`, userID, problemID, code, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to save draft: %w", err)
	}
	return nil
}

// GetDraft は保存済みの編集中コードを返します。保存されていない場合は ok=false を返します。
func (s *Store) GetDraft(userID int64, problemID string) (code string, ok bool, err error) {
	err = s.db.QueryRow(
		`SELECT code FROM code_drafts WHERE user_id = ? AND problem_id = ?`, userID, problemID,
	).Scan(&code)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("failed to get draft: %w", err)
	}
	return code, true, nil
}
