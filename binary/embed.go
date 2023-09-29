package binary

import (
	"embed"
)

//go:embed stf-binaries/node_modules/minicap-prebuilt/prebuilt
var Minicap embed.FS

//go:embed stf-binaries/node_modules/minitouch-prebuilt/prebuilt
var Minitouch embed.FS
