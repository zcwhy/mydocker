package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	LogInit()

	logger.Sugar().Info("1111111111")
}
