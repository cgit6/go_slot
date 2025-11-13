#!/usr/bin/env bash
set -e

go version

# 安裝 delve
go install github.com/go-delve/delve/cmd/dlv@latest

echo "Workspace ready"