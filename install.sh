#!/bin/sh
#
# This script downloads the latest Dnote release from github
# into /usr/bin/local.
#

set -eu

not_supported() {
  echo "OS not supported: ${UNAME}"
  echo "Please compile manually from https://github.com/dnote-io/cli"
  exit 1
}

install() {
  UNAME=$(uname)

  if [ "$UNAME" != "Linux" -a "$UNAME" != "Darwin" -a "$UNAME" != "OpenBSD" ] ; then
    not_supported
  fi

  if [ "$UNAME" = "Darwin" ]; then
    OSX_ARCH=$(uname -m)
    if [ "${OSX_ARCH}" = "x86_64" ]; then
      PLATFORM="darwin-amd64"
    else
      not_supported
    fi
  elif [ "$UNAME" = "Linux" ]; then
    LINUX_ARCH=$(uname -m)
    if [ "${LINUX_ARCH}" = "x86_64" ]; then
      PLATFORM="linux-amd64"
    elif [ "${LINUX_ARCH}" = "i686" ]; then
      PLATFORM="linux-386"
    else
      not_supported
    fi
  elif [ "$UNAME" = "OpenBSD" ]; then
    OPENBSD_ARCH=$(uname -m)
    if [ "${OPENBSD_ARCH}" = "x86_64" ]; then
      PLATFORM="openbsd-amd64"
    elif [ "${OPENBSD_ARCH}" = "i686" ]; then
      PLATFORM="openbsd-386"
    else
      not_supported
    fi
  fi

  LATEST=$(curl -s https://api.github.com/repos/dnote-io/cli/tags | grep -Eo '"name":[ ]*"v[0-9]*\.[0-9]*\.[0-9]*",' | head -n 1 | sed 's/[," ]//g' | cut -d ':' -f 2)
  URL="https://github.com/dnote-io/cli/releases/download/$LATEST/dnote-$PLATFORM"
  DEST=${DEST:-/usr/local/bin/dnote}

  if [ -z $LATEST ]; then
    echo "Error fetching latest version. Please try again."
    exit 1
  fi

  echo "Downloading Dnote binary from $URL to $DEST"
  if curl -L --progress-bar $URL -o $DEST; then
    chmod +x $DEST
    echo "Successfully installed Dnote"
  else
    echo "Installation failed. You might need elevated permission."
  fi
}

install
