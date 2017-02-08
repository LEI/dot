# dot-tmux

## Requirements

- [git](https://git-scm.com/)
- [tmux](https://tmux.github.io/)
- [Solarized](http://ethanschoonover.com/solarized) terminal colors

## Manual installation

    mkdir -p "$HOME/.tmux"
    ln -isv "$DOT/.tmux/*" "$HOME/.tmux"
    echo 'source-file $HOME/.tmux/tmux.conf' >> "$HOME/.tmux.conf"

## Resources

- [Tmux Cheat Sheet](http://tmuxcheatsheet.com)
- [Example .tmux.conf](https://github.com/tmux/tmux/blob/master/example_tmux.conf)

### Tmux Plugin Manager

- [tpm](https://github.com/tmux-plugins/tpm)
- [tmux-sensible](https://github.com/tmux-plugins/tmux-sensible)
- [tmux-resurrect](https://github.com/tmux-plugins/tmux-resurrect)

### Usage

The default bind prefix is `C-b` (Control+b)

This configuration sets the prefix to `C-a` when running tmux locally.

List all key bindings

    prefix + ?

Enter the tmux command prompt

    prefix + :

#### Sessions

Create a new sesion that can be named or detached

    tmux new [-s <session-name>] [-d]

    prefix + :new

Attach to a running session

    tmux attach -t <target-session>

    tmux a

List sessions

    tmux ls

    prefix + s

Rename the current session

    prefix + $

Switch to the next session

    prefix + )

Switch to the previous session

    prefix + (

Detach from the session

    prefix + d

Choose a client to detach

    prefix + D

Kill all sessions except one

    tmux kill-session -a -t <target-session>

#### Windows

Create a new window

    prefix + c

Rename the current window

    prefix + ,

Select a window by its index

    prefix + [0-9]


Change to the next window

    prefix + n

Change to the previous window

    prefix + p

Kill the current window

    prefix + &

#### Panes

Split the current pane horizontally

    prefix + "

Split the current pane vertically

    prefix + %

Navigate panes

    prefix + Up # Down, Left, Right

    prefix + k # j, h, l

Resize panes

    prefix + K # J, H, L

Arrange the current window in the next preset layout

    prefix + Space

Toggle pane zoom

    prefix + z

Kill the active pane

    prefix + x

#### Copy mode

Enter copy mode

    prefix + [

Page the last copied buffer of text

    prefix + ]

Enter copy mode and scroll one page up

    prefix + PageUp

Search forward

    prefix + /

Quit copy mode

    q
