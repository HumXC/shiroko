package screencap

import (
	"fmt"
	"strings"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/tools"
)

type Screencap struct {
	Base tools.Tool
}

func (s *Screencap) Png() ([]byte, error) {
	cmd := android.Command(s.Base.Exe(), "-p")
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to screencap: %w", err)
	}
	return b, nil
}
func (s *Screencap) PngWithDisplay(displayID string) ([]byte, error) {
	cmd := android.Command(s.Base.Exe(), "-p", "-d", displayID)
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to screencap with display id [%s]: %w", displayID, err)
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
		return nil, fmt.Errorf("can not get display-id: %w", err)
	}
	str := string(output)

	for _, line := range strings.Split(str, "\n") {
		id := strings.Split(line, " ")[1]
		result = append(result, id)
	}
	return result, nil
}
func NewScreencap() *Screencap {
	s := &Screencap{
		Base: tools.NewCommandTool("screencap"),
	}
	s.Base.Init()
	return s
}
