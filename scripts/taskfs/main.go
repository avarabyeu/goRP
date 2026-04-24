// Small helpers for Taskfile: mkdir, remove tree, test directory exists (works on Windows without rm/mkdir/test).
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: taskfs <mkdir|remove|exists|nonempty> <path>")
		os.Exit(2)
	}
	op, path := os.Args[1], os.Args[2]
	switch op {
	case "mkdir":
		if err := os.MkdirAll(path, 0o755); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "remove":
		if err := os.RemoveAll(path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "exists":
		st, err := os.Stat(path)
		if err != nil || !st.IsDir() {
			os.Exit(1)
		}
	case "nonempty":
		st, err := os.Stat(path)
		if err != nil || st.Size() == 0 {
			fmt.Fprintf(os.Stderr, "taskfs nonempty: missing or empty file: %s\n", path)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "unknown op:", op)
		os.Exit(2)
	}
}
