package cmd

func GetSystemInfo() (string, error)     { return getSystemInfo() }
func GetIpConfigInfo() (string, error)   { return getIpConfigInfo() }
func GetUsbInfo() (string, error)        { return getUsbInfo() }
func GetServices() (string, error)       { return getServices() }
func GetUsersInfo() (string, error)      { return getUsersInfo() }
func GetVpnConnections() (string, error) { return getVpnConnections() }
func GetProductsInfo() (string, error)   { return getProductsInfo() }
func GetNetInfo() (string, error)        { return getNetInfo() }
func CheckHealth() (string, error)       { return checkHealth() }
func CopyToClipboard(content string)     { copyToClipboard(content) }
func CopyToClipboardForce(content string) error { return copyToClipboardForce(content) }
func ExportReport(c, f, p string) error  { return exportReport(c, f, p) }
func SetRemoteTarget(t string)           { remoteTarget = t }
func TraceRoute(target string) (string, error) { return traceRoute(target) }

// Systemcare
func FlushDns() (string, error)      { return flushDns() }
func WingetUpdate() (string, error)  { return wingetUpdate() }
func ScanHealth() (string, error)    { return scanHealth() }
func RestoreHealth() (string, error) { return restoreHealth() }
func SearchUninstallStringAndMSI(query string) (string, error) { return searchUninstallStringAndMSI(query) }
func ResetWindowsUpdate() (string, error) { return resetWindowsUpdate() }