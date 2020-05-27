package wsep

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"cdr.dev/slog/sloggers/slogtest/assert"
)

func TestLocalExec(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	testExecer(ctx, t, LocalExecer{})
}

func testExecer(ctx context.Context, t *testing.T, execer Execer) {
	process, err := execer.Start(ctx, Command{
		Command: "pwd",
	})
	assert.Success(t, "start local cmd", err)
	var (
		stderr = process.Stderr()
		stdout = process.Stdout()
		wg     sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		stdoutByt, err := ioutil.ReadAll(stdout)
		assert.Success(t, "read stdout", err)
		wd, err := os.Getwd()
		assert.Success(t, "get real working dir", err)

		assert.Equal(t, "stdout", wd, strings.TrimSuffix(string(stdoutByt), "\n"))
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()

		stderrByt, err := ioutil.ReadAll(stderr)
		assert.Success(t, "read stderr", err)
		assert.True(t, "len stderr", len(stderrByt) == 0)
	}()

	wg.Wait()
	err = process.Wait()
	assert.Success(t, "wait for process to complete", err)
}