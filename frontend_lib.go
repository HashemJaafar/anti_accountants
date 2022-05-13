package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PercentageHBoxLayout fyne.Size
type PercentageVBoxLayout fyne.Size
type StretchHBoxLayout struct {
	Width           float32
	Height          float32
	ObjectToStertch fyne.CanvasObject
}
type StretchVBoxLayout struct {
	Width           float32
	Height          float32
	ObjectToStertch fyne.CanvasObject
}

func (s *PercentageHBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
}

func (s *PercentageVBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
}

func (s *StretchHBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
}

func (s *StretchVBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
}

func (s *PercentageHBoxLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
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

func (s *PercentageVBoxLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
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

func (s *StretchHBoxLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
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

func (s *StretchVBoxLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
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

	l1 := container.New(&PercentageHBoxLayout{Width: 50, Height: 80}, b1, b2, b3)
	l2 := container.New(&PercentageVBoxLayout{Width: 50, Height: 80}, b4, b5, b6)
	l3 := container.New(&PercentageVBoxLayout{Width: 50, Height: 80}, l1, l2)

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

	l1 := container.New(&StretchHBoxLayout{Width: 50, Height: 80, ObjectToStertch: b2}, b1, b2, b3)
	l2 := container.New(&StretchVBoxLayout{Width: 50, Height: 80, ObjectToStertch: b5}, b4, b5, b6)
	l3 := container.New(&PercentageHBoxLayout{Width: 20, Height: 20}, l1, l2)

	w.SetContent(l3)
	w.ShowAndRun()
}
