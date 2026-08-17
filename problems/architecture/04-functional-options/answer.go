//go:build ignore

package main

type Server struct {
	Host    string
	Port    int
	Timeout int
}

type Option func(*Server)

func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

func WithTimeout(timeout int) Option {
	return func(s *Server) {
		s.Timeout = timeout
	}
}

func NewServer(opts ...Option) *Server {
	s := &Server{Host: "localhost", Port: 8080, Timeout: 30}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func main() {}
