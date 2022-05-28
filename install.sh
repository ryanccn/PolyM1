#!/bin/sh

set -e

download_url="https://github.com/PolyM1/PolyM1/releases/latest/download/polym1"
download_dir="$(mktemp -d)"
download_path="${download_dir}/polym1"

echo "downloading..."

curl --fail --location --progress-bar --output "$download_path" "$download_url"
chmod +x "$download_path"

$download_path install
