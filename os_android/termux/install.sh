# case "$OS" in
#   android) apt install termux-tools -qq -y ;;
# esac

for p in $HOME/{.android.d,.termux}
do [[ -d "$p" ]] || mkdir -p "$p"
done

# https://termux.com/storage.html
if [[ ! -d "$HOME/storage" ]]
then log "Termux setup storage..."
  termux-setup-storage
fi
