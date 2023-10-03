package manager

import (
	"fmt"
	"strings"

	"github.com/HumXC/shiroko/logs"
	"github.com/HumXC/shiroko/tools/common"
	"github.com/spf13/cobra"
)

// 管理所有的 tool
type IManager interface {
	// 注册一个 tool
	Register(common.Tool)
	// 返回所有工具的名字
	List() []string
	// 对应 tools.common.Base
	// 约定 name 必须在 List 的结果之中，否则返回默认值
	Health(name string) error
	Install(name string) error
	Uninstall(name string) error
	Env(name string) []string
	Exe(name string) string
	Args(name string) []string
	Files(name string) []string
}

var log = logs.Get()

var _ IManager = &managerImpl{}
var Manager IManager = nil

type managerImpl struct {
	rootCmd  *cobra.Command
	allTools map[string]common.BaseTool
}

func (m *managerImpl) getTool(name string) common.BaseTool {
	t := m.allTools[name]
	if t == nil {
		log.Warn("Tool not found", "name", name)
	}
	return t
}

// Register implements IManager.
func (m *managerImpl) Register(tool common.Tool) {
	base := tool.Base()
	name := base.Name()
	log.Info("Register tool", "name", name)
	base.Init()
	m.allTools[base.Name()] = base

	to, ok := tool.(common.UseCommand)
	if ok {
		log.Info("Register command", "command", base.Name())
		subCmd := &cobra.Command{
			Use:   base.Name(),
			Short: base.Description(),
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		m.rootCmd.AddCommand(subCmd)
		to.RegCommand(subCmd)
	}
}

// Args implements IManager.
func (m *managerImpl) Args(name string) []string {
	t := m.getTool(name)
	if t == nil {
		return []string{}
	}
	return t.Args()
}

// Env implements IManager.
func (m *managerImpl) Env(name string) []string {
	t := m.getTool(name)
	if t == nil {
		return []string{}
	}
	return t.Env()
}

// Exe implements IManager.
func (m *managerImpl) Exe(name string) string {
	t := m.getTool(name)
	if t == nil {
		return ""
	}
	return t.Exe()
}

// Files implements IManager.
func (m *managerImpl) Files(name string) []string {
	t := m.getTool(name)
	if t == nil {
		return []string{}
	}
	return t.Files()
}

// Health implements IManager.
func (m *managerImpl) Health(name string) error {
	t := m.getTool(name)
	if t == nil {
		return nil
	}
	return t.Health()
}

// Install implements IManager.
func (m *managerImpl) Install(name string) error {
	t := m.getTool(name)
	if t == nil {
		return nil
	}
	return t.Install()
}

// List implements IManager.
func (m *managerImpl) List() []string {
	result := make([]string, len(m.allTools))
	for _, t := range m.allTools {
		result = append(result, t.Name())
	}
	return result
}

// Uninstall implements IManager.
func (m *managerImpl) Uninstall(name string) error {
	t := m.getTool(name)
	if t == nil {
		return nil
	}
	return t.Uninstall()
}

func New(rootCmd *cobra.Command) IManager {
	m := &managerImpl{
		rootCmd:  rootCmd,
		allTools: make(map[string]common.BaseTool),
	}
	m.setCommand()
	return m
}

func (m *managerImpl) setCommand() {
	type flagsSet struct {
		health, env, files, args, exe bool
	}
	setFlags := func(cmd *cobra.Command) {
		flags := cmd.Flags()
		flags.BoolP("health", "H", false, "health check")
		flags.Bool("env", false, "show used env")
		flags.BoolP("files", "F", false, "show used files")
		flags.BoolP("args", "A", false, "show default args")
		flags.BoolP("exe", "E", false, "show executable file")
	}
	parseFlags := func(cmd *cobra.Command) flagsSet {
		health, err := cmd.Flags().GetBool("health")
		if err != nil {
			panic(err)
		}
		exe, err := cmd.Flags().GetBool("exe")
		if err != nil {
			panic(err)
		}
		env, err := cmd.Flags().GetBool("env")
		if err != nil {
			panic(err)
		}
		_args, err := cmd.Flags().GetBool("args")
		if err != nil {
			panic(err)
		}
		files, err := cmd.Flags().GetBool("files")
		if err != nil {
			panic(err)
		}

		return flagsSet{health, env, files, _args, exe}
	}
	cmdTools := &cobra.Command{
		Use:   "tools",
		Short: "Manager tools",
		Run: func(cmd *cobra.Command, args []string) {
			flags := parseFlags(cmd)
			genStr := func(arr [][2]string) string {
				result := ""
				for i, v := range arr {
					if i == 0 {
						result += fmt.Sprintf("  %2d. %-10s %s", i, v[0], v[1])
						continue
					}
					result += fmt.Sprintf("\n  %2d. %-10s %s", i, v[0], v[1])
				}
				return result
			}
			if flags.health {
				result := make([][2]string, 0)
				for _, tool := range m.allTools {
					err := tool.Health()
					msg := "OK"
					if err != nil {
						msg = err.Error()
					}
					result = append(result, [2]string{tool.Name(), msg})
				}
				fmt.Println("Health Check:")
				fmt.Println(genStr(result))
			}

			if flags.exe {
				result := make([][2]string, 0)
				for _, tool := range m.allTools {
					exe := tool.Exe()
					result = append(result, [2]string{tool.Name(), exe})
				}
				fmt.Println("Executable: ")
				fmt.Println(genStr(result))
			}

			if flags.env {
				result := make([][2]string, 0)
				for _, tool := range m.allTools {
					env := tool.Env()
					msg := ""
					for _, v := range env {
						msg += "\n        " + v
					}
					result = append(result, [2]string{tool.Name(), msg})
				}
				fmt.Println("Environment: ")
				fmt.Println(genStr(result))
			}
			if flags.args {
				result := make([][2]string, 0)
				for _, tool := range m.allTools {
					_args := tool.Args()
					msg := strings.Join(_args, " ")
					result = append(result, [2]string{tool.Name(), msg})
				}
				fmt.Println("Arguments: ")
				fmt.Println(genStr(result))
			}
			if flags.files {
				result := make([][2]string, 0)
				for _, tool := range m.allTools {
					files := tool.Files()
					msg := ""
					for _, v := range files {
						msg += "\n        " + v
					}
					result = append(result, [2]string{tool.Name(), msg})
				}
				fmt.Println("Files: ")
				fmt.Println(genStr(result))
			}

			if !(flags.health ||
				flags.exe ||
				flags.env ||
				flags.args ||
				flags.files) {
				cmd.Help()
			}
		},
	}
	for _, tool := range m.allTools {
		t := tool
		c := &cobra.Command{
			Use:   t.Name(),
			Short: t.Description(),
			Run: func(cmd *cobra.Command, args []string) {
				flags := parseFlags(cmd)
				install, err := cmd.Flags().GetBool("install")
				if err != nil {
					panic(err)
				}
				uninstall, err := cmd.Flags().GetBool("uninstall")
				if err != nil {
					panic(err)
				}
				if flags.health {
					fmt.Println("Health Check: ")
					err := t.Health()
					if err != nil {
						fmt.Println("  " + err.Error())
					} else {
						fmt.Println("  OK")
					}
				}

				if flags.exe {
					fmt.Println("Executable: ")
					fmt.Println("  " + t.Exe())
				}

				if flags.env {
					fmt.Println("Environment: ")
					for _, env := range t.Env() {
						fmt.Println("  " + env)
					}
				}
				if flags.args {
					fmt.Println("Arguments: ")
					for _, arg := range t.Args() {
						fmt.Println("  " + arg)
					}
				}
				if flags.files {
					fmt.Println("Files: ")
					for _, file := range t.Files() {
						fmt.Println("  " + file)
					}
				}
				if install && !uninstall {
					err := t.Install()
					if err != nil {
						fmt.Println("Install: Failed - " + err.Error())
					} else {
						fmt.Println("Install: Successful")
					}
				}
				if uninstall && !install {
					err := t.Uninstall()
					if err != nil {
						fmt.Println("Uninstall: Failed - " + err.Error())
					} else {
						fmt.Println("Uninstall: Successful")
					}
				}
				if !(flags.health ||
					flags.exe ||
					flags.env ||
					flags.args ||
					flags.files ||
					install ||
					uninstall) {
					cmd.Help()
				}
			},
		}
		setFlags(c)
		c.Flags().BoolP("install", "I", false, "install")
		c.Flags().BoolP("uninstall", "U", false, "uninstall")
		cmdTools.AddCommand(c)
	}
	setFlags(cmdTools)
	m.rootCmd.AddCommand(cmdTools)
}
