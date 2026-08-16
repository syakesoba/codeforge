package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

// fakeShutdowner は Shutdowner の偽実装です。
// immediate が nil でなければすぐにそのエラー（nilならnil）を返します。
// immediate が設定されておらず block が true の場合は、
// 本物の *http.Server.Shutdown と同じように、ctx.Done() を待ってから
// ctx.Err() を返します（処理に時間がかかりすぎてタイムアウトした状況を再現する）。
type fakeShutdowner struct {
	immediate   error
	returnImmed bool
	block       bool
}

func (f *fakeShutdowner) Shutdown(ctx context.Context) error {
	if f.returnImmed {
		return f.immediate
	}
	if f.block {
		<-ctx.Done()
		return ctx.Err()
	}
	return nil
}

func TestGracefulShutdownSuccess(t *testing.T) {
	s := &fakeShutdowner{returnImmed: true, immediate: nil}
	if err := GracefulShutdown(s, 100*time.Millisecond); err != nil {
		t.Fatalf("GracefulShutdownがエラーを返しました: %v", err)
	}
}

func TestGracefulShutdownPropagatesError(t *testing.T) {
	wantErr := errors.New("listener already closed")
	s := &fakeShutdowner{returnImmed: true, immediate: wantErr}
	err := GracefulShutdown(s, 100*time.Millisecond)
	if !errors.Is(err, wantErr) {
		t.Fatalf("GracefulShutdownがエラーを伝播していません。実際: %v", err)
	}
}

func TestGracefulShutdownTimesOut(t *testing.T) {
	// shutdownerがcontextの期限に反応して初めて終了するケース。
	// timeoutを短く設定し、fakeShutdownerがctx.Done()で解放されて
	// ctx.Err()（DeadlineExceeded）を返すことを確認する。
	s := &fakeShutdowner{block: true}

	start := time.Now()
	err := GracefulShutdown(s, 50*time.Millisecond)
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("タイムアウトなのにDeadlineExceededが返りませんでした。実際: %v", err)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("GracefulShutdownの完了に時間がかかりすぎています: %v（timeoutが正しく設定されていますか？）", elapsed)
	}
}
