package screencap

import (
	"github.com/HumXC/shiroko/tools/common"
)

var _ common.Tool = &screenBase{}

type screenBase struct {
	cmd string
}

// Description implements common.Tool.
func (*screenBase) Description() string {
	return "Android command screencap"
}

// Args implements tools.Tool.
func (*screenBase) Args() []string {
	return []string{}
}

// Env implements tools.Tool.
func (*screenBase) Env() []string {
	return []string{}
}

// Exe implements tools.Tool.
func (s *screenBase) Exe() string {
	return s.cmd
}

// Files implements tools.Tool.
func (*screenBase) Files() []string {
	return []string{}
}

// Health implements tools.Tool.
func (s *screenBase) Health() error {
	return common.CommandHealth(s.cmd)
}

// Init implements tools.Tool.
func (*screenBase) Init() {}

// Install implements tools.Tool.
func (*screenBase) Install() error {
	return nil
}

// Name implements tools.Tool.
func (s *screenBase) Name() string {
	return s.cmd
}

// Uninstall implements tools.Tool.
func (*screenBase) Uninstall() error {
	return nil
}
