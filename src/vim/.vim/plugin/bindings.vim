" Key bindings

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
