// Unrated Coder t.me/Unrated_Coder

package bot

import (
	"bytes"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogCapturer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

var GlobalLogCapturer = &LogCapturer{}

func (lc *LogCapturer) Write(p []byte) (n int, err error) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	_, _ = os.Stdout.Write(p)
	return lc.buf.Write(p)
}

func (lc *LogCapturer) GetLogs() []byte {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	b := make([]byte, lc.buf.Len())
	copy(b, lc.buf.Bytes())
	return b
}

func NewZapLogger() *zap.Logger {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(GlobalLogCapturer),
		zap.DebugLevel,
	)
	return zap.New(core)
}
