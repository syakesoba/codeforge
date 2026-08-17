package main

import "testing"

func TestNewServerDefaults(t *testing.T) {
	s := NewServer()
	if s.Host != "localhost" {
		t.Errorf("Host = %q, want %q", s.Host, "localhost")
	}
	if s.Port != 8080 {
		t.Errorf("Port = %d, want 8080", s.Port)
	}
	if s.Timeout != 30 {
		t.Errorf("Timeout = %d, want 30", s.Timeout)
	}
}

func TestNewServerWithPort(t *testing.T) {
	s := NewServer(WithPort(9000))
	if s.Port != 9000 {
		t.Errorf("Port = %d, want 9000", s.Port)
	}
	if s.Timeout != 30 {
		t.Errorf("Timeout = %d, want 30（変更していないので既定値のはず）", s.Timeout)
	}
	if s.Host != "localhost" {
		t.Errorf("Host = %q, want %q", s.Host, "localhost")
	}
}

func TestNewServerWithMultipleOptions(t *testing.T) {
	s := NewServer(WithPort(9000), WithTimeout(5))
	if s.Port != 9000 {
		t.Errorf("Port = %d, want 9000", s.Port)
	}
	if s.Timeout != 5 {
		t.Errorf("Timeout = %d, want 5", s.Timeout)
	}
}

func TestNewServerOptionsAppliedInOrder(t *testing.T) {
	// 同じ項目を複数回設定した場合、最後に渡したOptionが勝つはず
	s := NewServer(WithPort(1111), WithPort(2222))
	if s.Port != 2222 {
		t.Errorf("Port = %d, want 2222（後勝ち）", s.Port)
	}
}

func TestWithPortReturnsIndependentOption(t *testing.T) {
	// WithPort(port)自体が呼ばれた時点でServerを変更せず、
	// 返された関数が適用されて初めて反映されることを確認する
	opt := WithPort(3000)
	s := &Server{Host: "example", Port: 1, Timeout: 1}
	opt(s)
	if s.Port != 3000 {
		t.Errorf("Port = %d, want 3000", s.Port)
	}
	if s.Host != "example" {
		t.Errorf("Host が意図せず変更されました: %q", s.Host)
	}
}
