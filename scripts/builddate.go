// Prints the current time in RFC3339 for Taskfile BUILD_DATE (works on Windows, Linux, macOS).
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Print(time.Now().Format(time.RFC3339))
}
