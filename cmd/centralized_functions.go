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
	if err := clipboard.WriteAll(content); err != nil {
		fmt.Println("Error copying to clipboard:", err)
	} else {
		fmt.Println("✅ Output copied to clipboard!")
	}
}

func exportReport(content, format, path string) error {
	if path == "" {
		// Desktop of the current user
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("Couldnt determine current user: %v", err)
		}
		path = filepath.Join(usr.HomeDir, "Desktop")
	}

	// Create filename
	filename := filepath.Join(path, fmt.Sprintf("report_%s.%s", time.Now().Format("2006-01-02_15-04-05"), format))

	switch format {
	case "md":
		return os.WriteFile(filename, []byte("```markdown\n"+content+"\n```"), 0644)
	case "html":
		htmlOut := "<html><body><pre>" + html.EscapeString(content) + "</pre></body></html>"
		return os.WriteFile(filename, []byte(htmlOut), 0644)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
