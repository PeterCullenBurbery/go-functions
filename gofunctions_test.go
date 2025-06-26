package gofunctions

import (
	"fmt"
	"testing"
)

func TestFullWorkflow(t *testing.T) {
	// Test SayHello
	fmt.Println("➡️ Calling SayHello:")
	SayHello("Peter")

	// Test MultiplyBy2718
	result := MultiplyBy2718(3)
	fmt.Printf("➡️ MultiplyBy2718(3) = %d\n", result)

	// Test DateTimeStamp
	timestamp, err := DateTimeStamp()
	if err != nil {
		t.Errorf("❌ DateTimeStamp failed: %v", err)
		return
	}
	fmt.Println("➡️ DateTimeStamp output:")
	fmt.Print(timestamp)
}