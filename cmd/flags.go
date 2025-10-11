package cmd

// Flags
var (
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

	// DISM / Windows Health
	scanHealthFlag    bool
	checkHealthFlag   bool
	restoreHealthFlag bool
)
