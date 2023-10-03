package common

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type Tool interface {
	Base() BaseTool
}
type BaseTool interface {
	Name() string
	// 简短的描述
	Description() string
	// 检查工具是否可用
	// error 为 nil 时表示可用
	Health() error
	// 将工具部署到设备
	Install() error
	Uninstall() error
	// 返回运行时需要的环境变量，例如: LD_LIBRARY_PATH=/data/local/tmp/
	Env() []string
	// 返回可执行文件的路径
	Exe() string
	// 返回需要的默认参数
	Args() []string
	// 返回所有部署在设备上的文件
	Files() []string
	// 在使用前进行初始化，由 tools 包调用
	Init()
}
type UseCommand interface {
	RegCommand(*cobra.Command)
}

// 检查命令是否在 $PATH 中
func CommandHealth(cmd string) error {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}
	return nil
}

func FullCommand(cmd *exec.Cmd, env ...string) string {
	r := strings.Join(cmd.Args, " ")
	if len(env) != 0 {
		r = strings.Join(env, " ") + " " + r
	}
	return r
}

// 检查所有文件是否存在，如果都存在则返回 nil
// 如果有文件不存在则返回一个包含所有不存在文件的错误
func HealthWithFiles(files []string) error {
	notFiends := []string{}
	for _, file := range files {
		_, err := os.Stat(file)
		if err == nil {
			continue
		}
		if errors.Is(err, os.ErrNotExist) {
			notFiends = append(notFiends, file)
		} else {
			return fmt.Errorf("file stat error: %s: %w", file, err)
		}
	}
	if len(notFiends) == 0 {
		return nil
	}
	return fmt.Errorf("file not exist: [%s]", strings.Join(notFiends, "] ["))
}
