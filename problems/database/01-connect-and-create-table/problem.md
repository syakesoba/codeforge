# Lesson 1: データベースに接続してテーブルを作る

## 今回学ぶこと

- Goの標準パッケージ `database/sql` の役割
- ドライバを「インポートするだけ」で登録する書き方（ブランクインポート）
- `sql.Open` で接続を開き、`db.Exec` でテーブルを作る方法

## 解説

### database/sql とドライバ

Goでは、データベース操作の**共通インターフェース**を標準パッケージ `database/sql` が提供し、**実際に特定のDBと通信する部分**は「ドライバ」と呼ばれる外部パッケージが担当します。

このコースではSQLiteを使います。ドライバは `github.com/glebarez/go-sqlite` です。

```go
import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)
```

ドライバのインポートに付いている `_`（アンダースコア）は**ブランクインポート**といいます。「このパッケージの関数は直接呼ばないが、パッケージの初期化処理（`init`関数）だけは実行してほしい」という意味です。ドライバは `init` の中で自分自身を `database/sql` に登録するので、この書き方が必要になります。

`_` を付けずに普通にインポートすると「importしたのに使っていない」というコンパイルエラーになります。

### 接続を開く

```go
db, err := sql.Open("sqlite", ":memory:")
```

- 第1引数はドライバ名（`github.com/glebarez/go-sqlite` は `"sqlite"` という名前で登録されます）
- 第2引数はデータソース名。`":memory:"` はSQLite特有の指定で、**メモリ上に一時的なDBを作る**という意味です。プログラムが終われば消えるので、テストや学習にちょうどよい指定です

> **注意**: `sql.Open` は名前に反して、実際にはまだ接続していません（遅延接続）。接続できるか今すぐ確認したい場合は `db.Ping()` を呼びます。

### テーブルを作る

結果の行を受け取らないSQL（`CREATE TABLE` / `INSERT` / `UPDATE` / `DELETE`）は `db.Exec` で実行します。

```go
func setupDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE users (
			id   INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
```

SQL文はバッククォート（`` ` ``）で囲むと複数行にまたがって書けるので読みやすくなります。

## 演習

`setupDB` 関数を実装してください。次の3つを行います。

1. `sql.Open("sqlite", ":memory:")` でDBを開く（エラーなら `nil, err` を返す）
2. 次のテーブルを作成する

```sql
CREATE TABLE books (
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    title  TEXT NOT NULL,
    author TEXT NOT NULL
)
```

3. 成功したら `db, nil` を返す（テーブル作成に失敗したら `nil, err` を返す）

## ヒント

- ドライバのブランクインポート `_ "github.com/glebarez/go-sqlite"` を忘れないでください（雛形には最初から書いてあります）
- `db.Exec` の戻り値は `(sql.Result, error)` の2つです。今回は結果を使わないので `_, err = db.Exec(...)` と書きます
