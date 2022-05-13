package main

import (
	// "image/color"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/app"
	// "fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	WindowWidth  = 500
	WindowHeight = 500
)

type menuButtons struct {
	pageName string
	page     *fyne.Container
}

func main() {
	// here i want to close the database after the app is closed
	defer DbClose()

	a := app.New()
	w := a.NewWindow("ANTI ACCOUNTANTS")
	height := widget.NewLabel("").MinSize().Height
	w.Resize(fyne.Size{Width: WindowWidth, Height: WindowHeight})

	menuButton := widget.NewButton("menu", nil)
	var menu *widget.List
	var appLayout *fyne.Container
	pages := []menuButtons{
		{"SIMPLE JOURNAL ENTRY", PageJournalEntry(height)},
		{"LOGIN", PageLogin()},
	}

	menu = widget.NewList(
		func() int {
			return len(pages)
		},
		func() fyne.CanvasObject {
			button := widget.NewButton("", nil)
			button.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)
			return button
		},
		func(k1 widget.ListItemID, v1 fyne.CanvasObject) {
			v1.(*widget.Button).SetText(pages[k1].pageName)
			v1.(*widget.Button).OnTapped = func() {
				appLayout = container.New(layout.NewBorderLayout(nil, menuButton, nil, nil), pages[k1].page, menuButton)
				w.SetContent(appLayout)
			}
		},
	)

	appLayout = container.New(layout.NewBorderLayout(nil, menuButton, nil, nil), menu, menuButton)

	menuButton.OnTapped = func() {
		appLayout = container.New(layout.NewBorderLayout(nil, menuButton, nil, nil), menu, menuButton)
		w.SetContent(appLayout)
	}

	w.SetContent(appLayout)
	w.ShowAndRun()
}

func Widget2(height float32) *fyne.Container {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, height))

	wsAccount := widget.NewSelect([]string{"cash", "book"}, nil)
	wsAccount.PlaceHolder = " "
	wsAccount.Resize(fyne.NewSize(20, height))

	weBarcode := widget.NewEntry()
	weBarcode.SetPlaceHolder("barcode")
	weBarcode.Resize(fyne.NewSize(20, height))

	wePrice := widget.NewEntry()
	wePrice.SetPlaceHolder("price")
	wePrice.Resize(fyne.NewSize(20, height))

	weQuantity := widget.NewEntry()
	weQuantity.SetPlaceHolder("quantity")
	weQuantity.Resize(fyne.NewSize(20, height))

	wbX := widget.NewButton("x", nil)
	wbX.Resize(fyne.NewSize(10, height))

	return container.New(&PercentageHBoxLayout{Width: WindowWidth, Height: height}, wc, wsAccount, weBarcode, wePrice, weQuantity, wbX)
}

func PageLogin() *fyne.Container {

	w1 := widget.NewEntry()
	w2 := widget.NewEntry()
	w3 := widget.NewPasswordEntry()
	w5 := widget.NewButton("login", nil)
	w6 := widget.NewButton("create new employee", nil)
	w7 := widget.NewButton("create new company", nil)

	w1.SetPlaceHolder("company name")
	w2.SetPlaceHolder("employee name")
	w3.SetPlaceHolder("password")

	return container.New(layout.NewVBoxLayout(), w1, w2, w3, w5, w6, w7)
}

func PageJournalEntry(height float32) *fyne.Container {
	var c2 *fyne.Container
	var wlEntries *widget.List
	lenEntries := 1

	wlEntries = widget.NewList(
		func() int {
			return lenEntries
		},
		func() fyne.CanvasObject {
			return Widget2(height)
		},
		func(k1 widget.ListItemID, v1 fyne.CanvasObject) {
			v1.(*fyne.Container).Objects[5].(*widget.Button).OnTapped = func() {}
		},
	)
	wlEntries.Resize(fyne.Size{Width: WindowWidth, Height: height})

	wcName := widget.NewCheck("", nil)
	wcName.Resize(fyne.NewSize(10, height))
	weName := widget.NewEntry()
	weName.SetPlaceHolder("name")
	weName.Resize(fyne.NewSize(80, height))
	wbName := widget.NewButton("x", func() {
		weName.SetText("")
	})
	wbName.Resize(fyne.NewSize(10, height))
	cName := container.New(&PercentageHBoxLayout{Width: WindowWidth, Height: height}, wcName, weName, wbName)

	wcNote := widget.NewCheck("", nil)
	wcNote.Resize(fyne.NewSize(10, height))
	weNote := widget.NewEntry()
	weNote.SetPlaceHolder("note")
	weNote.Resize(fyne.NewSize(80, height))
	wbNote := widget.NewButton("x", func() {
		weNote.SetText("")
	})
	wbNote.Resize(fyne.NewSize(10, height))
	cNote := container.New(&PercentageHBoxLayout{Width: WindowWidth, Height: height}, wcNote, weNote, wbNote)

	wcType := widget.NewCheck("", nil)
	wcType.Resize(fyne.NewSize(10, height))
	weType := widget.NewEntry()
	weType.SetPlaceHolder("type")
	weType.Resize(fyne.NewSize(80, height))
	wbType := widget.NewButton("x", func() {
		weType.SetText("")
	})
	wbType.Resize(fyne.NewSize(10, height))
	cType := container.New(&PercentageHBoxLayout{Width: WindowWidth, Height: height}, wcType, weType, wbType)

	wbSave := widget.NewButton("save", nil)
	wbSave.Resize(fyne.Size{Width: wbSave.MinSize().Width, Height: height})

	wbClearChecked := widget.NewButton("clear checked", func() {
		if wcName.Checked {
			weName.SetText("")
		}
		if wcNote.Checked {
			weNote.SetText("")
		}
		if wcType.Checked {
			weType.SetText("")
		}
	})
	wbClearChecked.Resize(fyne.Size{Width: wbClearChecked.MinSize().Width, Height: height})

	wbClearNotChecked := widget.NewButton("clear not checked", func() {
		if !wcName.Checked {
			weName.SetText("")
		}
		if !wcNote.Checked {
			weNote.SetText("")
		}
		if !wcType.Checked {
			weType.SetText("")
		}
	})
	wbClearNotChecked.Resize(fyne.Size{Width: wbClearNotChecked.MinSize().Width, Height: height})

	wbOk := widget.NewButton("ok", func() {
		if !wcName.Checked {
			weName.SetText("")
		}
		if !wcNote.Checked {
			weNote.SetText("")
		}
		if !wcType.Checked {
			weType.SetText("")
		}
	})
	wbOk.Resize(fyne.Size{Width: wbOk.MinSize().Width, Height: height})

	wbAdd := widget.NewButton("add", func() {
		lenEntries++
		c2.Refresh()
	})
	wbAdd.Resize(fyne.Size{Width: wbAdd.MinSize().Width, Height: height})

	wlError := widget.NewLabel("error massage")
	wlError.Resize(fyne.Size{Width: WindowWidth, Height: height})

	c1 := container.New(&PercentageHBoxLayout{Width: WindowWidth, Height: height}, wbOk, wbSave, wbClearChecked, wbClearNotChecked, wbAdd)
	c2 = container.New(&StretchVBoxLayout{Width: WindowWidth, Height: height, ObjectToStertch: wlEntries}, c1, cName, cNote, cType, wlError, wlEntries)

	return c2
}
