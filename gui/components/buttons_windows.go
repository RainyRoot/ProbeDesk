//go:build windows

package components

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

func isAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	return err == nil && member
}

func relaunchAsAdmin() {
	exe, _ := os.Executable()
	verbPtr, _ := windows.UTF16PtrFromString("runas")
	exePtr, _ := windows.UTF16PtrFromString(exe)
	cwd, _ := os.Getwd()
	cwdPtr, _ := windows.UTF16PtrFromString(cwd)

	err := windows.ShellExecute(0, verbPtr, exePtr, nil, cwdPtr, windows.SW_HIDE)
	if err != nil {
		fmt.Println("Failed to relaunch as admin:", err)
	}
	os.Exit(0)
}
