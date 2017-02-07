# dot-vim

## Requirements

- curl
- [git](https://git-scm.com)
- [vim](https://github.com/vim/vim) or [neovim](https://neovim.io)

## Manual installation

    mkdir -p "$HOME/.vim"
    ln -isv "$DOT/.vim/*" "$HOME/.vim"
    echo 'source ~/.vim/init.vim' >> "$HOME/.vim/vimrc"
