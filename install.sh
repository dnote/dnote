#!/bin/sh
#
# This script downloads the latest Dnote release from github
# into /usr/bin/local.
#

set -eu

not_supported() {
  echo "OS not supported: ${UNAME}"
  echo "Please compile manually from https://github.com/dnote/cli"
  exit 1
}

get_platform() {
  UNAME=$(uname)

  if [ "$UNAME" != "Linux" -a "$UNAME" != "Darwin" -a "$UNAME" != "OpenBSD" ] ; then
    not_supported
  fi

  if [ "$UNAME" = "Darwin" ]; then
    OSX_ARCH=$(uname -m)
    if [ "${OSX_ARCH}" = "x86_64" ]; then
      platform="darwin_amd64"
    else
      not_supported
    fi
  elif [ "$UNAME" = "Linux" ]; then
    LINUX_ARCH=$(uname -m)
    if [ "${LINUX_ARCH}" = "x86_64" ]; then
      platform="linux_amd64"
    elif [ "${LINUX_ARCH}" = "i686" ]; then
      platform="linux_386"
    else
      not_supported
    fi
  elif [ "$UNAME" = "OpenBSD" ]; then
    OPENBSD_ARCH=$(uname -m)
    if [ "${OPENBSD_ARCH}" = "x86_64" ]; then
      platform="openbsd_amd64"
    elif [ "${OPENBSD_ARCH}" = "i686" ]; then
      platform="openbsd_386"
    else
      not_supported
    fi
  fi

  echo $platform
}

get_version() {
  LATEST=$(curl -s https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/tags | grep -Eo '"name":[ ]*"v[0-9]*\.[0-9]*\.[0-9]*",' | head -n 1 | sed 's/[," ]//g' | cut -d ':' -f 2)
  if [ -z $LATEST ]; then
    echo "Error fetching latest version. Please try again."
    exit 1
  fi

  # remove the preceding 'v'
  echo ${LATEST#v}
}

execute() {
  echo "downloading Dnote v${LATEST}..."
  echo ${URL}
  if curl -L --progress-bar $URL -o "${TMPDIR}/${TARBALL}"; then
    (cd "${TMPDIR}" && tar -xzf "${TARBALL}")

    install -d "${BINDIR}"
    install "${TMPDIR}/${BINARY}" "${BINDIR}/"

    echo "Successfully installed Dnote"
  else
    echo "Installation failed. You might need elevated permission."
    exit 1
  fi
}

REPO_OWNER=dnote
REPO_NAME=cli
PLATFORM=$(get_platform)
LATEST=$(get_version)
TARBALL="dnote_${LATEST}_${PLATFORM}.tar.gz"
URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/v${LATEST}/${TARBALL}"
TMPDIR="$(mktemp -d)"
BINDIR=${BINDIR:-/usr/local/bin}
BINARY=dnote

execute
