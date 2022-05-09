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

type menuButtons struct {
	pageName string
	page     *fyne.Container
}

func main1() {
	// here i want to close the database after the app is closed
	defer DbClose()

	a := app.New()
	w := a.NewWindow("ANTI ACCOUNTANTS")
	w.Resize(fyne.Size{Width: 500, Height: 500})
	w.SetFixedSize(true)

	page := PageJournalEntry()

	menuButton := widget.NewButton("menu", nil)
	appLayout := container.New(layout.NewBorderLayout(nil, menuButton, nil, nil), page, menuButton)

	menuButton.OnTapped = func() {
		page, newPage := PageMenu()
		println(newPage)
		appLayout = container.New(layout.NewBorderLayout(nil, menuButton, nil, nil), page, menuButton)
		w.SetContent(appLayout)
		page.Refresh()
	}

	w.SetContent(appLayout)
	w.ShowAndRun()
}

func PageJournalEntry() *fyne.Container {
	height := widget.NewLabel("").MinSize().Height

	return container.NewVBox(
		Widget4(height, "save"),
		Widget4(height, "clear checked"),
		Widget4(height, "clear nan checked"),
		Widget3(height),
		Widget1(height, "name"),
		Widget1(height, "note"),
		Widget1(height, "type of compound entry"),
		Widget2(height),
		Widget2(height),
		Widget2(height),
		Widget2(height),
	)
}

func Widget3(height float32) *fyne.Container {
	w1 := widget.NewButton("ok", nil)
	w1.Resize(fyne.NewSize(50, height))

	w2 := widget.NewButton("add", nil)
	w2.Resize(fyne.NewSize(50, height))

	return container.New(&DynamicHBoxLayout{Width: 500, Height: height}, w1, w2)
}

func Widget4(height float32, placeHolder string) *fyne.Container {
	w1 := widget.NewButton(placeHolder, nil)
	w1.Resize(fyne.NewSize(50, height))

	w2 := widget.NewCheck("do with ok", nil)
	w2.Resize(fyne.NewSize(50, height))

	return container.New(&DynamicHBoxLayout{Width: 500, Height: height}, w1, w2)
}

func Widget2(height float32) *fyne.Container {
	w1 := widget.NewCheck("", nil)
	w1.Resize(fyne.NewSize(10, height))

	w2 := widget.NewSelect([]string{"cash", "book"}, nil)
	w2.PlaceHolder = " "
	w2.Resize(fyne.NewSize(20, height))

	w3 := widget.NewEntry()
	w3.SetPlaceHolder("barcode")
	w3.Resize(fyne.NewSize(20, height))

	w4 := widget.NewEntry()
	w4.SetPlaceHolder("price")
	w4.Resize(fyne.NewSize(20, height))

	w5 := widget.NewEntry()
	w5.SetPlaceHolder("quantity")
	w5.Resize(fyne.NewSize(20, height))

	w6 := widget.NewButton("x", nil)
	w6.Resize(fyne.NewSize(10, height))

	return container.New(&DynamicHBoxLayout{Width: 500, Height: height}, w1, w2, w3, w4, w5, w6)
}

func Widget1(height float32, placeHolder string) *fyne.Container {
	w1 := widget.NewCheck("", nil)
	w1.Resize(fyne.NewSize(10, height))

	w2 := widget.NewEntry()
	w2.SetPlaceHolder(placeHolder)
	w2.Resize(fyne.NewSize(80, height))

	w3 := widget.NewButton("x", nil)
	w3.Resize(fyne.NewSize(10, height))

	return container.New(&DynamicHBoxLayout{Width: 500, Height: height}, w1, w2, w3)
}

func PageMenu() (*fyne.Container, *fyne.Container) {

	pages := []menuButtons{
		{"SIMPLE JOURNAL ENTRY", PageJournalEntry()},
		{"REVERSE ENTRIES", PageJournalEntry()},
		{"JOURNAL FILTER", PageJournalEntry()},
		{"STATEMENT FILTER", PageJournalEntry()},
		{"ADD ACCOUNT", PageJournalEntry()},
		{"EDIT ACCOUNT", PageJournalEntry()},
		{"LOGIN", PageLogin()},
	}

	var page *fyne.Container
	menu := widget.NewList(
		func() int {
			return len(pages)
		},
		func() fyne.CanvasObject {
			button := widget.NewButton("", nil)
			button.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)
			return button
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Button).SetText(pages[i].pageName)
		},
	)
	PageMenu := container.New(layout.NewMaxLayout(), menu)
	return PageMenu, page
}

func PageLogin() *fyne.Container {

	w1 := widget.NewEntry()
	w1.SetPlaceHolder("company name")

	w2 := widget.NewEntry()
	w2.SetPlaceHolder("employee name")

	w3 := widget.NewPasswordEntry()
	w3.SetPlaceHolder("password")

	w4 := widget.NewLabel("")

	w5 := widget.NewButton("login", nil)
	w6 := widget.NewButton("create new employee", nil)
	w7 := widget.NewButton("create new company", nil)

	return container.New(layout.NewVBoxLayout(), w1, w2, w3, w4, w5, w6, w7)
}
