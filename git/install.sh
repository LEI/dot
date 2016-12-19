source $1/lib/template.bash

template "$1/git/.gitconfig.local.tpl" "$HOME/.gitconfig.local" \
  "GIT_AUTHOR_NAME" "GIT_AUTHOR_USERNAME" "GIT_AUTHOR_EMAIL"
