" Plugins

" https://github.com/skwp/dotfiles/blob/master/vim/settings.vim

if get(g:, 'loaded_plugins', 0)
  finish
endif
let g:loaded_plugins = 1

for s:path in split(globpath('~/.vim/settings', '*.vim'), '\n')
  let s:name = fnamemodify(s:path, ':t:r')
  let s:check = !empty(glob('~/.vim/plugged/*' . s:name . '*'))
  if s:check
    execute 'source ' . s:path
  endif
endfor
