" Vim

if filereadable(expand('~/.vim/before.vim'))
  source ~/.vim/before.vim
endif

" Auto download Vim Plug
let g:vim_plug = expand('~/.vim/autoload/plug.vim')
let g:vim_plug_url = 'https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim'
if !filereadable(g:vim_plug)
  execute 'silent !curl -sfLo ' . g:vim_plug
    \ . '  --create-dirs ' . g:vim_plug_url
endif

let g:vim_plugins = expand('~/.vim/plugged')
call plug#begin(g:vim_plugins)

" Plug 'altercation/vim-colors-solarized'
Plug 'kien/ctrlp.vim'
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

" try
"   set background=dark
"   colorscheme solarized
"   call togglebg#map('<F5>')
" catch /E185:/
"   colorscheme default
" endtry

if filereadable(expand('~/.vimrc.local'))
  source ~/.vimrc.local
endif
