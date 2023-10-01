package minicap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/tools/common"
)

type Info struct {
	Id       int32   `json:"id"`
	Width    int32   `json:"width"`
	Height   int32   `json:"height"`
	Xdpi     float32 `json:"xdpi"`
	Ydpi     float32 `json:"ydpi"`
	Size     float32 `json:"size"`
	Density  float32 `json:"density"`
	Fps      float32 `json:"fps"`
	Secure   bool    `json:"secure"`
	Rotation int32   `json:"rotation"`
}

type IMinicap interface {
	Info() (Info, error)
	Start(RWidth, RHeight, VWidth, VHeight, Orientation int32) error
	Stop() error
	Cat() (io.ReadCloser, error)
}

var Minicap *MinicapImpl = New()

type MinicapImpl struct {
	Base common.BaseTool
	proc *os.Process
	conn net.Conn
}

var _ IMinicap = &MinicapImpl{}

func (m *MinicapImpl) Info() (Info, error) {
	result := Info{}
	cmd := android.Command(m.Base.Exe(), append(m.Base.Args(), "-i")...)
	cmd.SetEnv(m.Base.Env())
	output, err := cmd.Output()
	if err != nil {
		return result, fmt.Errorf("failed to get info: %w", err)
	}
	err = json.Unmarshal(output, &result)
	if err != nil {
		return result, fmt.Errorf("failed to get info: %w", err)
	}
	return result, nil
}

func (m *MinicapImpl) Start(RWidth, RHeight, VWidth, VHeight, Orientation int32) error {
	if m.proc != nil {
		return fmt.Errorf("minicap already running")
	}
	args := append(m.Base.Args(), "-P", fmt.Sprintf("%dx%d@%dx%d/%d", RWidth, RHeight, VWidth, VHeight, Orientation))
	cmd := exec.Command(m.Base.Exe(), args...)
	cmd.Env = m.Base.Env()
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("minicap start error: %s: %w", common.FullCommand(cmd), err)
	}
	m.proc = cmd.Process
	return nil
}

func (m *MinicapImpl) Stop() error {
	if m.conn != nil {
		conn := m.conn
		m.conn = nil
		conn.Close()
	}
	if m.proc != nil {
		proc := m.proc
		m.proc = nil
		err := proc.Kill()
		if err != nil {
			return fmt.Errorf("minicap kill error: %d: %w", m.proc.Pid, err)
		}
	}
	return nil
}

func (m *MinicapImpl) Cat() (io.ReadCloser, error) {
	if m.proc == nil {
		return nil, errors.New("minicap not running")
	}

	conn, err := net.Dial("unix", "@minicap")
	if err != nil {
		return nil, fmt.Errorf("failed to connect minicap: %w", err)
	}
	return conn, nil
}

func New() *MinicapImpl {
	m := &MinicapImpl{
		Base: &minicapBase{},
	}
	return m
}
