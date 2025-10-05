// Copyright (c) 2025 RainyRoot
// MIT License
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	systemFlag            bool
	ipconfigFlag          bool
	netuseFlag            bool
	biosFlag              bool
	productsFlag          bool
	getUsersFlag          bool
	getVpnConnectionsFlag bool
	getServicesFlag       bool
	getPasswordInfoFlag   bool
	pingRequest           bool
	getUsbInfoFlag        bool
)

// winCmd represents "win get"
var winCmd = &cobra.Command{
	Use: "win",
	Long: `ProbeDesk can collect various information about a Windows system,
including system details, network configuration, BIOS info, and installed products.`,
	Short: "Collect Windows system and network information",
	Run: func(cmd *cobra.Command, args []string) {
		// If no flags are set, get all info
		if !systemFlag && !ipconfigFlag && !netuseFlag && !biosFlag && !productsFlag && !getUsersFlag && !getVpnConnectionsFlag && !getServicesFlag && !getPasswordInfoFlag && !pingRequest && !getUsbInfoFlag {
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
		if getUsersFlag {
			getUsers()
		}
		if getVpnConnectionsFlag {
			getVpnConnections()
		}
		if getServicesFlag {
			getServices()
		}
		if getPasswordInfoFlag {
			getPasswordInfo()
		}
		if pingRequest {
			if len(args) < 1 {
				fmt.Println("Please specify a host or IP to ping, e.g.: probedesk win --ping 8.8.8.8")
				return
			}
			pingHost(args[0])
		}
		if getUsbInfoFlag {
			getUsbInfo()
		}
	},
}

func init() {
	rootCmd.AddCommand(winCmd)

	// Flags
	winCmd.Flags().BoolVar(&systemFlag, "system", false, "Get system info")
	winCmd.Flags().BoolVar(&ipconfigFlag, "ipconfig", false, "Get IP configuration info")
	winCmd.Flags().BoolVar(&netuseFlag, "netuse", false, "Get network use info") //TODO testing
	winCmd.Flags().BoolVar(&biosFlag, "bios", false, "Get BIOS info")
	winCmd.Flags().BoolVar(&productsFlag, "products", false, "Get installed products info") //TODO weird output
	winCmd.Flags().BoolVar(&getUsersFlag, "users", false, "Get user accounts info")
	winCmd.Flags().BoolVar(&getVpnConnectionsFlag, "vpn", false, "Get VPN connections info") //TODO testing
	winCmd.Flags().BoolVar(&getServicesFlag, "services", false, "Get running services info")
	winCmd.Flags().BoolVar(&getPasswordInfoFlag, "passwords", false, "Get user password info") //TODO fixing
	winCmd.Flags().BoolVar(&pingRequest, "ping", false, "Ping a host (add host as argument)")
	winCmd.Flags().BoolVar(&getUsbInfoFlag, "usb", false, "Get connected USB devices info")

	// Custom HelpFunc for Flags as Modules
	winCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println()
		fmt.Println("ProbeDesk — List of available flags/modules for 'win' command")
		fmt.Println("-------------------------------")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			fmt.Printf("  --%-12s %s\n", f.Name, f.Usage)
		})
		fmt.Println()
		fmt.Println("Usage examples:")
		fmt.Println("  probedesk win --system --ipconfig   # run specific probes")
	})
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
	getVpnConnections()
	getServices()
	getPasswordInfo()
	pingHost("srv-fls-001.ad.adler-group.com")
	pingHost("8.8.8.8")
	getUsbInfo()
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
			strings.HasPrefix(line, "OS Version") {
			fmt.Println(line)
		}
	}
}

func getIpConfigInfo() {
	fmt.Println("\n=== IP Configuration Info ===")
	runPowershell("ipconfig /all")
}

func getNetInfo() {
	fmt.Println("\n=== Network Info ===")
	runPowershell("net use")
}

func getBiosInfo() {
	fmt.Println("\n=== BIOS Info ===")
	runPowershell("Get-CimInstance Win32_BIOS | Select-Object SerialNumber,Manufacturer,Version")
}

func getProductsInfo() {
	fmt.Println("\n=== Products Info ===")
	runPowershell("Get-ItemProperty HKLM:\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\* | Select-Object DisplayName,DisplayVersion")
}

func getUsers() {
	fmt.Println("\n=== User Accounts Names ===")
	runPowershell("Get-LocalUser | Select-Object Name")
}

func getVpnConnections() {
	fmt.Println("\n=== VPN Connections ===")
	runPowershell("Get-VpnConnection")
}

func getServices() {
	fmt.Println("\n=== Running Services ===")
	runPowershell("Get-Service | Where-Object {$_.Status -eq 'Running'} | Select-Object DisplayName,Name,StartType")
}

// filter out default user, Guest, admin
func getPasswordInfo() {
	fmt.Println("\n=== User Password Info ===")
	runPowershell("Get-LocalUser | Select-Object Name,Enabled,PasswordExpires,PasswordLastSet,LastLogon")
}

func pingHost(target string) {
	fmt.Printf("\n=== Pinging Host: %s ===\n", target)
	runPowershell(fmt.Sprintf("ping -n 4 %s", target))
}

func getUsbInfo() {
	fmt.Println("\n=== USB Devices Info ===")
	runPowershell("Get-WmiObject Win32_PnPEntity | Where-Object { $_.PNPClass -eq 'USB' -or $_.Name -match 'USB' } | Select-Object Name, Manufacturer")
}

// keep for later
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

// function to run powershell commands
func runPowershell(command string) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", command)
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))

	if err != nil {
		fmt.Printf("Error running PowerShell command: %v\n", err)
		return
	}

	if output == "" {
		fmt.Println("No output (possibly no data found).")
	} else {
		fmt.Println(output)
	}
}
