package screencap

import (
	"github.com/HumXC/shiroko/tools/common"
)

var _ common.BaseTool = &screencapBase{}

type screencapBase = common.AndroidCommandBase
