package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// winCmd represents "win get"
var winCmd = &cobra.Command{
	Use:   "win",
	Short: "Windows system info commands",
	Long:  `Collect various Windows system information for support purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "all" {
			getAllWindowsInfo()
		} else {
			fmt.Println("Unknown subcommand. Use 'win get all'")
		}
	},
}

func init() {
	rootCmd.AddCommand(winCmd)
}

func getAllWindowsInfo() {
	fmt.Println("=== Windows System Info ===")
	runCommand("systeminfo get OS Name")                          // Win-Version, Hersteller, RAM, etc.
	runCommand("ipconfig /all")                                   // IPv4 / IPv6
	runCommand("net use")                                         // Netzlaufwerke
	runCommand("wmic bios get serialnumber,manufacturer,version") // BIOS
	runCommand("wmic product get name,version")                   // Installed programs
}

func runCommand(command string) {
	fmt.Printf("\n> %s\n", command)
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		return
	}
	fmt.Println(string(out))
}
