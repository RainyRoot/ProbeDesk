// Copyright (c) 2025 RainyRoot
// MIT License
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Root command definition
var rootCmd = &cobra.Command{
	Use:   "probedesk",
	Short: "Collect Windows system and network information",
	Long: `ProbeDesk collects Windows system info, network configuration, 
	and installed products for auditing or support purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		var report strings.Builder

		// Handle one-off flags first
		if autocompleteInstallFlag {
			script := installAutocomplete()
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

		// If no flags set → run full collection
		if !anyFlagsSet() {
			getAllWindowsInfo()
			return
		}

		// Run the selected flags
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

		// TraceRoute example
		if traceRouteRequest {
			if len(args) < 1 {
				fmt.Println("Please specify a host or IP to trace, e.g.: probedesk --trace 8.8.8.8")
			} else {
				host := args[0]
				fmt.Printf("\n=== TraceRoute (%s) ===\n", host)
				out, _ := traceRoute(host)
				fmt.Println(out)
				report.WriteString(fmt.Sprintf("=== TraceRoute (%s) ===\n%s\n\n", host, out))
			}
		}

		// Export report
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

func Execute() {
	configureHelpAndUsage()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func configureHelpAndUsage() {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println()
		fmt.Println("ProbeDesk — quick help")
		fmt.Println("----------------------")
		fmt.Println("ProbeDesk collects system & network information to help with support or auditing.")
		fmt.Println("Usage examples:")
		fmt.Println("  probedesk         # --system --ipconfig   # run specific probes")
		fmt.Println()
		fmt.Println("Available commands:")
		printCommandsSummary(cmd)
		fmt.Println()
		cmd.Flags().PrintDefaults()
	})

	// usage function (called when cobra prints usage)
	rootCmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Println()
		fmt.Printf("Usage: %s\n\n", cmd.UseLine())
		fmt.Println("If you need help, run with `--help` or `-h`.")
		fmt.Println()
		fmt.Println("Commands:")
		printCommandsSummary(cmd)
		return nil
	})

	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		fmt.Printf("Error: %s\n\n", err.Error())
		_ = cmd.Help()
		return err
	})
	// Override the default help command
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "Show help for command",
		Run: func(cmd *cobra.Command, args []string) {
			_ = rootCmd.Help()
		},
	})
}

func printCommandsSummary(cmd *cobra.Command) {
	commands := cmd.Commands()
	if cmd.HasParent() {
		commands = cmd.Commands()
	}

	if len(commands) == 0 {
		fmt.Println("  (no commands available)")
		return
	}

	// Determine padding
	maxNameLen := 0
	for _, c := range commands {
		if c.Hidden {
			continue
		}
		if l := len(c.Name()); l > maxNameLen {
			maxNameLen = l
		}
	}
	if maxNameLen < 10 {
		maxNameLen = 10
	}

	for _, c := range commands {
		if c.Hidden {
			continue
		}
		// repeatable formatting: "  cmdName    \t short description"
		name := c.Name()
		short := strings.TrimSpace(c.Short)
		if short == "" {
			short = "(no short description)"
		}
		padding := strings.Repeat(" ", maxNameLen-len(name))
		fmt.Printf("  %s%s  %s\n", name, padding, short)
	}
}
