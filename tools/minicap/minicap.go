package minicap

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/logs"
	"github.com/HumXC/shiroko/tools/common"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var log *slog.Logger

func init() {
	log = logs.Get()
}

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
	Start(rWidth, rHeight, vWidth, vHeight, orientation, rate int32) error
	Jpg(rWidth, rHeight, vWidth, vHeight, orientation, quality int32) ([]byte, error)
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
		Short: "Start minicap, if either vw or vh is set to 0, they will be automatically set.",
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
			rate, err := flags.GetInt32("r")
			if err != nil {
				panic(err)
			}
			if vw == 0 {
				vw = rw
			}
			if vh == 0 {
				vh = rh
			}
			if rw == 0 || rh == 0 || vw == 0 || vh == 0 {
				cmd.Help()
				return nil
			}
			_ = m.Start(rw, rh, vw, vh, o, rate)
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
	flags.Int32("o", 0, "Orientation (0-3)")
	flags.Int32("r", 30, "Frame rate (frames/s)")

	cmdJpg := &cobra.Command{
		Use:   "jpg",
		Short: "Get screenshot and output to JPEG, if either vw or vh is set to 0, they will be automatically set.",
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
			quality, err := flags.GetInt32("q")
			if err != nil {
				panic(err)
			}
			if vw == 0 {
				vw = rw
			}
			if vh == 0 {
				vh = rh
			}
			if rw == 0 || rh == 0 || vw == 0 || vh == 0 {
				cmd.Help()
				return nil
			}
			data, err := m.Jpg(rw, rh, vw, vh, o, quality)
			if err != nil {
				return err
			}
			os.Stdout.Write(data)
			return nil
		},
	}
	flags = cmdJpg.Flags()
	flags.Int32("rw", 0, "Real wigth")
	flags.Int32("rh", 0, "Real height")
	flags.Int32("vw", 0, "Virtual wigth")
	flags.Int32("vh", 0, "Virtual height")
	flags.Int32("o", 0, "Orientation")
	flags.Int32("q", 100, "Jpg quality")

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
	cmd.AddCommand(cmdJpg)
	cmd.AddCommand(cmdInfo)
	cmd.AddCommand(cmdCat)
}

var _ IMinicap = &MinicapImpl{}
var _ common.UseCommand = &MinicapImpl{}

func (m *MinicapImpl) Info() (Info, error) {
	log.Info("Get info")
	result := Info{}
	cmd := android.Command(m.Base.Exe(), append(m.Base.Args(), "-i")...)
	cmd.SetEnv(m.Base.Env())
	log.Debug("Run command", "command", cmd.FullCmd())
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

// FIXME：ERROR: invalid format for -P, need <w>x<h>@<w>x<h>/{0|90|180|270}
// 貌似 apk 的参数和 cpp 二进制程序的 orientation 参数不一样
func (m *MinicapImpl) Start(rWidth, rHeight, vWidth, vHeight, orientation, rate int32) error {
	log.Info("Start minicap", "rWidth", rWidth, "rHeight", rHeight, "vWidth", vWidth, "vHeight", vHeight, "orientation", orientation, "rate", rate)
	if m.proc != nil {
		return fmt.Errorf("minicap already running")
	}
	args := append(
		m.Base.Args(),
		"-P",
		fmt.Sprintf("%dx%d@%dx%d/%d", rWidth, rHeight, vWidth, vHeight, orientation),
		"-r",
		strconv.Itoa(int(rate)),
	)
	cmd := android.Command(m.Base.Exe(), args...)
	cmd.SetEnv(m.Base.Env())
	log.Debug("Run command", "command", cmd.FullCmd())
	out := logs.File("nimicap")
	cmd.Stderr = out
	cmd.Stdout = out
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("minicap start error: %s: %w", cmd.FullCmd(), err)
	}
	m.proc = cmd.Process
	return nil
}

// Jpg implements IMinicap.
func (m *MinicapImpl) Jpg(rWidth int32, rHeight int32, vWidth int32, vHeight int32, orientation int32, quality int32) ([]byte, error) {
	log.Info("Minicap screenshot", "rWidth", rWidth, "rHeight", rHeight, "vWidth", vWidth, "vHeight", vHeight, "orientation", orientation, "quality", quality)
	args := append(
		m.Base.Args(),
		"-P",
		fmt.Sprintf("%dx%d@%dx%d/%d", rWidth, rHeight, vWidth, vHeight, orientation),
		"-s",
		"-Q",
		strconv.Itoa(int(quality)),
	)
	cmd := android.Command(m.Base.Exe(), args...)
	cmd.SetEnv(m.Base.Env())
	log.Debug("Run command", "command", cmd.FullCmd())
	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("minicap start error: %s: %w", cmd.FullCmd(), err)
	}
	return data, nil
}

func (m *MinicapImpl) Stop() error {
	log.Info("Minicap stop")
	if m.conn != nil {
		conn := m.conn
		m.conn = nil
		log.Debug("Close connection", "conn", conn.LocalAddr())
		_ = conn.Close()
	}
	if m.proc != nil {
		proc := m.proc
		m.proc = nil
		log.Debug("Kill minicap process", "pid", proc.Pid)
		err := proc.Kill()
		if err != nil {
			return fmt.Errorf("minicap kill error: %d: %w", m.proc.Pid, err)
		}
	}
	return nil
}

func (m *MinicapImpl) Cat() (io.ReadCloser, error) {
	log.Info("Connect minicap socket")
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
