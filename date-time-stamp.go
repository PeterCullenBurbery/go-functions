package gofunctions

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// DateTimeStamp runs a temporary Java program that returns a formatted date-time string
// with time zone, ISO week/year/day, and day-of-year.
func DateTimeStamp() (string, error) {
	// Ensure javac is available
	if _, err := exec.LookPath("javac"); err != nil {
		return "", fmt.Errorf("❌ 'javac' not found in PATH. Please ensure JDK is installed")
	}
	// Ensure java is available
	if _, err := exec.LookPath("java"); err != nil {
		return "", fmt.Errorf("❌ 'java' not found in PATH. Please ensure JRE is installed")
	}

	// Create a temporary folder
	tempDir, err := os.MkdirTemp("", "date_time_stamp")
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	const javaFileName = "date_time_stamp.java"
	const className = "date_time_stamp"
	javaFilePath := filepath.Join(tempDir, javaFileName)

	// Java source code
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

	// Write the Java file
	if err := os.WriteFile(javaFilePath, []byte(javaCode), 0644); err != nil {
		return "", fmt.Errorf("❌ Failed to write Java file: %w", err)
	}

	// Compile it
	cmdCompile := exec.Command("javac", javaFileName)
	cmdCompile.Dir = tempDir
	if err := cmdCompile.Run(); err != nil {
		return "", fmt.Errorf("❌ Failed to compile Java file: %w", err)
	}

	// Run it and capture the output
	cmdRun := exec.Command("java", className)
	cmdRun.Dir = tempDir

	var out bytes.Buffer
	cmdRun.Stdout = &out
	cmdRun.Stderr = &out

	if err := cmdRun.Run(); err != nil {
		return "", fmt.Errorf("❌ Failed to run Java class: %w\nOutput:\n%s", err, out.String())
	}

	return out.String(), nil
}
