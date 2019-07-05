package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Text struct {
	renderer *sdl.Renderer
	font     *ttf.Font
	color    sdl.Color
	text     string
	x        int32
	y        int32
	w        int32
	h        int32
}

func newText(renderer *sdl.Renderer) *Text {
	text := new(Text)

	text.renderer = renderer

	font, err := ttf.OpenFont("asset/font/RedRose-Bold.ttf", 32)
	if err != nil {
		panic(err)
	}

	text.font = font
	text.color = sdl.Color{R: 255, G: 255, B: 255, A: 255}

	return text
}

func (self *Text) paint() {
	sdl.Main(func() {
		surface, err := self.font.RenderUTF8Solid(self.text, self.color)
		checkErr(err)
		defer surface.Free()
		tex, err := self.renderer.CreateTextureFromSurface(surface)
		checkErr(err)
		defer tex.Destroy()
		err = self.renderer.Copy(tex, nil, &sdl.Rect{X: self.x, Y: self.y, W: self.w, H: self.h})
		checkErr(err)
		self.renderer.Present()
	})
}
