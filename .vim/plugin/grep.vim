" Grep

if executable('ag')
  " Use Ag over Grep
  set grepprg=ag\ --vimgrep\ --nogroup\ --nocolor
  " Command output format (default "%f:%l:%m,%f:%l%m,%f  %l%m")
  set grepformat=%f:%l:%c:%m

  " Bind \ (backward slash) to grep shortcut
  " command -nargs=+ -complete=file -bar Ag silent! grep! <args>|cwindow|redraw!
  " nnoremap \ :Ag<SPACE>
endif
