#!/usr/bin/env bash

# https://termux.com/storage.html
if [[ ! -d "$HOME/storage" ]]
then log "Termux setup storage..."
  termux-setup-storage
fi
