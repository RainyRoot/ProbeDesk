//go:build !windows

package components

func isAdmin() bool {
	return true
}

func relaunchAsAdmin() {
	// no-op on non-Windows
}
