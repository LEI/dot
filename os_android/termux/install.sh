#!/usr/bin/env bash

# https://termux.com/storage.html
if [[ ! -d "$HOME/storage" ]]
then log "Setup storage symlink"
  termux-setup-storage
fi
