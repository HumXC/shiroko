package screencap

import (
	"fmt"
	"strings"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/tools"
)

// screencap 工具特有的接口
type IScreencap interface {
	Png() ([]byte, error)
	PngWithDisplay(displayID string) ([]byte, error)
	Displays() ([]string, error)
}
type Screencap struct {
	Base tools.Tool
	IScreencap
}

func (s *Screencap) Png() ([]byte, error) {
	cmd := android.Command(s.Base.Exe(), "-p")
	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return b, nil
}
func (s *Screencap) PngWithDisplay(displayID string) ([]byte, error) {
	cmd := android.Command(s.Base.Exe(), "-p", "-d", displayID)
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%w: with display id [%s]", err, displayID)
	}
	return b, nil
}

// 返回所需的 Display ID
// 通过运行 dumpsys SurfaceFlinger --display-id 获取
func (s *Screencap) Displays() ([]string, error) {
	result := make([]string, 0, 1)
	cmd := android.Command("dumpsys", "SurfaceFlinger", "--display-id")
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
func New() *Screencap {
	s := &Screencap{
		Base: tools.NewCommandTool("screencap"),
	}
	s.Base.Init()
	return s
}
