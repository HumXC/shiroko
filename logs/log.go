package logs

import (
	"io"
	"os"

	"golang.org/x/exp/slog"
)

// TODO： 更好的 log
// 增加 tag 的功能
// 增加命令行不输出日志到终端
// 增加文件输出
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
