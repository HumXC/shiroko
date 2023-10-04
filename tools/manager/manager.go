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
	// 返回所有工具的名字
	List() []string
	// 对应 tools.common.Base
	Health(name string) error
	Install(name string) error
	Uninstall(name string) error
	Env(name string) ([]string, error)
	Exe(name string) (string, error)
	Args(name string) ([]string, error)
	Files(name string) ([]string, error)
}

var log = logs.Get()

var _ IManager = &ManagerImpl{}
var Manager *ManagerImpl = nil

type ManagerImpl struct {
	rootCmd  *cobra.Command
	allTools map[string]common.BaseTool
}

func (m *ManagerImpl) getTool(name string) (common.BaseTool, error) {
	t := m.allTools[name]
	if t == nil {
		return nil, fmt.Errorf("can not found tool with name: %s", name)
	}
	return t, nil
}

// Register implements IManager.
func (m *ManagerImpl) Register(tool common.Tool) {
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
func (m *ManagerImpl) Args(name string) ([]string, error) {
	t, err := m.getTool(name)
	if err != nil {
		return nil, err
	}
	return t.Args(), nil
}

// Env implements IManager.
func (m *ManagerImpl) Env(name string) ([]string, error) {
	t, err := m.getTool(name)
	if err != nil {
		return nil, err
	}
	return t.Env(), nil
}

// Exe implements IManager.
func (m *ManagerImpl) Exe(name string) (string, error) {
	t, err := m.getTool(name)
	if err != nil {
		return "", err
	}
	return t.Exe(), nil
}

// Files implements IManager.
func (m *ManagerImpl) Files(name string) ([]string, error) {
	t, err := m.getTool(name)
	if err != nil {
		return nil, err
	}
	return t.Files(), nil
}

// Health implements IManager.
func (m *ManagerImpl) Health(name string) error {
	t, err := m.getTool(name)
	if err != nil {
		return err
	}
	return t.Health()
}

// Install implements IManager.
func (m *ManagerImpl) Install(name string) error {
	t, err := m.getTool(name)
	if err != nil {
		return err
	}
	return t.Install()
}

// List implements IManager.
func (m *ManagerImpl) List() []string {
	result := make([]string, 0, len(m.allTools))
	for _, t := range m.allTools {
		result = append(result, t.Name())
	}
	return result
}

// Uninstall implements IManager.
func (m *ManagerImpl) Uninstall(name string) error {
	t, err := m.getTool(name)
	if err != nil {
		return err
	}
	return t.Uninstall()
}

func New(rootCmd *cobra.Command) *ManagerImpl {
	m := &ManagerImpl{
		rootCmd:  rootCmd,
		allTools: make(map[string]common.BaseTool),
	}
	return m
}

func (m *ManagerImpl) SetCommand() {
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
