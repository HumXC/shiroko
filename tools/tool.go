package tools

import (
	"github.com/HumXC/shiroko/tools/input"
	"github.com/HumXC/shiroko/tools/manager"
	"github.com/HumXC/shiroko/tools/minicap"
	"github.com/HumXC/shiroko/tools/screencap"
	"github.com/HumXC/shiroko/tools/shell"
	"github.com/HumXC/shiroko/tools/window"
	"github.com/spf13/cobra"
)

// 由 main 包调用
func Init(cmd *cobra.Command) {
	manager.Init(cmd)
	input.Init()
	shell.Init()
	screencap.Init()
	minicap.Init()
	window.Init()

	manager.Manager.Register(input.Input)
	manager.Manager.Register(shell.Shell)
	manager.Manager.Register(screencap.Screencap)
	manager.Manager.Register(minicap.Minicap)
	manager.Manager.Register(window.Window)
	manager.Manager.SetCommand()
}
