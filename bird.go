package main

import (
	sdlImg "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"gonn"
	"math/rand"
	"strconv"
	"time"
)

type bird struct {
	renderer          *sdl.Renderer
	i                 int
	x                 int32
	y                 int32
	w                 int32
	h                 int32
	angle             float64
	jumpY             int32
	velocity          int32
	isDead            bool
	score             int
	birdScore         int
	nearestPipeX      int32
	nearestPipeHeight int32
	brain             *gonn.MyNetwork
	lifeTime          int64
	inputs            [100000][]float64
	inputIndex        int
}

func newBird(scene *Scene) *bird {
	var b = new(bird)
	b.renderer = scene.renderer

	if scene.birdTextures == nil {
		scene.birdTextures = make([]*sdl.Texture, 0)
		for i := 1; i <= 4; i++ {
			texture, err := sdlImg.LoadTexture(b.renderer, "asset/player2/player"+strconv.Itoa(i)+".png")
			checkErr(err)
			scene.birdTextures = append(scene.birdTextures, texture)
		}
	}

	b.w = 75
	b.h = 70
	b.isDead = false
	b.jumpY = 15
	nn := new(gonn.MyNetwork)
	//girdiler: en yakın borunun yüksekliği,hızı,kuşun velocitysi,kuşun y değeri
	nn.Init([]int{3, 13, 1})
	b.brain = nn
	return b
}

func (self *bird) paint(scene *Scene) {
	sdl.Main(func() {
		var err error
		var hiz = 3
		err = self.renderer.CopyEx(scene.birdTextures[self.i/hiz], nil, &sdl.Rect{X: self.x, Y: self.y, W: self.w, H: self.h},
			self.angle, nil, sdl.FLIP_HORIZONTAL)
		checkErr(err)
		if self.i == 3*hiz {
			self.i = -1
		}
		self.i += 1
	})
}

func (self *bird) jump() {
	self.velocity -= self.jumpY
	if self.velocity < -int32(float64(self.jumpY)) {
		self.velocity = -int32(float64(self.jumpY))
	}
}

func (self *bird) impactAndScoreControl(pipes []*pipe, floorHeight int32, pipeSpeed int32, frameCount int64) bool {
	// kırpma işlemi
	kusX := self.x + self.w/5
	kusY := self.y + self.h/10
	kusW := self.w - self.w/3
	kusH := self.h - self.h/3

	self.nearestPipeX = WIDTH

	//boru kontrolü
	for _, pipe := range pipes {

		// en yakın olan boruyu bulmak
		if pipe.x+pipe.w > self.x {
			if pipe.x < self.nearestPipeX {
				self.nearestPipeX = pipe.x
				self.nearestPipeHeight = pipe.h
			}
		}

		//üst boruya çarpma kontrolü
		if (kusX > pipe.x || kusX+kusW > pipe.x) && (kusX < pipe.x+pipe.w || kusX+kusW < pipe.x+pipe.w) && kusY < pipe.h {
			self.lifeTime = frameCount
			return true
		}

		//üst boruya çarpma kontrolü
		var pipeY = pipe.h + pipe.pipeSpace
		// var pipeH = HEIGHT - pipeY
		if (kusX > pipe.x || kusX+kusW > pipe.x) && (kusX < pipe.x+pipe.w || kusX+kusW < pipe.x+pipe.w) && kusY+kusH > pipeY {
			self.lifeTime = frameCount
			return true
		}

		// Skor kontrolü yapılıyor
		if pipe.x+(pipe.w/2) > self.x+(self.w/2)-pipeSpeed && pipe.x+(pipe.w/2) <= self.x+(self.w/2) {
			self.score++
		}

		// yz için skor kontrolü yapılıyor
		if pipe.x+pipe.w > (self.x+5)-pipeSpeed && pipe.x+pipe.w <= (self.x+5) {
			self.birdScore++
		}
	}

	//ekran sınırı kontrolü
	if kusY < 0 || kusY+kusH > HEIGHT-floorHeight {
		self.lifeTime = frameCount
		return true
	}

	return false
}

func (self *bird) addInput(input []float64) {
	self.inputIndex++
	if self.inputIndex == len(self.inputs) {
		self.inputIndex = 0
	}
	self.inputs[self.inputIndex] = input
}

func (self *bird) reward(inputCount int) {
	var trainingSets [][][]float64

	var index int = self.inputIndex

	for j := 0; j < inputCount; j++ {

		if self.inputs[index] == nil {
			break
		}

		var output []float64
		output = self.brain.Predict(self.inputs[index])

		if output[0] > 0.5 {
			output[0] = 1
		} else {
			output[0] = 0
		}

		trainingSets = append(trainingSets, [][]float64{self.inputs[index], output})

		index--
		if index < 0 {
			index = len(self.inputs) - 1
		}
	}

	self.brain.Train(trainingSets, 1, 0.05, false)
}

func (self *bird) panish(inputCount int) {
	var trainingSets [][][]float64

	var index = self.inputIndex

	for j := 0; j < inputCount; j++ {

		if self.inputs[index] == nil {
			break
		}

		var output []float64
		output = self.brain.Predict(self.inputs[index])

		if output[0] > 0.5 {
			output[0] = 0
		} else {
			output[0] = 1
		}

		trainingSets = append(trainingSets, [][]float64{self.inputs[index], output})

		index--
		if index < 0 {
			index = len(self.inputs) - 1
		}
	}

	self.brain.Train(trainingSets, 1, 0.05, false)
}

func (self *bird) mutation(inputCount int) {
	var trainingSets [][][]float64

	var index int = self.inputIndex

	for j := 0; j < inputCount; j++ {

		if self.inputs[index] == nil {
			break
		}

		var output []float64
		output = self.brain.Predict(self.inputs[index])

		rand.Seed(time.Now().UnixNano())
		output[0] += (rand.Float64() * 2) - 1
		if output[0] > 1 {
			output[0] = 1
		} else if output[0] < 0 {
			output[0] = 0
		}

		trainingSets = append(trainingSets, [][]float64{self.inputs[index], output})

		index--
		if index < 0 {
			index = len(self.inputs) - 1
		}
	}

	self.brain.Train(trainingSets, 1, 0.05, false)
}

func (self *bird) rewardLocal(input []float64) {
	var trainingSets [][][]float64

	var output []float64
	output = self.brain.Predict(input)

	if output[0] > 0.5 {
		output[0] += 0.2
		if output[0] > 1 {
			output[0] = 1
		}
	} else {
		output[0] -= 0.2
		if output[0] < 0 {
			output[0] = 0
		}
	}

	trainingSets = append(trainingSets, [][]float64{input, output})

	self.brain.Train(trainingSets, 1, 0.05, false)
}

func (self *bird) panishLocal(input []float64) {
	var trainingSets [][][]float64

	var output []float64
	output = self.brain.Predict(input)

	if output[0] > 0.5 {
		output[0] -= 0.2
	} else {
		output[0] += 0.2
	}

	trainingSets = append(trainingSets, [][]float64{input, output})

	self.brain.Train(trainingSets, 1, 0.05, false)
}
