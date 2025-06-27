package gofunctions

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// AddDefenderExclusion adds a directory to Microsoft Defender's exclusion list using PowerShell,
// but only if it's not already excluded.
func AddDefenderExclusion(path string) {
	// First, check if the path is already excluded
	checkCmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass",
		"-Command", "Get-MpPreference | Select-Object -ExpandProperty ExclusionPath")

	var checkOut bytes.Buffer
	checkCmd.Stdout = &checkOut
	checkCmd.Stderr = &checkOut

	if err := checkCmd.Run(); err != nil {
		log.Printf("⚠️ Could not check existing exclusions: %v", err)
	} else {
		existing := strings.Split(checkOut.String(), "\n")
		for _, line := range existing {
			// ✅ Trim before comparing
			if strings.EqualFold(strings.TrimSpace(line), path) {
				log.Printf("ℹ️ Defender exclusion already exists:\n↳ %s", path)
				return
			}
		}
	}

	// Add the exclusion if it doesn't exist
	psCommand := fmt.Sprintf(`Add-MpPreference -ExclusionPath "%s"`, path)
	addCmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psCommand)
	addCmd.Stdout = os.Stdout
	addCmd.Stderr = os.Stderr

	if err := addCmd.Run(); err != nil {
		log.Printf("⚠️ Failed to exclude from Defender: %s\n↳ %v", path, err)
	} else {
		log.Printf("🛡️ Added Defender exclusion:\n↳ %s", path)
	}
}
