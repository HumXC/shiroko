package minicap

import (
	"testing"

	"github.com/HumXC/shiroko/android"
)

var mc minicapBase = minicapBase{
	env:   make([]string, 0, 2),
	files: make([]string, 0, 2),
	args:  make([]string, 0, 2),
}

func TestMain(m *testing.M) {
	android.TestProp = make(map[string]string)
	android.TestProp["ro.product.cpu.abi"] = "arm64-v8a"
	android.TestProp["ro.build.version.sdk"] = "30"
	android.TestProp["ro.build.version.release"] = "11"
	android.TestProp["ro.build.version.preview_sdk"] = "0"
	mc.Init()
}
func TestGetBin(t *testing.T) {
	abi, sdk, _ := mc.Getprop()
	target := "arm64-v8a/bin/minicap"
	bin := mc.getBin(abi, sdk)
	if bin != target {
		t.Errorf("want %s, got %s", target, bin)
	}
}

func TestGetLib(t *testing.T) {
	abi, sdk, rel := mc.Getprop()
	target := "arm64-v8a/lib/android-30/minicap.so"
	lib := mc.getLib(abi, sdk, rel)
	if lib != target {
		t.Errorf("want %s, got %s", target, lib)
	}
}
