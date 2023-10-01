package tools

import (
	"reflect"

	"github.com/HumXC/shiroko/tools/common"
	minicap "github.com/HumXC/shiroko/tools/minicap"
	"github.com/HumXC/shiroko/tools/screencap"
	"github.com/spf13/cobra"
)

// TODO: 日志
var allTools map[string]common.BaseTool = make(map[string]common.BaseTool)

// 由 main 包调用
func Init(cmd *cobra.Command) {
	register(cmd, minicap.Minicap)
	register(cmd, screencap.Screencap)
	setCommand(cmd)
}

// 添加工具到全局的 allTools
// 会验证 tool 是否符合要求：具有 Base 字段且实现了 common.BaseTool
// 如果实现了 common.UseCommand 接口, 则调用此接口的函数
func register(cmd *cobra.Command, tool any) {
	val := reflect.ValueOf(tool).Elem()
	typ := val.Type()

	// 检查是否存在 "Base" 字段
	baseField := val.FieldByName("Base")
	if !baseField.IsValid() {
		panic("Base field not found in: " + typ.String())
	}

	base, ok := baseField.Interface().(common.BaseTool)
	if !ok {
		panic("failed Base is not implements tools.Tool in: " + typ.String())
	}
	base.Init()
	allTools[base.Name()] = base

	to, ok := tool.(common.UseCommand)
	if ok {
		subCmd := &cobra.Command{
			Use:   base.Name(),
			Short: base.Description(),
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		cmd.AddCommand(subCmd)
		to.RegCommand(subCmd)
	}
}
