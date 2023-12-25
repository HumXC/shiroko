package screencap

import (
	"os"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/logs"
	"github.com/HumXC/shiroko/tools/common"
	"github.com/spf13/cobra"
)

var log *logs.Logger = nil
var Screencap *ScreencapImpl = nil

func Init() {
	log = logs.Get("screencap")
	Screencap = New()
}

// screencap 工具特有的接口
type IScreencap interface {
	// 相当于 screencap -p
	Png() ([]byte, error)
}

type ScreencapImpl struct {
	base common.BaseTool
}

var _ IScreencap = &ScreencapImpl{}
var _ common.Tool = &ScreencapImpl{}
var _ common.UseCommand = &ScreencapImpl{}

// Base implements common.Tool.
func (s *ScreencapImpl) Base() common.BaseTool {
	return s.base
}
func (s *ScreencapImpl) RegCommand(c *cobra.Command) {
	cmdPng := &cobra.Command{
		Use:   "png",
		Short: "Get screenshot and write to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := s.Png()
			if err != nil {
				return err
			}
			_, _ = os.Stdout.Write(result)
			return nil
		},
	}

	cmdPng.Flags().StringP("display-id", "d", "", "display id")

	c.AddCommand(cmdPng)
}

func (s *ScreencapImpl) Png() ([]byte, error) {
	cmd := android.Command(s.base.Exe(), "-p")
	log.Info("Get screenshot and write to stdout")
	log.Debug("Run command", "command", cmd.FullCmd())
	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func New() *ScreencapImpl {
	s := &ScreencapImpl{
		base: &screencapBase{Cmd: "screencap"},
	}
	return s
}
