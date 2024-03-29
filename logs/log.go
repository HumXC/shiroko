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

var level = slog.LevelInfo

var commonOutput io.Writer = os.Stderr

type Logger struct {
	output io.Writer
	*slog.Logger
}

func (l *Logger) Output() io.Writer {
	return l.output
}

var inited = false

func Init(logLevel slog.Level, outputToFile bool) {
	if inited {
		return
	}
	inited = true
	level = logLevel
	if !outputToFile {
		return
	}
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
	commonOutput = io.MultiWriter(os.Stderr, f)
}

func Get(tag string) *Logger {
	if !inited {
		panic("logs.Init() not called")
	}
	h := slog.NewTextHandler(
		commonOutput,
		&slog.HandlerOptions{
			Level: level,
		},
	)

	return &Logger{
		output: commonOutput,
		Logger: slog.New(h).With("tag", tag),
	}
}
