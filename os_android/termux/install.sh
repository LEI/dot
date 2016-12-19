for f in $1/lib/{directory,utils}.bash; do source "$f"; done

directory present $HOME/{.android.d,.termux}

# https://termux.com/storage.html
if [[ ! -d "$HOME/storage" ]]
then log "Termux setup storage..."
  termux-setup-storage
fi
