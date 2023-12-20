package template

import (
	"github.com/HumXC/shiroko/tools/common"
	"github.com/spf13/cobra"
)

var Tool *ToolImpl = nil

func Init() {
	Tool = New()
}

type ITool interface{}
type ToolImpl struct {
	base common.BaseTool
}

var _ ITool = &ToolImpl{}
var _ common.Tool = &ToolImpl{}
var _ common.UseCommand = &ToolImpl{} // 可选

func New() *ToolImpl {
	return &ToolImpl{
		// make some data, like base, slice and map.
	}
}

// Base implements common.Tool.
func (x *ToolImpl) Base() common.BaseTool {
	return x.base
}

// RegCommand implements common.UseCommand.
func (*ToolImpl) RegCommand(*cobra.Command) {
	panic("unimplemented")
}
