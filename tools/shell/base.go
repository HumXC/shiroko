package shell

import "github.com/HumXC/shiroko/tools/common"

var _ common.BaseTool = &shellBase{}

type shellBase = common.AndroidCommandBase
