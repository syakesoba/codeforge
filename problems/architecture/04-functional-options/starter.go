package main

// Server はいくつかの設定項目を持つ構造体です。
type Server struct {
	Host    string
	Port    int
	Timeout int // 秒
}

// Option は Server の設定を変更する関数です。
// この型を使い、「設定したい項目だけ」を呼び出し側が指定できるようにします
// （関数オプションパターン）。
type Option func(*Server)

// WithPort はポート番号を上書きする Option を返します。
func WithPort(port int) Option {
	// TODO: return func(s *Server) { s.Port = port } を返す
	return nil
}

// WithTimeout はタイムアウト秒数を上書きする Option を返します。
func WithTimeout(timeout int) Option {
	// TODO: WithPortと同様に、s.Timeout = timeout を行うOptionを返す
	return nil
}

// NewServer は、デフォルト値（Host: "localhost", Port: 8080, Timeout: 30）を持つ
// Server を作り、渡された opts を順番に適用して返します。
func NewServer(opts ...Option) *Server {
	// TODO:
	// s := &Server{Host: "localhost", Port: 8080, Timeout: 30} でデフォルト値を用意する
	// for _, opt := range opts { opt(s) } で1つずつ適用する
	// s を返す
	return nil
}

func main() {}
