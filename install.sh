#!/bin/sh

set -e

platform="$(uname -s)"
processor="$(uname -p)"

if [[ $platform != "Darwin" || $processor != "arm" ]]; then
  echo "\x1b[31mYou're not supposed to run this on an non-M1 machine!\x1b[39m"
  exit 1
fi

download_url="https://github.com/ryanccn/PolyM1/releases/latest/download/polym1"
download_dir="${HOME}/.polym1"
download_path="${download_dir}/polym1"

if [ -d "$download_dir" ]; then rm -rf "$download_dir"; fi
mkdir "$download_dir"

echo "downloading..."

curl --fail --location --progress-bar --output "$download_path" "$download_url"
chmod +x "$download_path"

$download_path install
