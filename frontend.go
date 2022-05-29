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
	Vheight                  float32 = 50
	VfcEntries               *fyne.Container
	VAllLowLevelAccountNames = FAllLowLevelAccountNames()
	VWindow                  fyne.Window
)

func main() {
	defer FDbCloseAll()

	a := app.New()
	VWindow = a.NewWindow(CNameOfTheApp)
	VWindow.Resize(fyne.Size{Width: CWindowWidth, Height: CWindowHeight})

	Vheight = widget.NewLabel("").MinSize().Height

	VWindow.SetContent(FPageLogin())
	VWindow.ShowAndRun()
}

func FPageMenu() fyne.CanvasObject {
	VWindow.SetTitle(CNameOfTheApp + " - " + CMenu)

	wbPageJournalEntry := FSetTheWindow("ENTRY", FPageJournalEntry)
	wbPageJournal := FSetTheWindow("JOURNAL", FPageJournal)
	wbPageAccounts := FSetTheWindow("ACCOUNTS", FPageAccounts)
	wbPageChangePassword := FSetTheWindow("CHANGE PASSWORD", FPageChangePassword)
	wbPageLogin := widget.NewButton("LOGOUT", func() { VWindow.SetContent(FPageLogin()) })
	wbPageLogin.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)

	fc := container.New(layout.NewVBoxLayout(), wbPageJournalEntry, wbPageJournal, wbPageAccounts, wbPageChangePassword, wbPageLogin)
	if VEmployeeName == CManger {
		fc.Add(widget.NewLabel("this below options are just for " + CManger))
		fc.Add(FSetTheWindow("ADD NEW EMPLOYEE", FPageAddNewEmployee))
		fc.Add(FSetTheWindow("CHANGE EMPLOYEE NAME", FPageChangeEmployeeName))
		fc.Add(FSetTheWindow("DELETE EMPLOYEE", FPageDeleteEmployee))
	}
	return container.NewVScroll(fc)
}

func FPageJournalEntry() fyne.CanvasObject {

	wlError := widget.NewLabel("")
	wlError.Wrapping = fyne.TextWrapWord
	wlError.Resize(fyne.Size{Width: CWindowWidth, Height: Vheight})

	go func() {
		for range time.Tick(time.Second / 4) {
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

	return container.New(&SStretchVBoxLayout{Width: CWindowWidth, Height: Vheight, ObjectToStertch: wsEntries}, fc1, wsEntries, wlError)
}

func FPageJournal() fyne.CanvasObject {
	keys, journal := FDbRead[SJournal](VDbJournal)
	dates := FConvertByteSliceToTime(keys)

	fcJournal := container.New(layout.NewVBoxLayout(), container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight},
		FLabelToDisplayJournal("dates"),
		FLabelToDisplayJournal("IsReverse"),
		FLabelToDisplayJournal("IsReversed"),
		FLabelToDisplayJournal("ReverseEntryNumberCompound"),
		FLabelToDisplayJournal("ReverseEntryNumberSimple"),
		FLabelToDisplayJournal("EntryNumberCompound"),
		FLabelToDisplayJournal("EntryNumberSimple"),
		FLabelToDisplayJournal("Value"),
		FLabelToDisplayJournal("PriceDebit"),
		FLabelToDisplayJournal("PriceCredit"),
		FLabelToDisplayJournal("QuantityDebit"),
		FLabelToDisplayJournal("QuantityCredit"),
		FLabelToDisplayJournal("AccountDebit"),
		FLabelToDisplayJournal("AccountCredit"),
		FLabelToDisplayJournal("Notes"),
		FLabelToDisplayJournal("Name"),
		FLabelToDisplayJournal("Employee"),
		FLabelToDisplayJournal("TypeOfCompoundEntry"),
	))
	for k1, v1 := range journal {
		fcJournal.Add(container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight},
			FLabelToDisplayJournal(dates[k1]),
			FLabelToDisplayJournal(v1.IsReverse),
			FLabelToDisplayJournal(v1.IsReversed),
			FLabelToDisplayJournal(v1.ReverseEntryNumberCompound),
			FLabelToDisplayJournal(v1.ReverseEntryNumberSimple),
			FLabelToDisplayJournal(v1.EntryNumberCompound),
			FLabelToDisplayJournal(v1.EntryNumberSimple),
			FLabelToDisplayJournal(v1.Value),
			FLabelToDisplayJournal(v1.PriceDebit),
			FLabelToDisplayJournal(v1.PriceCredit),
			FLabelToDisplayJournal(v1.QuantityDebit),
			FLabelToDisplayJournal(v1.QuantityCredit),
			FLabelToDisplayJournal(v1.AccountDebit),
			FLabelToDisplayJournal(v1.AccountCredit),
			FLabelToDisplayJournal(v1.Notes),
			FLabelToDisplayJournal(v1.Name),
			FLabelToDisplayJournal(v1.Employee),
			FLabelToDisplayJournal(v1.TypeOfCompoundEntry),
		))
	}

	return container.NewVScroll(fcJournal)
}

func FPageAccounts() fyne.CanvasObject {

	fcAccounts := container.New(layout.NewVBoxLayout(), container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}))

	for _, v1 := range VAccounts {
		fcAccounts.Add(container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight},
			FLabelToDisplayJournal(v1.TIsCredit),
			FLabelToDisplayJournal(v1.TCostFlowType),
			FLabelToDisplayJournal(v1.TAccountName),
			FLabelToDisplayJournal(v1.TAccountNotes),
			FLabelToDisplayJournal(v1.TAccountImage),
			FLabelToDisplayJournal(v1.TAccountBarcode),
			FLabelToDisplayJournal(v1.TAccountNumber),
			FLabelToDisplayJournal(v1.TAccountLevels),
			FLabelToDisplayJournal(v1.TAccountFathersName),
		))
	}

	wsAccounts := container.NewVScroll(fcAccounts)

	wbIsCredit := widget.NewSelect([]string{CCredit, CDebit}, nil)
	wbIsCredit.Resize(fyne.NewSize(1, 0))
	wbIsCredit.SetSelected(CCredit)

	wbCostFlowType := widget.NewSelect([]string{CHighLevelAccount, CWma, CFifo, CLifo}, nil)
	wbCostFlowType.Resize(fyne.NewSize(1, 0))
	wbCostFlowType.SetSelected(CHighLevelAccount)

	wbAccountName := widget.NewEntry()
	wbAccountName.Resize(fyne.NewSize(1, 0))
	wbAccountName.SetPlaceHolder("Name")

	wbAccountNotes := widget.NewEntry()
	wbAccountNotes.Resize(fyne.NewSize(1, 0))
	wbAccountNotes.SetPlaceHolder("Notes")

	wbAccountImage := widget.NewEntry()
	wbAccountImage.Resize(fyne.NewSize(1, 0))
	wbAccountImage.SetPlaceHolder("Image")

	wbAccountBarcode := widget.NewEntry()
	wbAccountBarcode.Resize(fyne.NewSize(1, 0))
	wbAccountBarcode.SetPlaceHolder("Barcode")

	wbAccountNumber := widget.NewEntry()
	wbAccountNumber.Resize(fyne.NewSize(1, 0))
	wbAccountNumber.SetPlaceHolder("Number")

	fc1 := container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight},
		wbIsCredit,
		wbCostFlowType,
		wbAccountName,
		wbAccountNotes,
		wbAccountImage,
		wbAccountBarcode,
		wbAccountNumber,
	)
	return container.New(&SStretchVBoxLayout{Width: CWindowWidth, Height: Vheight, ObjectToStertch: wsAccounts}, fc1, wsAccounts)
}

func FPageChangePassword() fyne.CanvasObject {
	wePassword := widget.NewPasswordEntry()
	wePassword.SetPlaceHolder("password")

	wlError := widget.NewLabel("")
	wlError.Alignment = fyne.TextAlignCenter

	wb := widget.NewButton("ok", func() {
		wlError.SetText("")

		if wePassword.Text == "" {
			wlError.SetText(wePassword.PlaceHolder + " is empty")
			return
		}

		FDbUpdate(VDbEmployees, []byte(VEmployeeName), wePassword.Text)
		wePassword.SetText("")
		VWindow.SetContent(FPageMenu())
	})

	return container.New(layout.NewVBoxLayout(), wePassword, wlError, wb)
}

func FPageLogin() fyne.CanvasObject {
	VWindow.SetTitle(CNameOfTheApp + " - " + "LOGIN")
	FDbCloseAll()

	companies, _ := FFilesName(CPathDataBase)

	weCompanyName := widget.NewSelectEntry(companies)
	weCompanyName.SetPlaceHolder("company name")

	weEmployeeName := widget.NewSelectEntry(nil)
	weEmployeeName.SetPlaceHolder("employee name")

	weCompanyName.OnChanged = func(s string) {
		_, isIn := FFind(weCompanyName.Text, companies)
		if isIn && weCompanyName.Text != "" {
			VCompanyName = weCompanyName.Text
			employees, _ := FDbOpenAndReadEmployees()
			weEmployeeName.SetOptions(employees)
		} else {
			FDbCloseAll()
			weEmployeeName.SetOptions(nil)
		}
	}

	wePassword := widget.NewPasswordEntry()
	wePassword.SetPlaceHolder("password")

	wlError := widget.NewLabel("")
	wlError.Alignment = fyne.TextAlignCenter

	wbLogin := widget.NewButton("login", func() {
		wlError.SetText("")

		_, isIn := FFind(weCompanyName.Text, companies)
		if !isIn {
			wlError.SetText(weCompanyName.PlaceHolder + " is wrong")
			return
		}

		employees, passwords := FDbOpenAndReadEmployees()

		indexOfEmployee, isIn := FFind(weEmployeeName.Text, employees)
		if !isIn {
			wlError.SetText(weEmployeeName.PlaceHolder + " is wrong")
			return
		}

		if passwords[indexOfEmployee] != wePassword.Text {
			wlError.SetText(wePassword.PlaceHolder + " is wrong")
			return
		}

		VCompanyName = weCompanyName.Text
		VEmployeeName = weEmployeeName.Text

		FDbOpenAll()
		VWindow.SetContent(FPageMenu())
	})

	wbCreateNewCompany := widget.NewButton("create new company", func() {
		wlError.SetText("")
		weEmployeeName.SetText(CManger)

		_, isIn := FFind(weCompanyName.Text, companies)
		if isIn {
			wlError.SetText(weCompanyName.PlaceHolder + " is used")
			return
		}

		if weCompanyName.Text == "" {
			wlError.SetText(weCompanyName.PlaceHolder + " is empty")
			return
		}

		if wePassword.Text == "" {
			wlError.SetText(wePassword.PlaceHolder + " is empty")
			return
		}

		VCompanyName = weCompanyName.Text
		VEmployeeName = weEmployeeName.Text

		VDbEmployees = FDbOpen(VDbEmployees, CPathDataBase+VCompanyName+CPathEmployees)
		FDbUpdate(VDbEmployees, []byte(VEmployeeName), wePassword.Text)

		FDbOpenAll()
		VWindow.SetContent(FPageMenu())
	})

	return container.New(layout.NewVBoxLayout(), weCompanyName, weEmployeeName, wePassword, wlError, wbLogin, wbCreateNewCompany)
}

func FPageAddNewEmployee() fyne.CanvasObject {
	employees, _ := FDbOpenAndReadEmployees()

	weEmployeeName := widget.NewSelectEntry(employees)
	weEmployeeName.SetPlaceHolder("employee name")

	wePassword := widget.NewPasswordEntry()
	wePassword.SetPlaceHolder("password")

	wlError := widget.NewLabel("")
	wlError.Alignment = fyne.TextAlignCenter

	wb := widget.NewButton("ok", func() {
		wlError.SetText("")

		if weEmployeeName.Text == "" {
			wlError.SetText(weEmployeeName.PlaceHolder + " is empty")
			return
		}

		if weEmployeeName.Text == CManger {
			wlError.SetText("you can't use " + CManger + " as " + weEmployeeName.PlaceHolder)
			return
		}

		employees, _ := FDbOpenAndReadEmployees()
		_, isIn := FFind(weEmployeeName.Text, employees)
		if isIn {
			wlError.SetText(weEmployeeName.PlaceHolder + " is used")
			return
		}

		if wePassword.Text == "" {
			wlError.SetText(wePassword.PlaceHolder + " is empty")
			return
		}

		FDbUpdate(VDbEmployees, []byte(weEmployeeName.Text), wePassword.Text)
		weEmployeeName.SetText("")
		wePassword.SetText("")
	})

	return container.New(layout.NewVBoxLayout(), weEmployeeName, wePassword, wlError, wb)
}

func FPageChangeEmployeeName() fyne.CanvasObject {
	employees, _ := FDbOpenAndReadEmployees()

	weEmployeeName := widget.NewSelectEntry(employees)
	weEmployeeName.SetPlaceHolder("employee name")

	weNewEmployeeName := widget.NewEntry()
	weNewEmployeeName.SetPlaceHolder("new employee name")

	wlError := widget.NewLabel("")
	wlError.Alignment = fyne.TextAlignCenter

	wb := widget.NewButton("ok", func() {
		wlError.SetText("")

		if weEmployeeName.Text == CManger {
			wlError.SetText("you can't use " + CManger + " as " + weEmployeeName.PlaceHolder)
			return
		}

		if weNewEmployeeName.Text == "" {
			wlError.SetText(weNewEmployeeName.PlaceHolder + " is empty")
			return
		}

		if weNewEmployeeName.Text == CManger {
			wlError.SetText("you can't use " + CManger + " as " + weNewEmployeeName.PlaceHolder)
			return
		}

		employees, passwords := FDbOpenAndReadEmployees()
		indexOfEmployee, isIn := FFind(weEmployeeName.Text, employees)
		if !isIn {
			wlError.SetText(weEmployeeName.PlaceHolder + " is wrong")
			return
		}

		employees, _ = FDbOpenAndReadEmployees()
		_, isIn = FFind(weNewEmployeeName.Text, employees)
		if isIn {
			wlError.SetText(weNewEmployeeName.PlaceHolder + " is used")
			return
		}

		FDbUpdate(VDbEmployees, []byte(weEmployeeName.Text), passwords[indexOfEmployee])
		FChangeEmployeeName(weEmployeeName.Text, weNewEmployeeName.Text)
		weEmployeeName.SetText("")
		weNewEmployeeName.SetText("")
	})

	return container.New(layout.NewVBoxLayout(), weEmployeeName, weNewEmployeeName, wlError, wb)
}

func FPageDeleteEmployee() fyne.CanvasObject {
	employees, _ := FDbOpenAndReadEmployees()

	weEmployeeName := widget.NewSelectEntry(employees)
	weEmployeeName.SetPlaceHolder("employee name")

	wlError := widget.NewLabel("")
	wlError.Alignment = fyne.TextAlignCenter

	wb := widget.NewButton("ok", func() {
		wlError.SetText("")

		if weEmployeeName.Text == CManger {
			wlError.SetText("you can't use " + CManger + " as " + weEmployeeName.PlaceHolder)
			return
		}

		employees, _ := FDbOpenAndReadEmployees()
		_, isIn := FFind(weEmployeeName.Text, employees)
		if !isIn {
			wlError.SetText(weEmployeeName.PlaceHolder + " is wrong")
			return
		}

		FDbDelete(VDbEmployees, []byte(weEmployeeName.Text))
		weEmployeeName.SetText("")
	})

	return container.New(layout.NewVBoxLayout(), weEmployeeName, wlError, wb)
}

func FDbOpenAndReadEmployees() ([]string, []string) {
	VDbEmployees = FDbOpen(VDbEmployees, CPathDataBase+VCompanyName+CPathEmployees)
	keys, passwords := FDbRead[string](VDbEmployees)
	employees := FConvertFromByteSliceToString(keys)
	return employees, passwords
}

func FLabelToDisplayJournal(x any) *widget.Label {
	wl := widget.NewLabel(fmt.Sprint(x))
	wl.Resize(fyne.NewSize(10, Vheight))
	return wl
}

func FAllLowLevelAccountNames() []string {
	var allAccountNames []string
	for _, v1 := range VAccounts {
		if v1.TCostFlowType != CHighLevelAccount {
			allAccountNames = append(allAccountNames, v1.TAccountName)
		}
	}
	return allAccountNames
}

func FFcEntryABPQ() *fyne.Container {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, Vheight))

	wsAccount := widget.NewSelectEntry(VAllLowLevelAccountNames)
	wsAccount.Resize(fyne.NewSize(80/3, Vheight))

	wePrice := widget.NewEntry()
	wePrice.Resize(fyne.NewSize(80/3, Vheight))
	wePrice.SetPlaceHolder("price")

	weQuantity := widget.NewEntry()
	weQuantity.Resize(fyne.NewSize(80/3, Vheight))
	weQuantity.SetPlaceHolder("quantity")
	weQuantity.SetText("1")

	wbX := widget.NewButton("x", nil)
	wbX.Resize(fyne.NewSize(10, Vheight))

	fc := container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wc, wsAccount, wePrice, weQuantity, wbX)

	wbX.OnTapped = func() {
		index, _ := FFindObject(fc, VfcEntries.Objects)
		VfcEntries.Objects = FRemove(VfcEntries.Objects, index)
	}
	return fc
}

func FDisplayTheError(errorMessage error, wlError *widget.Label) {
	if errorMessage != nil {
		wlError.SetText(errorMessage.Error())
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
	wb := widget.NewButton("x", func() { we.SetText("") })
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

func FSetTheWindow(label string, page func() fyne.CanvasObject) *widget.Button {
	wbMenu := widget.NewButton(CMenu, func() {
		VWindow.SetContent(FPageMenu())
	})
	wb := widget.NewButton(label, func() {
		VWindow.SetTitle(CNameOfTheApp + " - " + label)
		VWindow.SetContent(container.New(layout.NewBorderLayout(nil, wbMenu, nil, nil), wbMenu, page()))
	})
	wb.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)
	return wb
}
