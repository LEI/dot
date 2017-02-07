# dot-tmux

## Requirements

- [git](https://git-scm.com/)
- [tmux](https://tmux.github.io/)
- [Solarized](http://ethanschoonover.com/solarized) terminal colors

## Manual installation

    mkdir -p "$HOME/.tmux"
    ln -isv "$DOT/.tmux/*" "$HOME/.tmux"
    echo 'source-file $HOME/.tmux/tmux.conf' >> "$HOME/.tmux.conf"

## Tmux Plugin Manager

- [tpm](https://github.com/tmux-plugins/tpm)
- [tmux-sensible](https://github.com/tmux-plugins/tmux-sensible)
- [tmux-resurrect](https://github.com/tmux-plugins/tmux-resurrect)

## Resources

- [Example .tmux.conf](https://github.com/tmux/tmux/blob/master/example_tmux.conf)
