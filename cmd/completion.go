// completion.go
package cmd

import (
	"os"

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

func init() {
	rootCmd.AddCommand(completionCmd)
}
