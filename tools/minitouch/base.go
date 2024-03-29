package minitouch

import (
	"errors"
	"os"
	"path"
	"strconv"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/blob"
	"github.com/HumXC/shiroko/tools/common"
)

var _ common.BaseTool = &minitouchBase{}

type minitouchBase struct {
	embedBin string // 在 embed 里的 bin 文件路径
	bin      string // 设备里的 bin 文件路径
}

// Args implements common.BaseTool.
func (*minitouchBase) Args() []string {
	return []string{}
}

// Description implements common.BaseTool.
func (*minitouchBase) Description() string {
	return "https://github.com/DeviceFarmer/minitouch"
}

// Env implements common.BaseTool.
func (*minitouchBase) Env() []string {
	return []string{}
}

// Exe implements common.BaseTool.
func (m *minitouchBase) Exe() string {
	return m.bin
}

// Files implements common.BaseTool.
func (m *minitouchBase) Files() []string {
	return []string{m.bin}
}

// Health implements common.BaseTool.
func (m *minitouchBase) Health() error {
	return common.HealthWithFiles([]string{m.bin})
}

// Init implements common.BaseTool.
func (m *minitouchBase) Init() {
	// https://github.com/DeviceFarmer/minitouch/blob/master/run.sh
	abi := android.Getprop("ro.product.cpu.abi")
	sdk := android.Getprop("ro.build.version.sdk")
	_sdk, _ := strconv.Atoi(sdk)
	bin := "minitouch-nopie"
	if _sdk < 16 {
		bin = "minitouch"
	}
	m.embedBin = path.Join(abi, "bin", bin)
	m.bin = path.Join(android.TMP_DIR, bin)
}

// Install implements common.BaseTool.
func (m *minitouchBase) Install() error {
	if !blob.Minitouch.IsExist(m.embedBin) {
		return errors.New("binary not found: " + m.embedBin)
	}
	b, err := blob.Minitouch.ReadFile(m.embedBin)
	if err != nil {
		return err
	}
	return os.WriteFile(m.bin, b, 0755)
}

// Name implements common.BaseTool.
func (*minitouchBase) Name() string {
	return "minitouch"
}

// Uninstall implements common.BaseTool.
func (m *minitouchBase) Uninstall() error {
	return os.Remove(m.bin)
}
