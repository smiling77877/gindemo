package syntax

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestLoggerFile(t *testing.T) {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"logs/info.log"}
	l, err := cfg.Build()
	require.NoError(t, err)
	l.Error("hello")
}
