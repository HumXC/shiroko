package minicap

import (
	"encoding/json"
	"fmt"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/tools"
)

type Info struct {
	Id       int8    `json:"id"`
	Width    int32   `json:"width"`
	Height   int32   `json:"height"`
	Xdpi     float32 `json:"xdpi"`
	Ydpi     float32 `json:"ydpi"`
	Size     float32 `json:"size"`
	Density  float32 `json:"density"`
	Fps      float32 `json:"fps"`
	Secure   bool    `json:"secure"`
	Rotation int16   `json:"rotation"`
}
type Minicap struct {
	Base tools.Tool
}

func (m *Minicap) Info() (Info, error) {
	result := Info{}
	cmd := android.Command(m.Base.Exe(), append(m.Base.Args(), "-i")...)
	cmd.SetEnv(m.Base.Env())
	output, err := cmd.Output()
	if err != nil {
		return result, fmt.Errorf("failed to get info: %w", err)
	}
	err = json.Unmarshal(output, &result)
	if err != nil {
		return result, fmt.Errorf("failed to get info: %w", err)
	}
	return result, nil
}

// TODO
func (m *Minicap) Run(RWidth, RHeight, VWidth, VHeight, Orientation int32) {
	// args := append(m.Base.Args(), "-P", fmt.Sprintf("%dx%d@%dx%d/%d", RWidth, RHeight, VWidth, VHeight, Orientation))
	// cmd := exec.Command(m.Base.Exe(), args...)
	// cmd.Start()
	// ctx, _ := context.WithTimeout(context.Background(), 1000*time.Millisecond)

}
func NewMinicap() *Minicap {
	m := &Minicap{
		Base: &minicap{},
	}
	m.Base.Init()
	return m
}
