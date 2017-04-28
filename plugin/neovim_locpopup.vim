autocmd CursorMoved,CursorMovedI,CursorHold,CursorHoldI,InsertLeave * call neovim_locpopup#show()

function! neovim_locpopup#show()
    call rpcnotify(0, "LocPopup", "show", getloclist(winnr("$")))
endfunction
