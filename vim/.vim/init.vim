" Vim configuration

if filereadable(expand('~/.vim/before.vim'))
  source ~/.vim/before.vim
endif

" Auto download Vim Plug
let g:vim_plug = expand('~/.vim/autoload/plug.vim')
let g:vim_plug_url = 'https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim'
if filereadable(g:vim_plug)
  execute '!curl -sfLo --create-dirs "' . g:vim_plug . '"'
    \ . ' "' . g:vim_plug_url . '"'
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

if &term =~# '256color'
  " Disable Background Color Erase (BCE) so that color schemes
  " work properly when Vim is used inside tmux and GNU screen.
  " See also http://snk.tuxfamily.org/log/vim-256color-bce.html
  set t_ut=
endif

try
  set background=dark
  colorscheme solarized
  call togglebg#map('<F5>')
catch /E185:/
  colorscheme default
endtry

set synmaxcol=500

" Relative to textwidth
if exists('+colorcolumn')
  set colorcolumn=+1
endif

set number
if exists('&relativenumber')
  set relativenumber
endif

if has('mouse')
  set mouse+=a
endif

if !has('nvim')
  " Fix mouse inside screen and tmux
  if &term =~# '^screen' || strlen($TMUX) > 0
    set ttymouse=xterm2
  endif
  " Fast terminal connection
  set ttyfast
endif

" Current mode in status line
set showmode

" Display incomplete commands
set showcmd

" Highlight previous matches
set hlsearch

" Ignore case in search patterns
set ignorecase

" Do not ignore when the pattern containes upper case characters
set smartcase

if !exists('g:loaded_sleuth')
  set expandtab
  set shiftwidth=4
  set softtabstop=4
  set tabstop=4
endif

" Show invisible characters
set list
" set listchars=tab:>\ ,trail:-,extends:>,precedes:<,nbsp:+
let &listchars = 'tab:' . nr2char(0x25B8) . ' '
  \ . ',trail:' . nr2char(0x00B7)
  \ . ',extends:' . nr2char(0x276F)
  \ . ',precedes:' . nr2char(0x276E)
  \ . ',nbsp:' . nr2char(0x005F)
  \ . ',eol:' . nr2char(0x00AC)

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

" Change leader
let g:mapleader = "\<Space>"

" Yank from the cursor to the end of the line
map Y y$

" Move vertically on wrapped lines
nnoremap j gj
nnoremap k gk

" Split navigation shortcuts
nnoremap <C-h> <C-w>h
nnoremap <C-j> <C-w>j
nnoremap <C-k> <C-w>k
nnoremap <C-l> <C-w>l

" Clear highlighted search results
nnoremap <Space> :nohlsearch<CR>

" command W :execute ':silent w !sudo tee % > /dev/null' | :edit!
command Write :w !sudo tee % > /dev/null
" noremap <Leader>W :w !sudo tee % > /dev/null<CR>
