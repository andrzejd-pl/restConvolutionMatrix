package main

import (
	"image"
	"image/color"
	"sync"
)

type Apm interface {
	BeginConvolution()
	EndConvolution()
}

type Task struct {
	waitGroup    *sync.WaitGroup
	lines, start int
	oldImage     *image.Image
	newImage     *image.Gray
}

func (t *Task) BeginConvolution() {
	go t.call(t.oldImage, t.newImage, t.lines, t.start)
}

func (t *Task) EndConvolution() {
	t.waitGroup.Wait()
}

func (t *Task) call(oldImage *image.Image, newImage *image.Gray, lines int, start int) {
	defer t.waitGroup.Done()
	for i := start; i < lines+start; i++ {
		for j := 0; j < 1024; j++ {
			value, _, _, _ := (*oldImage).At(i, j).RGBA()
			var up, down, left, right uint32

			if i > 0 {
				up, _, _, _ = (*oldImage).At(i-1, j).RGBA()
			} else {
				up = 0
			}

			if i < 1023 {
				down, _, _, _ = (*oldImage).At(i+1, j).RGBA()
			} else {
				down = 0
			}

			if j > 0 {
				left, _, _, _ = (*oldImage).At(i, j-1).RGBA()
			} else {
				left = 0
			}

			if j < 1023 {
				right, _, _, _ = (*oldImage).At(i, j+1).RGBA()
			} else {
				right = 0
			}

			value = ((value * 6) + (up) + (right) + (left) + (down)) / 10

			newImage.SetGray(i, j, color.Gray{Y: (uint8)(value >> 8)})
		}
	}
}
