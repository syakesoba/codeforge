//go:build ignore

package main

import "net/http"

type Config struct {
	Port  string
	Debug bool
}

type HealthChecker func() error

func HealthHandler(check HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := check(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func DebugHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("debug mode enabled"))
}

func NewApp(cfg Config, check HealthChecker, routes map[string]http.HandlerFunc) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", HealthHandler(check))
	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}
	if cfg.Debug {
		mux.HandleFunc("/debug", DebugHandler)
	}
	return mux
}

func main() {}
