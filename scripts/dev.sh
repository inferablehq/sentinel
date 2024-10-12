#!/bin/bash

# Install reflex if not already installed
if ! command -v reflex &> /dev/null; then
    echo "Installing reflex..."
    go install github.com/cespare/reflex@latest

    # Add the Go bin directory to PATH
    export PATH=$PATH:$(go env GOPATH)/bin

    # Verify installation
    which reflex || echo "Reflex installation failed or not in PATH"
fi

# Run all Go files in the current directory and watch for changes
reflex -r '\.go$' -s -- go run .
