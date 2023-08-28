#!/bin/bash

# Config
USERNAME="orvn"
REPO_NAME="cidls"
VERSION="v1.0.0"
APP="cidls"

BASE_URL="https://github.com/${USERNAME}/${REPO_NAME}/releases/download/${VERSION}"

OS="$(uname)"
ARCH="$(uname -m)"

if [ "$OS" == "Darwin" ]; then
    if [ "$ARCH" == "x86_64" ]; then
        BINARY_URL="${BASE_URL}/${APP}_darwin_amd64"
    elif [ "$ARCH" == "arm64" ]; then
        BINARY_URL="${BASE_URL}/${APP}_darwin_arm64"
    else
        echo "Unsupported architecture: $ARCH"
        exit 1
    fi
elif [ "$OS" == "Linux" ]; then
    if [ "$ARCH" == "x86_64" ]; then
        BINARY_URL="${BASE_URL}/${APP}_linux_amd64"
    elif [[ "$ARCH" == "arm"* ]]; then
        BINARY_URL="${BASE_URL}/${APP}_linux_arm64"
    else
        echo "Unsupported architecture: $ARCH"
        exit 1
    fi
else
    echo "Unsupported OS: $OS"
    exit 1
fi

curl -L $BINARY_URL -o /tmp/$APP
chmod +x /tmp/$APP
sudo mv /tmp/$APP /usr/local/bin/$APP
