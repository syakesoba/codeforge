package main

// Config はアプリケーションの設定です。
type Config struct {
	Port   string
	DBPath string
	Debug  bool
}

// LoadConfig は環境変数から設定を読み込みます。
// テストしやすくするため、os.Getenv を直接呼ぶのではなく、
// 「環境変数名を受け取って値を返す関数」を引数として受け取ります
// （本番では os.Getenv を渡し、テストでは偽の関数を渡せます）。
//
// 各環境変数が空文字列（未設定）の場合は、以下のデフォルト値を使ってください。
//   - PORT: デフォルト "8080"
//   - DB_PATH: デフォルト "data/app.db"
//   - DEBUG: "true" ならtrue、それ以外（未設定を含む）はfalse
func LoadConfig(getenv func(string) string) Config {
	// TODO:
	// port := getenv("PORT") 。空文字列なら "8080" を使う
	// dbPath := getenv("DB_PATH") 。空文字列なら "data/app.db" を使う
	// debug := getenv("DEBUG") == "true"
	// Config{Port: port, DBPath: dbPath, Debug: debug} を返す
	return Config{}
}

func main() {}
