//go:build !windows

package cmd

import "syscall"

func hiddenProcess() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func psCommand() string {
	return "pwsh"
}
