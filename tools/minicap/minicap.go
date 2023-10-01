package minicap

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/tools/common"
	"github.com/spf13/cobra"
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

// RegCommand implements common.UseCommand.
func (m *MinicapImpl) RegCommand(cmd *cobra.Command) {
	cmdStart := &cobra.Command{
		Use:   "start",
		Short: "Start minicap",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			rw, err := flags.GetInt32("rw")
			if err != nil {
				panic(err)
			}
			rh, err := flags.GetInt32("rh")
			if err != nil {
				panic(err)
			}
			vw, err := flags.GetInt32("vw")
			if err != nil {
				panic(err)
			}
			vh, err := flags.GetInt32("vh")
			if err != nil {
				panic(err)
			}
			o, err := flags.GetInt32("o")
			if err != nil {
				panic(err)
			}
			if rw == 0 || rh == 0 || vw == 0 || vh == 0 {
				cmd.Help()
				return nil
			}
			_ = m.Start(rw, rh, vw, vh, o)
			for {
				time.Sleep(1 * time.Second)
			}
		},
	}
	flags := cmdStart.Flags()
	flags.Int32("rw", 0, "Real wigth")
	flags.Int32("rh", 0, "Real height")
	flags.Int32("vw", 0, "Virtual wigth")
	flags.Int32("vh", 0, "Virtual height")
	flags.Int32("o", 0, "Orientation")

	cmdInfo := &cobra.Command{
		Use:   "info",
		Short: "Show display info",
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := m.Info()
			if err != nil {
				return err
			}
			_info, _ := json.MarshalIndent(info, "", "    ")
			fmt.Println(string(_info))
			return nil
		},
	}

	cmdCat := &cobra.Command{
		Use:   "cat",
		Short: "Connect and output minicap socket",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader, err := m.Cat()
			if err != nil {
				return err
			}
			_, _ = io.Copy(os.Stdout, reader)
			return nil
		},
	}
	cmd.AddCommand(cmdStart)
	cmd.AddCommand(cmdInfo)
	cmd.AddCommand(cmdCat)
}

var _ IMinicap = &MinicapImpl{}
var _ common.UseCommand = &MinicapImpl{}

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
