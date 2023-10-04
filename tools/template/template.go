package template

import (
	"github.com/HumXC/shiroko/tools/common"
	"github.com/spf13/cobra"
)

var XXX *XXXImpl = New()

type IXXX interface{}
type XXXImpl struct {
	base common.BaseTool
}

var _ IXXX = &XXXImpl{}
var _ common.Tool = &XXXImpl{}
var _ common.UseCommand = &XXXImpl{} // 可选

func New() *XXXImpl {
	return &XXXImpl{
		// make some data, like base, slice and map.
	}
}

// Base implements common.Tool.
func (x *XXXImpl) Base() common.BaseTool {
	return x.base
}

// RegCommand implements common.UseCommand.
func (*XXXImpl) RegCommand(*cobra.Command) {
	panic("unimplemented")
}
