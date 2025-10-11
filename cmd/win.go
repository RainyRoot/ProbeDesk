// Copyright (c) 2025 RainyRoot
// MIT License
package cmd

import (
	"fmt"
	"html"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Flags
var (
	systemFlag            bool
	ipconfigFlag          bool
	netuseFlag            bool
	productsFlag          bool
	getVpnConnectionsFlag bool
	getServicesFlag       bool
	getUserInfoFlag       bool
	traceRouteRequest     bool
	getUsbInfoFlag        bool
	autocompleteFlag      bool
	remoteTarget          string
	reportFormat          string
	confirmationFlag      bool
	flushDnsFlag          bool
	wingetUpdateFlag      bool
	scanHealthFlag        bool
	checkHealthFlag       bool
	restoreHealthFlag     bool
)

// Struct: Flag + Name + Action
type WinAction struct {
	flag *bool
	name string
	run  func() (string, error)
}

// Define all Windows-related actions
var winActions = []WinAction{
	{&systemFlag, "system", getSystemInfo},
	{&ipconfigFlag, "ipconfig", getIpConfigInfo},
	{&netuseFlag, "netuse", getNetInfo},
	{&productsFlag, "products", getProductsInfo},
	{&getVpnConnectionsFlag, "vpn", getVpnConnections},
	{&getServicesFlag, "services", getServices},
	{&getUserInfoFlag, "users", getUsersInfo},
	{&getUsbInfoFlag, "usb", getUsbInfo},
	{&checkHealthFlag, "check-health", checkHealth},

	//---additional flags defined in init():---
	// traceRoute -> trace <ip/host>
	// remoteTarget -> remote <host>
	// reportFormat  -> report <html|md>
	// autocomplete -> --win <tab>

	//Cases requiring confirmation:
	// confirmationFlag -> ----yes to confirm prompts like flushing DNS and winget update
	// flushDns -> --yes
	// wingetUpdate -> --yes
	// restoreHealth -> --yes
	// scanHealth -> --yes

}

// Command definition
var winCmd = &cobra.Command{
	Use:   "win",
	Short: "Collect Windows system and network information",
	Long: `ProbeDesk can collect various information about a Windows system,
including system details, network configuration, BIOS info, and installed products.`,
	Run: func(cmd *cobra.Command, args []string) {
		var report strings.Builder

		if autocompleteFlag {
			script := autocompleteScript()
			fmt.Println(script)
			return
		}

		if flushDnsFlag {
			out, _ := flushDns()
			fmt.Println(out)
			return
		}

		if wingetUpdateFlag {
			out, _ := wingetUpdate()
			fmt.Println(out)
			return
		}

		if scanHealthFlag {
			out, _ := scanHealth()
			fmt.Println(out)
			return
		}

		if restoreHealthFlag {
			out, _ := restoreHealth()
			fmt.Println(out)
			return
		}

		if !anyFlagsSet() {
			getAllWindowsInfo()
			return
		}

		// Report Header
		for _, a := range winActions {
			if *a.flag {
				fmt.Printf("\n=== %s ===\n", strings.Title(a.name))

				out, err := a.run()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error running %s: %v\n", a.name, err)
				}

				fmt.Println(out)
				report.WriteString(fmt.Sprintf("=== %s ===\n%s\n\n", strings.Title(a.name), out))
			}
		}

		if traceRouteRequest {

			if len(args) < 1 {
				fmt.Println("Please specify a host or IP to trace, e.g.: probedesk win --trace 8.8.8.8")
			} else {
				host := args[0]
				fmt.Printf("\n=== TraceRoute (%s) ===\n", host)
				out, _ := traceRoute(host)
				fmt.Println(out)
				report.WriteString(fmt.Sprintf("=== TraceRoute (%s) ===\n%s\n\n", host, out))
			}
		}

		// Export report if format specified
		finalReport := report.String()
		if finalReport != "" {
			copyToClipboard(finalReport)
			if reportFormat != "" {
				if err := exportReport(finalReport, reportFormat, ""); err != nil {
					fmt.Println("Error exporting report:", err)
				} else {
					fmt.Printf("✅ Report exported successfully as %s\n", reportFormat)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(winCmd)

	for _, a := range winActions {
		winCmd.Flags().BoolVar(a.flag, a.name, false, fmt.Sprintf("Get %s info", a.name))
	}

	// Additional flags
	winCmd.Flags().BoolVar(&traceRouteRequest, "trace", false, "Trace a host (add host as argument)")
	winCmd.Flags().StringVar(&remoteTarget, "remote", "", "Run commands remotely on target host (requires PS Remoting)")
	winCmd.Flags().StringVar(&reportFormat, "report", "", "Export collected data to report (html or md)")
	winCmd.Flags().BoolVar(&autocompleteFlag, "autocomplete", false, "Enable PowerShell autocomplete for this session")
	winCmd.Flags().BoolVar(&confirmationFlag, "yes", false, "confirmation flag")
	winCmd.Flags().BoolVar(&flushDnsFlag, "flush", false, "Flush DNS cache (requires --yes)")
	winCmd.Flags().BoolVar(&wingetUpdateFlag, "winget-update", false, "Update installed packages using winget (requires --yes)")
	winCmd.Flags().BoolVar(&scanHealthFlag, "scan-health", false, "Scan system health (requires --yes)")
	winCmd.Flags().BoolVar(&restoreHealthFlag, "restore-health", false, "Restore system health (requires --yes)")

	winCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("\nProbeDesk — List of available flags/modules for 'win' command")
		fmt.Println("-------------------------------")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			fmt.Printf("  --%-12s %s\n", f.Name, f.Usage)
		})
		fmt.Println("\nUsage examples:")
		fmt.Println("  probedesk win --system --ipconfig")
		fmt.Println("  probedesk win --bios --remote server01")
		fmt.Println("  probedesk win --system --report html")
	})
}

func anyFlagsSet() bool {
	for _, a := range winActions {
		if *a.flag {
			return true
		}
	}
	return traceRouteRequest
}

// ========================
// Centralized Functions
// ========================

func runPowershellReturnOutput(command string) (string, error) {
	// Force Powershell UTF-8 output
	psCmd := fmt.Sprintf("[Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8; %s", command)
	if remoteTarget != "" {
		psCmd = fmt.Sprintf(`Invoke-Command -ComputerName %s -ScriptBlock { [Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8; %s }`, remoteTarget, command)
	}

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psCmd)

	// CombinedOutput []byte UTF-8
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))

	if output == "" {
		if err != nil {
			return fmt.Sprintf("⚠️ Error executing: %v", err), nil
		}
		return "No output (possibly no data found).\n", nil
	}
	return output, nil
}

func copyToClipboard(content string) {
	if content == "" {
		fmt.Println("Nothing to copy.")
		return
	}
	if err := clipboard.WriteAll(content); err != nil {
		fmt.Println("Error copying to clipboard:", err)
	} else {
		fmt.Println("✅ Output copied to clipboard!")
	}
}

func exportReport(content, format, path string) error {
	if path == "" {
		// Desktop of the current user
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("Couldnt determine current user: %v", err)
		}
		path = filepath.Join(usr.HomeDir, "Desktop")
	}

	// Create filename
	filename := filepath.Join(path, fmt.Sprintf("report_%s.%s", time.Now().Format("2006-01-02_15-04-05"), format))

	switch format {
	case "md":
		return os.WriteFile(filename, []byte("```markdown\n"+content+"\n```"), 0644)
	case "html":
		htmlOut := "<html><body><pre>" + html.EscapeString(content) + "</pre></body></html>"
		return os.WriteFile(filename, []byte(htmlOut), 0644)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// PowerShell Autocomplete Script
func autocompleteScript() string {
	return `
$flags = @("system","ipconfig","netuse","products","vpn","services","users","usb","trace","remote","report")

Register-ArgumentCompleter -CommandName "probedesk" -ScriptBlock {
    param($commandName, $parameterName, $wordToComplete, $commandAst, $fakeBoundParameter)

    # Prüfen, ob der Befehl 'win' verwendet wird
    if ($commandAst.CommandElements.Count -gt 1 -and $commandAst.CommandElements[1].Value -eq "win") {
        $flags | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { 
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_) 
            }
    }
}`
}

// ========================
// Windows Actions
// ========================

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

func traceRoute(target string) (string, error) {
	cmd := fmt.Sprintf("tracert -d -h 10 %s", target)

	out, err := runPowershellReturnOutput(cmd)
	if err != nil {
		return fmt.Sprintf("⚠️ TraceRoute could not reach %s or error:\n%s", target, out), nil
	}
	return out, nil
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

// ========================
// All-in-one Collector
// ========================

func getAllWindowsInfo() {
	fmt.Println("=== Collecting All Windows Info ===")
	var report strings.Builder

	// Report Header current User
	usr, err := user.Current()
	username := "Unknown"
	if err == nil {
		username = usr.Username
	}
	report.WriteString(fmt.Sprintf("Report generated by: %s\nDate: %s\n\n", username, time.Now().Format("2006-01-02 15:04:05")))

	for _, a := range winActions {
		fmt.Printf("\n=== %s ===\n", strings.Title(a.name))
		out, _ := a.run()
		fmt.Println(out)
		report.WriteString(fmt.Sprintf("=== %s ===\n%s\n\n", strings.Title(a.name), out))
	}

	// Tracing example
	examples := []string{"srv-fls-001.ad.adler-group.com", "8.8.8.8"}
	for _, host := range examples {
		fmt.Printf("\n=== TraceRoute (%s) ===\n", host)
		out, _ := traceRoute(host)
		fmt.Println(out)
		report.WriteString(fmt.Sprintf("=== TraceRoute (%s) ===\n%s\n\n", host, out))
	}

	// Export report
	finalReport := report.String()
	copyToClipboard(finalReport)
	if err := exportReport(finalReport, "html", ""); err != nil {
		fmt.Println("Error exporting report:", err)
	} else {
		fmt.Printf("✅ Report exported successfully as %s\n", "html")
	}
}
