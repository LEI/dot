#!/bin/sh

set -e

echo "SETUP USER: $USER"

USER="${USER:-$1}"

# Add /usr/local/bin to PATH while sudoing
#sed -e 's#Defaults[[:blank:]]+secure_path = /sbin:/bin:/usr/sbin:/usr/bin#Defaults[[:blank:]]+secure_path = /sbin:/bin:/usr/sbin:/usr/bin:/usr/local/bin#' /etc/sudoers
# Deb: Defaults        secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
# sed -e '/secure_path/ s[=.*[&:/usr/local/bin[' /etc/sudoers
# sed -r -e '/^\s*Defaults\s+secure_path/ s[=(.*)[=\1:/usr/local/bin[' /etc/sudoers

# groupadd $GROUP
# echo "%$GROUP ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers.d/$GROUP
# usermod -aG sudo $USER

# # Allow members of group to execute any command without password
echo "%$USER ALL=(ALL) NOPASSWD: ALL" >> "/etc/sudoers.d/$USER"

#shopt -s extglob
chmod 0440 /etc/sudoers.d/*

# Create user: --password '' --groups sudo,...
useradd --create-home --shell /bin/bash --user-group "$@"
# echo "$USER:$USER" | chpasswd
