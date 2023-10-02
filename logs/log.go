package logs

import (
	"io"
	"os"

	"golang.org/x/exp/slog"
)

var programLevel = new(slog.LevelVar)

func init() {
	programLevel.Set(slog.LevelDebug)
}
func Get() *slog.Logger {
	h := slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{
			Level: programLevel,
		},
	)

	return slog.New(h)
}

func File(subdir string) io.Writer {
	return os.Stderr
}
