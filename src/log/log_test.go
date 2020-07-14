package log

import (
	"go.uber.org/zap"
	"testing"
)

func TestInitLogger(t *testing.T) {
	s := []string{
		"hello info",
		"hello error",
		"hello debug",
		"hello fatal",
	}
	Logger.Info("info:", zap.String("s", s[0]))
	Logger.Error("info:", zap.String("s", s[1]))
	Logger.Debug("info:", zap.String("s", s[2]))
	Logger.Fatal("info:", zap.String("s", s[3]))
}
