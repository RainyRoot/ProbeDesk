package components

import "github.com/RainyRoot/ProbeDesk/cmd"

var FlagTooltips = map[string]string{
	"system":               cmd.GetFlagDescription("system"),
	"ipconfig":             cmd.GetFlagDescription("ipconfig"),
	"netuse":               cmd.GetFlagDescription("netuse"),
	"products":             cmd.GetFlagDescription("products"),
	"vpn":                  cmd.GetFlagDescription("vpn"),
	"services":             cmd.GetFlagDescription("services"),
	"users":                cmd.GetFlagDescription("users"),
	"usb":                  cmd.GetFlagDescription("usb"),
	"check-health":         cmd.GetFlagDescription("check-health"),
	"trace":                cmd.GetFlagDescription("trace"),
	"autocomplete-install": cmd.GetFlagDescription("autocomplete-install"),
	"remote":               cmd.GetFlagDescription("remote"),
	"report":               cmd.GetFlagDescription("report"),
	"yes":                  cmd.GetFlagDescription("yes"),
	"flush":                cmd.GetFlagDescription("flush"),
	"winget-update":        cmd.GetFlagDescription("winget-update"),
	"scan-health":          cmd.GetFlagDescription("scan-health"),
	"restore-health":       cmd.GetFlagDescription("restore-health"),
}
