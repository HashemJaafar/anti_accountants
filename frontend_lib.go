package main

import (
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TPercentageHBoxLayout fyne.Size
type TPercentageVBoxLayout fyne.Size
type SStretchHBoxLayout struct {
	Width           float32
	Height          float32
	ObjectToStertch fyne.CanvasObject
}
type SStretchVBoxLayout struct {
	Width           float32
	Height          float32
	ObjectToStertch fyne.CanvasObject
}

func (s *TPercentageHBoxLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
}

func (s *TPercentageVBoxLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
}

func (s *SStretchHBoxLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
}

func (s *SStretchVBoxLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
}

func (s *TPercentageHBoxLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	var total float32
	for _, v1 := range objects {
		total += v1.Size().Width
	}
	var pos fyne.Position
	for _, v1 := range objects {
		width := v1.Size().Width / total * containerSize.Width
		v1.Resize(fyne.Size{Width: width, Height: containerSize.Height})
		v1.Move(pos)
		pos = pos.Add(fyne.NewPos(width, 0))
	}
}

func (s *TPercentageVBoxLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	var total float32
	for _, v1 := range objects {
		total += v1.Size().Height
	}
	var pos fyne.Position
	for _, v1 := range objects {
		height := v1.Size().Height / total * containerSize.Height
		v1.Resize(fyne.Size{Width: containerSize.Width, Height: height})
		v1.Move(pos)
		pos = pos.Add(fyne.NewPos(0, height))
	}
}

func (s *SStretchHBoxLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	var index int
	var total float32
	for k1, v1 := range objects {
		if v1 == s.ObjectToStertch {
			index = k1
			continue
		}
		total += v1.Size().Width
	}
	objects[index].Resize(fyne.NewSize(containerSize.Width-total, containerSize.Height))
	var pos fyne.Position
	for _, v1 := range objects {
		width := v1.Size().Width
		v1.Resize(fyne.Size{Width: width, Height: containerSize.Height})
		v1.Move(pos)
		pos = pos.Add(fyne.NewPos(width, 0))
	}
}

func (s *SStretchVBoxLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	var index int
	var total float32
	for k1, v1 := range objects {
		if v1 == s.ObjectToStertch {
			index = k1
			continue
		}
		total += v1.Size().Height
	}
	objects[index].Resize(fyne.NewSize(containerSize.Width, containerSize.Height-total))
	var pos fyne.Position
	for _, v1 := range objects {
		height := v1.Size().Height
		v1.Resize(fyne.Size{Width: containerSize.Width, Height: height})
		v1.Move(pos)
		pos = pos.Add(fyne.NewPos(0, height))
	}
}

func main2() {
	a := app.New()
	w := a.NewWindow("")

	b1 := widget.NewButton("b1", nil)
	b2 := widget.NewButton("b2", nil)
	b3 := widget.NewButton("b3", nil)
	b4 := widget.NewButton("b4", nil)
	b5 := widget.NewButton("b5", nil)
	b6 := widget.NewButton("b6", nil)

	b1.Resize(fyne.NewSize(50, 100))
	b2.Resize(fyne.NewSize(100, 100))
	b3.Resize(fyne.NewSize(50, 100))
	b4.Resize(fyne.NewSize(100, 100))
	b5.Resize(fyne.NewSize(100, 100))
	b6.Resize(fyne.NewSize(100, 200))

	l1 := container.New(&TPercentageHBoxLayout{Width: 50, Height: 80}, b1, b2, b3)
	l2 := container.New(&TPercentageVBoxLayout{Width: 50, Height: 80}, b4, b5, b6)
	l3 := container.New(&TPercentageVBoxLayout{Width: 50, Height: 80}, l1, l2)

	w.SetContent(l3)
	w.ShowAndRun()
}

func main3() {
	a := app.New()
	w := a.NewWindow("")

	b1 := widget.NewButton("b1", nil)
	b2 := widget.NewButton("b2", nil)
	b3 := widget.NewButton("b3", nil)
	b4 := widget.NewButton("b4", nil)
	b5 := widget.NewButton("b5", nil)
	b6 := widget.NewButton("b6", nil)

	b1.Resize(fyne.NewSize(50, 100))
	b2.Resize(fyne.NewSize(100, 100))
	b3.Resize(fyne.NewSize(50, 100))
	b4.Resize(fyne.NewSize(100, 100))
	b5.Resize(fyne.NewSize(100, 100))
	b6.Resize(fyne.NewSize(100, 200))

	l1 := container.New(&SStretchHBoxLayout{Width: 50, Height: 80, ObjectToStertch: b2}, b1, b2, b3)
	l2 := container.New(&SStretchVBoxLayout{Width: 50, Height: 80, ObjectToStertch: b5}, b4, b5, b6)
	l3 := container.New(&TPercentageHBoxLayout{Width: 20, Height: 20}, l1, l2)

	w.SetContent(l3)
	w.ShowAndRun()
}

func FFindObject(object fyne.CanvasObject, objects []fyne.CanvasObject) (int, bool) {
	for k1, v1 := range objects {
		if v1 == object {
			return k1, true
		}
	}
	return 0, false
}

func main1() {
	// VCompanyName = "anti_accountants"
	// FDbOpenAll()
	// _, value := FDbRead[*fyne.Container](VDbJournalDrafts)
	// FDbCloseAll()

	// fmt.Println(value)
	a := app.New()
	w := a.NewWindow("")

	wb1 := widget.NewButton("b1", nil)
	wb1.Text = "ihasiodjodjaoi"
	fc1 := container.NewHBox(wb1)

	v, err := json.Marshal(fc1)
	fmt.Println(v)
	fmt.Println(err)

	var fc2 *fyne.Container
	err = json.Unmarshal(v, &fc2)
	fmt.Println(err)

	w.SetContent(fc2)
	w.ShowAndRun()
}
