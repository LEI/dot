" ~/.vimrc

" Auto download Vim Plug
let g:vim_plug = expand('~/.vim/autoload/plug.vim')
let g:vim_plug_url = 'https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim'
if !filereadable(g:vim_plug)
  execute '!curl -sfLo ' . g:vim_plug
    \ . ' --create-dirs '
    \ . g:vim_plug_url
endif

let g:vim_plugins = expand('~/.vim/plugged')
call plug#begin(g:vim_plugins)

Plug 'altercation/vim-colors-solarized'
" kien/ctrlp.vim
Plug 'tpope/vim-commentary'
Plug 'tpope/vim-repeat'
Plug 'tpope/vim-sensible'
Plug 'tpope/vim-sleuth'
Plug 'tpope/vim-surround'

" Add plugins to &runtimepath
call plug#end()

" Install plugins
if !isdirectory(g:vim_plugins)
  PlugInstall
endif

" Settings

set background=dark
colorscheme solarized

if has('mouse')
  set mouse+=a
endif

" Show invisible characters
set list
set listchars=tab:>\ ,trail:-,extends:>,precedes:<,nbsp:+

" Do not capture all global options
set sessionoptions-=options

"set viminfo='10,\"100,:20,%,n~/.viminfo
if !empty(&viminfo)
  set viminfo^=!
endif

" Disable swapfiles and backups
set noswapfile
set nobackup
set nowritebackup

" Keep undo history across sessions
if has('persistent_undo')
  let g:vim_backups = expand('~/.vim/backups')
  if !isdirectory(g:vim_backups)
    execute 'silent !mkdir ' . g:vim_backups . ' > /dev/null 2>&1'
  endif
  let &undodir = g:vim_backups
  set undofile
endif

function! RestoreCursor()
  if line("'\"") > 0 && line("'\"") <= line("$")
    normal! g`"
    return 1
  endif
endfunction

augroup RestoreCursor
  autocmd!
  autocmd BufReadPost * call RestoreCursor()
augroup END
