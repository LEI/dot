#!/bin/sh

set -e

# USER="${USER:-$1}"

if [ -z "$USER" ]; then
  exit 1
fi

# groupadd $GROUP
# echo "%$GROUP ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers.d/$GROUP
# usermod -aG sudo $USER

# Create user: --password '' --groups sudo,...
useradd --create-home --shell /bin/bash --user-group "$@"
# echo "$USER:$USER" | chpasswd

# Allow user to execute any command without password
echo "$USER ALL=(ALL) NOPASSWD: ALL" >>"/etc/sudoers.d/$USER"

chmod 0440 "/etc/sudoers.d/$USER"
