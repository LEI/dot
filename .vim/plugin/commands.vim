" Commands

command! Write :execute ':silent w !sudo tee % > /dev/null' | :edit!
