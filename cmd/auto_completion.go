// completion.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for your shell.
For PowerShell: 
  probedesk completion powershell > probedesk.ps1
  . ./probedesk.ps1
`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"powershell"},
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "powershell":
			if err := cmd.Root().GenPowerShellCompletion(os.Stdout); err != nil {
				os.Exit(1)
			}
		}
	},
}

func installAutocomplete() error {
	profilePath := filepath.Join(os.Getenv("USERPROFILE"), "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")

	script := autocompleteScript()

	// Check: Create profile file if it doesn't exist
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(profilePath), 0755); err != nil {
			return fmt.Errorf("failed to create profile directory: %v", err)
		}
		if _, err := os.Create(profilePath); err != nil {
			return fmt.Errorf("failed to create profile file: %v", err)
		}
	}

	// Append script if it doesn't already exist
	f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open profile file: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString("\n" + script + "\n")
	if err != nil {
		return fmt.Errorf("failed to write autocomplete script: %v", err)
	}

	return nil
}

// PowerShell Autocomplete Script
func autocompleteScript() string {
	return `
$flags = @("system","ipconfig","netuse","products","vpn","services","users","usb","trace","remote","report","flushdns","winget-update","scan-health","check-health","restore-health","autocomplete-install")

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

func init() {
	rootCmd.AddCommand(completionCmd)
}
