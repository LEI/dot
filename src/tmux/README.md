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

### Usage

List shortcuts

    prefix + ?

Command prompt

    prefix + :

#### Sessions

Create a new sesion that can be named or detached

    tmux new [-s session-name] [-d]

    prefix + :new

Attach to a running session

    tmux attach -t target-session

    tmux a

List sessions

    tmux ls

    prefix + s

Rename the current session

    prefix + $

Detach from the session

    prefix + d

Kill a named session

    tmux kill-session -t target-session

#### Windows

Create a new window

    prefix + c

Rename the current window

    prefix + ,

Next window

    prefix + n

Previous window

    prefix + p

Kill the current window

    prefix + &

#### Panes

Create an horizontal split

    prefix + "

Create a vertical split

    prefix + %

Toggle pane zoom

    prefix + z

Kill pane

    prefix + x

#### Copy mode

    prefix + [
