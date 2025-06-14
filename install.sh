#!/bin/bash
set -e

echo "Building docker-ai..."
go build -o docker-ai ./cmd/docker-ai

echo "Installing to /usr/local/bin (may require sudo)..."
sudo cp docker-ai /usr/local/bin/
echo "Done. Run 'docker-ai --help' to get started." 