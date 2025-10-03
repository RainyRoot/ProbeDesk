// Copyright (c) 2025 RainyRoot
// MIT License
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
	Use: "win",
	Long: `ProbeDesk can collect various information about a Windows system,
including system details, network configuration, BIOS info, and installed products.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no flags are set, get all info
		if !systemFlag && !ipconfigFlag && !netuseFlag && !biosFlag && !productsFlag {
			getAllWindowsInfo()
			return
		}

		if systemFlag {
			getSystemInfo()
		}
		if ipconfigFlag {
			getIpConfigInfo()
		}
		if netuseFlag {
			getNetInfo()
		}
		if biosFlag {
			getBiosInfo()
		}
		if productsFlag {
			getProductsInfo()
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

// Collect all Windows information
func getAllWindowsInfo() {
	fmt.Println("=== Windows System Info ===")
	getSystemInfo()
	getIpConfigInfo()
	getNetInfo()
	getBiosInfo()
	getProductsInfo()
}

// Different functions to get specific information
func getSystemInfo() {
	fmt.Println("=== System Info ===")
	runCommand(`systeminfo | findstr /B /C:"OS Name" /C:"OS Version" /C:"Total Physical Memory"`)
}

func getNetInfo() {
	fmt.Println("=== Network Info ===")
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

func getIpConfigInfo() {
	fmt.Println("=== IP Configuration Info ===")
	runCommand("ipconfig /all")
}

// Executes a command and prints its output
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
