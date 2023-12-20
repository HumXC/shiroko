package shell

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/tools/common"
)

// 一些常用的命令
var Shell *ShellImpl = nil

func Init() {
	Shell = New()
}

type IShell interface {
	// 在安卓里执行命令，超时之后会主动终止，如果 timeout 为 0，则持续等待执行
	Run(cmd string, timeoutMs int32) ([]byte, error)
	Push(filename string, data io.Reader) error
	Pull(filename string) (io.ReadCloser, error)
	HttpGet(url, dist string, timeoutMs int32) error
	Install(apkpath string) error
	Uninstall(pkgname string) error
	// pm list packages
	ListApps() ([]string, error)
	// am start -n
	StartApp(active string) error
	// am force-stop
	StopApp(pkgname string) error
	// getprop
	Getprop(key string) (string, error)
}
type ShellImpl struct {
	base common.BaseTool
}

// HttpGet implements IShell.
func (*ShellImpl) HttpGet(url string, dist string, timeout int32) error {
	client := http.Client{Timeout: time.Duration(timeout) * time.Millisecond}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(dist)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// Getprop implements IShell.
func (s *ShellImpl) Getprop(key string) (string, error) {
	c := android.Command("getprop", key)
	out, err := c.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Install implements IShell.
func (s *ShellImpl) Install(apkpath string) error {
	c := android.Command("pm", "install", apkpath)
	_, err := c.Output()
	return err
}

// ListApps implements IShell.
func (s *ShellImpl) ListApps() ([]string, error) {
	c := android.Command("pm", "list", "packages")
	out, err := c.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	for i := 0; i < len(lines); i++ {
		// 去除开头 "package:"
		lines[i] = lines[i][8:]
	}
	return lines, nil
}

// Pull implements IShell.
func (s *ShellImpl) Pull(filename string) (io.ReadCloser, error) {
	src, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return src, nil
}

// Push implements IShell.
func (s *ShellImpl) Push(filename string, data io.Reader) error {
	if data == nil {
		return errors.New("data is nil")
	}
	if filename == "" {
		return errors.New("filename is empty")
	}
	dist, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dist.Close()
	_, err = io.Copy(dist, data)
	return err
}

// Run implements IShell.
func (s *ShellImpl) Run(cmd string, timeout int32) ([]byte, error) {
	c := android.Command("sh", "-c", cmd)
	stderr := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	c.Stderr = stderr
	c.Stdout = stdout
	err := c.Cmd.Start()
	if err != nil {
		return nil, err
	}
	parseErr := func(err error) error {
		if err != nil {
			err = fmt.Errorf("%w: stderr: %s", err, strings.TrimRight(stderr.String(), "\n"))
			return err
		}
		return nil
	}
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				if c.Cmd.Process != nil {
					c.Cmd.Process.Kill()
				}
				return stdout.Bytes(), errors.New("timeout")
			default:
				if c.Cmd.ProcessState == nil || !c.Cmd.ProcessState.Exited() {
					time.Sleep(10 * time.Millisecond)
					continue
				} else {
					return stdout.Bytes(), parseErr(c.Cmd.Err)
				}
			}
		}
	} else {
		err = c.Cmd.Wait()
		return stdout.Bytes(), parseErr(err)
	}
}

// StartApp implements IShell.
func (s *ShellImpl) StartApp(active string) error {
	c := android.Command("am", "start", "-n", active)
	out, err := c.Output()
	fmt.Println(out)
	return err
}

// StopApp implements IShell.
func (s *ShellImpl) StopApp(pkgname string) error {
	c := android.Command("am", "force-stop", pkgname)
	_, err := c.Output()
	return err
}

// Uninstall implements IShell.
func (s *ShellImpl) Uninstall(pkgname string) error {
	c := android.Command("pm", "uninstall", pkgname)
	_, err := c.Output()
	if err != nil {
		return err
	}
	return nil
}

var _ IShell = &ShellImpl{}
var _ common.Tool = &ShellImpl{}

func New() *ShellImpl {
	return &ShellImpl{
		base: &shellBase{Cmd: "sh"},
	}
}

// Base implements common.Tool.
func (s *ShellImpl) Base() common.BaseTool {
	return s.base
}
