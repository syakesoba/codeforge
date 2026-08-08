package judge

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

var errTimeout = errors.New("execution timed out")

// dockerRunner はサンドボックスコンテナ内で `go test` を実行します。
type dockerRunner struct {
	Image string
}

func randomContainerName() (string, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return "judge-" + hex.EncodeToString(buf), nil
}

// run はワークスペースディレクトリをマウントしたコンテナで `go test -json ./...` を実行します。
// 戻り値の stdout/stderr はタイムアウト時も可能な限り取得したものを返します。
func (d dockerRunner) run(ctx context.Context, workdir string, timeout time.Duration) (stdout, stderr []byte, runErr error) {
	name, err := randomContainerName()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate container name: %w", err)
	}

	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// --cpus は当初0.5だったが、Gin/GORM等リンクが重いレッスンで
	// 実測7〜23秒とばらつき、まれにタイムアウトしていた。1.0にしたところ
	// 安定して6〜7秒になったため引き上げた（メモリも余裕を持たせて768mに）。
	args := []string{
		"run", "--rm", "--name", name,
		"--network", "none",
		"--memory", "768m", "--memory-swap", "768m",
		"--cpus", "1.0",
		"--pids-limit", "64",
		"--security-opt", "no-new-privileges",
		"--cap-drop", "ALL",
		"-v", workdir + ":/workspace:rw",
		"-w", "/workspace",
		d.Image,
		"go", "test", "-vet=off", "-json", "./...",
	}

	cmd := exec.CommandContext(runCtx, "docker", args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	execErr := cmd.Run()

	if runCtx.Err() == context.DeadlineExceeded {
		// docker run --rm がタイムアウトのkillで後始末できていない可能性があるため、明示的に削除する
		killCtx, killCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer killCancel()
		_ = exec.CommandContext(killCtx, "docker", "rm", "-f", name).Run()
		return outBuf.Bytes(), errBuf.Bytes(), errTimeout
	}

	return outBuf.Bytes(), errBuf.Bytes(), execErr
}
