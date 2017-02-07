# dot-bash

## Requirements

- bash

## Manual installation

    mkdir -p "$HOME/.bash.d"
    ln -isv "$DOT/.{bash_*,bash.d/*,inputrc}" "$HOME"
    echo 'source ~/.bash_profile' >> "$HOME/.bashrc"
