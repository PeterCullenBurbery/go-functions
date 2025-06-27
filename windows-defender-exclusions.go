package gofunctions

import (
	"fmt"
	"os/exec"
	"os"
	"log"
)

// addDefenderExclusion adds a directory to Microsoft Defender's exclusion list using PowerShell.
// This prevents Defender from scanning or interfering with the specified path.
//
// Parameters:
//   path - the full file system path to exclude from Defender scans.
//
// Example:
//   addDefenderExclusion("C:\\my\\safe\\tools")
func addDefenderExclusion(path string) {
	// Build the PowerShell command to add the exclusion.
	// Note: path is quoted to support spaces or special characters.
	psCommand := fmt.Sprintf(`Add-MpPreference -ExclusionPath "%s"`, path)

	// Construct the exec.Command to run PowerShell with flags to avoid profile interference
	// and bypass execution policy restrictions.
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psCommand)

	// Attach the current process's standard output and error to the PowerShell command,
	// so any messages from PowerShell are shown in the terminal.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and handle success or failure.
	if err := cmd.Run(); err != nil {
		// Log a warning if the exclusion fails, including the error message.
		log.Printf("⚠️ Failed to exclude from Defender: %s\n↳ %v", path, err)
	} else {
		// Log a success message when the exclusion is successfully added.
		log.Printf("🛡️ Added Defender exclusion:\n↳ %s", path)
	}
}
