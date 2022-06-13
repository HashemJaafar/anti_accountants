package main

import (
	"fyne.io/fyne/v2"
)

type TPercentageHBoxLayout fyne.Size

func (s *TPercentageHBoxLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
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

type TPercentageVBoxLayout fyne.Size

func (s *TPercentageVBoxLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
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

type SStretchHBoxLayout struct {
	Width           float32
	Height          float32
	ObjectToStertch fyne.CanvasObject
}

func (s *SStretchHBoxLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
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

type SStretchVBoxLayout struct {
	Width           float32
	Height          float32
	ObjectToStertch fyne.CanvasObject
}

func (s *SStretchVBoxLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: s.Width, Height: s.Height}
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
