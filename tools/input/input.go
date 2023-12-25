package input

import (
	"strconv"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/tools/common"
)

var Input *InputImpl = nil

func Init() {
	Input = New()
}

// Default: -1 for key event, 0 for motion event if not specified.
type IInput interface {
	// text <string> (Default: touchscreen)
	Text(text string) error
	// keyevent [--longpress] <key code number or name> ... (Default: keyboard)
	Keyevent(longpress bool, keycode ...string) error
	// tap <x> <y> (Default: touchscreen)
	Tap(x, y int32) error
	// swipe <x1> <y1> <x2> <y2> [duration(ms)] (Default: touchscreen)
	Swipe(x1, y1, x2, y2 int32, duration int32) error
	// draganddrop <x1> <y1> <x2> <y2> [duration(ms)] (Default: touchscreen)
	Draganddrop(x1, y1, x2, y2 int32, duration int32) error
	// motionevent <DOWN|UP|MOVE> <x> <y> (Default: touchscreen)
	Motionevent(event string, x, y int32) error
}
type InputImpl struct {
	base common.BaseTool
}

// Draganddrop implements IInput.
func (i *InputImpl) Draganddrop(x1 int32, y1 int32, x2 int32, y2 int32, duration int32) error {
	return i.run([]string{
		"draganddrop",
		strconv.Itoa(int(x1)),
		strconv.Itoa(int(y1)),
		strconv.Itoa(int(x2)),
		strconv.Itoa(int(y2)),
		strconv.Itoa(int(duration)),
	})
}

// Keyevent implements IInput.
func (i *InputImpl) Keyevent(longpress bool, keycode ...string) error {
	args := []string{"keyevent"}
	if longpress {
		args = append(args, "--longpress")
	}
	args = append(args, keycode...)
	return i.run(args)
}

// Motionevent implements IInput.
func (i *InputImpl) Motionevent(event string, x int32, y int32) error {
	return i.run([]string{"motionevent", event, strconv.Itoa(int(x)), strconv.Itoa(int(y))})
}

// Swipe implements IInput.
func (i *InputImpl) Swipe(x1 int32, y1 int32, x2 int32, y2 int32, duration int32) error {
	return i.run([]string{
		"swipe",
		strconv.Itoa(int(x1)),
		strconv.Itoa(int(y1)),
		strconv.Itoa(int(x2)),
		strconv.Itoa(int(y2)),
		strconv.Itoa(int(duration)),
	})
}

// Tap implements IInput.
func (i *InputImpl) Tap(x int32, y int32) error {
	return i.run([]string{"tap", strconv.Itoa(int(x)), strconv.Itoa(int(y))})
}

// Text implements IInput.
func (i *InputImpl) Text(text string) error {
	return i.run([]string{"text", text})
}

func (i *InputImpl) run(args []string) error {
	cmd := android.Command("input", args...)
	_, err := cmd.Output()
	return err
}

var _ IInput = &InputImpl{}
var _ common.Tool = &InputImpl{}

func New() *InputImpl {
	return &InputImpl{
		base: &inputBase{Cmd: "input"},
	}
}

// Base implements common.Tool.
func (i *InputImpl) Base() common.BaseTool {
	return i.base
}
