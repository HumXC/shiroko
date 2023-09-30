package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/process"
)

type Daemonize struct {
	IAmDaemon bool
	Pid       int
}

func Daemon() (Daemonize, error) {
	// 检查程序是否已经是守护进程
	if os.Getenv("DAEMON") == "true" {
		return Daemonize{
			IAmDaemon: true,
		}, nil
	}
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Env = append(os.Environ(), "DAEMON=true")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // 创建新的会话和新的进程组
	}

	if err := cmd.Start(); err != nil {
		return Daemonize{}, fmt.Errorf("failed to daemonize: %w", err)
	}

	return Daemonize{
		IAmDaemon: false,
		Pid:       cmd.Process.Pid,
	}, nil
}

func Kill() {
	selfExe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// 获取所有进程
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
		// 以 Daemon 运行的进程名可能会是 "/data/local/tmp/shiroko (deleted)"
		exe = strings.TrimSuffix(exe, " (deleted)")
		if exe == selfExe {
			pid := proc.Pid
			// 检查 PID 是否是当前进程
			if pid != int32(os.Getpid()) {
				fmt.Printf("Killing process %d with name %s\n", pid, exe)
				_ = proc.Kill()
			}
		}
	}
}
