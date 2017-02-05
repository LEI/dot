" Disable swapfiles and backups
set noswapfile
set nobackup
set nowritebackup

" Keep undo history across sessions
if has('persistent_undo')
  let g:vim_backups = expand(get(g:, 'vim_backups', '~/.vim/backups'))
  if !isdirectory(g:vim_backups)
    execute 'silent !mkdir ' . g:vim_backups . ' > /dev/null 2>&1'
  endif
  let &undodir = g:vim_backups
  set undofile
endif
