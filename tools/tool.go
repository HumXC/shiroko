package tools

type Tool interface {
	// 检查工具是否可用
	// error 为 nil 时表示可用
	Health() error
	// 将工具部署到设备
	Install() error
	Uninstall() error
	// 返回运行时需要的环境变量，例如: LD_LIBRARY_PATH=/data/local/tmp/
	Env() []string
	// 返回可执行文件的路径
	Exe() string
	// 返回需要的默认参数
	Args() []string
	// 返回所有部署在设备上的文件
	Files() []string
	// 在被调用时初始化，由 tools 包调用
	Init()
}
