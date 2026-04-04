/// Centralized Functions
/// powershell execution, clipboard, report export

package cmd

import (
	"fmt"
	"html"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

func runPowershellReturnOutput(command string) (string, error) {
	// Force Powershell UTF-8 output
	psCmd := fmt.Sprintf("[Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8; %s", command)
	if remoteTarget != "" {
		psCmd = fmt.Sprintf(`Invoke-Command -ComputerName %s -ScriptBlock { [Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8; %s }`, remoteTarget, command)
	}

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psCmd)
	cmd.SysProcAttr = hiddenSysProcAttr()

	// CombinedOutput []byte UTF-8
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))

	if output == "" {
		if err != nil {
			return fmt.Sprintf("⚠️ Error executing: %v", err), nil
		}
		return "No output (possibly no data found).\n", nil
	}
	return output, nil
}

func copyToClipboard(content string) {
	if content == "" {
		fmt.Println("Nothing to copy.")
		return
	}

	// CLI-only auto-copy
	if !guiMode {
		if err := clipboard.WriteAll(content); err != nil {
			fmt.Println("Error copying to clipboard:", err)
		} else {
			fmt.Println("✅ Output copied to clipboard!")
		}
	}
}

func copyToClipboardForce(content string) error {
	if content == "" {
		return fmt.Errorf("nothing to copy")
	}
	return clipboard.WriteAll(content)
}

func exportReport(content, format, path string) error {
	if path == "" {
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("Could not determine current user: %v", err)
		}
		path = filepath.Join(usr.HomeDir, "Desktop",
			fmt.Sprintf("report_%s.%s", time.Now().Format("2006-01-02_15-04-05"), format))
		return writeReportFile(content, format, path)
	}

	// If path already ends with .html or .md -> write directly
	if strings.HasSuffix(strings.ToLower(path), "."+format) {
		return writeReportFile(content, format, path)
	}

	// Otherwise, assume it's a directory
	filename := filepath.Join(path,
		fmt.Sprintf("report_%s.%s", time.Now().Format("2006-01-02_15-04-05"), format))
	return writeReportFile(content, format, filename)
}

func writeReportFile(content, format, filename string) error {

	switch format {
	case "md":
		return os.WriteFile(filename, []byte("```markdown\n"+content+"\n```"), 0644)
	case "html":
		htmlOut := "<!DOCTYPE html>\n<html lang=\"de\">\n<head>\n<meta charset=\"UTF-8\">\n<title>ProbeDesk Report</title>\n<style>body{font-family:monospace;background:#1e1e1e;color:#d4d4d4;padding:20px;}pre{white-space:pre-wrap;word-wrap:break-word;}</style>\n</head>\n<body><pre>" + html.EscapeString(content) + "</pre></body></html>"
		return os.WriteFile(filename, []byte(htmlOut), 0644)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func GetFlagDescription(name string) string {
	f := rootCmd.Flags().Lookup(name)
	if f != nil {
		return f.Usage
	}
	return ""
}

// check for admin privileges
func IsAdmin() bool {
	cmd := exec.Command("net", "session")
	cmd.SysProcAttr = hiddenSysProcAttr()
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// relaunch app as admin after prompt
func RelaunchAsAdmin() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command("powershell", "-Command",
		"Start-Process", "'"+exe+"'", "-Verb", "runAs")
	cmd.SysProcAttr = hiddenSysProcAttr()
	return cmd.Start()
}
