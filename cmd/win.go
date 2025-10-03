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
	getUsersFlag bool
)

// winCmd represents "win get"
var winCmd = &cobra.Command{
	Use: "win",
	Long: `ProbeDesk can collect various information about a Windows system,
including system details, network configuration, BIOS info, and installed products.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no flags are set, get all info
		if !systemFlag && !ipconfigFlag && !netuseFlag && !biosFlag && !productsFlag && !getUsersFlag {
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

//Ideas
/*
vpn? vlan netz name
pingen?
lizenz keys?
how old user pw? last change
how many user profiles + name
welche dienste laufen?
welche geräte sind angeschlossen
*/

func init() {
	rootCmd.AddCommand(winCmd)

	// Flags
	winCmd.Flags().BoolVar(&systemFlag, "system", false, "Get system info")
	winCmd.Flags().BoolVar(&ipconfigFlag, "ipconfig", false, "Get IP configuration info") //TODO filter for specific fields
	winCmd.Flags().BoolVar(&netuseFlag, "netuse", false, "Get network use info")          //TODO testing
	winCmd.Flags().BoolVar(&biosFlag, "bios", false, "Get BIOS info")
	winCmd.Flags().BoolVar(&productsFlag, "products", false, "Get installed products info") //TODO filter for specific fields
	winCmd.Flags().BoolVar(&getUsersFlag, "users", false, "Get user accounts info")
}

// Collect all Windows information
func getAllWindowsInfo() {
	fmt.Println("=== Windows System Info ===")
	getSystemInfo()
	getIpConfigInfo()
	getNetInfo()
	getBiosInfo()
	getProductsInfo()
	getUsers()
}

// Different functions to get specific information
func getSystemInfo() {
	fmt.Println("=== System Info ===")

	out, err := exec.Command("systeminfo").Output()
	if err != nil {
		fmt.Printf("Error running systeminfo: %v\n", err)
		return
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "OS Name") ||
			strings.HasPrefix(line, "OS Version") ||
			strings.HasPrefix(line, "Total Physical Memory") {
			fmt.Println(line)
		}
	}
}

func getIpConfigInfo() {
	fmt.Println("=== IP Configuration Info ===")
	runCommand("ipconfig /all")
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

func getUsers() {
	fmt.Println("=== User Accounts Info ===")
	runCommand("wmic useraccount get name")
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
