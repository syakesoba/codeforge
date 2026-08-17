//go:build ignore

package main

type Config struct {
	Port   string
	DBPath string
	Debug  bool
}

func LoadConfig(getenv func(string) string) Config {
	port := getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/app.db"
	}
	debug := getenv("DEBUG") == "true"

	return Config{Port: port, DBPath: dbPath, Debug: debug}
}

func main() {}
