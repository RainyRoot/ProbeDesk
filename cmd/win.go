package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	systemFlag   bool
	ipconfigFlag bool
	netuseFlag   bool
	biosFlag     bool
	productsFlag bool
)

// winCmd represents "win get"
var winCmd = &cobra.Command{
	Use:   "win",
	Short: "Windows system info commands",
	Long:  `Collect various Windows system information for support purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !systemFlag && !ipconfigFlag && !netuseFlag && !biosFlag && !productsFlag {
			getAllWindowsInfo()
		} else if systemFlag {
			getSystemInfo()
		} else if ipconfigFlag {
			getNetInfo()
		} else if netuseFlag {
			getNetInfo()
		} else if biosFlag {
			getBiosInfo()
		} else if productsFlag {
			getProductsInfo()
		} else {
			fmt.Println("Unknown subcommand. Use '--help' for available options.")
		}
	},
}

func init() {
	rootCmd.AddCommand(winCmd)

	// Flags
	winCmd.Flags().BoolVar(&systemFlag, "system", false, "Get system info")
	winCmd.Flags().BoolVar(&ipconfigFlag, "ipconfig", false, "Get IP configuration info")
	winCmd.Flags().BoolVar(&netuseFlag, "netuse", false, "Get network use info")
	winCmd.Flags().BoolVar(&biosFlag, "bios", false, "Get BIOS info")
	winCmd.Flags().BoolVar(&productsFlag, "products", false, "Get installed products info")
}

func getAllWindowsInfo() {
	fmt.Println("=== Windows System Info ===")
	runCommand(`systeminfo | findstr /B /C:"OS Name" /C:"OS Version" /C:"BIOS Version" /C:"Total Physical Memory" /C:"Available Physical Memory" /C:"Domain" /C:"Logon Server"`) // Win-Version, Hersteller, RAM, etc.
	runCommand("ipconfig /all")                                                                                                                                                  // IPv4 / IPv6
	runCommand("net use")                                                                                                                                                        // Netzlaufwerke
	runCommand("wmic bios get serialnumber,manufacturer,version")                                                                                                                // BIOS
	runCommand("wmic product get name,version")                                                                                                                                  // Installed programs
}

// Different functions to get specific information
func getSystemInfo() {
	fmt.Println("=== System Info ===")
	runCommand(`systeminfo | findstr /B /C:"OS Name" /C:"OS Version" /C:"Total Physical Memory"`)
}

func getNetInfo() {
	fmt.Println("=== Network Info ===")
	runCommand("ipconfig /all")
	runCommand("net use")
}

func getBiosInfo() {
	fmt.Println("=== BIOS Info ===")
	runCommand("wmic bios get serialnumber,manufacturer,version")
}

func getProductsInfo() {
	fmt.Println("=== Products Info ===")
	runCommand("wmic product get name,version")
}

// Default run command function
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
