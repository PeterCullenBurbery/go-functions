package gofunctions

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DateTimeStamp returns a timestamp string formatted via a temporary Java program.
// It supports optional overrides for javac/java paths.
func DateTimeStamp(args ...string) (string, error) {
	var javacCmd, javaCmd string

	switch len(args) {
	case 0:
		// Default: look in PATH
		var err error
		javacCmd, err = exec.LookPath("javac")
		if err != nil {
			return "", fmt.Errorf("❌ 'javac' not found in PATH. Please ensure JDK is installed")
		}
		javaCmd, err = exec.LookPath("java")
		if err != nil {
			return "", fmt.Errorf("❌ 'java' not found in PATH. Please ensure JRE is installed")
		}
	case 2:
		javacCmd = args[0]
		javaCmd = args[1]
	default:
		return "", fmt.Errorf("❌ DateTimeStamp() expects 0 or 2 arguments (javacPath, javaPath)")
	}

	tempDir, err := os.MkdirTemp("", "date_time_stamp")
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	const javaFileName = "date_time_stamp.java"
	const className = "date_time_stamp"
	javaFilePath := filepath.Join(tempDir, javaFileName)

	javaCode := `import java.time.*;
import java.time.format.DateTimeFormatter;
import java.time.temporal.WeekFields;

public class date_time_stamp {
    public static void main(String[] args) {
        ZonedDateTime now = ZonedDateTime.now();
        ZoneId tz = now.getZone();
        String date_part = now.format(DateTimeFormatter.ofPattern("yyyy-0MM-0dd"));
        String time_part = now.format(DateTimeFormatter.ofPattern("0HH.0mm.0ss.nnnnnnn"));
        WeekFields wf = WeekFields.ISO;
        int week = now.get(wf.weekOfWeekBasedYear());
        int weekday = now.get(wf.dayOfWeek());
        int iso_year = now.get(wf.weekBasedYear());
        int day_of_year = now.getDayOfYear();
        String output = String.format(
            "%s %s %04d-W%03d-%03d %04d-%03d",
            date_part, time_part, iso_year, week, weekday, now.getYear(), day_of_year
        );
        output = output.replace(time_part, time_part + " " + tz);
        System.out.println(output);
    }
}`

	if err := os.WriteFile(javaFilePath, []byte(javaCode), 0644); err != nil {
		return "", fmt.Errorf("❌ Failed to write Java file: %w", err)
	}

	// Compile
	cmdCompile := exec.Command(javacCmd, javaFileName)
	cmdCompile.Dir = tempDir
	if err := cmdCompile.Run(); err != nil {
		return "", fmt.Errorf("❌ Failed to compile Java file: %w", err)
	}

	// Run
	cmdRun := exec.Command(javaCmd, className)
	cmdRun.Dir = tempDir
	var out bytes.Buffer
	cmdRun.Stdout = &out
	cmdRun.Stderr = &out
	if err := cmdRun.Run(); err != nil {
		return "", fmt.Errorf("❌ Failed to run Java class: %w\nOutput:\n%s", err, out.String())
	}

	return out.String(), nil
}

// SafeTimeStamp optionally replaces "/" with " slash " if mode == 1.
func SafeTimeStamp(timestamp string, mode int) string {
	if mode == 1 {
		return strings.ReplaceAll(timestamp, "/", " slash ")
	}
	return timestamp
}