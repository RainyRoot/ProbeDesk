// Windows Actions

package cmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Struct: Flag + Name + Action
type WinAction struct {
	flag *bool
	name string
	run  func() (string, error)
}

func getSystemInfo() (string, error) {
	return runPowershellReturnOutput("systeminfo | Select-String 'OS Name','OS Version'")
}

func getIpConfigInfo() (string, error) {
	return runPowershellReturnOutput("ipconfig /all")
}

func getNetInfo() (string, error) {
	return runPowershellReturnOutput("net use")
}

func getProductsInfo() (string, error) {
	return runPowershellReturnOutput("Get-ItemProperty HKLM:\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\* | Select-Object DisplayName,DisplayVersion")
}

func getVpnConnections() (string, error) {
	return runPowershellReturnOutput("Get-VpnConnection")
}

func getServices() (string, error) {
	return runPowershellReturnOutput("Get-Service | Where-Object {$_.Status -eq 'Running'} | Select-Object DisplayName,Name,StartType")
}

func getUsersInfo() (string, error) {
	return runPowershellReturnOutput("Get-LocalUser | Select-Object Name,Enabled,PasswordExpires,PasswordLastSet,LastLogon")
}

func traceRoute(target string) (string, error) {
	// Validate the target: only allow letters, numbers, dots, hyphens
	if !isValidHost(target) {
		return "Invalid target: only letters, digits, dots, and hyphens are allowed.", nil
	}

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"-Command", "tracert", "-d", "-h", "10", target)

	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))

	if err != nil && output == "" {
		return fmt.Sprintf("⚠️ TraceRoute failed for %s: %v", target, err), nil
	}
	return output, nil
}

func isValidHost(input string) bool {
	// regex: allows letters, digits, dots, hyphens, but not empty or spaces
	re := regexp.MustCompile(`^[a-zA-Z0-9.\-]+$`)
	return re.MatchString(input)
}

func flushDns() (string, error) {
	if !confirmationFlag {
		return "Flushing DNS requires --yes flag to confirm.", nil
	}
	cmd := "ipconfig /flushdns"
	return runPowershellReturnOutput(cmd)
}

func wingetUpdate() (string, error) {
	if !confirmationFlag {
		return "Running winget upgrade requires --yes flag to confirm.", nil
	}
	cmd := "winget upgrade --accept-source-agreements --accept-package-agreements"
	return runPowershellReturnOutput(cmd)
}

func scanHealth() (string, error) {
	if !confirmationFlag {
		return "Scanning health requires --yes flag to confirm.", nil
	}
	return runPowershellReturnOutput("Dism /Online /Cleanup-Image /ScanHealth")
}

func restoreHealth() (string, error) {
	if !confirmationFlag {
		return "Restoring health requires --yes flag to confirm.", nil
	}
	return runPowershellReturnOutput("Dism /Online /Cleanup-Image /RestoreHealth")
}

func checkHealth() (string, error) {
	return runPowershellReturnOutput("Dism /Online /Cleanup-Image /CheckHealth")
}

func getUsbInfo() (string, error) {
	psCmd := `
	$usbDevices = Get-PnpDevice -PresentOnly |
		Where-Object {
			$_.InstanceId -match '^USB' -and
			$_.FriendlyName -and
			$_.Manufacturer -and
			$_.Manufacturer -notmatch 'Standard system devices' -and
			$_.Manufacturer -notmatch 'Standard USB Host Controller' -and
			$_.Manufacturer -notmatch 'Standard USB HUBs' -and
			$_.Manufacturer -notmatch 'Generic USB Audio' -and
			$_.Class -notmatch 'HIDClass'
		} |
		Select-Object FriendlyName, Manufacturer, Class

	if (!$usbDevices) {
		Write-Host "No external USB devices detected."
	} else {
		$usbDevices | ForEach-Object {
			Write-Host ("• " + $_.FriendlyName)
			Write-Host ("    Manufacturer: " + $_.Manufacturer)
			if ($_.Class) { Write-Host ("    Type:         " + $_.Class) }
			Write-Host ""
		}
	}
	`
	return runPowershellReturnOutput(psCmd)
}
