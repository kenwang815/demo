package log_test

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github/demo/utils/log"
)

func TestShowFileInfo(t *testing.T) {
	wantMsg := "\x1b[37mDEBU\x1b[0m[0000] debug" +
		"                                         " +
		"\x1b[37mfile\x1b[0m=logger_test.go \x1b[37mfunc\x1b[0m=log_test.TestShowFileInfo.func1 \x1b[37mline\x1b[0m=22\n"
	log.Init("development", "", "debug")
	log.SetFormat(true, false)
	re := captureOutput(func() {
		log.Debug("debug")
	})
	assert.Equal(t, wantMsg, re)
}

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}
