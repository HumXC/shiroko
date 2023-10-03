package android

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/HumXC/shiroko/tools/common"
)

var Sysenv = os.Environ()

const TMP_DIR = "/data/local/tmp"

type Cmd struct {
	*exec.Cmd
	CustomEnv []string
}

func (c *Cmd) Output() ([]byte, error) {
	reader, err := c.OutputReader()
	b, _ := io.ReadAll(reader)
	return TrimEnd(b), err
}
func (c *Cmd) OutputReader() (io.Reader, error) {
	stderr := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	c.Stderr = stderr
	c.Stdout = stdout
	err := c.Cmd.Run()
	if err != nil {
		err = fmt.Errorf("%w: stderr: %s", err, strings.TrimRight(stderr.String(), "\r\n"))
		return stdout, err
	}
	return stdout, nil
}
func (c *Cmd) SetEnv(env []string) {
	c.CustomEnv = env
	c.Cmd.Env = append(Sysenv, env...)
}
func (c *Cmd) FullCmd() string {
	return common.FullCommand(c.Cmd, c.CustomEnv...)
}
func Command(cmd string, args ...string) *Cmd {
	_cmd := exec.Command(cmd, args...)
	return &Cmd{
		Cmd: _cmd,
	}
}

// 去除结尾换行
func TrimEnd[T string | []byte](data T) T {
	if len(data) > 1 && (data[len(data)-1] == '\n' || data[len(data)-1] == '\r') {
		return data[:len(data)-1]
	}
	return data
}
