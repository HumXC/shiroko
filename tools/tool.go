package tools

import (
	"reflect"

	"github.com/HumXC/shiroko/tools/common"
	minicap "github.com/HumXC/shiroko/tools/minicap"
	"github.com/HumXC/shiroko/tools/screencap"
	"github.com/spf13/cobra"
)

var allTools map[string]common.Tool = make(map[string]common.Tool)

func Init(cmd *cobra.Command) {
	register(cmd, minicap.Minicap)
	register(cmd, screencap.Screencap)
}

// 由 tool 的 init 函数内调用
func register(cmd *cobra.Command, t any) {
	val := reflect.ValueOf(t).Elem()
	typ := val.Type()

	// 检查是否存在 "Base" 字段
	baseField := val.FieldByName("Base")
	if !baseField.IsValid() {
		panic("Base field not found in: " + typ.String())
	}

	base, ok := baseField.Interface().(common.Tool)
	if !ok {
		panic("failed Base is not implements tools.Tool in: " + typ.String())
	}
	base.Init()
	allTools[base.Name()] = base

	to, ok := t.(common.UseCommand)
	if ok {
		subCmd := &cobra.Command{
			Use:   base.Name(),
			Short: base.Description(),
		}
		to.RegCommand(subCmd)
		cmd.AddCommand(subCmd)
	}
}

func Names() []string {
	names := make([]string, 0, len(allTools))
	for name := range allTools {
		names = append(names, name)
	}
	return names
}
