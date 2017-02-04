" CtrlP
" https://github.com/skwp/dotfiles/blob/master/vim/settings/ctrlp.vim

if executable('ag')
  " Use The Silver Searcher in CtrlP for listing files
  " Respect .gitignore and .agignore, ignores hidden files by default
  let g:ctrlp_user_command = 'ag --nogroup --nocolor --files-with-matches --hidden --ignore .git -g "" --literal %s'
  " Disable per-session caching
  let g:ctrlp_use_caching = 0
else
  " Exclude .gitignore patterns
  let g:ctrlp_user_command = ['.git/', 'git --git-dir=%s/.git ls-files -oc --exclude-standard']
endif
