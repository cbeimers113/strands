package component

import (
	"sort"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/gui/assets/icon"

	"cbeimers113/strands/internal/gui/color"
)

type fileItem struct {
	*gui.Panel
	label        *gui.Label
	openButton   *gui.Button
	deleteButton *gui.Button
}

func newFileItem(text string, width float32, viewType int) *fileItem {
	fi := &fileItem{}

	fi.openButton = gui.NewButton("")
	fi.openButton.SetUserData(viewType)
	fi.openButton.SetIcon(icon.FolderOpen)
	fi.openButton.SetPosition(width-2*fi.openButton.Width()-10, 0)

	fi.deleteButton = gui.NewButton("")
	fi.deleteButton.SetUserData(viewType)
	fi.deleteButton.SetIcon(icon.Delete)
	fi.deleteButton.SetPosition(width-fi.deleteButton.Width()-5, 0)

	fi.label = gui.NewLabel(text)
	fi.label.SetUserData(viewType)
	fi.label.SetPosition(5, (fi.openButton.Height()-fi.label.Height())/2)

	fi.Panel = gui.NewPanel(width, fi.openButton.Height())
	fi.Panel.SetUserData(viewType)
	fi.Panel.Add(fi.label)
	fi.Panel.Add(fi.openButton)
	fi.Panel.Add(fi.deleteButton)
	fi.Panel.SetVisible(true)

	// Panel hover color
	fi.Panel.Subscribe(gui.OnCursor, func(evname string, ev interface{}) {
		if !fi.Enabled() {
			return
		}

		fi.Panel.SetColor(color.Focus)
	})
	fi.Panel.Subscribe(gui.OnCursorLeave, func(evname string, ev interface{}) {
		fi.Panel.SetColor(color.Background)
	})

	// open button hover color
	fi.openButton.Subscribe(gui.OnCursor, func(evname string, ev interface{}) {
		if !fi.Enabled() {
			return
		}

		fi.Panel.SetColor(color.Green)
	})

	// delete button hover color
	fi.deleteButton.Subscribe(gui.OnCursor, func(evname string, ev interface{}) {
		if !fi.Enabled() {
			return
		}

		fi.Panel.SetColor(color.Red)
	})

	return fi
}

type FileList struct {
	*gui.ItemScroller
	Selected string // Track which file has been selected for the primary action (opening, etc)
	Deleted  string // Track which file has been selected for deletion
}

func NewFileList(filepaths map[string]string, width, height float32, viewType int) *FileList {
	f := &FileList{
		ItemScroller: gui.NewVScroller(width, height),
	}

	var (
		y     float32
		order = make([]string, 0)
	)

	for name := range filepaths {
		order = append(order, name)
	}

	// List files alphabetically
	sort.Strings(order)
	for _, name := range order {
		item := newFileItem(name, width, viewType)
		item.SetPosition(0, y)
		item.openButton.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
			f.Selected = filepaths[item.label.Text()]
			f.Deleted = ""
		})
		item.deleteButton.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
			f.Deleted = filepaths[item.label.Text()]
			f.Selected = ""
		})

		y += item.Height() + 5
		f.Add(item)
	}

	return f
}

func (f *FileList) SetEnabled(active bool) {
	for _, child := range f.Children() {
		if fi, ok := child.(*fileItem); ok {
			fi.openButton.SetEnabled(active)
			fi.deleteButton.SetEnabled(active)
			fi.SetEnabled(active)
		}
	}
}
