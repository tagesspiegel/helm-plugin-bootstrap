#!/bin/sh -e

# convert architecture of the target system to a compatible GOARCH value.
# Otherwise failes to download of the plugin from github, because the provided
# architecture by `uname -m` is not part of the github release.
arch=""
case $(uname -m) in
  x86_64)
    arch="amd64"
    ;;
  armv6*)
    arch="armv6"
    ;;
  # match every arm processor version like armv7h, armv7l and so on.
  armv7*)
    arch="armv7"
    ;;
  aarch64 | arm64)
    arch="arm64"
    ;;
  *)
    echo "Failed to detect target architecture"
    exit 1
    ;;
esac

# detect the operating system of the target system.
os=""
binExtension=""
case "$(uname)" in
    Darwin)
        os="Darwin"
        binExtension=""
        ;;
    Linux)
        os="Linux"
        binExtension=""
        ;;
    Windows)
        os="Windows"
        binExtension=".exe"
        ;;
    *)
        echo "Failed to detect target operating system"
        exit 1
        ;;
esac

echo "Installing Helm boostrap plugin for ${os} ${arch}"
url="https://github.com/tagesspiegel/helm-plugin-bootstrap/releases/latest/download/helm-plugin-bootstrap_${os}_${arch}${binExtension}"
echo "Downloading from ${url}"
binDir="$(dirname "$0")/../bin"
echo "Storing binary in: ${binDir}"

mkdir -p "${binDir}"
# Download with curl if possible.
if [ -x "$(which curl 2>/dev/null)" ]; then
    curl -sSL "${url}" -o "${binDir}/bootstrap${binExtension}"
else
    wget -q "${url}" -O "${binDir}/bootstrap${binExtension}"
fi

chmod a+x "${binDir}/bootstrap${binExtension}"