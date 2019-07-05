package main

import (
	"fmt"
	sdlImg "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
	"strconv"
)

const (
	WIDTH        = 1000
	HEIGHT       = 600
	FPS          = 70
	PIPE_COUNT   = 5
	PIPE_SPACING = WIDTH / (PIPE_COUNT - 1)
	FLOOR_COUNT  = 6
)

type Scene struct {
	window            *sdl.Window
	renderer          *sdl.Renderer
	bg                *sdl.Texture
	frameCount        int64
	floors            []*floor
	bird              *bird
	pipes             []*pipe
	pipeAndFloorSpeed int32
	text              *Text
	birdTextures      []*sdl.Texture
}

var epoch = 0
var previousLifeTime int64 = 0
var previousBirdScore = 0

var lifeTimes []int

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	checkErr(err)
	defer sdl.Quit()

	err = ttf.Init()
	checkErr(err)
	defer ttf.Quit()

	var scene Scene

	scene.window, scene.renderer, err = sdl.CreateWindowAndRenderer(WIDTH, HEIGHT, sdl.WINDOW_SHOWN)
	checkErr(err)
	defer scene.window.Destroy()
	defer scene.renderer.Destroy()

	start(&scene)
	restart(&scene)

	for {
		scene.frameCount++
		update(&scene)
		sdl.Delay(1000 / FPS)
	}

}

func drawBackground(scene *Scene) {
	err := scene.renderer.Copy(scene.bg, nil, nil)
	checkErr(err)
}

func start(scene *Scene) {
	scene.window.SetTitle("Yapay Sinir Ağları İle Flappy Bird")
	iconSurface, err := sdlImg.Load("asset/player/player1.png")
	checkErr(err)
	defer iconSurface.Free()
	scene.window.SetIcon(iconSurface)

	scene.bg, err = sdlImg.LoadTexture(scene.renderer, "asset/background.png")
	checkErr(err)

	scene.bird = newBird(scene)
	if _, err := os.Stat("brain.json"); !os.IsNotExist(err) {
		scene.bird.brain.Load("brain.json")
	}

	for i := 0; i < PIPE_COUNT; i++ {
		scene.pipes = append(scene.pipes, newPipe(scene.renderer))
	}

	for i := 0; i < FLOOR_COUNT; i++ {
		scene.floors = append(scene.floors, newFloor(scene.renderer))
	}

	scene.text = newText(scene.renderer)
}

func restart(scene *Scene) {
	scene.frameCount = 0
	scene.pipeAndFloorSpeed = 3
	var newInputs [100000][]float64
	scene.bird.inputs = newInputs
	scene.bird.inputIndex = 0

	scene.renderer.Clear()

	scene.bird.x = 150
	scene.bird.y = HEIGHT / 2 //int32(((HEIGHT / 2) - 250) + (i * 6))
	scene.bird.isDead = false
	scene.bird.velocity = 4
	scene.bird.score = 0
	scene.bird.inputIndex = 0
	scene.bird.lifeTime = 0

	rand.Seed(time.Now().UTC().UnixNano() /*int64(math.Pow(float64(i), 2))*/)
	for i, pipe := range scene.pipes {
		pipe.x = WIDTH + (int32(i) * PIPE_SPACING)
		pipe.y = 0
		//random
		pipe.h = int32(rand.Intn(250)) + 50
	}

	for i, floor := range scene.floors {
		floor.x = (int32(i) * floor.w)
	}

	scene.text.x = WIDTH - (12*7 + 20)
	scene.text.y = 0
	scene.text.w = 12 * 7
	scene.text.h = 24
	scene.text.text = "SKOR: 0"

	scene.renderer.Present()
}

func update(scene *Scene) {
	switch event := sdl.PollEvent().(type) {
	case *sdl.QuitEvent:
		fmt.Println("byeee")
		sdl.Quit()
		destroyAll(scene) //bellek temizleme
		scene.bird.brain.Save("brain.json")
		os.Exit(0)
	case *sdl.MouseButtonEvent:
		if event.GetType() == sdl.MOUSEBUTTONDOWN && event.Button == sdl.BUTTON_LEFT {
			//scene.bird.jump()
		}
	case *sdl.KeyboardEvent:
		if event.GetType() == 768 && event.Keysym.Sym == 32 {
			//scene.bird.jump()
		}
	}

	scene.renderer.Clear()

	drawBackground(scene)

	//boruların hızını 40 saniyede bir arttırıyoruz. max hız 6 olacak.
	/*if scene.frameCount%2400 == 0 && scene.pipeAndFloorSpeed != 6 {
		scene.pipeAndFloorSpeed += 1
	}*/

	for _, pipe := range scene.pipes {
		pipe.x -= scene.pipeAndFloorSpeed
		if pipe.x+pipe.w < 0 {
			pipe.x = WIDTH + PIPE_SPACING - pipe.w
			//random
			rand.Seed(time.Now().UTC().UnixNano() /*scene.frameCount*/)
			pipe.h = int32(rand.Intn(250) + 50)
		}
		pipe.paint()
	}

	for _, floor := range scene.floors {
		floor.x -= scene.pipeAndFloorSpeed
		if floor.x+floor.w < 0 {
			floor.x = (floor.w * (FLOOR_COUNT - 1)) - scene.pipeAndFloorSpeed
		}
		floor.paint()
	}

	//Yapay zeka oyun kontrolü

	//kuşun y değer aralığı:0-600,borunun x değeri: 150-350,borunun yükseklik aralığı: 50-300
	input1 := normalize(float64(scene.bird.y), 0, 600)
	input2 := normalize(float64(scene.bird.nearestPipeHeight), 50, 300)
	input3 := normalize(float64(scene.bird.velocity), float64(scene.bird.jumpY), 5)

	scene.bird.addInput([]float64{input1, input2, input3})
	out := scene.bird.brain.Predict([]float64{input1, input2, input3})

	if out[0] > 0.5 {
		scene.bird.jump()
	}
	//Yapay zeka oyun kontrolü son

	scene.bird.y += scene.bird.velocity //y değeri güncellendi
	if scene.bird.velocity < 5 { //aşağıya max 5 hızla düşmesi için
		scene.bird.velocity += 1
	}

	scene.bird.angle = float64(scene.bird.velocity * 5)
	if math.Abs(scene.bird.angle) > 20 {
		if scene.bird.angle < 0 {
			scene.bird.angle = -20
		} else {
			scene.bird.angle = 20
		}
	}
	scene.bird.paint(scene)

	scene.text.text = "SKOR: " + strconv.Itoa(scene.bird.score)
	scene.text.paint()

	scene.renderer.Present()

	if scene.bird.birdScore > previousBirdScore { //eğer boru geçmişse iki defa ödüllendiriyoruz.
		if previousBirdScore == 0 {
			scene.bird.reward(len(scene.bird.inputs))
			scene.bird.reward(len(scene.bird.inputs))
		} else {
			scene.bird.reward(int(PIPE_SPACING / scene.pipeAndFloorSpeed))
			scene.bird.reward(int(PIPE_SPACING / scene.pipeAndFloorSpeed))
		}
		previousBirdScore = scene.bird.birdScore //güncelleme işlemi
	}

	if scene.bird.impactAndScoreControl(scene.pipes, scene.floors[0].h, scene.pipeAndFloorSpeed, scene.frameCount) {
		scene.bird.isDead = true
		//kuş ölmüşse ceza ödül işlemleri yapacağız.
		if scene.bird.lifeTime > previousLifeTime { //daha fazla yaşamışsa ödüllendiriyoruz.
			scene.bird.reward(len(scene.bird.inputs))
		}

		if scene.bird.lifeTime <= previousLifeTime { //daha az yaşamışsa cezalandırıyoruz.
			if scene.bird.score == 0 {
				scene.bird.panish(len(scene.bird.inputs)) //çarptığı için cezalandırıyoruz.
			} else {
				scene.bird.panish(int(PIPE_SPACING / scene.pipeAndFloorSpeed)) //çarptığı için cezalandırıyoruz.
			}
		}

		previousLifeTime = scene.bird.lifeTime //güncelleme işlemi
		previousBirdScore = 0                  //güncelleme işlemi
	}

	if scene.bird.isDead {
		epoch++
		lifeTimes = append(lifeTimes, int(scene.bird.lifeTime))
		if epoch%100 == 0 {
			sort.Ints(lifeTimes)
			//fmt.Println(epoch, scene.bird.lifeTime)
			fmt.Println(epoch, "min:", lifeTimes[0], "max:", lifeTimes[len(lifeTimes)-1])
			lifeTimes = make([]int, 0)
		}
		restart(scene)
	}

}

func normalize(x, min, max float64) float64 {
	if x < min {
		x = min
	} else if x > max {
		x = max
	}
	return math.Floor(((x-min)/(max-min))*100) / 100
}

func destroyAll(scene *Scene) {
	scene.bg.Destroy()

	for _, t := range scene.birdTextures {
		t.Destroy()
	}

	for _, pipe := range scene.pipes {
		pipe.topTex.Destroy()
		pipe.bottomTex.Destroy()
	}

	for _, floor := range scene.floors {
		floor.texture.Destroy()
	}
}
