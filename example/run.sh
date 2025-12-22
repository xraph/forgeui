#!/bin/bash
# Run ForgeUI example from the correct directory
cd "$(dirname "$0")"
echo "Running from: $(pwd)"
go run *.go

