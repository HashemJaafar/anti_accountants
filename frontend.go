package main

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	CWindowWidth  = 500
	CWindowHeight = 500
)

var (
	Vheight          float32 = 50
	VfcEntries       *fyne.Container
	VAllAccountNames = FAllAccountNames()
)

func main() {
	defer FDbClose()

	a := app.New()
	w := a.NewWindow("ANTI ACCOUNTANTS")
	Vheight = widget.NewLabel("").MinSize().Height
	w.Resize(fyne.Size{Width: CWindowWidth, Height: CWindowHeight})

	wbMenu := widget.NewButton("menu", nil)
	var cApp *fyne.Container

	wbPageJournalEntry := widget.NewButton("SIMPLE JOURNAL ENTRY", func() {
		cApp = container.New(layout.NewBorderLayout(nil, wbMenu, nil, nil), wbMenu, FPageJournalEntry())
		w.SetContent(cApp)
	})
	wbPageJournalEntry.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)

	wbPageJournal := widget.NewButton("JOURNAL", func() {
		cApp = container.New(layout.NewBorderLayout(nil, wbMenu, nil, nil), wbMenu, FPageJournal())
		w.SetContent(cApp)
	})
	wbPageJournal.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)

	wbLogin := widget.NewButton("LOGIN", func() {
		w.SetContent(FPageLogin())
	})
	wbLogin.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)

	wsMenu := container.NewVScroll(container.New(layout.NewVBoxLayout(), wbPageJournalEntry, wbPageJournal, wbLogin))

	wbMenu.OnTapped = func() {
		cApp = container.New(layout.NewBorderLayout(nil, wbMenu, nil, nil), wbMenu, wsMenu)
		w.SetContent(cApp)
	}

	cApp = container.New(layout.NewBorderLayout(nil, wbMenu, nil, nil), wbMenu, wsMenu)
	w.SetContent(cApp)

	w.ShowAndRun()
}

func FPageJournal() fyne.CanvasObject {
	keys, journal := FDbRead[SJournal](VDbJournal)
	dates := FConvertByteSliceToTime(keys)

	fcJournal := container.New(layout.NewVBoxLayout(), container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight},
		FLabelToDeiplayJournal("dates"),
		FLabelToDeiplayJournal("IsReverse"),
		FLabelToDeiplayJournal("IsReversed"),
		FLabelToDeiplayJournal("ReverseEntryNumberCompound"),
		FLabelToDeiplayJournal("ReverseEntryNumberSimple"),
		FLabelToDeiplayJournal("EntryNumberCompound"),
		FLabelToDeiplayJournal("EntryNumberSimple"),
		FLabelToDeiplayJournal("Value"),
		FLabelToDeiplayJournal("PriceDebit"),
		FLabelToDeiplayJournal("PriceCredit"),
		FLabelToDeiplayJournal("QuantityDebit"),
		FLabelToDeiplayJournal("QuantityCredit"),
		FLabelToDeiplayJournal("AccountDebit"),
		FLabelToDeiplayJournal("AccountCredit"),
		FLabelToDeiplayJournal("Notes"),
		FLabelToDeiplayJournal("Name"),
		FLabelToDeiplayJournal("Employee"),
		FLabelToDeiplayJournal("TypeOfCompoundEntry"),
	))
	for k1, v1 := range journal {
		fcJournal.Add(container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight},
			FLabelToDeiplayJournal(dates[k1]),
			FLabelToDeiplayJournal(v1.IsReverse),
			FLabelToDeiplayJournal(v1.IsReversed),
			FLabelToDeiplayJournal(v1.ReverseEntryNumberCompound),
			FLabelToDeiplayJournal(v1.ReverseEntryNumberSimple),
			FLabelToDeiplayJournal(v1.EntryNumberCompound),
			FLabelToDeiplayJournal(v1.EntryNumberSimple),
			FLabelToDeiplayJournal(v1.Value),
			FLabelToDeiplayJournal(v1.PriceDebit),
			FLabelToDeiplayJournal(v1.PriceCredit),
			FLabelToDeiplayJournal(v1.QuantityDebit),
			FLabelToDeiplayJournal(v1.QuantityCredit),
			FLabelToDeiplayJournal(v1.AccountDebit),
			FLabelToDeiplayJournal(v1.AccountCredit),
			FLabelToDeiplayJournal(v1.Notes),
			FLabelToDeiplayJournal(v1.Name),
			FLabelToDeiplayJournal(v1.Employee),
			FLabelToDeiplayJournal(v1.TypeOfCompoundEntry),
		))
	}

	return container.NewVScroll(fcJournal)
}

func FLabelToDeiplayJournal(x any) *widget.Label {
	wl := widget.NewLabel(fmt.Sprint(x))
	wl.Resize(fyne.NewSize(10, Vheight))
	return wl
}

func FAllAccountNames() []string {
	var allAccountNames []string
	for _, v1 := range VAccounts {
		if v1.TCostFlowType != "" {
			allAccountNames = append(allAccountNames, v1.TAccountName)
		}
	}
	return allAccountNames
}

func FFcEntryABPQ() *fyne.Container {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, Vheight))

	wsAccount := widget.NewSelectEntry(VAllAccountNames)
	wsAccount.Resize(fyne.NewSize(80/3, Vheight))

	wePrice := widget.NewEntry()
	wePrice.SetPlaceHolder("price")
	wePrice.Resize(fyne.NewSize(80/3, Vheight))

	weQuantity := widget.NewEntry()
	weQuantity.SetPlaceHolder("quantity")
	weQuantity.Resize(fyne.NewSize(80/3, Vheight))

	wbX := widget.NewButton("x", nil)
	wbX.Resize(fyne.NewSize(10, Vheight))

	fc := container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wc, wsAccount, wePrice, weQuantity, wbX)

	wbX.OnTapped = func() {
		index, _ := FFindObject(fc, VfcEntries.Objects)
		VfcEntries.Objects = FRemove(VfcEntries.Objects, index)
	}
	return fc
}

func FPageLogin() *fyne.Container {

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

func FPageJournalEntry() *fyne.Container {

	wlError := widget.NewLabel("")
	wlError.Wrapping = fyne.TextWrapWord
	wlError.Resize(fyne.Size{Width: CWindowWidth, Height: Vheight})

	go func() {
		for range time.Tick(time.Second) {
			_, err := FSaveEntries(false)
			FDisplayTheError(err, wlError)
		}
	}()

	wbOk := widget.NewButton("ok", func() {
		_, err := FSaveEntries(true)
		FDisplayTheError(err, wlError)
		if err == nil {
			FClearEntries(false)
		}
	})
	wbOk.Resize(fyne.Size{Width: wbOk.MinSize().Width, Height: Vheight})

	wbSave := widget.NewButton("save", func() {
		_, err := FSaveEntries(true)
		FDisplayTheError(err, wlError)
	})
	wbSave.Resize(fyne.Size{Width: wbSave.MinSize().Width, Height: Vheight})

	wbClearChecked := widget.NewButton("clear checked", func() {
		FClearEntries(true)
	})
	wbClearChecked.Resize(fyne.Size{Width: wbClearChecked.MinSize().Width, Height: Vheight})

	wbClearNotChecked := widget.NewButton("clear not checked", func() {
		FClearEntries(false)
	})
	wbClearNotChecked.Resize(fyne.Size{Width: wbClearNotChecked.MinSize().Width, Height: Vheight})

	VfcEntries = container.New(layout.NewVBoxLayout(), FFcEntryInfo("note"), FFcEntryInfo("name"), FFcEntryInfo("type"), FFcEntryABPQ())
	wsEntries := container.NewVScroll(VfcEntries)

	wbAdd := widget.NewButton("+", func() {
		VfcEntries.Add(FFcEntryABPQ())
		VfcEntries.Refresh()
	})
	wbAdd.Resize(fyne.Size{Width: wbAdd.MinSize().Width, Height: Vheight})

	fc1 := container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wbOk, wbSave, wbClearChecked, wbClearNotChecked, wbAdd)
	fc2 := container.New(&SStretchVBoxLayout{Width: CWindowWidth, Height: Vheight, ObjectToStertch: wsEntries}, fc1, wsEntries, wlError)

	return fc2
}

func FDisplayTheError(err error, wlError *widget.Label) {
	if err != nil {
		wlError.SetText(err.Error())
	} else {
		wlError.SetText("")
	}
}

func FFcEntryInfo(label string) *fyne.Container {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, Vheight))
	we := widget.NewEntry()
	we.SetPlaceHolder(label)
	we.Resize(fyne.NewSize(80, Vheight))
	wb := widget.NewButton("x", func() {
		we.SetText("")
	})
	wb.Resize(fyne.NewSize(10, Vheight))
	return container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wc, we, wb)
}

func FClearEntries(isChecked bool) {
	for k1, v1 := range VfcEntries.Objects {
		if k1 == 3 {
			break
		}
		if v1.(*fyne.Container).Objects[0].(*widget.Check).Checked == isChecked {
			v1.(*fyne.Container).Objects[1].(*widget.Entry).SetText("")
		}
	}

	k1 := 3
	for k1 < len(VfcEntries.Objects) {
		if VfcEntries.Objects[k1].(*fyne.Container).Objects[0].(*widget.Check).Checked == isChecked {
			VfcEntries.Objects = FRemove(VfcEntries.Objects, k1)
		} else {
			k1++
		}
	}
}

func FSaveEntries(insert bool) ([]SAPQ, error) {
	entryInfo := SEntry{
		TEntryNotes:          VfcEntries.Objects[0].(*fyne.Container).Objects[1].(*widget.Entry).Text,
		TPersonName:          VfcEntries.Objects[1].(*fyne.Container).Objects[1].(*widget.Entry).Text,
		TTypeOfCompoundEntry: VfcEntries.Objects[2].(*fyne.Container).Objects[1].(*widget.Entry).Text,
	}

	var entries []SAPQ
	for _, v1 := range VfcEntries.Objects[3:] {
		price, _ := strconv.ParseFloat(v1.(*fyne.Container).Objects[2].(*widget.Entry).Text, 64)
		quantity, _ := strconv.ParseFloat(v1.(*fyne.Container).Objects[3].(*widget.Entry).Text, 64)

		entries = append(entries, SAPQ{
			TAccountName: v1.(*fyne.Container).Objects[1].(*widget.SelectEntry).Text,
			TPrice:       price,
			TQuantity:    quantity,
		})
	}

	return FSimpleJournalEntry(entries, entryInfo, insert)
}
