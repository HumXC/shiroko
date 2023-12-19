package tools

import (
	"github.com/HumXC/shiroko/tools/input"
	"github.com/HumXC/shiroko/tools/manager"
	"github.com/HumXC/shiroko/tools/minicap"
	"github.com/HumXC/shiroko/tools/minitouch"
	"github.com/HumXC/shiroko/tools/screencap"
	"github.com/spf13/cobra"
)

// 由 main 包调用
func Init(cmd *cobra.Command) {
	manager.Init(cmd)
	input.Init()
	minicap.Init()
	minicap.Init()
	minitouch.Init()
	screencap.Init()

	manager.Manager.Register(minicap.Minicap)
	manager.Manager.Register(screencap.Screencap)
	manager.Manager.SetCommand()
}
