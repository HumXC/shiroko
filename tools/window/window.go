package window

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/tools/common"
)

// wm 的命令
var Window *WindowImpl = nil

func Init() {
	Window = New()
}

type IWindow interface {
	GetSize() (int32, int32, error)
	// wm size xxx
	SetSize(width, height int32) error
	// wm size reset
	ResetSize() error
	// wm set-user-rotation user-rotation
	SetRotation(lock bool, rotation int32) error
	GetDensity() (int32, error)
	SetDensity(density int32) error
	ResetDensity() error
}
type WindowImpl struct {
	base common.BaseTool
}

// GetDensity implements IWindow.
func (w *WindowImpl) GetDensity() (int32, error) {
	c := android.Command("wm", "density")
	out, err := c.Output()
	if err != nil {
		return 0, err
	}
	r := strings.Split(string(out), " ")
	dens, err := strconv.Atoi(r[len(r)-1])
	if err != nil {
		return 0, err
	}
	return int32(dens), nil
}

// GetSize implements IWindow.
func (w *WindowImpl) GetSize() (int32, int32, error) {
	c := android.Command("wm", "size")
	out, err := c.Output()
	if err != nil {
		return 0, 0, err
	}
	r := strings.Split(string(out), " ")
	size := strings.Split(r[len(r)-1], "x")
	sw, err := strconv.Atoi(size[0])
	if err != nil {
		return 0, 0, err
	}
	sh, err := strconv.Atoi(size[1])
	if err != nil {
		return 0, 0, err
	}
	return int32(sw), int32(sh), nil
}

// ResetDensity implements IWindow.
func (w *WindowImpl) ResetDensity() error {
	c := android.Command("wm", "density", "reset")
	_, err := c.Output()
	return err
}

// ResetSize implements IWindow.
func (w *WindowImpl) ResetSize() error {
	c := android.Command("wm", "size", "reset")
	_, err := c.Output()
	return err
}

// SetDensity implements IWindow.
func (w *WindowImpl) SetDensity(density int32) error {
	c := android.Command("wm", "density", strconv.Itoa(int(density)))
	_, err := c.Output()
	return err
}

// SetRotation implements IWindow.
func (w *WindowImpl) SetRotation(lock bool, rotation int32) error {
	// 根据版本不同可能会有两个不同的子命令 set-user-rotation user-rotation
	cmds := []string{"user-rotation", "set-user-rotation"}
	// 第一个元素留给 subcmd
	args := []string{"", "free", strconv.Itoa(int(rotation))}
	if lock {
		args[1] = "lock"
	}
	errs := []error{}
	for _, subcmd := range cmds {
		args[0] = subcmd
		c := android.Command("wm", args...)
		_, err := c.Output()
		if err == nil {
			return nil
		}
		errs = append(errs, err)
	}
	return fmt.Errorf("failed to set rotation: %v", errs)
}

// SetSize implements IWindow.
func (w *WindowImpl) SetSize(width int32, height int32) error {
	c := android.Command(
		"wm", "size",
		strconv.Itoa(int(width))+"x"+strconv.Itoa(int(height)),
	)
	_, err := c.Output()
	return err
}

var _ IWindow = &WindowImpl{}
var _ common.Tool = &WindowImpl{}

func New() *WindowImpl {
	return &WindowImpl{
		base: &windowBase{Cmd: "wm"},
	}
}

// Base implements common.Tool.
func (w *WindowImpl) Base() common.BaseTool {
	return w.base
}
