package component

import (
	"sort"

	"github.com/g3n/engine/gui"
)

type FileList struct {
	*gui.ItemScroller
	Selected string
}

func NewFileList(filepaths map[string]string, width, height float32, viewType int) *FileList {
	f := &FileList{
		ItemScroller: gui.NewVScroller(width, height),
	}
	var y float32
	var order = make([]string, 0)

	for name := range filepaths {
		order = append(order, name)
	}

	// List files alphabetically
	sort.Strings(order)
	for _, name := range order {
		button := gui.NewButton(name)
		button.SetWidth(width)
		button.SetHeight(button.ContentHeight() * 1.1)
		button.SetPosition(0, y)
		button.SetUserData(viewType)
		button.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
			f.Selected = filepaths[button.Label.Text()]
		})

		y += button.Height() + 5
		f.Add(button)
	}

	return f
}
