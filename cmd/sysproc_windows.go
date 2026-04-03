//go:build windows

package cmd

import "syscall"

func hiddenProcess() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{HideWindow: true}
}

func psCommand() string {
	return "powershell"
}
