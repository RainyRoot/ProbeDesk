package gui

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	//"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/RainyRoot/ProbeDesk/cmd"
	"github.com/RainyRoot/ProbeDesk/gui/components"
)

func Run() {
	a := app.New()
	w := a.NewWindow("ProbeDesk GUI")
	w.Resize(fyne.NewSize(900, 600))
	cmd.SetGuiMode(true)

	resourceIcon, _ := fyne.LoadResourceFromPath("icon.ico")
	w.SetIcon(resourceIcon)

	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapWord
	output.SetPlaceHolder("Output will appear here...")
	output.TextStyle.Monospace = true

	status := widget.NewLabel("Ready")

	// // Check adminrights
	// if !cmd.IsAdmin() {
	// 	dialog.NewConfirm(
	// 		"Admin rights are recommend",
	// 		"Some actions (e.g. Flush DNS, DISM, Windows Update) require admin privileges.\n\nDo you like to restart ProbeDesk as admin?",
	// 		func(ok bool) {
	// 			if ok {
	// 				cmd.RelaunchAsAdmin()
	// 				a.Quit()
	// 			}
	// 		},
	// 		w,
	// 	).Show()
	// }

	// Checkboxes
	systemCheck := widget.NewCheck("System Info", nil)
	ipCheck := widget.NewCheck("IP Config", nil)
	usbCheck := widget.NewCheck("USB Devices", nil)
	servicesCheck := widget.NewCheck("Running Services", nil)
	usersCheck := widget.NewCheck("Local Users", nil)
	traceRouteRequest := widget.NewCheck("Trace Route", nil)
	vpnCheck := widget.NewCheck("VPN Connections", nil)
	productsCheck := widget.NewCheck("Installed Products", nil)
	netuseCheck := widget.NewCheck("Mapped Drives", nil)
	checkHealthCheck := widget.NewCheck("Check Health", nil)
	scanHealth := widget.NewCheck("Scan System Health", nil)
	flushDns := widget.NewCheck("Flush DNS", nil)
	winGetUpdate := widget.NewCheck("Check for Windows updates", nil)
	restoreHealth := widget.NewCheck("Restore broken Windows image", nil)

	// Toggle all
	actionChecks := []*widget.Check{
		systemCheck, ipCheck, usbCheck, servicesCheck, usersCheck,
		vpnCheck, productsCheck, netuseCheck,
		traceRouteRequest,
	}

	// Entry fields
	tracerouteEntry := widget.NewEntry()
	tracerouteEntry.SetPlaceHolder("Optional: IP/DNS")

	remoteEntry := widget.NewEntry()
	remoteEntry.SetPlaceHolder("Optional: Remote host")

	formatSelect := widget.NewSelect([]string{"html", "md"}, func(string) {})
	formatSelect.SetSelected("html")

	var mu sync.Mutex

	// Buttons
	runButton := components.NewRunButton(output, status,
		systemCheck, ipCheck, usbCheck, servicesCheck, usersCheck, traceRouteRequest,
		vpnCheck, productsCheck, netuseCheck, checkHealthCheck,
		flushDns, winGetUpdate, scanHealth, restoreHealth,
		remoteEntry, tracerouteEntry, formatSelect, &mu)

	copyButton := components.NewCopyButton(output, status)
	exportButton := components.NewExportButton(output, status, w, formatSelect)

	// Info Labels
	healthInfo := widget.NewLabel("🛈 Some actions (like Health Checks or Windows Restore) can take several minutes.")
	healthInfo.Wrapping = fyne.TextWrapWord
	healthInfo.Alignment = fyne.TextAlignLeading

	flushDnsInfo := widget.NewLabel("⚠️ Flushing DNS may disrupt static IP configurations temporarily.")
	flushDnsInfo.Wrapping = fyne.TextWrapWord
	flushDnsInfo.Alignment = fyne.TextAlignLeading

	toggleAllButton := widget.NewButton("Toggle All", func() {
		anyUnchecked := false
		for _, c := range actionChecks {
			if !c.Checked {
				anyUnchecked = true
				break
			}
		}

		for _, c := range actionChecks {
			c.SetChecked(anyUnchecked)
		}
	})

	// Layout
	scrollContent := container.NewVBox(
		widget.NewLabelWithStyle("Select actions to run:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		toggleAllButton,
		systemCheck, ipCheck, usbCheck, servicesCheck, usersCheck,
		vpnCheck, productsCheck, netuseCheck,
		traceRouteRequest, tracerouteEntry,
		widget.NewSeparator(),

		widget.NewLabelWithStyle("System Maintenance & Health Checks", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		checkHealthCheck, flushDns, winGetUpdate, scanHealth, restoreHealth,
		healthInfo,
		flushDnsInfo,

		widget.NewSeparator(),
		widget.NewLabelWithStyle("Remote Target:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		remoteEntry,
		widget.NewLabelWithStyle("Report Format:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		formatSelect,
	)

	scrollOptions := container.NewVScroll(scrollContent)
	scrollOptions.SetMinSize(fyne.NewSize(250, 0))

	buttonBar := container.NewVBox(
		widget.NewSeparator(),
		runButton,
		copyButton,
		exportButton,
	)

	// combine scrollarea + fix Buttons vertical
	leftPanel := container.NewBorder(nil, buttonBar, nil, nil, scrollOptions)

	// Output in Card
	outputCard := widget.NewCard("Output", "", output)

	// Main
	main := container.NewHSplit(
		leftPanel,
		container.NewBorder(nil, status, nil, nil, outputCard),
	)
	main.Offset = 0.3

	w.SetContent(main)
	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}
