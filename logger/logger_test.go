package logger

import "testing"

func TestWarn(t *testing.T) {
	Warn("warn %d", 123)
}

func TestInfo(t *testing.T) {
	Info("info %d", 123)
}
