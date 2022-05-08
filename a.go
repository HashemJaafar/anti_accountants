package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type DynamicHBoxLayout fyne.Size
type DynamicVBoxLayout fyne.Size

func (s *DynamicHBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
}

func (s *DynamicVBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
}

func (s *DynamicHBoxLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	var total float32
	for _, v1 := range objects {
		total += v1.Size().Width
	}
	var pos fyne.Position
	for _, v1 := range objects {
		width := v1.Size().Width / total * s.Width
		v1.Resize(fyne.Size{Width: width, Height: s.Height})
		v1.Move(pos)
		pos = pos.Add(fyne.NewPos(width, 0))
	}
}

func (s *DynamicVBoxLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	var total float32
	for _, v1 := range objects {
		total += v1.Size().Height
	}
	var pos fyne.Position
	for _, v1 := range objects {
		height := v1.Size().Height / total * s.Height
		v1.Resize(fyne.Size{Width: s.Width, Height: height})
		v1.Move(pos)
		pos = pos.Add(fyne.NewPos(0, height))
	}
}

func Main() {
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
	b6.Resize(fyne.NewSize(100, 100))

	l1 := container.New(&DynamicHBoxLayout{Width: 400, Height: 250}, b1, b2, b3)
	l2 := container.New(&DynamicVBoxLayout{Width: 400, Height: 250}, b4, b5, b6)
	l3 := container.New(&DynamicVBoxLayout{Width: 400, Height: 500}, l1, l2)

	w.SetContent(l3)
	w.ShowAndRun()
}
