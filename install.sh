#!/usr/bin/env bash

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}=== Building and Installing GoFetch ===${NC}\n"

# Verify Go installation
if ! command -v go &>/dev/null; then
  echo -e "${RED}Error: Go is not installed on your system.${NC}"
  exit 1
fi

# Ensure target binary directory exists
mkdir -p ~/.local/bin

# Build and install GoFetch
echo -e "${BLUE}--> Compiling GoFetch...${NC}"
go build -o gofetch .

echo -e "${BLUE}--> Moving binary to ~/.local/bin...${NC}"
mv gofetch ~/.local/bin

echo -e "${GREEN}GoFetch installed successfully in ~/.local/bin/gofetch!${NC}\n"

# PATH warning check
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
  echo -e "${RED}Warning: ~/.local/bin is not in your \$PATH.${NC}"
  echo -e "Add the following line to your shell configuration file:"
  echo -e "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi
