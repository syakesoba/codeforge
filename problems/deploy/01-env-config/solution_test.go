package main

import "testing"

func fakeGetenv(values map[string]string) func(string) string {
	return func(key string) string {
		return values[key]
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	cfg := LoadConfig(fakeGetenv(map[string]string{}))

	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
	if cfg.DBPath != "data/app.db" {
		t.Errorf("DBPath = %q, want %q", cfg.DBPath, "data/app.db")
	}
	if cfg.Debug != false {
		t.Errorf("Debug = %v, want false", cfg.Debug)
	}
}

func TestLoadConfigOverridesPort(t *testing.T) {
	cfg := LoadConfig(fakeGetenv(map[string]string{"PORT": "3000"}))
	if cfg.Port != "3000" {
		t.Errorf("Port = %q, want %q", cfg.Port, "3000")
	}
	if cfg.DBPath != "data/app.db" {
		t.Errorf("DBPath = %q, want %q（他は既定値のはず）", cfg.DBPath, "data/app.db")
	}
}

func TestLoadConfigOverridesDBPath(t *testing.T) {
	cfg := LoadConfig(fakeGetenv(map[string]string{"DB_PATH": "/var/lib/app/data.db"}))
	if cfg.DBPath != "/var/lib/app/data.db" {
		t.Errorf("DBPath = %q, want %q", cfg.DBPath, "/var/lib/app/data.db")
	}
}

func TestLoadConfigDebugTrue(t *testing.T) {
	cfg := LoadConfig(fakeGetenv(map[string]string{"DEBUG": "true"}))
	if cfg.Debug != true {
		t.Errorf("Debug = %v, want true", cfg.Debug)
	}
}

func TestLoadConfigDebugInvalidValue(t *testing.T) {
	// "1" のような、trueでもfalseでもない値は false 扱いにする
	cfg := LoadConfig(fakeGetenv(map[string]string{"DEBUG": "1"}))
	if cfg.Debug != false {
		t.Errorf("Debug = %v, want false（\"true\"以外はfalse扱い）", cfg.Debug)
	}
}

func TestLoadConfigAllOverridden(t *testing.T) {
	cfg := LoadConfig(fakeGetenv(map[string]string{
		"PORT":    "9090",
		"DB_PATH": "custom.db",
		"DEBUG":   "true",
	}))
	if cfg.Port != "9090" || cfg.DBPath != "custom.db" || !cfg.Debug {
		t.Errorf("cfg = %+v, 期待した値と異なります", cfg)
	}
}
