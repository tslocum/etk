package etk

import (
	"image/color"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// FilePickerMode represents a FilePicker selection mode.
type FilePickerMode int

// FilePicker modes.
const (
	ModeCreateDir  FilePickerMode = 0
	ModeCreateFile FilePickerMode = 1
	ModeSelectDir  FilePickerMode = 2
	ModeSelectFile FilePickerMode = 3
)

// FilePicker is a file and directory creation and selection dialog.
type FilePicker struct {
	*Grid
	List        *List
	mode        FilePickerMode
	dir         string
	entries     []string
	extensions  []string
	dirLabel    *Text
	inputField  *Input
	cancelLabel string
	selectLabel string
	onResult    func(path string) error
	needRebuild bool
	sync.Mutex
}

// NewFilePicker returns a new FilePicker.
func NewFilePicker(mode FilePickerMode, dir string, extensions []string, onResult func(path string) error) *FilePicker {
	itemHeight := int(float64(Style.TextSize) * 1.5)

	f := &FilePicker{
		Grid:        NewGrid(),
		mode:        mode,
		dir:         dir,
		extensions:  extensions,
		dirLabel:    NewText(dir),
		cancelLabel: "Cancel",
		selectLabel: "Select",
		onResult:    onResult,
		needRebuild: true,
	}

	f.dirLabel.SetVertical(AlignCenter)
	f.dirLabel.SetAutoResize(true)

	f.List = NewList(itemHeight, f.onListSelected, f.onListConfirmed)

	f.inputField = NewInput("", f.onInputSelected)
	f.inputField.SetVertical(AlignCenter)
	return f
}

func (f *FilePicker) SetFocus(focus bool) (accept bool) {
	if focus {
		SetFocus(f.inputField)
	}
	return false
}

func (f *FilePicker) handleResult(index int) {
	var path string
	modeCreate := f.mode == ModeCreateDir || f.mode == ModeCreateFile
	if modeCreate {
		path = f.dir
		name := f.inputField.Text()
		if name == "" {
			name = f.entries[index]
		} else if len(f.extensions) == 1 && !strings.HasSuffix(strings.ToLower(name), f.extensions[0]) {
			name += f.extensions[0]
		}
		path = filepath.Join(path, name)
	} else {
		var err error
		path, err = filepath.Abs(filepath.Join(f.dir, f.entries[index]))
		if err != nil {
			log.Fatal(err)
		}
	}
	err := f.onResult(path)
	if err != nil {
		log.Fatal(err)
	}
}

func (f *FilePicker) handleSelected(index int) {
	path, err := filepath.Abs(filepath.Join(f.dir, f.entries[index]))
	if err != nil {
		log.Fatal(err)
	}
	modeCreate := f.mode == ModeCreateDir || f.mode == ModeCreateFile
	if modeCreate {
		if strings.TrimSpace(f.inputField.Text()) == "" {
			isDir := strings.HasSuffix(f.entries[index], "/")
			if isDir != (f.mode == ModeCreateDir) {
				if isDir {
					f.dir = path
					f.needRebuild = true
				}
				return
			}
		}
		f.handleResult(index)
		return
	}
	selectFile := f.mode == ModeCreateFile || f.mode == ModeSelectFile
	if selectFile && strings.HasSuffix(f.entries[index], "/") {
		f.dir = path
		f.needRebuild = true
		return
	}
	f.handleResult(index)
}

func (f *FilePicker) onListSelected(index int) (accept bool) {
	return true
}

func (f *FilePicker) onListConfirmed(index int) {
	entry := f.entries[index]
	_, selected := f.List.SelectedItem()
	if index == selected && strings.HasSuffix(entry, "/") {
		abs, err := filepath.Abs(filepath.Join(f.dir, entry))
		if err == nil {
			f.dir = abs
			f.needRebuild = true
		}
		return
	}
	f.handleResult(index)
}

func (f *FilePicker) onInputSelected(text string) (handled bool) {
	_, index := f.List.SelectedItem()
	f.handleSelected(index)
	return true
}

func (f *FilePicker) onButtonSelected() error {
	_, index := f.List.SelectedItem()
	f.handleSelected(index)
	return nil
}

func (f *FilePicker) onCancel() error {
	return f.onResult("")
}

func (f *FilePicker) walkDir(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	label := strings.TrimPrefix(path, f.dir)
	if len(label) > 1 && label[0] == '/' {
		label = label[1:]
	}

	if d.IsDir() {
		if path == f.dir {
			return nil
		}
		label += "/"
		f.entries = append(f.entries, label)
		return filepath.SkipDir
	}

	if len(f.extensions) > 0 {
		var found bool
		for i := range f.extensions {
			if strings.HasSuffix(strings.ToLower(d.Name()), f.extensions[i]) {
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	f.entries = append(f.entries, label)
	return nil
}

func (f *FilePicker) rebuild() {
	f.Grid.Clear()
	f.List.Clear()

	path := f.dir
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	f.dirLabel.SetText(path)

	f.entries = f.entries[:0]
	filepath.WalkDir(f.dir, f.walkDir)
	sort.Slice(f.entries, func(i, j int) bool {
		if strings.HasSuffix(f.entries[i], "/") != strings.HasSuffix(f.entries[j], "/") {
			return strings.HasSuffix(f.entries[i], "/")
		}
		return strings.ToLower(f.entries[i]) < strings.ToLower(f.entries[j])
	})

	selectDir := f.mode == ModeCreateDir || f.mode == ModeSelectDir
	if selectDir {
		f.entries = append([]string{"./"}, f.entries...)
	}
	if f.dir != "/" {
		f.entries = append([]string{"../"}, f.entries...)
	}

	var y int
	for _, entry := range f.entries {
		t := NewText(entry)
		t.SetPadding(0)
		t.SetAutoResize(true)
		t.SetVertical(AlignCenter)

		g := NewGrid()
		g.SetColumnSizes(5, -1, 5)
		g.AddChildAt(&WithoutMouse{t}, 1, 0, 1, 1)

		f.List.AddChildAt(&WithoutMouse{g}, 0, y)
		y++
	}
	f.List.SetSelectedItem(0, 0)

	dividerA := NewBox()
	dividerA.SetBackground(color.RGBA{255, 255, 255, 255})
	dividerB := NewBox()
	dividerB.SetBackground(color.RGBA{255, 255, 255, 255})

	showInput := f.mode == ModeCreateDir || f.mode == ModeCreateFile

	rowSizes := []int{Style.TextSize * 2, 2, -1, 2}
	if showInput {
		rowSizes = append(rowSizes, Style.TextSize*2)
	}
	rowSizes = append(rowSizes, Style.TextSize*2)

	f.Grid.SetRowSizes(rowSizes...)
	f.Grid.AddChildAt(f.dirLabel, 0, 0, 2, 1)
	f.Grid.AddChildAt(dividerA, 0, 1, 2, 1)
	f.Grid.AddChildAt(&WithoutFocus{f.List}, 0, 2, 2, 1)
	f.Grid.AddChildAt(dividerB, 0, 3, 2, 1)
	y = 4
	if showInput {
		nameText := "Name"
		nameSize := 150
		if len(f.extensions) == 1 {
			nameText += " (" + f.extensions[0] + ")"
			nameSize = 300
		}
		nameLabel := NewText(nameText)
		nameLabel.SetVertical(AlignCenter)
		nameLabel.SetAutoResize(true)
		g := NewGrid()
		g.SetColumnPadding(5)
		g.SetColumnSizes(nameSize, -1)
		g.AddChildAt(nameLabel, 0, 0, 1, 1)
		g.AddChildAt(f.inputField, 1, 0, 1, 1)
		f.Grid.AddChildAt(g, 0, y, 2, 1)
		y++
	}
	f.Grid.AddChildAt(NewButton(f.cancelLabel, f.onCancel), 0, y, 1, 1)
	f.Grid.AddChildAt(NewButton(f.selectLabel, f.onButtonSelected), 1, y, 1, 1)
}

// SetMode sets the FilePicker mode.
func (f *FilePicker) SetMode(mode FilePickerMode) {
	f.Lock()
	defer f.Unlock()

	f.mode = mode
	f.needRebuild = true
}

// SetExtensions sets the desired file extensions, if any. When set, only files
// with matching extensions are shown. When creating a file and only one
// extension is set, the file will be created with the specified extension.
func (f *FilePicker) SetExtensions(extensions []string) {
	f.Lock()
	defer f.Unlock()

	f.extensions = make([]string, len(extensions))
	for i := range extensions {
		f.extensions[i] = strings.ToLower(extensions[i])
	}
	f.needRebuild = true
}

// SetOnResult sets the FilePicker result handler. When a file or directory is
// selected, depending on the FilePicker mode, the path to the file or directory
// is provided. When the FilePicker is canceled, a blank path is provided.
func (f *FilePicker) SetOnResult(onResult func(path string) error) {
	f.Lock()
	defer f.Unlock()

	f.onResult = onResult
	f.needRebuild = true
}

// SetButtonLabels sets the FilePicker cancel and confirm button labels.
func (f *FilePicker) SetButtonLabels(cancel string, confirm string) {
	f.Lock()
	defer f.Unlock()

	f.cancelLabel = cancel
	f.selectLabel = confirm
	f.needRebuild = true
}

// Draw draws the FilePicker on the screen.
func (f *FilePicker) Draw(screen *ebiten.Image) error {
	f.Lock()
	defer f.Unlock()

	if f.needRebuild {
		f.rebuild()
		f.needRebuild = false
	}

	return f.Grid.Draw(screen)
}
