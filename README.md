# dot-bash

## Requirements

- bash

## Manual installation

    mkdir -p "$HOME/.bash.d"
    ln -isv "$DOT/.{bash_*,bash.d/*,inputrc}" "$HOME"
    echo '[[ -n "$PS1" ]] && [[ -f ~/.bash_profile ]] && source ~/.bash_profile' >> "$HOME/.bashrc"
