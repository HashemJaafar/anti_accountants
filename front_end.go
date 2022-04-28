package main

import (
	// "image/color"

	"fmt"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/app"
	// "fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// here i want to close the database after the app is closed
	defer DB_CLOSE()

	a := app.New()
	w := a.NewWindow("ANTI ACCOUNTANTS")

	// p_journal_entry := page_journal_entry()
	p_menu, _ := page_menu()
	p_login := page_login()
	page := p_login

	menu_button := widget.NewButton("menu", func() {
		page = p_menu
	})

	app_layout := container.New(layout.NewBorderLayout(nil, menu_button, nil, nil), page, menu_button)

	// page := container.New(layout.NewMaxLayout(), p_menu)
	w.SetContent(app_layout)
	w.ShowAndRun()
}

func page_journal_entry() *fyne.Container {
	name := widget.NewEntry()
	name.SetPlaceHolder("name")
	notes := widget.NewEntry()
	notes.SetPlaceHolder("notes")
	check_name := widget.NewCheck("", func(bool) {})
	check_notes := widget.NewCheck("", func(bool) {})
	cancel_name := widget.NewButton("x", func() {})
	cancel_notes := widget.NewButton("x", func() {})
	c := container.New(layout.NewGridLayout(3), check_name, name, cancel_name, check_notes, notes, cancel_notes)

	entry1 := widget.NewSelect([]string{"cash", "book"}, func(string) {})
	entry1.PlaceHolder = "cash"
	entry2 := widget.NewEntry()
	entry2.SetPlaceHolder("barcode")
	entry3 := widget.NewEntry()
	entry3.SetPlaceHolder("price")
	entry4 := widget.NewEntry()
	entry4.SetPlaceHolder("quantity")
	e := container.New(layout.NewGridLayout(4), entry1, entry2, entry3, entry4)

	check := widget.NewCheck("", func(bool) {})
	check.MinSize()
	cancel := widget.NewButton("x", func() {})
	cancel.MinSize()
	entries := container.New(layout.NewGridLayout(3), check, e, cancel)

	button_ok := widget.NewButton("ok", func() {})
	button_add := widget.NewButton("add", func() {})
	b := container.New(layout.NewGridLayout(2), button_ok, button_add)

	return container.New(layout.NewVBoxLayout(), c, entries, layout.NewSpacer(), b)
}

func page_menu() (*fyne.Container, *fyne.Container) {

	pages := []menu_buttons{
		{"SIMPLE JOURNAL ENTRY", page_journal_entry()},
		{"REVERSE ENTRIES", page_journal_entry()},
		{"JOURNAL FILTER", page_journal_entry()},
		{"STATEMENT FILTER", page_journal_entry()},
		{"ADD ACCOUNT", page_journal_entry()},
		{"EDIT ACCOUNT", page_journal_entry()},
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
			o.(*widget.Button).SetText(pages[i].page_name)
		},
	)
	page_menu := container.New(layout.NewMaxLayout(), menu)
	return page_menu, page
}

type menu_buttons struct {
	page_name string
	page      *fyne.Container
}

func page_login() *fyne.Container {

	employee_new := widget.NewCheck("new", func(bool) {})
	employee_name := widget.NewEntry()
	employee_name.SetPlaceHolder("employee name")
	employee_entry := container.New(layout.NewHBoxLayout(), employee_name, employee_new)

	company_new := widget.NewCheck("new", func(i bool) {
		if i {
			fmt.Println("yeeeeeeeeees")
			employee_new.SetChecked(true)
			employee_new.Refresh()
		}
	})
	company_name := widget.NewEntry()
	company_name.SetPlaceHolder("company name")
	company_entry := container.New(layout.NewHBoxLayout(), company_name, company_new)

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("password")

	err_label := widget.NewLabel("")

	login := widget.NewButton("login", func() {
		COMPANY_NAME = company_name.Text
		EMPLOYEE_NAME = employee_name.Text
	})

	create_new := widget.NewButton("create new", func() {
		COMPANY_NAME = company_name.Text
		EMPLOYEE_NAME = employee_name.Text
	})

	return container.New(layout.NewVBoxLayout(), company_entry, employee_entry, password, err_label, login, create_new)
}
