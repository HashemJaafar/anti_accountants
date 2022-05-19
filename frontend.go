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

var (
	height    float32 = 50
	fcEntries *fyne.Container
)

type menuButtons struct {
	pageName string
	page     *fyne.Container
}

func main() {
	defer DbClose()

	a := app.New()
	w := a.NewWindow("ANTI ACCOUNTANTS")
	height = widget.NewLabel("").MinSize().Height
	w.Resize(fyne.Size{Width: WindowWidth, Height: WindowHeight})

	wbMenu := widget.NewButton("menu", nil)
	var cApp *fyne.Container

	wbJournal := widget.NewButton("SIMPLE JOURNAL ENTRY", func() {
		cApp = container.New(layout.NewBorderLayout(nil, wbMenu, nil, nil), wbMenu, PageJournalEntry())
		w.SetContent(cApp)
	})
	wbJournal.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)

	wbLogin := widget.NewButton("LOGIN", func() {
		w.SetContent(PageLogin())
	})
	wbLogin.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)

	wsMenu := container.NewVScroll(container.New(layout.NewVBoxLayout(), wbJournal, wbLogin))

	wbMenu.OnTapped = func() {
		cApp = container.New(layout.NewBorderLayout(nil, wbMenu, nil, nil), wbMenu, wsMenu)
		w.SetContent(cApp)
	}

	cApp = container.New(layout.NewBorderLayout(nil, wbMenu, nil, nil), wbMenu, wsMenu)
	w.SetContent(cApp)

	w.ShowAndRun()
}

func FcEntryABPQ() *fyne.Container {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, height))

	wsAccount := widget.NewSelectEntry([]string{"cash", "book"})
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

	fc := container.New(&PercentageHBoxLayout{Width: WindowWidth, Height: height}, wc, wsAccount, weBarcode, wePrice, weQuantity, wbX)

	wbX.OnTapped = func() {
		index, _ := FindObject(fc, fcEntries.Objects)
		fcEntries.Objects = Remove(fcEntries.Objects, index)
	}
	return fc
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

func PageJournalEntry() *fyne.Container {

	wlError := widget.NewLabel("")
	wlError.Wrapping = fyne.TextWrapWord
	wlError.Resize(fyne.Size{Width: WindowWidth, Height: height})

	wbOk := widget.NewButton("ok", func() {
		_, err := SaveEntries(true)
		DisplayTheError(err, wlError)
		ClearEntries(false)
	})
	wbOk.Resize(fyne.Size{Width: wbOk.MinSize().Width, Height: height})

	wbSave := widget.NewButton("save", func() {
		_, err := SaveEntries(true)
		DisplayTheError(err, wlError)
	})
	wbSave.Resize(fyne.Size{Width: wbSave.MinSize().Width, Height: height})

	wbClearChecked := widget.NewButton("clear checked", func() {
		ClearEntries(true)
	})
	wbClearChecked.Resize(fyne.Size{Width: wbClearChecked.MinSize().Width, Height: height})

	wbClearNotChecked := widget.NewButton("clear not checked", func() {
		ClearEntries(false)
	})
	wbClearNotChecked.Resize(fyne.Size{Width: wbClearNotChecked.MinSize().Width, Height: height})

	fcEntries = container.New(layout.NewVBoxLayout(), FcEntryInfo("note"), FcEntryInfo("name"), FcEntryInfo("type"), FcEntryABPQ())
	wsEntries := container.NewVScroll(fcEntries)

	wbAdd := widget.NewButton("+", func() {
		fcEntries.Add(FcEntryABPQ())
		fcEntries.Refresh()
	})
	wbAdd.Resize(fyne.Size{Width: wbAdd.MinSize().Width, Height: height})

	fc1 := container.New(&PercentageHBoxLayout{Width: WindowWidth, Height: height}, wbOk, wbSave, wbClearChecked, wbClearNotChecked, wbAdd)
	fc2 := container.New(&StretchVBoxLayout{Width: WindowWidth, Height: height, ObjectToStertch: wsEntries}, fc1, wsEntries, wlError)

	return fc2
}

func DisplayTheError(err error, wlError *widget.Label) {
	if err != nil {
		wlError.SetText(err.Error())
	} else {
		wlError.SetText("")
	}
}

func FcEntryInfo(label string) *fyne.Container {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, height))
	we := widget.NewEntry()
	we.SetPlaceHolder(label)
	we.Resize(fyne.NewSize(80, height))
	wb := widget.NewButton("x", func() {
		we.SetText("")
	})
	wb.Resize(fyne.NewSize(10, height))
	return container.New(&PercentageHBoxLayout{Width: WindowWidth, Height: height}, wc, we, wb)
}

func ClearEntries(isChecked bool) {
	for k1, v1 := range fcEntries.Objects {
		if k1 == 3 {
			break
		}
		if v1.(*fyne.Container).Objects[0].(*widget.Check).Checked == isChecked {
			v1.(*fyne.Container).Objects[1].(*widget.Entry).SetText("")
		}
	}

	k1 := 3
	for k1 < len(fcEntries.Objects) {
		if fcEntries.Objects[k1].(*fyne.Container).Objects[0].(*widget.Check).Checked == isChecked {
			fcEntries.Objects = Remove(fcEntries.Objects, k1)
		} else {
			k1++
		}
	}
}

func SaveEntries(insert bool) ([]APQB, error) {
	entryInfo := EntryInfo{
		Notes:               fcEntries.Objects[0].(*fyne.Container).Objects[1].(*widget.Entry).Text,
		Name:                fcEntries.Objects[1].(*fyne.Container).Objects[1].(*widget.Entry).Text,
		TypeOfCompoundEntry: fcEntries.Objects[2].(*fyne.Container).Objects[1].(*widget.Entry).Text,
	}

	var entries []APQB
	for _, v1 := range fcEntries.Objects[3:] {
		entries = append(entries, APQB{
			Name: v1.(*fyne.Container).Objects[1].(*widget.Select).Selected,
			// Price:    v1.(*fyne.Container).Objects[3].(*widget.Label).Text,
			// Quantity: v1.(*fyne.Container).Objects[4].(*widget.Label).Text,
			Barcode: v1.(*fyne.Container).Objects[2].(*widget.Entry).Text,
		})
	}

	return SimpleJournalEntry(entries, entryInfo, insert)
}
