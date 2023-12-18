package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/sevlyar/go-daemon"
	"github.com/shirou/gopsutil/process"
)

func Daemon() (pid int, err error) {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	workdir := path.Dir(exe)
	cntxt := &daemon.Context{
		PidFileName: path.Join(workdir, "shiroko.pid"),
		PidFilePerm: 0644,
		WorkDir:     workdir,
		Umask:       027,
		Args:        os.Args,
	}
	defer cntxt.Release()
	proc, err := cntxt.Reborn()
	if err != nil {
		err = fmt.Errorf("unable to run: %s", err)
		return
	}
	if proc == nil {
		return
	}
	pid = proc.Pid
	return
}
func List() {
	target, err := os.Executable()
	if err != nil {
		panic(err)
	}
	for _, proc := range findProcess(target) {
		cmdl, _ := proc.Cmdline()
		fmt.Println(proc.Pid, cmdl)
	}
}
func findProcess(target string) []*process.Process {
	result := make([]*process.Process, 0, 1)
	processes, err := process.Processes()
	if err != nil {
		panic(err)
	}

	// 遍历所有进程
	for _, proc := range processes {
		exe, err := proc.Exe()
		if err != nil {
			continue
		}
		if strings.HasPrefix(exe, target) && proc.Pid != int32(os.Getpid()) {
			result = append(result, proc)
		}
	}
	return result
}
func Kill() {
	target, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// 遍历所有进程
	for _, proc := range findProcess(target) {
		exe, _ := proc.Exe()
		fmt.Printf("Killing process %s pid is %d\n", exe, proc.Pid)
		_ = proc.Kill()
	}
}
