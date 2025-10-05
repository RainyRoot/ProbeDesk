// Copyright (c) 2025 RainyRoot
// MIT License
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "probedesk",
	Short: "ProbeDesk collects system and network information",
	Long: `ProbeDesk is a tool to collect various Windows system information,
network configuration, and installed software details for support purposes.

Use subcommands like 'win' to perform specific tasks.`,
}

// Execute adds all child commands to the root command and sets flags
func Execute() {
	configureHelpAndUsage()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func configureHelpAndUsage() {
	// prints a default blurb + commands list
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println()
		fmt.Println("ProbeDesk — quick help")
		fmt.Println("----------------------")
		fmt.Println("ProbeDesk collects system & network information to help with support or auditing.")
		fmt.Println("Usage examples:")
		fmt.Println("  probedesk win     # run Windows-specific probes")
		fmt.Println("  probedesk win     # --system --ipconfig   # run specific probes")
		fmt.Println()
		fmt.Println("Available commands:")
		printCommandsSummary(cmd)
		fmt.Println()
		// Print flags for the current command (if any)
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

	// When an unknown flag or invalid flag usage occurs
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		fmt.Printf("Error: %s\n\n", err.Error())
		_ = cmd.Help()
		return err
	})
	// Override the default help command to show our custom help
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "Show help for command",
		Run: func(cmd *cobra.Command, args []string) {
			_ = rootCmd.Help()
		},
	})
}

func printCommandsSummary(cmd *cobra.Command) {
	// We want top-level commands from root; if cmd is root, use its children.
	commands := cmd.Commands()
	if cmd.HasParent() {
		commands = cmd.Commands()
	}

	if len(commands) == 0 {
		fmt.Println("  (no commands available)")
		return
	}

	// Determine padding for nice alignment
	maxNameLen := 0
	for _, c := range commands {
		// skip hidden commands
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
