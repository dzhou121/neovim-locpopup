package locpopup

import (
	"sort"
	"sync"

	"github.com/neovim/go-client/nvim"
)

// Locpop is
type Locpop struct {
	nvim      *nvim.Nvim
	mutex     sync.Mutex
	lastShown bool
	lastText  string
	lastType  string
}

// RegisterPlugin registers this remote plugin
func RegisterPlugin(nvim *nvim.Nvim) {
	locpop := &Locpop{
		nvim: nvim,
	}
	nvim.Subscribe("LocPopup")
	nvim.RegisterHandler("LocPopup", func(args ...interface{}) {
		if len(args) < 1 {
			return
		}
		event, ok := args[0].(string)
		if !ok {
			return
		}
		switch event {
		case "show":
			locpop.show(args[1:])
		}
	})
}

func (l *Locpop) show(args []interface{}) {
	l.mutex.Lock()
	shown := false
	defer func() {
		if shown {
			l.lastShown = true
		} else {
			if l.lastShown {
				l.lastShown = false
				l.lastText = ""
				l.lastType = ""
				l.nvim.Call("rpcnotify", nil, 0, "Gui", "locpopup_hide")
			}
		}
		l.mutex.Unlock()
	}()
	curWin, err := l.nvim.CurrentWindow()
	if err != nil {
		return
	}
	pos, err := l.nvim.WindowCursor(curWin)
	if err != nil {
		return
	}
	result := new([]map[string]interface{})
	err = l.nvim.Call("getloclist", result, "winnr(\"$\")")
	if err != nil {
		return
	}
	locs := []map[string]interface{}{}
	for _, loc := range *result {
		lnumInterface := loc["lnum"]
		if lnumInterface == nil {
			continue
		}
		lnum := reflectToInt(lnumInterface)
		if lnum == pos[0] {
			locs = append(locs, loc)
		}
	}
	if len(locs) == 0 {
		return
	}
	if len(locs) > 1 {
		sort.Sort(ByCol(locs))
	}
	var loc map[string]interface{}
	for _, loc = range locs {
		if pos[1] >= reflectToInt(loc["col"])-1 {
			break
		}
	}

	locType := loc["type"].(string)
	text := loc["text"].(string)
	if locType != l.lastType || text != l.lastText {
		l.lastText = text
		l.lastType = locType
		l.nvim.Call("rpcnotify", nil, 0, "Gui", "locpopup_show", loc)
	}
	shown = true
}

// ByCol sorts locations by column
type ByCol []map[string]interface{}

// Len of locations
func (s ByCol) Len() int {
	return len(s)
}

// Swap locations
func (s ByCol) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less than
func (s ByCol) Less(i, j int) bool {
	return reflectToInt(s[i]["col"]) > reflectToInt(s[j]["col"])
}

func reflectToInt(iface interface{}) int {
	o, ok := iface.(int64)
	if ok {
		return int(o)
	}
	return int(iface.(uint64))
}
