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

func Main1() {
	// here i want to close the database after the app is closed
	defer DbClose()

	a := app.New()
	w := a.NewWindow("ANTI ACCOUNTANTS")
	w.Resize(fyne.Size{Width: 500, Height: 500})

	// p_journal_entry := PageJournalEntry()
	pMenu, _ := PageMenu()
	pLogin := PageLogin()
	page := pLogin

	menuButton := widget.NewButton("menu", func() {
		page = pMenu
	})

	appLayout := container.New(layout.NewBorderLayout(nil, menuButton, nil, nil), page, menuButton)

	// page := container.New(layout.NeWMAxLayout(), pMenu)
	w.SetContent(appLayout)
	w.ShowAndRun()
}

func PageJournalEntry() *fyne.Container {
	name := widget.NewEntry()
	name.SetPlaceHolder("name")
	notes := widget.NewEntry()
	notes.SetPlaceHolder("notes")
	checkName := widget.NewCheck("", func(bool) {})
	checkNotes := widget.NewCheck("", func(bool) {})
	cancelName := widget.NewButton("x", func() {})
	cancelNotes := widget.NewButton("x", func() {})
	c := container.New(layout.NewGridLayout(3), checkName, name, cancelName, checkNotes, notes, cancelNotes)

	entry1 := widget.NewSelect([]string{"cash", "book"}, func(string) {})
	entry1.PlaceHolder = "cash"
	entry2 := widget.NewEntry()
	entry2.SetPlaceHolder("barcode")
	entry3 := widget.NewEntry()
	entry3.SetPlaceHolder("PRICE")
	entry4 := widget.NewEntry()
	entry4.SetPlaceHolder("QUANTITY")
	e := container.New(layout.NewGridLayout(4), entry1, entry2, entry3, entry4)

	check := widget.NewCheck("", func(bool) {})
	check.MinSize()
	cancel := widget.NewButton("x", func() {})
	cancel.MinSize()
	entries := container.New(layout.NewGridLayout(3), check, e, cancel)

	ButtonOk := widget.NewButton("ok", func() {})
	ButtonAdd := widget.NewButton("add", func() {})
	b := container.New(layout.NewGridLayout(2), ButtonOk, ButtonAdd)

	return container.New(layout.NewVBoxLayout(), c, entries, layout.NewSpacer(), b)
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
			button := widget.NewButton("", func() {})
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

type menuButtons struct {
	pageName string
	page     *fyne.Container
}

func PageLogin() *fyne.Container {

	employeeNew := widget.NewCheck("new", func(bool) {})
	employeeName := widget.NewEntry()
	employeeName.SetPlaceHolder("employee name")
	employeeEntry := container.New(layout.NewGridLayout(2), employeeName, employeeNew)

	companyNew := widget.NewCheck("new", func(i bool) {
		if i {
			employeeNew.SetChecked(true)
			employeeNew.Disable()
			employeeNew.Refresh()
		} else {
			employeeNew.SetChecked(false)
			employeeNew.Enable()
			employeeNew.Refresh()
		}
	})
	companyName := widget.NewEntry()
	companyName.SetPlaceHolder("company name")
	companyEntry := container.New(layout.NewGridLayout(2), companyName, companyNew)

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("password")

	errLabel := widget.NewLabel("")

	login := widget.NewButton("login", func() {
		CompanyName = companyName.Text
		EmployeeName = employeeName.Text
	})

	createNew := widget.NewButton("create new", func() {
		CompanyName = companyName.Text
		EmployeeName = employeeName.Text
	})

	return container.New(layout.NewVBoxLayout(), companyEntry, employeeEntry, password, errLabel, login, createNew)
}
