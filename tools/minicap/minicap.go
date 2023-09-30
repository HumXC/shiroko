package minicap

import "github.com/HumXC/shiroko/tools"

type Minicap struct {
	Base tools.Tool
}

func NewMinicap() *Minicap {
	m := &Minicap{
		Base: &minicap{},
	}
	m.Base.Init()
	return m
}
