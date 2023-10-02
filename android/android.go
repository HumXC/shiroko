package android

import (
	"bytes"
	"fmt"
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
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	c.Stderr = &stderr
	c.Stdout = &stdout
	err := c.Cmd.Run()
	if err != nil {
		err = fmt.Errorf("%w: stderr: %s", err, strings.TrimRight(stderr.String(), "\r\n"))
		return stdout.Bytes(), err
	}
	b := stdout.Bytes()
	if len(b) > 1 && (b[len(b)-1] == '\n' || b[len(b)-1] == '\r') {
		b = b[:len(b)-1]
	}
	return b, nil
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
