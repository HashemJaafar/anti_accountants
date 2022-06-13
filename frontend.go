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
	badger "github.com/dgraph-io/badger/v3"
)

const (
	CNameOfTheApp           = "ANTI ACCOUNTANTS"
	CPageMenu               = "MENU"
	CPageInvoiceEntry       = "INVOICE ENTRY"
	CPageJournalEntry       = "JOURNAL ENTRY"
	CPageJournalDraft       = "JOURNAL DRAFT"
	CPageJournal            = "JOURNAL"
	CPageAccounts           = "ACCOUNTS"
	CPageAddAccount         = "ADD ACCOUNT"
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
	Vheight        float32 = 50
	VCurrentWindow fyne.Window
)

func main() {
	defer FDbCloseAll()

	a := app.New()
	VCurrentWindow = a.NewWindow("")
	VCurrentWindow.Resize(fyne.Size{Width: CWindowWidth, Height: CWindowHeight})

	Vheight = widget.NewLabel("").MinSize().Height

	VCurrentWindow.SetContent(SPage{}.FLogin())
	VCurrentWindow.ShowAndRun()
}

type SPage struct{}

func (s SPage) FMenu() fyne.CanvasObject {
	FSetTitle(CPageMenu)

	FButtonToSetContentForCurrentWindow := func(label string, page func() fyne.CanvasObject) fyne.CanvasObject {
		wb := widget.NewButton(label, func() { SWidget{}.FSetContentForCurrentWindowWithBackButton(page) })
		wb.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)
		return wb
	}

	wbPageInvoiceEntry := FButtonToSetContentForCurrentWindow(CPageInvoiceEntry, s.FInvoiceEntry)
	wbPageJournalEntry := FButtonToSetContentForCurrentWindow(CPageJournalEntry, s.FJournalEntry)
	wbPageJournalDraft := FButtonToSetContentForCurrentWindow(CPageJournalDraft, s.FJournalDraft)
	wbPageJournal := FButtonToSetContentForCurrentWindow(CPageJournal, s.FJournal)
	wbPageAccounts := FButtonToSetContentForCurrentWindow(CPageAccounts, s.FAccounts)
	wbPageAddAccount := FButtonToSetContentForCurrentWindow(CPageAddAccount, s.FAddAccount)
	wbPageChangePassword := FButtonToSetContentForCurrentWindow(CPageChangePassword, s.FChangePassword)
	wbPageLogin := widget.NewButton("LOGOUT", func() { VCurrentWindow.SetContent(s.FLogin()) })
	wbPageLogin.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)

	fc := container.NewVBox(wbPageInvoiceEntry, wbPageJournalEntry, wbPageJournalDraft, wbPageJournal, wbPageAccounts, wbPageAddAccount, wbPageChangePassword, wbPageLogin)
	if VEmployeeName == CManger {
		fc.Add(widget.NewLabel("this below options are just for " + CManger))
		fc.Add(FButtonToSetContentForCurrentWindow(CPageAddNewEmployee, s.FAddNewEmployee))
		fc.Add(FButtonToSetContentForCurrentWindow(CPageChangeEmployeeName, s.FChangeEmployeeName))
		fc.Add(FButtonToSetContentForCurrentWindow(CPageDeleteEmployee, s.FDeleteEmployee))
		fc.Add(FButtonToSetContentForCurrentWindow(CPageChangeCompanyName, s.FChangeCompanyName))
		fc.Add(FButtonToSetContentForCurrentWindow(CPageDeleteCompany, s.FDeleteCompany))
	}
	return container.NewVScroll(fc)
}

func (SPage) FInvoiceEntry() fyne.CanvasObject {
	FSetTitle(CPageInvoiceEntry)
	return FSetTheEntries(SPageInvoiceEntries{container.NewVBox(SWidget{}.FWCEB("draft"), SWidget{}.FWCEB("note"), SWidget{}.FWCEB("name"), SWidget{}.FWCEB("type"))})
}

func (SPage) FJournalEntry() fyne.CanvasObject {
	FSetTitle(CPageJournalEntry)
	return FSetTheEntries(SPageJournalEntries{container.NewVBox(SWidget{}.FWCEB("draft"), SWidget{}.FWCEB("note"), SWidget{}.FWCEB("name"), SWidget{}.FWCEB("type"))})
}

func (SPage) FJournalDraft() fyne.CanvasObject {
	return FSetDraft(CPageJournalDraft, VDbJournalDrafts)
}

func (SPage) FJournal() fyne.CanvasObject {
	FSetTitle(CPageJournal)

	keys, journal := FDbRead[SJournal1](VDbJournal)
	dates := FConvertByteSliceToTime(keys)

	fcJournal := container.NewVBox(container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight},
		SWidget{}.FWL("dates"),
		SWidget{}.FWL("IsReverse"),
		SWidget{}.FWL("IsReversed"),
		SWidget{}.FWL("ReverseEntryNumberCompound"),
		SWidget{}.FWL("ReverseEntryNumberSimple"),
		SWidget{}.FWL("EntryNumberCompound"),
		SWidget{}.FWL("EntryNumberSimple"),
		SWidget{}.FWL("Value"),
		SWidget{}.FWL("PriceDebit"),
		SWidget{}.FWL("PriceCredit"),
		SWidget{}.FWL("QuantityDebit"),
		SWidget{}.FWL("QuantityCredit"),
		SWidget{}.FWL("AccountDebit"),
		SWidget{}.FWL("AccountCredit"),
		SWidget{}.FWL("Notes"),
		SWidget{}.FWL("Name"),
		SWidget{}.FWL("Employee"),
		SWidget{}.FWL("TypeOfCompoundEntry"),
	))
	for k1, v1 := range journal {
		fcJournal.Add(container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight},
			SWidget{}.FWL(dates[k1]),
			SWidget{}.FWL(v1.IsReverse),
			SWidget{}.FWL(v1.IsReversed),
			SWidget{}.FWL(v1.ReverseEntryNumberCompound),
			SWidget{}.FWL(v1.ReverseEntryNumberSimple),
			SWidget{}.FWL(v1.EntryNumberCompound),
			SWidget{}.FWL(v1.EntryNumberSimple),
			SWidget{}.FWL(v1.Value),
			SWidget{}.FWL(v1.DebitPrice),
			SWidget{}.FWL(v1.CreditPrice),
			SWidget{}.FWL(v1.DebitQuantity),
			SWidget{}.FWL(v1.CreditQuantity),
			SWidget{}.FWL(v1.DebitAccountName),
			SWidget{}.FWL(v1.CreditAccountName),
			SWidget{}.FWL(v1.Notes),
			SWidget{}.FWL(v1.Name),
			SWidget{}.FWL(v1.Employee),
			SWidget{}.FWL(v1.TypeOfCompoundEntry),
		))
	}

	return container.NewVScroll(fcJournal)
}

func (SPage) FAccounts() fyne.CanvasObject {
	FSetTitle(CPageAccounts)

	fcAccounts := container.NewVBox(container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}))

	for _, v1 := range VAccounts {
		fcAccounts.Add(container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight},
			SWidget{}.FWL(v1.IsCredit),
			SWidget{}.FWL(v1.CostFlowType),
			SWidget{}.FWL(v1.Name),
			SWidget{}.FWL(v1.Notes),
			SWidget{}.FWL(v1.Image),
			SWidget{}.FWL(v1.Barcode),
			SWidget{}.FWL(v1.Number),
			SWidget{}.FWL(v1.Levels),
			SWidget{}.FWL(v1.FathersName),
		))
	}

	return container.NewVScroll(fcAccounts)
}

func (SPage) FAddAccount() fyne.CanvasObject {
	FSetTitle(CPageAddAccount)

	wsIsCredit := SWidget{}.FWCSB([]string{CCredit, CDebit})
	wsCostFlowType := SWidget{}.FWCSB([]string{CHighLevelAccount, CWma, CFifo, CLifo})
	weInventory := SWidget{}.FWCSEB("Inventory", func() []string {
		var inventories []string
		for _, v1 := range VAccounts {
			inventories = append(inventories, v1.Inventory)
		}
		return inventories
	}())
	weAccountName := SWidget{}.FWCEB("Name")
	weAccountNotes := SWidget{}.FWCEB("Notes")
	weAccountImage := SWidget{}.FFWE("Image")
	weAccountBarcode := SWidget{}.FFWE("Barcode")
	weAccountNumber := SWidget{}.FFWE("Number")

	clearSelect := func(fc fyne.CanvasObject, isClearChecked bool) {
		if fc.(*fyne.Container).Objects[0].(*widget.Check).Checked == isClearChecked {
			fc.(*fyne.Container).Objects[1].(*widget.Select).SetSelectedIndex(0)
		}
	}

	clearSelectEntry := func(fc fyne.CanvasObject, isClearChecked bool) {
		if fc.(*fyne.Container).Objects[0].(*widget.Check).Checked == isClearChecked {
			fc.(*fyne.Container).Objects[1].(*widget.SelectEntry).SetText("")
		}
	}

	clearEntry := func(fc fyne.CanvasObject, isClearChecked bool) {
		if fc.(*fyne.Container).Objects[0].(*widget.Check).Checked == isClearChecked {
			fc.(*fyne.Container).Objects[1].(*widget.Entry).SetText("")
		}
	}

	clear := func(isClearChecked bool) {
		clearSelect(wsIsCredit, isClearChecked)
		clearSelect(wsCostFlowType, isClearChecked)
		clearSelectEntry(weInventory, isClearChecked)
		clearEntry(weAccountName, isClearChecked)
		clearEntry(weAccountNotes, isClearChecked)
		SDelete{}.FLines(weAccountImage, 1, isClearChecked)
		SDelete{}.FLines(weAccountBarcode, 1, isClearChecked)
		SDelete{}.FLines(weAccountNumber, 1, isClearChecked)
	}

	getWsText := func(fc fyne.CanvasObject) string {
		return fc.(*fyne.Container).Objects[1].(*widget.Select).Selected
	}

	getWeText := func(fc fyne.CanvasObject) string {
		return fc.(*fyne.Container).Objects[1].(*widget.Entry).Text
	}

	getWseText := func(fc fyne.CanvasObject) string {
		return fc.(*fyne.Container).Objects[1].(*widget.SelectEntry).Text
	}

	getFweText := func(fc fyne.CanvasObject) []string {
		var slice []string
		for _, v1 := range fc.(*fyne.Container).Objects[1:] {
			slice = append(slice, getWeText(v1))
		}
		return slice
	}

	save := func(isSave bool) {
		err := FAddAccount(isSave, SAccount1{
			IsCredit:     getWsText(wsIsCredit) == CCredit,
			CostFlowType: getWsText(wsCostFlowType),
			Inventory:    getWseText(weInventory),
			Name:         getWeText(weAccountName),
			Notes:        getWeText(weAccountNotes),
			Image:        getFweText(weAccountImage),
			Barcode:      getFweText(weAccountBarcode),
			Number: func() [][]uint {
				var slice1 [][]uint
				for _, v1 := range getFweText(weAccountNumber) {
					if v1 == "" {
						continue
					}
					slice2 := []uint{}
					for _, v2 := range strings.Split(v1, ",") {
						i, err := strconv.Atoi(v2)
						log.Println(err)
						slice2 = append(slice2, uint(i))
					}
					slice1 = append(slice1, slice2)
				}
				return slice1
			}(),
		})
		log.Println(err)
		FPrintFormatedAccounts()
	}

	wbOk := SWidget{}.FButton("ok", func() {
		save(true)
		clear(false)
	})
	wbSave := SWidget{}.FButton("save", func() { save(true) })
	wbClearChecked := SWidget{}.FButton("clear checked", func() { clear(true) })
	wbClearNotChecked := SWidget{}.FButton("clear not checked", func() { clear(false) })

	fc1 := container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wbOk, wbSave, wbClearChecked, wbClearNotChecked)

	return container.New(layout.NewBorderLayout(fc1, nil, nil, nil), fc1,
		container.NewVScroll(container.NewVBox(
			wsIsCredit,
			wsCostFlowType,
			weInventory,
			weAccountName,
			weAccountNotes,
			weAccountImage,
			weAccountBarcode,
			weAccountNumber,
		)))
}

func (SPage) FChangePassword() fyne.CanvasObject {
	FSetTitle(CPageChangePassword)

	wePassword := SWidget{}.FWEPassword()
	wlError := SWidget{}.FWlError()

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

func (s SPage) FLogin() fyne.CanvasObject {
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

	wePassword := SWidget{}.FWEPassword()
	wlError := SWidget{}.FWlError()

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
		VCurrentWindow.SetContent(s.FMenu())
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
		VCurrentWindow.SetContent(s.FMenu())
	})

	return container.NewVBox(weCompanyName, weEmployeeName, wePassword, wlError, wbLogin, wbCreateNewCompany)
}

func (SPage) FAddNewEmployee() fyne.CanvasObject {
	FSetTitle(CPageAddNewEmployee)

	weEmployeeName := widget.NewEntry()
	weEmployeeName.SetPlaceHolder("employee name")

	wePassword := SWidget{}.FWEPassword()
	wlError := SWidget{}.FWlError()

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

func (SPage) FChangeEmployeeName() fyne.CanvasObject {
	FSetTitle(CPageChangeEmployeeName)

	employees, _ := FDbOpenAndReadEmployees()

	weEmployeeName := widget.NewSelectEntry(employees)
	weEmployeeName.SetPlaceHolder("employee name")

	weNewEmployeeName := widget.NewEntry()
	weNewEmployeeName.SetPlaceHolder("new employee name")

	wlError := SWidget{}.FWlError()

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

func (SPage) FDeleteEmployee() fyne.CanvasObject {
	FSetTitle(CPageDeleteEmployee)

	employees, _ := FDbOpenAndReadEmployees()

	weEmployeeName := widget.NewSelectEntry(employees)
	weEmployeeName.SetPlaceHolder("employee name")

	wlError := SWidget{}.FWlError()

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

func (SPage) FChangeCompanyName() fyne.CanvasObject {
	FSetTitle(CPageChangeCompanyName)

	weCompanyName := widget.NewEntry()
	weCompanyName.SetPlaceHolder("company name")

	wlError := SWidget{}.FWlError()

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

func (s SPage) FDeleteCompany() fyne.CanvasObject {
	FSetTitle(CPageDeleteCompany)

	wePassword1 := SWidget{}.FWEPassword()
	wePassword2 := SWidget{}.FWEPassword()
	wePassword3 := SWidget{}.FWEPassword()
	wePassword4 := SWidget{}.FWEPassword()
	wePassword5 := SWidget{}.FWEPassword()
	wePassword6 := SWidget{}.FWEPassword()
	wePassword7 := SWidget{}.FWEPassword()
	wePassword8 := SWidget{}.FWEPassword()
	wePassword9 := SWidget{}.FWEPassword()
	wePassword10 := SWidget{}.FWEPassword()

	wlError := SWidget{}.FWlError()

	wb := widget.NewButton("delete "+VCompanyName+" company", func() {
		wlError.SetText("")

		employees, passwords := FDbOpenAndReadEmployees()

		indexOfEmployee, _ := FFind(CManger, employees)
		if passwords[indexOfEmployee] != wePassword1.Text {
			wlError.SetText(wePassword1.PlaceHolder + " is wrong")
			return
		}

		if wePassword1.Text != wePassword2.Text ||
			wePassword1.Text != wePassword3.Text ||
			wePassword1.Text != wePassword4.Text ||
			wePassword1.Text != wePassword5.Text ||
			wePassword1.Text != wePassword6.Text ||
			wePassword1.Text != wePassword7.Text ||
			wePassword1.Text != wePassword8.Text ||
			wePassword1.Text != wePassword9.Text ||
			wePassword1.Text != wePassword10.Text {
			wlError.SetText("not all the passwords are the same")
			return
		}

		fmt.Println(os.RemoveAll(CPathDataBase + VCompanyName))
		VCurrentWindow.SetContent(s.FLogin())
	})

	return container.NewVScroll(container.NewVBox(wePassword1, wePassword2, wePassword3, wePassword4,
		wePassword5, wePassword6, wePassword7, wePassword8, wePassword9, wePassword10, wlError, wb))
}

type IPageEntries interface {
	FGetEntries() *fyne.Container
	FTheIndexOfTheFirstEntry() int
	FSave(bool) error
	FClear(bool)
	FNewLine() fyne.CanvasObject
	FDraft(string)
}

type SPageInvoiceEntries struct{ *fyne.Container }

func (s SPageInvoiceEntries) FGetEntries() *fyne.Container { return s.Container }
func (SPageInvoiceEntries) FTheIndexOfTheFirstEntry() int  { return 4 }

func (SPageInvoiceEntries) FAllAccounts() []string {
	var allAccountNames []string
	for _, v1 := range VAutoCompletionEntries {
		allAccountNames = append(allAccountNames, v1.TAccountName)
	}
	return allAccountNames
}

func (s SPageInvoiceEntries) FSave(insert bool) error {
	entryInfo := SEntry{
		Notes:               s.FGetEntries().Objects[1].(*fyne.Container).Objects[1].(*widget.Entry).Text,
		Name:                s.FGetEntries().Objects[2].(*fyne.Container).Objects[1].(*widget.Entry).Text,
		TypeOfCompoundEntry: s.FGetEntries().Objects[3].(*fyne.Container).Objects[1].(*widget.Entry).Text,
	}

	var entries1 []SAPQ
	for _, v1 := range s.FGetEntries().Objects[s.FTheIndexOfTheFirstEntry():] {
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
	for k1, v1 := range s.FGetEntries().Objects {
		if k1 == s.FTheIndexOfTheFirstEntry() {
			break
		}
		if v1.(*fyne.Container).Objects[0].(*widget.Check).Checked == isChecked {
			v1.(*fyne.Container).Objects[1].(*widget.Entry).SetText("")
		}
	}

	SDelete{}.FLines(s.FGetEntries(), s.FTheIndexOfTheFirstEntry(), isChecked)
}

func (s SPageInvoiceEntries) FNewLine() fyne.CanvasObject {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, Vheight))

	wsAccount := SWidget{}.FSelectEntry(s.FAllAccounts())
	wsAccount.Resize(fyne.NewSize(80/3, Vheight))
	wsAccount.SetPlaceHolder("account")

	wlPrice := widget.NewLabel("")
	wlPrice.Resize(fyne.NewSize(80/3, Vheight))

	weQuantity := widget.NewEntry()
	weQuantity.Resize(fyne.NewSize(80/3, Vheight))
	weQuantity.SetPlaceHolder("quantity")
	weQuantity.SetText(fmt.Sprint(1))

	wbX := widget.NewButton("x", nil)
	wbX.Resize(fyne.NewSize(10, Vheight))

	fc := container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wc, wsAccount, wlPrice, weQuantity, wbX)

	wbX.OnTapped = func() { s.FGetEntries().Objects = SDelete{}.FLine(s.FGetEntries(), fc) }
	return fc
}

func (s SPageInvoiceEntries) FDraft(nameOfTheDraft string) {
	FDbUpdate(VDbInvoiceDrafts, []byte(nameOfTheDraft), s.FGetEntries())
}

type SPageJournalEntries struct{ *fyne.Container }

func (s SPageJournalEntries) FGetEntries() *fyne.Container { return s.Container }
func (SPageJournalEntries) FTheIndexOfTheFirstEntry() int  { return 4 }

func (SPageJournalEntries) FAllAccounts() []string {
	var allAccountNames []string
	for _, v1 := range VAccounts {
		if v1.CostFlowType != CHighLevelAccount {
			allAccountNames = append(allAccountNames, v1.Name)
		}
	}
	return allAccountNames
}

func (s SPageJournalEntries) FSave(insert bool) error {
	entryInfo := SEntry{
		Notes:               s.FGetEntries().Objects[1].(*fyne.Container).Objects[1].(*widget.Entry).Text,
		Name:                s.FGetEntries().Objects[2].(*fyne.Container).Objects[1].(*widget.Entry).Text,
		TypeOfCompoundEntry: s.FGetEntries().Objects[3].(*fyne.Container).Objects[1].(*widget.Entry).Text,
	}

	var entries1 []SAPQ
	for _, v1 := range s.FGetEntries().Objects[s.FTheIndexOfTheFirstEntry():] {
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
	for k1, v1 := range s.FGetEntries().Objects {
		if k1 == s.FTheIndexOfTheFirstEntry() {
			break
		}
		if v1.(*fyne.Container).Objects[0].(*widget.Check).Checked == isChecked {
			v1.(*fyne.Container).Objects[1].(*widget.Entry).SetText("")
		}
	}

	SDelete{}.FLines(s.FGetEntries(), s.FTheIndexOfTheFirstEntry(), isChecked)
}

func (s SPageJournalEntries) FNewLine() fyne.CanvasObject {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, Vheight))

	wsAccount := SWidget{}.FSelectEntry(s.FAllAccounts())
	wsAccount.Resize(fyne.NewSize(80/3, Vheight))
	wsAccount.SetPlaceHolder("account")

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

	wbX.OnTapped = func() { s.FGetEntries().Objects = SDelete{}.FLine(s.FGetEntries(), fc) }
	return fc
}

func (s SPageJournalEntries) FDraft(nameOfTheDraft string) {
	FDbUpdate(VDbJournalDrafts, []byte(nameOfTheDraft), s.FGetEntries())
}

type SWidget struct{}

func (SWidget) FButton(label string, tapped func()) *widget.Button {
	wb := widget.NewButton(label, tapped)
	wb.Resize(fyne.Size{Width: wb.MinSize().Width, Height: Vheight})
	return wb
}

func (SWidget) FWlError() *widget.Label {
	wl := widget.NewLabel("")
	wl.Resize(fyne.Size{Width: CWindowWidth, Height: Vheight})
	wl.Wrapping = fyne.TextWrapWord
	wl.Alignment = fyne.TextAlignCenter
	return wl
}

func (SWidget) FSelectEntry(options []string) *widget.SelectEntry {
	ws := widget.NewSelectEntry(options)
	ws.OnChanged = func(option string) {
		// if option == "" {
		// 	ws.SetOptions(options)
		// 	return
		// }
		var newOptions []string
		for _, v1 := range options {
			if strings.Contains(v1, option) {
				newOptions = append(newOptions, v1)
			}
		}
		ws.SetOptions(newOptions)
	}
	return ws
}

func (SWidget) FSetContentForCurrentWindowWithBackButton(page func() fyne.CanvasObject) {
	previousWindow := VCurrentWindow.Content()
	previousTitel := VCurrentWindow.Title()
	wbBack := widget.NewButton("back", func() {
		VCurrentWindow.SetTitle(previousTitel)
		VCurrentWindow.SetContent(previousWindow)
	})
	VCurrentWindow.SetContent(container.New(layout.NewBorderLayout(nil, wbBack, nil, nil), wbBack, page()))
}

func (SWidget) FWCSB(options []string) fyne.CanvasObject {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, Vheight))

	ws := widget.NewSelect(options, nil)
	ws.Resize(fyne.NewSize(80, 0))
	ws.SetSelected(options[0])

	wb := widget.NewButton("x", func() { ws.SetSelectedIndex(0) })
	wb.Resize(fyne.NewSize(10, Vheight))

	return container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wc, ws, wb)
}

func (SWidget) FWCSEB(placeHolder string, options []string) fyne.CanvasObject {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, Vheight))

	ws := SWidget{}.FSelectEntry(options)
	ws.Resize(fyne.NewSize(80, 0))
	ws.SetPlaceHolder(placeHolder)

	wb := widget.NewButton("x", func() { ws.SetText("") })
	wb.Resize(fyne.NewSize(10, Vheight))

	return container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wc, ws, wb)
}

func (SWidget) FWCEB(placeHolder string) fyne.CanvasObject {
	wc := widget.NewCheck("", nil)
	wc.Resize(fyne.NewSize(10, Vheight))

	we := widget.NewEntry()
	we.SetPlaceHolder(placeHolder)
	we.Resize(fyne.NewSize(80, Vheight))

	wb := widget.NewButton("x", func() { we.SetText("") })
	wb.Resize(fyne.NewSize(10, Vheight))

	return container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wc, we, wb)
}

func (SWidget) FWEPassword() *widget.Entry {
	we := widget.NewPasswordEntry()
	we.SetPlaceHolder("password")
	return we
}

func (SWidget) FWL(x any) *widget.Label {
	wl := widget.NewLabel(fmt.Sprint(x))
	wl.Resize(fyne.NewSize(10, Vheight))
	return wl
}

func (SWidget) FFWE(placeHolder string) fyne.CanvasObject {
	wbAdd := widget.NewButton("Add "+placeHolder, func() {})
	fc := container.NewVBox(wbAdd, SWidget{}.FWCEB(placeHolder))
	wbAdd.OnTapped = func() { fc.Add(SWidget{}.FWCEB(placeHolder)) }
	return fc
}

type SDelete struct{}

func (SDelete) FLines(fc fyne.CanvasObject, indexToStart int, isChecked bool) {
	k1 := indexToStart
	for k1 < len(fc.(*fyne.Container).Objects) {
		if fc.(*fyne.Container).Objects[k1].(*fyne.Container).Objects[0].(*widget.Check).Checked == isChecked {
			fc.(*fyne.Container).Objects = FRemove(fc.(*fyne.Container).Objects, k1)
		} else {
			k1++
		}
	}
}

func (SDelete) FLine(fcBig fyne.CanvasObject, fcSmall fyne.CanvasObject) []fyne.CanvasObject {
	for k1, v1 := range fcBig.(*fyne.Container).Objects {
		if v1 == fcSmall {
			return FRemove(fcBig.(*fyne.Container).Objects, k1)
		}
	}
	return fcBig.(*fyne.Container).Objects
}

////////////////////////////////////////////////////////////////////////////////////////////

func FDbOpenAndReadEmployees() ([]string, []string) {
	VDbEmployees = FDbOpen(VDbEmployees, CPathDataBase+VCompanyName+CPathEmployees)
	keys, passwords := FDbRead[string](VDbEmployees)
	employees := FConvertFromByteSliceToString(keys)
	return employees, passwords
}

func FDisplayTheError(errorMessage error, wlError *widget.Label) {
	if errorMessage != nil {
		wlError.SetText(errorMessage.Error())
	} else {
		wlError.SetText("")
	}
}

func FSetTitle(title string) {
	if title != CPageLogin {
		VCurrentWindow.SetTitle(CNameOfTheApp + " - " + VCompanyName + " - " + VEmployeeName + " - " + title)
	} else {
		VCurrentWindow.SetTitle(CNameOfTheApp + " - " + title)
	}
}

func FSetTheEntries(page IPageEntries) fyne.CanvasObject {
	wlError := SWidget{}.FWlError()

	go func() {
		for range time.Tick(time.Second / 4) {
			err := page.FSave(false)
			FDisplayTheError(err, wlError)
			if len(page.FGetEntries().Objects) == page.FTheIndexOfTheFirstEntry() {
				page.FGetEntries().Add(page.FNewLine())
			}
		}
	}()

	wePrint := widget.NewEntry()
	wePrint.SetPlaceHolder("print count")
	wePrint.Resize(fyne.NewSize(wePrint.MinSize().Width, Vheight))

	wbOk := SWidget{}.FButton("ok", func() {
		err := page.FSave(true)
		FDisplayTheError(err, wlError)
		if err == nil {
			page.FClear(false)
			i, err := strconv.ParseUint(wePrint.Text, 8, 8)
			if i > 0 {
				log.Println(i)
			} else {
				log.Println(err)
			}
		}
	})
	wbSave := SWidget{}.FButton("save", func() {
		err := page.FSave(true)
		FDisplayTheError(err, wlError)
	})
	wbClearChecked := SWidget{}.FButton("clear checked", func() { page.FClear(true) })
	wbClearNotChecked := SWidget{}.FButton("clear not checked", func() { page.FClear(false) })
	wbDraft := SWidget{}.FButton("draft", func() { page.FDraft(page.FGetEntries().Objects[0].(*fyne.Container).Objects[1].(*widget.Entry).Text) })
	wbAdd := SWidget{}.FButton("+", func() { page.FGetEntries().Add(page.FNewLine()) })

	page.FGetEntries().Add(page.FNewLine())
	wsEntries := container.NewVScroll(page.FGetEntries())
	fc1 := container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wbOk, wbSave, wbClearChecked, wbClearNotChecked, wbDraft, wbAdd, wePrint)
	return container.New(layout.NewBorderLayout(fc1, wlError, nil, nil), fc1, wsEntries, wlError)
}

func FSetDraft(title1 string, db *badger.DB) fyne.CanvasObject {
	FSetTitle(title1)
	keys, value := FDbRead[*fyne.Container](db)
	fc := container.NewVBox(widget.NewLabel(""))

	DeleteLine := func(wbName *widget.Button) {
		for k2, v2 := range fc.Objects[1:] {
			if v2.(*fyne.Container).Objects[0].(*widget.Button).Text == wbName.Text {
				FDbDelete(db, []byte(wbName.Text))
				fc.Objects = FRemove(fc.Objects, k2+1)
				fc.Refresh()
			}
		}
	}

	SetTheDraftedEntry := func(k1 int) {
		SWidget{}.FSetContentForCurrentWindowWithBackButton(func() fyne.CanvasObject { return FSetTheEntries(SPageJournalEntries{value[k1]}) })
	}

	for k1, v1 := range keys {

		wbName := widget.NewButton(string(v1), func() { SetTheDraftedEntry(k1) })
		wbName.Resize(fyne.NewSize(70, Vheight))
		wbName.Alignment = widget.ButtonAlign(widget.ButtonAlignLeading)

		wbEdite := widget.NewButton("edite", func() {
			DeleteLine(wbName)
			SetTheDraftedEntry(k1)
		})
		wbEdite.Resize(fyne.NewSize(20, Vheight))

		wbX := widget.NewButton("x", func() { DeleteLine(wbName) })
		wbX.Resize(fyne.NewSize(10, Vheight))

		fc.Add(container.New(&TPercentageHBoxLayout{Width: CWindowWidth, Height: Vheight}, wbName, wbEdite, wbX))
	}
	return container.NewVScroll(fc)
}
