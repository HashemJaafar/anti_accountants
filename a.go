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
	var totalWidth float32
	for _, v1 := range objects {
		totalWidth += v1.Size().Width
	}
	var pos fyne.Position
	for _, v1 := range objects {
		width := v1.Size().Width / totalWidth * s.Width
		v1.Resize(fyne.Size{Width: width, Height: s.Height})
		v1.Move(pos)
		pos = pos.Add(fyne.NewPos(width, 0))
	}
}

func (s *DynamicVBoxLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	var totalHeight float32
	for _, v1 := range objects {
		totalHeight += v1.Size().Height
	}
	var pos fyne.Position
	for _, v1 := range objects {
		height := v1.Size().Height / totalHeight * s.Height
		v1.Resize(fyne.Size{Width: s.Width, Height: height})
		v1.Move(pos)
		pos = pos.Add(fyne.NewPos(0, height))
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("")

	text1 := widget.NewButton("text1", nil)
	text2 := widget.NewButton("text2", nil)
	text3 := widget.NewButton("text3", nil)
	text4 := widget.NewButton("text4", nil)
	text5 := widget.NewButton("text5", nil)
	text6 := widget.NewButton("text6", nil)
	text1.Resize(fyne.NewSize(50, 100))
	text2.Resize(fyne.NewSize(100, 100))
	text3.Resize(fyne.NewSize(50, 100))
	text4.Resize(fyne.NewSize(100, 100))
	text5.Resize(fyne.NewSize(100, 100))
	text6.Resize(fyne.NewSize(100, 100))

	l1 := container.New(&DynamicHBoxLayout{Width: 500, Height: 250}, text1, text2, text3)
	l2 := container.New(&DynamicVBoxLayout{Width: 500, Height: 250}, text4, text5, text6)
	l3 := container.New(&DynamicVBoxLayout{Width: 500, Height: 500}, l1, l2)
	w.SetContent(l3)
	w.ShowAndRun()
}
