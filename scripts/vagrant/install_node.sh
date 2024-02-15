#!/usr/bin/env bash
# shellcheck disable=SC1090,SC1091
set -eux

VERSION=12.16.2
NVM_VERSION=v0.35.0

# Install nvm
wget -qO- https://raw.githubusercontent.com/nvm-sh/nvm/"$NVM_VERSION"/install.sh | bash
cat >> /home/vagrant/.bash_profile<< EOF
export NVM_DIR="\$([ -z "\${XDG_CONFIG_HOME-}" ] && printf %s "\${HOME}/.nvm" || printf %s "\${XDG_CONFIG_HOME}/nvm")"
[ -s "\$NVM_DIR/nvm.sh" ] && \. "\$NVM_DIR/nvm.sh" # This loads nvm
EOF
source /home/vagrant/.bash_profile

# Install a node and alias
nvm install --no-progress "$VERSION" 1>/dev/null
nvm alias default "$VERSION"
nvm use default
