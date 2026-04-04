//go:build !windows

package cmd

import "syscall"

func hiddenSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
