autocmd CursorMoved,CursorHold,InsertEnter,InsertLeave,BufEnter,BufLeave * call neovim_locpopup#show()

function! neovim_locpopup#show()
    call rpcnotify(0, "LocPopup", "show")
endfunction
