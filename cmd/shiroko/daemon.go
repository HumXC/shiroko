package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
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
