package components

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	fynedialog "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/RainyRoot/ProbeDesk/cmd"
	nativeDialog "github.com/sqweek/dialog"
	"golang.org/x/sys/windows"
)

// Check if admin privileges
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

// // Restart and prompt UAC
// func relaunchAsAdmin() {
// 	exe, _ := os.Executable()
// 	cmd := exec.Command("powershell", "-Command", "Start-Process", exe, "-Verb", "runAs", "-WindowStyle", "Hidden")
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	cmd.Run()
// 	os.Exit(0)
// }

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

func NewRunButton(
	output *widget.Entry, status *widget.Label,
	systemCheck, ipCheck, usbCheck, servicesCheck, usersCheck, traceRouteRequest,
	vpnCheck, productsCheck, netuseCheck, checkHealthCheck,
	flushDns, wingetUpdate, scanHealth, restoreHealth *widget.Check,
	remoteEntry *widget.Entry, traceRoute *widget.Entry,
	formatSelect *widget.Select, mu *sync.Mutex,
) *widget.Button {

	return widget.NewButton("Run Selected", func() {
		// Check if actions requires admin privileges
		adminChecks := []*widget.Check{checkHealthCheck, flushDns, wingetUpdate, scanHealth, restoreHealth}
		adminRequired := false
		for _, c := range adminChecks {
			if c.Checked {
				adminRequired = true
				break
			}
		}

		// if admin required but not given
		if adminRequired && !isAdmin() {
			parent := fyne.CurrentApp().Driver().AllWindows()[0]
			fynedialog.ShowConfirm(
				"Administrator Privileges Required",
				"One or more selected actions require administrator rights.\n\nWould you like to restart ProbeDesk with elevated permissions?",
				func(confirmed bool) {
					if confirmed {
						relaunchAsAdmin()
					}
				},
				parent,
			)
			return
		}

		mu.Lock()
		defer mu.Unlock()
		status.SetText("Running...")
		output.SetText("")

		go func() {
			var report strings.Builder
			cmd.SetRemoteTarget(remoteEntry.Text)

			run := func(name string, f func() (string, error)) {
				out, err := f()
				if err != nil {
					report.WriteString(fmt.Sprintf("=== %s ===\nError: %v\n\n", name, err))
				} else {
					report.WriteString(fmt.Sprintf("=== %s ===\n%s\n\n", name, out))
				}
			}

			// run actions
			if systemCheck.Checked {
				run("System Info", cmd.GetSystemInfo)
			}
			if ipCheck.Checked {
				run("IP Config", cmd.GetIpConfigInfo)
			}
			if usbCheck.Checked {
				run("USB Devices", cmd.GetUsbInfo)
			}
			if servicesCheck.Checked {
				run("Running Services", cmd.GetServices)
			}
			if usersCheck.Checked {
				run("Local Users", cmd.GetUsersInfo)
			}
			if vpnCheck.Checked {
				run("VPN Connections", cmd.GetVpnConnections)
			}
			if productsCheck.Checked {
				run("Installed Products", cmd.GetProductsInfo)
			}
			if netuseCheck.Checked {
				run("Mapped Drives", cmd.GetNetInfo)
			}

			if traceRouteRequest.Checked {
				target := traceRoute.Text
				if target == "" {
					target = "8.8.8.8"
				}
				run(fmt.Sprintf("Trace Route (%s)", target), func() (string, error) {
					return cmd.TraceRoute(target)
				})
			}

			if checkHealthCheck.Checked {
				run("Check Health", cmd.CheckHealth)
			}
			if flushDns.Checked {
				run("Flush DNS", cmd.FlushDns)
			}
			if wingetUpdate.Checked {
				run("Windows Update (winget)", cmd.WingetUpdate)
			}
			if scanHealth.Checked {
				run("Scan Health", cmd.ScanHealth)
			}
			if restoreHealth.Checked {
				run("Restore Health", cmd.RestoreHealth)
			}

			result := report.String()
			if result == "" {
				result = "No actions selected."
			}

			fyne.Do(func() {
				output.SetText(result)
				status.SetText("Done")
			})

			cmd.CopyToClipboard(result)
		}()
	})
}

func NewCopyButton(output *widget.Entry, status *widget.Label) *widget.Button {
	return widget.NewButton("Copy Output", func() {
		err := cmd.CopyToClipboardForce(output.Text)
		if err != nil {
			status.SetText("Nothing to copy")
		} else {
			status.SetText("Output copied.")
		}
	})
}

func NewExportButton(output *widget.Entry, status *widget.Label, w fyne.Window, formatSelect *widget.Select) *widget.Button {
	return widget.NewButton("Export Report", func() {
		content := output.Text
		if content == "" {
			fynedialog.ShowInformation("Export", "No report to export.", w)
			return
		}

		format := formatSelect.Selected
		defaultName := fmt.Sprintf("report_%s.%s", time.Now().Format("2006-01-02_15-04-05"), format)

		filePath, err := nativeDialog.File().
			SetStartFile(defaultName).
			Title("Save Report As...").
			Save()
		if err != nil {
			return
		}

		if err := cmd.ExportReport(content, format, filePath); err != nil {
			fynedialog.ShowError(err, w)
			return
		}

		fynedialog.ShowInformation("Export", fmt.Sprintf("✅ Report exported successfully:\n%s", filePath), w)
		status.SetText("Report exported.")
	})
}
