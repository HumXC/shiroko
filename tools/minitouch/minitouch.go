package minitouch

import (
	"io"

	"github.com/HumXC/shiroko/tools/common"
	"github.com/spf13/cobra"
)

var Minitouch *MinitouchImpl = New()

// 高版本要 root，先不实现
type IMinitouch interface {
}
type MinitouchImpl struct {
	base common.BaseTool
}

// RegCommand implements common.UseCommand.
func (*MinitouchImpl) RegCommand(*cobra.Command) {
	return
}

// Base implements common.Tool.
func (m *MinitouchImpl) Base() common.BaseTool {
	return m.base
}

// Open implements IMinitouch.
func (m *MinitouchImpl) Open() (io.ReadCloser, error) {
	// android.Command(m.base.Exe())
	return nil, nil
}

var _ IMinitouch = &MinitouchImpl{}

var _ common.Tool = &MinitouchImpl{}
var _ common.UseCommand = &MinitouchImpl{}

func New() *MinitouchImpl {
	return &MinitouchImpl{
		base: &minitouchBase{},
	}
}
