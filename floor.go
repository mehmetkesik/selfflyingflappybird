package main

import "github.com/veandco/go-sdl2/sdl"
import sdlImg "github.com/veandco/go-sdl2/img"

type floor struct {
	texture  *sdl.Texture
	renderer *sdl.Renderer
	x        int32
	y        int32
	w        int32
	h        int32
}

func newFloor(r *sdl.Renderer) *floor {
	var b = new(floor)
	b.renderer = r
	var err error
	b.texture, err = sdlImg.LoadTexture(r, "asset/floor.png")
	checkErr(err)
	b.w = 312
	b.h = 76
	b.x = 0
	b.y = HEIGHT - b.h
	return b
}

func (self *floor) paint() {
	sdl.Main(func() {
		var err error
		err = self.renderer.Copy(self.texture, nil, &sdl.Rect{X: self.x, Y: self.y, W: self.w, H: self.h})
		checkErr(err)
	})
}
