#!/usr/bin/env bash
set -e

## Tip: Use 'gofmt -w ./**/*go' to reformat all code

errors=$(gofmt -d ./**/*go; golint ./...; misspell .;)
if [[ ${errors} ]]; then
    echo "Checkstyle errors found: ${errors}";
    exit 1;
fi
