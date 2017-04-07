#!/bin/sh

not_supported() {
  echo "Sorry, OS is not supported: ${UNAME}. Please compile manually from https://github.com/dnote-io/cli"
  exit 1
}

install() {
  set -eu

  UNAME=$(uname)

  if [ "$UNAME" != "Linux" -a "$UNAME" != "Darwin" ] ; then
    not_supported
  fi

  if [ "$UNAME" = "Darwin" ]; then
    OSX_ARCH=$(uname -m)
    if [ "${OSX_ARCH}" = "x86_64" ]; then
      PLATFORM="darwin_amd64"
    else
      not_supported
    fi
  elif [ "$UNAME" = "Linux" ]; then
    LINUX_ARCH=$(uname -m)
    if [ "${LINUX_ARCH}" = "x86_64" ]; then
      PLATFORM="linux_amd64"
    else
      not_supported
    fi
  fi

  LATEST=$(curl -s https://api.github.com/repos/dnote-io/cli/tags | grep -Eo '"name":.*[^\\]",'  | head -n 1 | sed 's/[," ]//g' | cut -d ':' -f 2)
  URL="https://github.com/dnote-io/cli/releases/download/$LATEST/dnote_$PLATFORM"
  DEST=${DEST:-/usr/local/bin/dnote}

  if [ -z $LATEST ]; then
    echo "Error fetching. Please try again."
    exit 1
  else
    echo "Download Dnote binary from curl https://github.com/dnote-io/cli/releases/download/$LATEST/dnote_$PLATFORM to $DEST"
    if curl -sL https://github.com/dnote-io/cli/releases/download/$LATEST/dnote_$PLATFORM -o $DEST; then
      chmod +x $DEST
      echo "Dnote installation was successful"
    else
      echo "Installation failed"
    fi
  fi
}

install

