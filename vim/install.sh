source $1/lib/{directory,lineinfile}.bash

# if !filereadable(expand("~/.vim/init.vim"))
lineinfile present "$HOME/.vimrc" "source ~/.vim/init.vim"

directory present $HOME/.vim/{plugin,settings}

if hash nvim 2>/dev/null
then directory present $HOME/.config
  [[ -d "$HOME/.vim" ]] && [[ ! -e "$HOME/.config/nvim" ]] \
    && ln -s $HOME/{.vim,.config/nvim}
fi
