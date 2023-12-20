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
	if err != nil {
		return nil, err
	}
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
		err = fmt.Errorf("%w: stderr: %sstdout: %s", err, strings.TrimRight(stderr.String(), "\n"), strings.TrimRight(stdout.String(), "\n"))
		return nil, err
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
	c := &Cmd{
		Cmd: _cmd,
	}
	c.SetEnv(nil)
	return c
}

// 去除结尾换行
func TrimEnd[T string | []byte](data T) T {
	if len(data) > 1 && (data[len(data)-1] == '\n') {
		return data[:len(data)-1]
	}
	return data
}

func Model() string {
	c := Command("getprop", "ro.product.model")
	model, _ := c.Output()
	return string(model)
}
