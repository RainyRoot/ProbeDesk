// Windows Actions

package cmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

// Struct: Flag + Name + Action
type WinAction struct {
	flag *bool
	name string
	run  func() (string, error)
}

func init() {
	//var traceTarget string
	// System & Network flags
	rootCmd.Flags().BoolVar(&systemFlag, "system", false, "Collect system info")
	rootCmd.Flags().BoolVar(&ipconfigFlag, "ipconfig", false, "Collect IP configuration info")
	rootCmd.Flags().BoolVar(&netuseFlag, "netuse", false, "Show mapped network drives")
	rootCmd.Flags().BoolVar(&productsFlag, "products", false, "Show installed products")
	rootCmd.Flags().BoolVar(&getVpnConnectionsFlag, "vpn", false, "Show VPN connections")
	rootCmd.Flags().BoolVar(&getServicesFlag, "services", false, "Show running services")
	rootCmd.Flags().BoolVar(&getUserInfoFlag, "users", false, "Show local users")
	rootCmd.Flags().BoolVar(&getUsbInfoFlag, "usb", false, "Show connected USB devices")
	rootCmd.Flags().BoolVar(&checkHealthFlag, "check-health", false, "Check Windows health status")

	// One-off / special flags
	rootCmd.Flags().BoolVar(&traceRouteRequest, "trace", false, "Trace a host (provide host as argument)")
	rootCmd.Flags().BoolVar(&autocompleteInstallFlag, "autocomplete-install", false, "Install PowerShell autocomplete")
	rootCmd.Flags().StringVar(&remoteTarget, "remote", "", "Run commands remotely on target host (requires PS Remoting)")
	rootCmd.Flags().StringVar(&reportFormat, "report", "", "Export collected data to report (html or md)")
	rootCmd.Flags().BoolVar(&confirmationFlag, "yes", false, "Confirmation flag")
	rootCmd.Flags().BoolVar(&flushDnsFlag, "flush", false, "Flush DNS cache (requires --yes)")
	rootCmd.Flags().BoolVar(&wingetUpdateFlag, "winget-update", false, "Update installed packages using winget (requires --yes)")
	rootCmd.Flags().BoolVar(&scanHealthFlag, "scan-health", false, "Scan system health (requires --yes)")
	rootCmd.Flags().BoolVar(&restoreHealthFlag, "restore-health", false, "Restore system health (requires --yes)")
	rootCmd.Flags().BoolVar(&searchUninstallStringFlag, "ustring", false, "Search for uninstall strings of installed products")
	rootCmd.Flags().BoolVar(&resetWindowsUpdateFlag, "reset-windows-update", false, "Reset Windows Update components (requires --yes)")
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

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

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
	if !confirmationFlag && !guiMode {
		return "Flushing DNS requires --yes flag to confirm.", nil
	}
	return runPowershellReturnOutput("ipconfig /flushdns")
}

func wingetUpdate() (string, error) {
	if !confirmationFlag && !guiMode {
		return "Running winget upgrade requires --yes flag to confirm.", nil
	}

	// Run winget with silent/progress-suppressing flags
	out, err := runPowershellReturnOutput(
		"winget upgrade --accept-source-agreements --accept-package-agreements --silent",
	)
	if err != nil {
		return out, err
	}

	// Clean any remaining progress bar characters or lines
	re := regexp.MustCompile(`(?m)^[\s█▒]+$`) // lines with only progress characters
	cleanOutput := re.ReplaceAllString(out, "")

	reInline := regexp.MustCompile(`[█▒]+`) // remove inline block characters
	cleanOutput = reInline.ReplaceAllString(cleanOutput, "")

	return cleanOutput, nil
}

func scanHealth() (string, error) {
	if !confirmationFlag && !guiMode {
		return "Scanning health requires --yes flag to confirm.", nil
	}
	return runPowershellReturnOutput("Dism /Online /Cleanup-Image /ScanHealth")
}

func restoreHealth() (string, error) {
	if !confirmationFlag && !guiMode {
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

func searchUninstallStringAndMSI(query string) (string, error) {
	// Basic sanity check
	if strings.TrimSpace(query) == "" {
		return "Query cannot be empty or invalid.", nil
	}

	psScript := fmt.Sprintf(`
$paths = @(
    "HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall",
    "HKLM:\Software\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall",
    "HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall"
)

$result = @()

foreach ($path in $paths) {
    try {
        $keys = Get-ChildItem $path -ErrorAction SilentlyContinue
        foreach ($key in $keys) {
            $item = Get-ItemProperty $key.PSPath -ErrorAction SilentlyContinue

            if ($item.DisplayName -and ($item.DisplayName -match "(?i)%s")) {
                $prod = $key.PSChildName
                if ($prod -notmatch '^{[0-9A-Fa-f\-]+}$') {
                    $prod = "(N/A or non-MSI installer)"
                }

                $result += [PSCustomObject]@{
                    DisplayName     = $item.DisplayName
                    UninstallString = $item.UninstallString
                    ProductCode     = $prod
                }
            }
        }
    } catch {}
}

if ($result.Count -eq 0) {
    Write-Output "No matching uninstall entry found for '%s'."
} else {
    $result | ForEach-Object {
"{0}
{1}
{2}
---" -f $_.DisplayName, $_.UninstallString, $_.ProductCode
    }
}
`, query, query)

	return runPowershellReturnOutput(psScript)
}

func resetWindowsUpdate() (string, error) {

	if !confirmationFlag && !guiMode {
		return "Running winget upgrade requires --yes flag to confirm.", nil
	}
	if !IsAdmin() {
		return "Resetting Windows Update requires administrative privileges. Please run as administrator.", nil
	}
	psCmd := `
    Stop-Service wuauserv -Force
    Stop-Service bits -Force
    Stop-Service cryptsvc -Force

    Rename-Item "C:\Windows\SoftwareDistribution" "SoftwareDistribution.old" -Force
    Rename-Item "C:\Windows\System32\catroot2" "catroot2.old" -Force

    Start-Service wuauserv
    Start-Service bits
    Start-Service cryptsvc

    wuauclt /detectnow
    "Windows Update components reset successfully. Please check for updates again."
    `

	return runPowershellReturnOutput(psCmd)
}
