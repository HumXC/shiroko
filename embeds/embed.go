package embeds

import (
	"embed"
	"io/fs"
	"path"
)

//go:embed stf-binaries/node_modules/@devicefarmer/minicap-prebuilt/prebuilt
var minicap_embed embed.FS
var Minicap = achive{
	fs:      minicap_embed,
	baseDir: "stf-binaries/node_modules/@devicefarmer/minicap-prebuilt/prebuilt",
}

//go:embed stf-binaries/node_modules/@devicefarmer/minitouch-prebuilt/prebuilt
var minitouch_embed embed.FS
var Minitouch = achive{
	fs:      minitouch_embed,
	baseDir: "stf-binaries/node_modules/@devicefarmer/minitouch-prebuilt/prebuilt",
}

type achive struct {
	baseDir string
	fs      embed.FS
}

func (a *achive) Open(name string) (fs.File, error) {
	return a.fs.Open(path.Join(a.baseDir, name))
}
func (a *achive) ReadDir(name string) ([]fs.DirEntry, error) {
	return a.fs.ReadDir(path.Join(a.baseDir, name))
}
func (a *achive) ReadFile(name string) ([]byte, error) {
	return a.fs.ReadFile(path.Join(a.baseDir, name))
}
func (a *achive) IsExist(name string) bool {
	f, err := a.fs.Open(path.Join(a.baseDir, name))
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}
