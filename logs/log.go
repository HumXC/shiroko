package logs

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/exp/slog"
)

var dir string

// 返回保存日志的文件夹
func Dir() string {
	return dir
}

var programLevel = new(slog.LevelVar)

var commonOutput io.Writer

type Logger struct {
	output io.Writer
	*slog.Logger
}

func (l *Logger) Output() io.Writer {
	return l.output
}

func init() {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(filepath.Dir(exe), "logs")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
	t := time.Now().Local().Format("2006-01-02_15:04:05")
	logFile := filepath.Join(dir, t+".log")
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	commonOutput = f
	SetLevel(slog.LevelInfo)
}

func SetLevel(level slog.Level) {
	programLevel.Set(level)
}

func Get(tag string) *Logger {
	h := slog.NewTextHandler(
		commonOutput,
		&slog.HandlerOptions{
			Level: programLevel,
		},
	)

	return &Logger{
		output: commonOutput,
		Logger: slog.New(h).With("tag", tag),
	}
}

func File(subdir string) io.Writer {
	return os.Stderr
}

// 同时输出到 os.Stderr
func WriteToStderr() {
	commonOutput = io.MultiWriter(os.Stderr, commonOutput)
}
