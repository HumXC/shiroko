package tools

import (
	"github.com/HumXC/shiroko/logs"
	"github.com/HumXC/shiroko/tools/manager"
	minicap "github.com/HumXC/shiroko/tools/minicap"
	"github.com/HumXC/shiroko/tools/screencap"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var log *slog.Logger

func init() {
	log = logs.Get()
}

// 由 main 包调用
func Init(cmd *cobra.Command) {
	log.Info("Init tools")
	manager.Manager = manager.New(cmd)
	manager.Manager.Register(minicap.Minicap)
	manager.Manager.Register(screencap.Screencap)
	manager.Manager.SetCommand()
}
