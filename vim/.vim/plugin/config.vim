" Main configuration

if &term =~# '256color'
  " Disable Background Color Erase (BCE) so that color schemes
  " work properly when Vim is used inside tmux and GNU screen.
  " See also http://snk.tuxfamily.org/log/vim-256color-bce.html
  set t_ut=
endif

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

" Allow modified buffers in the background
set hidden

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
