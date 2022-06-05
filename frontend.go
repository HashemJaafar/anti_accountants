package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	CNameOfTheApp           = "ANTI ACCOUNTANTS"
	CPageMenu               = "MENU"
	CPageInvoiceEntry       = "INVOICE ENTRY"
	CPageJournalEntry       = "JOURNAL ENTRY"
	CPageJournalDraft       = "JOURNAL DRAFT"
	CPageJournal            = "JOURNAL"
	CPageAccounts           = "ACCOUNTS"
	CPageChangePassword     = "CHANGE PASSWORD"
	CPageLogin              = "LOGIN"
	CPageAddNewEmployee     = "ADD NEW EMPLOYEE"
	CPageChangeEmployeeName = "CHANGE EMPLOYEE NAME"
	CPageDeleteEmployee     = "DELETE EMPLOYEE"
	CPageChangeCompanyName  = "CHANGE COMPANY NAME"
	CPageDeleteCompany      = "DELETE COMPANY"

	CManger = "Manger"

	CWindowWidth  = 500
	CWindowHeight = 500
)

var (
	Vheight           float32 = 50
	VfcJournalEntries *fyne.Container
	VfcInvoiceEntries *fyne.Container
	VWindow           fyne.Window
)

func main() {
	defer FDbCloseAll()

	a := app.New()
	VWindow = a.NewWindow("")
	VWindow.Resize(fyne.Size{Width: CWindowWidth, Height: CWindowHeight})

	Vheight = widget.NewLabel("").MinSize().Height

	VWindow.SetContent(FPageLogin())
	VWindow.ShowAndRun()
}

func FPageMenu() fyne.CanvasObject {
	FSetTitle(CPageMenu)

	wbPageInvoiceEntry := FSetTheWindow(CPageInvoiceEntry, FPageInvoiceEntry)
	wbPageJournalEntry := FSetTheWindow(CPageJournalEntry, FPageJournalEntry)
	wbPageJournalDraft := FSetTheWindow(CPageJournalDraft, FPageJournalDraft)
	wbPageJournal := FSetTheWindow(CPageJournal, FPageJournal)
	wbPageAccounts := FSetTheWindow(CPageAccounts, FPageAccounts)
	wbPageChangePassword := FSetTheWindow(CPageChangePassword, FPageChangePassword)
	wbPageLogin := widget.NewButton("LOGOUT", func() { VWindow.SetContent(FPageLogin()) })
	wbPageLogin.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)

	fc := container.NewVBox(wbPageInvoiceEntry, wbPageJournalEntry, wbPageJournalDraft, wbPageJournal, wbPageAccounts, wbPageChangePassword, wbPageLogin)
	if VEmployeeName == CManger {
		fc.Add(widget.NewLabel("this below options are just for " + CManger))
		fc.Add(FSetTheWindow(CPageAddNewEmployee, FPageAddNewEmployee))
		fc.Add(FSetTheWindow(CPageChangeEmployeeName, FPageChangeEmployeeName))
		fc.Add(FSetTheWindow(CPageDeleteEmployee, FPageDeleteEmployee))
		fc.Add(FSetTheWindow(CPageChangeCompanyName, FPageChangeCompanyName))
		fc.Add(FSetTheWindow(CPageDeleteCompany, FPageDeleteCompany))
	}
	return container.NewVScroll(fc)
}

func FPageInvoiceEntry() fyne.CanvasObject {
	FSetTitle(CPageInvoiceEntry)
	VfcInvoiceEntries = container.NewVBox(FFcEntryInfo("draft"), FFcEntryInfo("note"), FFcEntryInfo("name"), FFcEntryInfo("type"))
	return FSetTheEntries(SPageInvoiceEntries{}, VfcInvoiceEntries)
}

func FPageJournalEntry() fyne.CanvasObject {
	FSetTitle(CPageJournalEntry)
	VfcJournalEntries = container.NewVBox(FFcEntryInfo("draft"), FFcEntryInfo("note"), FFcEntryInfo("name"), FFcEntryInfo("type"))
	return FSetTheEntries(SPageJournalEntries{}, VfcJournalEntries)
}

func FPageJournalDraft() fyne.CanvasObject {
	FSetTitle(CPageJournalDraft)

	SetTheDraftedEntry := func(values []fyne.Container, k1 int) {
		FSetTitle(CPageJournalEntry)
		VfcJournalEntries = &values[k1]
		VWindow.SetContent(FSetTheEntries(SPageJournalEntries{}, VfcJournalEntries))
	}

	DeleteLine := func(fc *fyne.Container, wbName *widget.Button) {
		for k2, v2 := range fc.Objects[1:] {
			if v2.(*fyne.Container).Objects[0].(*widget.Button).Text == wbName.Text {
				FDbDelete(VDbJournalDrafts, []byte(wbName.Text))
				fc.Objects = FRemove(fc.Objects, k2+1)
				fc.Refresh()
			}
		}
	}

	keys, values := FDbRead[fyne.Container](VDbJournalDrafts)
	fc := container.NewVBox(widget.NewLabel(""))
	for k1, v1 := range keys {

		wbName := widget.NewButton(string(v1), func() { SetTheDraftedEntry(values, k1) })
		wbName.Resize(fyne.NewSize(60, Vheight))

		wbEdit := widget.NewButton("Delete And Open", func() {
			DeleteLine(fc, wbName)
			SetTheDraftedEntry(values, k1)
		})
		wbEdit.Resize(fyne.NewSize(30, Vheight))

		wbX := widget.NewButton("x", func() { DeleteLine(fc, wbName) })
		wbX.Resize(fyne.NewSize(10, Vheight))

		fc.Add(container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wbName, wbEdit, wbX))
	}
	return container.NewVScroll(fc)
}

func FPageJournal() fyne.CanvasObject {
	FSetTitle(CPageJournal)

	keys, journal := FDbRead[SJournal](VDbJournal)
	dates := FConvertByteSliceToTime(keys)

	fcJournal := container.NewVBox(container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight},
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
			FLabelToDisplayJournal(v1.DebitPrice),
			FLabelToDisplayJournal(v1.CreditPrice),
			FLabelToDisplayJournal(v1.DebitQuantity),
			FLabelToDisplayJournal(v1.CreditQuantity),
			FLabelToDisplayJournal(v1.DebitAccountName),
			FLabelToDisplayJournal(v1.CreditAccountName),
			FLabelToDisplayJournal(v1.Notes),
			FLabelToDisplayJournal(v1.Name),
			FLabelToDisplayJournal(v1.Employee),
			FLabelToDisplayJournal(v1.TypeOfCompoundEntry),
		))
	}

	return container.NewVScroll(fcJournal)
}

func FPageAccounts() fyne.CanvasObject {
	FSetTitle(CPageAccounts)

	fcAccounts := container.NewVBox(container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}))

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
	FSetTitle(CPageChangePassword)

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
	})

	return container.NewVBox(wePassword, wlError, wb)
}

func FPageLogin() fyne.CanvasObject {
	FSetTitle(CPageLogin)
	FDbCloseAll()

	companies, _ := FFilesName(CPathDataBase)

	weCompanyName := widget.NewSelectEntry(companies)
	weCompanyName.SetPlaceHolder("company name")

	weEmployeeName := widget.NewSelectEntry(nil)
	weEmployeeName.SetPlaceHolder("employee name")

	weCompanyName.OnChanged = func(s string) {
		_, isIn := FFind(weCompanyName.Text, companies)
		if isIn {
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

	return container.NewVBox(weCompanyName, weEmployeeName, wePassword, wlError, wbLogin, wbCreateNewCompany)
}

func FPageAddNewEmployee() fyne.CanvasObject {
	FSetTitle(CPageAddNewEmployee)

	weEmployeeName := widget.NewEntry()
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

	return container.NewVBox(weEmployeeName, wePassword, wlError, wb)
}

func FPageChangeEmployeeName() fyne.CanvasObject {
	FSetTitle(CPageChangeEmployeeName)

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

		FDbDelete(VDbEmployees, []byte(weEmployeeName.Text))
		FDbUpdate(VDbEmployees, []byte(weNewEmployeeName.Text), passwords[indexOfEmployee])
		FChangeEmployeeName(weEmployeeName.Text, weNewEmployeeName.Text)
		weEmployeeName.SetText("")
		weNewEmployeeName.SetText("")
	})

	return container.NewVBox(weEmployeeName, weNewEmployeeName, wlError, wb)
}

func FPageDeleteEmployee() fyne.CanvasObject {
	FSetTitle(CPageDeleteEmployee)

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

	return container.NewVBox(weEmployeeName, wlError, wb)
}

func FPageChangeCompanyName() fyne.CanvasObject {
	FSetTitle(CPageChangeCompanyName)

	weCompanyName := widget.NewEntry()
	weCompanyName.SetPlaceHolder("company name")

	wlError := widget.NewLabel("")
	wlError.Alignment = fyne.TextAlignCenter

	wb := widget.NewButton("ok", func() {
		wlError.SetText("")

		if weCompanyName.Text == "" {
			wlError.SetText(weCompanyName.PlaceHolder + " is empty")
			return
		}

		companies, _ := FFilesName(CPathDataBase)
		_, isIn := FFind(weCompanyName.Text, companies)
		if isIn {
			wlError.SetText(weCompanyName.PlaceHolder + " is used")
			return
		}

		os.Rename(CPathDataBase+VCompanyName, CPathDataBase+weCompanyName.Text)

		VCompanyName = weCompanyName.Text
		weCompanyName.SetText("")
	})

	return container.NewVBox(weCompanyName, wlError, wb)
}

func FPageDeleteCompany() fyne.CanvasObject {
	FSetTitle(CPageDeleteCompany)

	wePassword := widget.NewPasswordEntry()
	wePassword.SetPlaceHolder("password")

	wePassword1 := widget.NewPasswordEntry()
	wePassword1.SetPlaceHolder("password")

	wePassword2 := widget.NewPasswordEntry()
	wePassword2.SetPlaceHolder("password")

	wePassword3 := widget.NewPasswordEntry()
	wePassword3.SetPlaceHolder("password")

	wePassword4 := widget.NewPasswordEntry()
	wePassword4.SetPlaceHolder("password")

	wePassword5 := widget.NewPasswordEntry()
	wePassword5.SetPlaceHolder("password")

	wePassword6 := widget.NewPasswordEntry()
	wePassword6.SetPlaceHolder("password")

	wePassword7 := widget.NewPasswordEntry()
	wePassword7.SetPlaceHolder("password")

	wePassword8 := widget.NewPasswordEntry()
	wePassword8.SetPlaceHolder("password")

	wePassword9 := widget.NewPasswordEntry()
	wePassword9.SetPlaceHolder("password")

	wlError := widget.NewLabel("")
	wlError.Alignment = fyne.TextAlignCenter

	wb := widget.NewButton("delete "+VCompanyName+" company", func() {
		wlError.SetText("")

		employees, passwords := FDbOpenAndReadEmployees()

		indexOfEmployee, _ := FFind(CManger, employees)
		if passwords[indexOfEmployee] != wePassword.Text {
			wlError.SetText(wePassword.PlaceHolder + " is wrong")
			return
		}

		if wePassword.Text != wePassword1.Text ||
			wePassword.Text != wePassword2.Text ||
			wePassword.Text != wePassword3.Text ||
			wePassword.Text != wePassword4.Text ||
			wePassword.Text != wePassword5.Text ||
			wePassword.Text != wePassword6.Text ||
			wePassword.Text != wePassword7.Text ||
			wePassword.Text != wePassword8.Text ||
			wePassword.Text != wePassword9.Text {
			wlError.SetText("not all the passwords are the same")
			return
		}

		fmt.Println(os.RemoveAll(CPathDataBase + VCompanyName))
		VWindow.SetContent(FPageLogin())
	})

	return container.NewVScroll(container.NewVBox(wePassword, wePassword1, wePassword2, wePassword3, wePassword4,
		wePassword5, wePassword6, wePassword7, wePassword8, wePassword9, wlError, wb))
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

func FDisplayTheError(errorMessage error, wlError *widget.Label) {
	if errorMessage != nil {
		wlError.SetText(errorMessage.Error())
	} else {
		wlError.SetText("")
	}
}

func FFcEntryInfo(placeHolder string) fyne.CanvasObject {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, Vheight))
	we := widget.NewEntry()
	we.SetPlaceHolder(placeHolder)
	we.Resize(fyne.NewSize(80, Vheight))
	wb := widget.NewButton("x", func() { we.SetText("") })
	wb.Resize(fyne.NewSize(10, Vheight))
	return container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wc, we, wb)
}

func FSetTheWindow(label string, page func() fyne.CanvasObject) fyne.CanvasObject {
	wb := widget.NewButton(label, func() {
		wbMenu := widget.NewButton(CPageMenu, func() { VWindow.SetContent(FPageMenu()) })
		VWindow.SetContent(container.New(layout.NewBorderLayout(nil, wbMenu, nil, nil), wbMenu, page()))
	})
	wb.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)
	return wb
}

func FSetTitle(Title string) {
	if Title != CPageLogin {
		VWindow.SetTitle(CNameOfTheApp + " - " + VCompanyName + " - " + VEmployeeName + " - " + Title)
	} else {
		VWindow.SetTitle(CNameOfTheApp + " - " + Title)
	}
}

func FSetTheEntries(page IPage, fcEntries *fyne.Container) fyne.CanvasObject {

	wlError := widget.NewLabel("")
	wlError.Wrapping = fyne.TextWrapWord
	wlError.Resize(fyne.Size{Width: CWindowWidth, Height: Vheight})

	go func() {
		for range time.Tick(time.Second / 4) {
			err := page.FSave(false)
			FDisplayTheError(err, wlError)
			if len(fcEntries.Objects) == page.FTheIndexOfTheFirstEntry() {
				fcEntries.Add(page.FNewLine())
			}
		}
	}()

	wbOk := widget.NewButton("ok", func() {
		err := page.FSave(true)
		FDisplayTheError(err, wlError)
		if err == nil {
			page.FClear(false)
		}
	})
	wbOk.Resize(fyne.Size{Width: wbOk.MinSize().Width, Height: Vheight})

	wbSave := widget.NewButton("save", func() {
		err := page.FSave(true)
		FDisplayTheError(err, wlError)
	})
	wbSave.Resize(fyne.Size{Width: wbSave.MinSize().Width, Height: Vheight})

	wbClearChecked := widget.NewButton("clear checked", func() { page.FClear(true) })
	wbClearChecked.Resize(fyne.Size{Width: wbClearChecked.MinSize().Width, Height: Vheight})

	wbClearNotChecked := widget.NewButton("clear not checked", func() { page.FClear(false) })
	wbClearNotChecked.Resize(fyne.Size{Width: wbClearNotChecked.MinSize().Width, Height: Vheight})

	wbDraft := widget.NewButton("draft", func() { page.FDraft(fcEntries.Objects[0].(*fyne.Container).Objects[1].(*widget.Entry).Text) })
	wbDraft.Resize(fyne.Size{Width: wbDraft.MinSize().Width, Height: Vheight})

	fcEntries.Add(page.FNewLine())
	log.Println("fcEntries:", fcEntries)
	log.Println("VfcJournalEntries:", VfcJournalEntries)
	log.Println("VfcInvoiceEntries:", VfcInvoiceEntries)
	wsEntries := container.NewVScroll(fcEntries)

	wbAdd := widget.NewButton("+", func() {
		fcEntries.Add(page.FNewLine())
		fcEntries.Refresh()
	})
	wbAdd.Resize(fyne.Size{Width: wbAdd.MinSize().Width, Height: Vheight})

	fc1 := container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wbOk, wbSave, wbClearChecked, wbClearNotChecked, wbDraft, wbAdd)

	return container.New(&SStretchVBoxLayout{Width: CWindowWidth, Height: Vheight, ObjectToStertch: wsEntries}, fc1, wsEntries, wlError)
}

type SPageInvoiceEntries struct{}

func (SPageInvoiceEntries) FTheIndexOfTheFirstEntry() int { return 4 }

func (SPageInvoiceEntries) FAllAccounts(account string) []string {
	var allAccountNames []string
	_, autoCompletion := FDbRead[SAutoCompletion](VDbAutoCompletionEntries)
	for _, v1 := range autoCompletion {
		if strings.Contains(v1.TAccountName, account) {
			allAccountNames = append(allAccountNames, v1.TAccountName)
		}
	}
	return allAccountNames
}

func (s SPageInvoiceEntries) FSave(insert bool) error {
	entryInfo := SEntry{
		Notes:               VfcInvoiceEntries.Objects[0].(*fyne.Container).Objects[1].(*widget.Entry).Text,
		Name:                VfcInvoiceEntries.Objects[1].(*fyne.Container).Objects[1].(*widget.Entry).Text,
		TypeOfCompoundEntry: VfcInvoiceEntries.Objects[2].(*fyne.Container).Objects[1].(*widget.Entry).Text,
	}

	var entries1 []SAPQ
	for _, v1 := range VfcInvoiceEntries.Objects[s.FTheIndexOfTheFirstEntry():] {
		quantity, _ := strconv.ParseFloat(v1.(*fyne.Container).Objects[3].(*widget.Entry).Text, 64)

		entries1 = append(entries1, SAPQ{
			TAccountName: v1.(*fyne.Container).Objects[1].(*widget.SelectEntry).Text,
			TQuantity:    quantity,
		})
	}

	_, err := FInvoiceJournalEntry("", 0, 0, entries1, entryInfo, insert)
	return err
}

func (s SPageInvoiceEntries) FClear(isChecked bool) {
	for k1, v1 := range VfcInvoiceEntries.Objects {
		if k1 == s.FTheIndexOfTheFirstEntry() {
			break
		}
		if v1.(*fyne.Container).Objects[0].(*widget.Check).Checked == isChecked {
			v1.(*fyne.Container).Objects[1].(*widget.Entry).SetText("")
		}
	}

	k1 := s.FTheIndexOfTheFirstEntry()
	for k1 < len(VfcInvoiceEntries.Objects) {
		if VfcInvoiceEntries.Objects[k1].(*fyne.Container).Objects[0].(*widget.Check).Checked == isChecked {
			VfcInvoiceEntries.Objects = FRemove(VfcInvoiceEntries.Objects, k1)
		} else {
			k1++
		}
	}
}

func (s SPageInvoiceEntries) FNewLine() fyne.CanvasObject {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, Vheight))

	wsAccount := widget.NewSelectEntry(s.FAllAccounts(""))
	wsAccount.Resize(fyne.NewSize(80/3, Vheight))
	wsAccount.SetPlaceHolder("account")
	wsAccount.OnChanged = func(account string) {
		wsAccount.SetOptions(s.FAllAccounts(account))
	}

	wlPrice := widget.NewLabel("")
	wlPrice.Resize(fyne.NewSize(80/3, Vheight))

	weQuantity := widget.NewEntry()
	weQuantity.Resize(fyne.NewSize(80/3, Vheight))
	weQuantity.SetPlaceHolder("quantity")
	weQuantity.SetText(fmt.Sprint(1))

	wbX := widget.NewButton("x", nil)
	wbX.Resize(fyne.NewSize(10, Vheight))

	fc := container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wc, wsAccount, wlPrice, weQuantity, wbX)

	wbX.OnTapped = func() {
		index, _ := FFindObject(fc, VfcInvoiceEntries.Objects)
		VfcInvoiceEntries.Objects = FRemove(VfcInvoiceEntries.Objects, index)
	}
	return fc
}

func (SPageInvoiceEntries) FDraft(nameOfTheDraft string) {
	FDbUpdate(VDbJournalDrafts, []byte(nameOfTheDraft), VfcInvoiceEntries)
}

type SPageJournalEntries struct{}

func (SPageJournalEntries) FTheIndexOfTheFirstEntry() int { return 4 }

func (SPageJournalEntries) FAllAccounts(account string) []string {
	var allAccountNames []string
	for _, v1 := range VAccounts {
		if v1.TCostFlowType != CHighLevelAccount && strings.Contains(v1.TAccountName, account) {
			allAccountNames = append(allAccountNames, v1.TAccountName)
		}
	}
	return allAccountNames
}

func (s SPageJournalEntries) FSave(insert bool) error {
	entryInfo := SEntry{
		Notes:               VfcJournalEntries.Objects[1].(*fyne.Container).Objects[1].(*widget.Entry).Text,
		Name:                VfcJournalEntries.Objects[2].(*fyne.Container).Objects[1].(*widget.Entry).Text,
		TypeOfCompoundEntry: VfcJournalEntries.Objects[3].(*fyne.Container).Objects[1].(*widget.Entry).Text,
	}

	var entries1 []SAPQ
	for _, v1 := range VfcJournalEntries.Objects[s.FTheIndexOfTheFirstEntry():] {
		price, _ := strconv.ParseFloat(v1.(*fyne.Container).Objects[2].(*widget.Entry).Text, 64)
		quantity, _ := strconv.ParseFloat(v1.(*fyne.Container).Objects[3].(*widget.Entry).Text, 64)

		entries1 = append(entries1, SAPQ{
			TAccountName: v1.(*fyne.Container).Objects[1].(*widget.SelectEntry).Text,
			TPrice:       price,
			TQuantity:    quantity,
		})
	}

	_, err := FSimpleJournalEntry(entries1, entryInfo, insert)
	return err
}

func (s SPageJournalEntries) FClear(isChecked bool) {
	for k1, v1 := range VfcJournalEntries.Objects {
		if k1 == s.FTheIndexOfTheFirstEntry() {
			break
		}
		if v1.(*fyne.Container).Objects[0].(*widget.Check).Checked == isChecked {
			v1.(*fyne.Container).Objects[1].(*widget.Entry).SetText("")
		}
	}

	k1 := s.FTheIndexOfTheFirstEntry()
	for k1 < len(VfcJournalEntries.Objects) {
		if VfcJournalEntries.Objects[k1].(*fyne.Container).Objects[0].(*widget.Check).Checked == isChecked {
			VfcJournalEntries.Objects = FRemove(VfcJournalEntries.Objects, k1)
		} else {
			k1++
		}
	}
}

func (s SPageJournalEntries) FNewLine() fyne.CanvasObject {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, Vheight))

	wsAccount := widget.NewSelectEntry(s.FAllAccounts(""))
	wsAccount.Resize(fyne.NewSize(80/3, Vheight))
	wsAccount.SetPlaceHolder("account")
	wsAccount.OnChanged = func(account string) {
		wsAccount.SetOptions(s.FAllAccounts(account))
	}

	wePrice := widget.NewEntry()
	wePrice.Resize(fyne.NewSize(80/3, Vheight))
	wePrice.SetPlaceHolder("price")
	wePrice.SetText(fmt.Sprint(1))

	weQuantity := widget.NewEntry()
	weQuantity.Resize(fyne.NewSize(80/3, Vheight))
	weQuantity.SetPlaceHolder("quantity")
	weQuantity.SetText(fmt.Sprint(1))

	wbX := widget.NewButton("x", nil)
	wbX.Resize(fyne.NewSize(10, Vheight))

	fc := container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wc, wsAccount, wePrice, weQuantity, wbX)

	wbX.OnTapped = func() {
		index, _ := FFindObject(fc, VfcJournalEntries.Objects)
		VfcJournalEntries.Objects = FRemove(VfcJournalEntries.Objects, index)
	}
	return fc
}

func (SPageJournalEntries) FDraft(nameOfTheDraft string) {
	FDbUpdate(VDbJournalDrafts, []byte(nameOfTheDraft), VfcJournalEntries)
}

type IPage interface {
	FTheIndexOfTheFirstEntry() int
	FSave(bool) error
	FClear(bool)
	FNewLine() fyne.CanvasObject
	FDraft(string)
}
