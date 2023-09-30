package tools

import "os/exec"

type commandTool struct {
	exe string
}

// Args implements tools.Tool.
func (*commandTool) Args() []string {
	return []string{}
}

// Env implements tools.Tool.
func (*commandTool) Env() []string {
	return []string{}
}

// Exe implements tools.Tool.
func (s *commandTool) Exe() string {
	return s.exe
}

// Files implements tools.Tool.
func (*commandTool) Files() []string {
	return []string{}
}

func (s *commandTool) Health() error {
	_, err := exec.LookPath(s.exe)
	if err != nil {
		return err
	}
	return nil
}

// Init implements tools.Tool.
func (s *commandTool) Init() {}

// Install implements tools.Tool.
func (*commandTool) Install() error {
	return nil
}

// Uninstall implements tools.Tool.
func (*commandTool) Uninstall() error {
	return nil
}

func NewCommandTool(cmd string) Tool {
	return &commandTool{exe: cmd}
}
