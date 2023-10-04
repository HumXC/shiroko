package template

import "github.com/HumXC/shiroko/tools/common"

func init() {
	panic("don`t import this package: tools.template, it`s just a template")
}

var _ common.BaseTool = &XXXBase{}

type XXXBase struct {
}

// Args implements common.BaseTool.
func (*XXXBase) Args() []string {
	panic("unimplemented")
}

// Description implements common.BaseTool.
func (*XXXBase) Description() string {
	panic("unimplemented")
}

// Env implements common.BaseTool.
func (*XXXBase) Env() []string {
	panic("unimplemented")
}

// Exe implements common.BaseTool.
func (*XXXBase) Exe() string {
	panic("unimplemented")
}

// Files implements common.BaseTool.
func (*XXXBase) Files() []string {
	panic("unimplemented")
}

// Health implements common.BaseTool.
func (*XXXBase) Health() error {
	panic("unimplemented")
}

// Init implements common.BaseTool.
func (*XXXBase) Init() {
	panic("unimplemented")
}

// Install implements common.BaseTool.
func (*XXXBase) Install() error {
	panic("unimplemented")
}

// Name implements common.BaseTool.
func (*XXXBase) Name() string {
	panic("unimplemented")
}

// Uninstall implements common.BaseTool.
func (*XXXBase) Uninstall() error {
	panic("unimplemented")
}
