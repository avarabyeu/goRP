// Runs the yq→jq Docker pipeline that normalises the upstream ReportPortal spec into
// openapi-modified.json without requiring a POSIX shell (works on Windows, Linux, macOS).
//
// Usage: go run ./scripts/preprocess <workdir> <spec-path> <output-path>
//
//	workdir     – absolute path mounted as /workdir inside both containers
//	spec-path   – path to the input spec, relative to workdir (passed to yq inside the container)
//	output-path – absolute path where the processed JSON is written on the host
package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: preprocess <workdir> <spec-path> <output-path>")
		os.Exit(2)
	}
	workdir, specPath, outputPath := os.Args[1], os.Args[2], os.Args[3]

	absWorkdir, err := filepath.Abs(workdir)
	if err != nil {
		fatal(fmt.Errorf("resolve workdir: %w", err))
	}

	absOutput, err := filepath.Abs(outputPath)
	if err != nil {
		fatal(fmt.Errorf("resolve output path: %w", err))
	}

	fmt.Fprintf(os.Stderr, "preprocess: %s → %s\n", specPath, absOutput)

	// yq: normalize YAML/JSON spec to JSON, using /workdir as the Docker working directory.
	yqCmd := exec.Command("docker", //nolint:gosec
		"run", "--rm",
		"-v", absWorkdir+":/workdir",
		"-w", "/workdir",
		"mikefarah/yq:4",
		"-o=json", ".", specPath,
	)
	yqCmd.Stderr = os.Stderr

	// jq: apply the preprocessing filter, reading JSON from yq's stdout.
	jqCmd := exec.Command("docker", //nolint:gosec
		"run", "--rm", "-i",
		"-v", absWorkdir+":/workdir",
		"ghcr.io/jqlang/jq",
		"-f", "/workdir/scripts/openapi-preprocess.jq",
	)
	jqCmd.Stderr = os.Stderr

	// Wire yq stdout → jq stdin via an in-process pipe.
	pr, pw := io.Pipe()
	yqCmd.Stdout = pw
	jqCmd.Stdin = pr

	out, err := os.Create(absOutput)
	if err != nil {
		fatal(fmt.Errorf("create output file: %w", err))
	}
	jqCmd.Stdout = out

	if err := yqCmd.Start(); err != nil {
		fatal(fmt.Errorf("start yq: %w", err))
	}
	if err := jqCmd.Start(); err != nil {
		fatal(fmt.Errorf("start jq: %w", err))
	}

	// Close the write end of the pipe once yq exits so jq receives EOF.
	yqErrCh := make(chan error, 1)
	go func() {
		yqErr := yqCmd.Wait()
		pw.CloseWithError(yqErr)
		yqErrCh <- yqErr
	}()

	jqErr := jqCmd.Wait()
	// If jq exited early (e.g. filter error), no one drains the pipe read-end and yq
	// blocks writing into a full buffer — yqCmd.Wait() never returns. Closing pr here
	// delivers a broken-pipe error to yq so it exits promptly.
	pr.CloseWithError(jqErr)
	yqErr := <-yqErrCh

	if yqErr != nil {
		fatal(fmt.Errorf("yq: %w", yqErr))
	}
	if jqErr != nil {
		fatal(fmt.Errorf("jq: %w", jqErr))
	}

	if err := out.Close(); err != nil {
		fatal(fmt.Errorf("close output file: %w", err))
	}

	fmt.Fprintf(os.Stderr, "preprocess: written %s\n", absOutput)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "preprocess:", err)
	os.Exit(1)
}
