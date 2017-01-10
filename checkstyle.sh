#!/usr/bin/env bash
set -e

errors=$(gofmt -d ./**/*go; golint ./...; misspell .;)
if [[ ${errors} ]]; then
    echo "Checkstyle errors found: ${errors}";
    exit 1;
fi
