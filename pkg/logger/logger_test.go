package logger

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("info logs to buffer", func(t *testing.T) {
		var buf bytes.Buffer
		Initialize(&buf)

		Info("Test message", "key", "value")

		output := buf.String()

		if !strings.Contains(output, `"level":"INFO"`) {
			t.Errorf("Expected level=INFO, got: %s", output)
		}
		if !strings.Contains(output, `"msg":"Test message"`) {
			t.Errorf("Expected msg 'Test message', got: %s", output)
		}
		if !strings.Contains(output, `"key":"value"`) {
			t.Errorf("Expected key value, got: %s", output)
		}
	})

	t.Run("Error logs with error key", func(t *testing.T) {
		var buf bytes.Buffer
		Initialize(&buf)

		testErr := errors.New("something went wrong")
		Error("Operation failed", testErr, "user_id", 123)

		output := buf.String()
		if !strings.Contains(output, `"level":"ERROR"`) {
			t.Errorf("Expected level=ERROR, got: %s", output)
		}
		if !strings.Contains(output, `"msg":"Operation failed"`) {
			t.Errorf("Expected msg 'Operation failed', got: %s", output)
		}
		if !strings.Contains(output, `"error":"something went wrong"`) {
			t.Errorf("Expected error message, got: %s", output)
		}
		if !strings.Contains(output, `"user_id":123`) {
			t.Errorf("Expected user_id 123, got: %s", output)
		}
	})

	t.Run("Warn logs with WARN level", func(t *testing.T) {
		var buf bytes.Buffer
		Initialize(&buf)
		Warn("Some issue", "resource_id", 123)
		output := buf.String()
		if !strings.Contains(output, `"level":"WARN"`) {
			t.Errorf("Expected level WARN, got: %s", output)
		}
		if !strings.Contains(output, `"msg":"Some issue"`) {
			t.Errorf("Expected msg 'Potential issue', got: %s", output)
		}
		if !strings.Contains(output, `"resource_id":123`) {
			t.Errorf("Expected resource_id 456, got: %s", output)
		}
	})

	t.Run("With creates child logger with context", func(t *testing.T) {
		// 1. Setup
		var buf bytes.Buffer
		Initialize(&buf)

		// 2. Act
		child := With("req_id", "abc-123")
		child.Info("Child log")

		// 3. Assert
		output := buf.String()
		if !strings.Contains(output, `"req_id":"abc-123"`) {
			t.Errorf("Expected context req_id, got: %s", output)
		}
	})
}

func TestLazyLoggerInitialization(t *testing.T) {
	originalLog := Log
	originalStdout := os.Stdout
	defer func() {
		Log = originalLog
		os.Stdout = originalStdout
	}()

	captureStdout := func(f func()) string {
		r, w, _ := os.Pipe()
		os.Stdout = w
		f()
		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		return buf.String()
	}

	t.Run("InitLogger defaults to stdout", func(t *testing.T) {
		Log = nil
		output := captureStdout(func() {
			InitLogger()
			Info("default init check")
		})

		if Log == nil {
			t.Error("InitLogger should initialize global Log")
		}
		if !strings.Contains(output, "default init check") {
			t.Error("InitLogger should write to stdout by default")
		}
	})

	t.Run("Info triggers lazy init", func(t *testing.T) {
		Log = nil
		output := captureStdout(func() {
			Info("lazy info")
		})
		if Log == nil {
			t.Error("Info should have triggered InitLogger")
		}
		if !strings.Contains(output, "lazy info") {
			t.Error("Info should have logged to stdout")
		}
	})

	t.Run("Error triggers lazy init", func(t *testing.T) {
		Log = nil
		output := captureStdout(func() {
			Error("lazy error", errors.New("oops"))
		})

		if Log == nil {
			t.Error("Error should have triggered InitLogger")
		}
		if !strings.Contains(output, "lazy error") {
			t.Error("Error should have logged to stdout")
		}
	})

	t.Run("Warn triggers lazy init", func(t *testing.T) {
		Log = nil
		output := captureStdout(func() {
			Warn("lazy warn")
		})
		if Log == nil {
			t.Error("Warn should have triggered InitLogger")
		}
		if !strings.Contains(output, "lazy warn") {
			t.Error("Warn should have logged to stdout")
		}
	})

	t.Run("With triggers lazy init", func(t *testing.T) {
		Log = nil

		_ = With("context", "lazy")

		if Log == nil {
			t.Error("With should have triggered InitLogger")
		}
	})
}
