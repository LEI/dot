# Termux

# https://termux.com/linux.html

PS1='\w\$ '

# Git work tree path
TERMUX_CFG="$HOME/storage/shared/termux-config"
# TERMUX_GIT="$HOME/termux-config.git"

termux-git() {
  git --git-dir="$TERMUX_CFG/.git" --work-tree="$TERMUX_CFG" "$@"
}

termux-stow() {
  stow -d "$TERMUX_CFG" -t "$HOME" "$@"
}

termux-fix-tpm() {
  termux-fix-shebang ~/.tmux/plugged/**/*.{sh,tmux}
}
