package android

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const TMP_DIR = "/data/local/tmp"

type Cmd interface {
	Output() ([]byte, error)
	SetEnv(env []string)
	FullCmd() string
}
type ccmd struct {
	*exec.Cmd
}

func (c *ccmd) Output() ([]byte, error) {
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	c.Cmd.Stderr = &stderr
	c.Cmd.Stdout = &stdout
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
func (c *ccmd) SetEnv(env []string) {
	c.Cmd.Env = env
}
func (c *ccmd) FullCmd() string {
	return fmt.Sprint(strings.Join(c.Cmd.Env, " "), " ", strings.Join(c.Cmd.Args, " "))
}

func Command(cmd string, args ...string) Cmd {
	_cmd := exec.Command(cmd, args...)
	return &ccmd{
		Cmd: _cmd,
	}
}
