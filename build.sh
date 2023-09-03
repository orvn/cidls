#!/usr/bin/env bash
# Build script to compile to multiple platforms and architectures

set -o errexit
set -o nounset
set -o pipefail

# Base name of the project
APP="cidls"

# version is the first argument passed to the script (in format n.n.n)
VERSION=$1

# Generate a build number as current date in the format YYYYMMDDhhmm
BUILD=$(date +"%Y%m%d%H%M")

# Setup the -ldflags option for go build, adding the version and build number
LDFLAGS="-X main.Version=${VERSION} -X main.Build=${BUILD}"

[ ! -d "build" ] && mkdir -p build

# Build for Darwin (MacOS) on both amd64 and arm64 architectures
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o build/${APP}_darwin_amd64
GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o build/${APP}_darwin_arm64

# Build for Linux on both amd64 and arm64 architectures
GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o build/${APP}_linux_amd64
GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o build/${APP}_linux_arm64

echo -e "\n 🎉 Build complete \n"
