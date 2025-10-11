// Copyright (c) 2025 RainyRoot
// MIT License
package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

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

func init() {
	rootCmd.AddCommand(winCmd)

	for _, a := range winActions {
		winCmd.Flags().BoolVar(a.flag, a.name, false, fmt.Sprintf("Get %s info", a.name))
	}

	// Additional flags
	winCmd.Flags().BoolVar(&traceRouteRequest, "trace", false, "Trace a host (add host as argument)")
	winCmd.Flags().StringVar(&remoteTarget, "remote", "", "Run commands remotely on target host (requires PS Remoting)")
	winCmd.Flags().StringVar(&reportFormat, "report", "", "Export collected data to report (html or md)")
	winCmd.Flags().BoolVar(&autocompleteInstallFlag, "autocomplete-install", false, "Install PowerShell autocomplete for 'win' command (persists in profile)")
	winCmd.Flags().BoolVar(&confirmationFlag, "yes", false, "confirmation flag")
	winCmd.Flags().BoolVar(&flushDnsFlag, "flush", false, "Flush DNS cache (requires --yes)")
	winCmd.Flags().BoolVar(&wingetUpdateFlag, "winget-update", false, "Update installed packages using winget (requires --yes)")
	winCmd.Flags().BoolVar(&scanHealthFlag, "scan-health", false, "Scan system health (requires --yes)")
	winCmd.Flags().BoolVar(&restoreHealthFlag, "restore-health", false, "Restore system health (requires --yes)")
	winCmd.Flags().BoolVar(&pingFlag, "ping", false, "Ping a host (add host as argument)")

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

// Command definition
var winCmd = &cobra.Command{
	Use:   "win",
	Short: "Collect Windows system and network information",
	Long: `ProbeDesk can collect various information about a Windows system,
including system details, network configuration, BIOS info, and installed products.`,
	Run: func(cmd *cobra.Command, args []string) {
		var report strings.Builder

		if autocompleteInstallFlag {
			script := installAutocomplete()
			fmt.Println(script)
			return
		}

		if pingFlag {
			out, _ := ping()
			fmt.Println(out)
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

func anyFlagsSet() bool {
	for _, a := range winActions {
		if *a.flag {
			return true
		}
	}
	return traceRouteRequest
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
