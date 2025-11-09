package cmd

// Flags
var (
	// Global
	guiMode bool

	// System & Network
	systemFlag            bool
	ipconfigFlag          bool
	netuseFlag            bool
	productsFlag          bool
	getVpnConnectionsFlag bool
	getServicesFlag       bool
	getUserInfoFlag       bool
	getUsbInfoFlag        bool
	traceRouteRequest     bool

	// Special Commands
	autocompleteInstallFlag bool
	remoteTarget            string
	reportFormat            string

	// Confirmation / Actions
	confirmationFlag bool
	flushDnsFlag     bool
	wingetUpdateFlag bool
	resetWindowsUpdateFlag bool

	// DISM / Windows Health
	scanHealthFlag    bool
	checkHealthFlag   bool
	restoreHealthFlag bool
	searchUninstallStringFlag bool
)

//>>>>>>>>>>>>>>>>TEST ME<<<<<<<<<<<<<<<<<<  verschieben?
func SetGuiMode(mode bool) {
	guiMode = mode
}
