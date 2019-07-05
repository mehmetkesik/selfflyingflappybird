package main

import (
	sdlImg "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipe struct {
	renderer  *sdl.Renderer
	topTex    *sdl.Texture
	bottomTex *sdl.Texture
	x         int32
	y         int32
	w         int32
	h         int32
	pipeSpace int32
}

func newPipe(r *sdl.Renderer) *pipe {
	var b = new(pipe)
	b.renderer = r
	var err error
	b.topTex, err = sdlImg.LoadTexture(r, "asset/topPipe.png")
	checkErr(err)
	b.bottomTex, err = sdlImg.LoadTexture(r, "asset/bottomPipe.png")
	checkErr(err)
	b.w = 70
	b.h = 200
	b.pipeSpace = 200 //200 normali
	return b
}

func (self *pipe) paint() {
	sdl.Main(func() {
		var err error
		err = self.renderer.Copy(self.topTex, nil, &sdl.Rect{X: self.x, Y: self.y, W: self.w, H: self.h})
		checkErr(err)
		var y = self.h + self.pipeSpace
		var h = HEIGHT - y
		err = self.renderer.Copy(self.bottomTex, nil, &sdl.Rect{X: self.x, Y: y, W: self.w, H: h})
		checkErr(err)
	})
}
