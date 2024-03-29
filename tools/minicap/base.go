package minicap

// 参考：https://github.com/DeviceFarmer/minicap/blob/0276fbeff6803c7ff4e39450f7c87a2ba59be25e/run.sh
import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/blob"
	"github.com/HumXC/shiroko/tools/common"
)

const PROP_ABI = "ro.product.cpu.abi"
const PROP_SDK = "ro.build.version.sdk"
const PROP_PRE = "ro.build.version.preview_sdk"
const PROP_REL = "ro.build.version.release"

type minicapBase struct {
	env      []string
	files    []string
	args     []string
	exe      string
	embedBin string // 在 embed 里的 bin 文件路径
	embedLib string // 在 embed 里的 lib 文件路径
	bin      string // 设备里的 bin 文件路径
	lib      string // 设备里的 lib 文件路径
}

// Description implements common.Tool.
func (*minicapBase) Description() string {
	return "https://github.com/DeviceFarmer/minicap"
}

var _ common.BaseTool = &minicapBase{}

// Name implements tools.Tool.
func (*minicapBase) Name() string {
	return "minicap"
}

// Uninstall implements tools.Tool.
func (m *minicapBase) Uninstall() error {
	for _, file := range m.files {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}

// Files implements tools.Tool.
func (m *minicapBase) Files() []string {
	return m.files
}

// Init implements tools.Tool.
func (m *minicapBase) Init() {
	m.env = make([]string, 0, 2)
	m.files = make([]string, 0, 2)
	m.args = make([]string, 0, 2)
	abi, sdk, rel := m.Getprop()
	log.Info("Device info", "ABI", abi, "SDK", sdk, "Release", rel)
	embedBin := m.getBin(abi, sdk)
	embedLib := m.getLib(abi, sdk, rel)
	// 可能会存在设备安卓版本太新而找不到对应 lib 的情况
	if !strings.HasPrefix(m.embedBin, "noarch") && embedLib == "" {
		embedBin = "noarch/minicap.apk"
	}
	m.embedBin = embedBin
	m.embedLib = embedLib
	if embedLib != "" {
		m.env = append(m.env, "LD_LIBRARY_PATH="+android.TMP_DIR)
		lib := path.Join(android.TMP_DIR, path.Base(embedLib))
		m.files = append(m.files, lib)
		m.lib = lib
	}

	bin := path.Join(android.TMP_DIR, path.Base(m.embedBin))
	m.bin = bin
	m.exe = bin
	// 如果 embedBin 以 noarch 开头
	if strings.HasPrefix(m.embedBin, "noarch") {
		m.exe = "app_process"
		m.env = append(m.env, "CLASSPATH="+bin)
		m.args = append(m.args, "/system/bin", "io.devicefarmer.minicap.Main")
	}
	m.files = append(m.files, bin)
	log.Info("Use minicap", "bin", m.embedBin, "lib", m.embedLib)
}

func (minicapBase) Getprop() (abi, sdk, rel string) {
	abi = android.Getprop(PROP_ABI)
	sdk = android.Getprop(PROP_SDK)
	_pre := android.Getprop(PROP_PRE)
	_rel := android.Getprop(PROP_REL)

	iPre, _ := strconv.Atoi(_pre)
	iRel, _ := strconv.Atoi(_rel)
	rel = strconv.Itoa(iRel + iPre)
	return
}

// Args implements tools.Tool.
func (m *minicapBase) Args() []string {
	return m.args
}

// Install implements tools.Tool.
func (m *minicapBase) Install() error {
	log.Info("Install minicap")
	if m.embedBin != "" {
		log.Info("Copy file", "src", m.embedBin, "dst", m.bin)
		b, err := blob.Minicap.ReadFile(m.embedBin)
		if err != nil {
			return err
		}
		err = os.WriteFile(m.bin, b, 0755)
		if err != nil {
			return err
		}
	}
	if m.embedLib != "" {
		log.Info("Copy file", "src", m.embedLib, "dst", m.lib)
		b, err := blob.Minicap.ReadFile(m.embedLib)
		if err != nil {
			return err
		}
		err = os.WriteFile(m.lib, b, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// Env implements tools.Tool.
func (m *minicapBase) Env() []string {
	return m.env
}

// Exe implements tools.Tool.
func (m *minicapBase) Exe() string {
	return m.exe
}

// Health implements tools.Tool.
func (m *minicapBase) Health() error {
	return common.HealthWithFiles(m.Files())
}

func (minicapBase) getBin(abi, sdk string) string {
	name := "minicap"
	if _sdk, err := strconv.Atoi(sdk); err == nil && _sdk >= 16 {
		name = "minicap-nopie"
	}
	dir := path.Join(abi, "bin")
	bin := path.Join(dir, name)
	if !blob.Minicap.IsExist(bin) {
		return "noarch/minicap.apk"
	}
	return bin
}

func (minicapBase) getLib(abi, sdk, rel string) string {
	libDir := path.Join(abi, "lib")
	lib := path.Join(libDir, "android-"+rel)
	if _, err := blob.Minicap.ReadDir(lib); err == nil {
		return path.Join(lib, "minicap.so")
	}
	lib = path.Join(libDir, "android-"+sdk)
	if _, err := blob.Minicap.ReadDir(lib); err == nil {
		return path.Join(lib, "minicap.so")
	}
	return ""
}
