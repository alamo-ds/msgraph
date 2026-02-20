#!/bin/sh
# Adapted from the Deno installer: Copyright 2019 the Deno authors. All rights reserved. MIT license.
# Ref: https://github.com/denoland/deno_install

set -e

os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

if [ "$arch" == "aarch64" ]; then
    arch="arm64"
fi

if [ $# -eq 0 ]; then
    uri="https://github.com/alamo-ds/msgraph/releases/latest/download/msgraph_${os}_${arch}"
else
    uri="https://github.com/alamo-ds/msgraph/releases/download/${1}/msgraph_${os}_${arch}"
fi

msgraph_install="/usr/local"
bin_dir="${msgraph_install}/bin"
exe="${bin_dir}/msgraph"

if [ ! -d "${bin_dir}" ]; then
    mkdir -p "${bin_dir}"
fi

curl --silent --show-error --location --fail --location --output "${exe}" "${uri}"
chmod +x "${exe}"

echo "msgraph was installed successfully to ${exe}"
if command -v msgraph >/dev/null; then
    echo "Run 'msgraph --help' to get started"
fi
