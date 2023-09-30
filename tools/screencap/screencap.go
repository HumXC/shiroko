package screencap

import "github.com/HumXC/shiroko/tools"

type Screencap struct {
	Base tools.Tool
}

func NewScreencap() *Screencap {
	s := &Screencap{
		Base: tools.NewCommandTool("screencap"),
	}
	s.Base.Init()
	return s
}
