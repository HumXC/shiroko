package screencap

import (
	"fmt"
	"os"
	"strings"

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

// screencap 工具特有的接口
type IScreencap interface {
	// 相当于 screencap -p -d displayID
	// displayID 可以为空
	Png(displayID string) ([]byte, error)
	Displays() ([]string, error)
}

var Screencap *ScreencapImpl = New()

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
	cmdDisplays := &cobra.Command{
		Use:   "displays",
		Short: "show displays ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			displays, err := s.Displays()
			if err != nil {
				return err
			}
			for _, display := range displays {
				fmt.Println(display)
			}
			return nil
		},
	}
	cmdPng := &cobra.Command{
		Use:   "png",
		Short: "Get screenshot and write to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			displayID, err := cmd.Flags().GetString("display-id")
			if err != nil {
				panic(err)
			}
			result, err := s.Png(displayID)
			if err != nil {
				return err
			}
			_, _ = os.Stdout.Write(result)
			return nil

		},
	}

	cmdPng.Flags().StringP("display-id", "d", "", "display id")

	c.AddCommand(cmdDisplays)
	c.AddCommand(cmdPng)
}

func (s *ScreencapImpl) Png(displayID string) ([]byte, error) {
	args := []string{"-p"}
	if displayID != "" {
		args = append(args, "-d", displayID)
	}
	cmd := android.Command(s.base.Exe(), args...)
	log.Info("Get screenshot and write to stdout", "displayID", displayID)
	log.Debug("Run command", "command", cmd.FullCmd())
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%w: with display id [%s]", err, displayID)
	}
	return b, nil
}

// 返回所需的 Display ID
// 通过运行 dumpsys SurfaceFlinger --display-id 获取
func (s *ScreencapImpl) Displays() ([]string, error) {
	result := make([]string, 0, 1)
	cmd := android.Command("dumpsys", "SurfaceFlinger", "--display-id")
	log.Info("Get displays")
	log.Debug("Run command", "command", cmd.FullCmd())
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command error [%s]: %w", cmd.FullCmd(), err)
	}
	str := string(output)

	for _, line := range strings.Split(str, "\n") {
		id := strings.Split(line, " ")[1]
		result = append(result, id)
	}
	return result, nil
}

func New() *ScreencapImpl {
	s := &ScreencapImpl{
		base: &screenBase{cmd: "screencap"},
	}
	return s
}
