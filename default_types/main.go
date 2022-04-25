package main

import (
	"fmt"
	"math"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width          = 1200
	height         = 800
	stepX  float64 = 2.0 / width
	stepY  float64 = 2.0 / height
)

var (
	square = []float32{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0}
	count  int64
	steps  = 10000000
	digits = 10
)

type block struct {
	x    float64
	y    int
	w    int
	v    float64
	m    int
	wall int
}

func (b *block) hitWall() bool {
	return b.x <= float64(b.wall)
}

func (b *block) reverse() {
	b.v *= -1
}

func (b *block) update() {
	b.x += b.v
}

func (b *block) collide(b2 block) bool {
	return !(b.x+float64(b.w) < b2.x || b.x > b2.x+float64(b2.w))
}

func (b block) bounce(b2 block) float64 {

	sumM := b.m + b2.m
	newV := float64(float64(b.m-b2.m)/float64(sumM)) * b.v
	newV += float64(float64(2*b2.m)/float64(sumM)) * b2.v
	return newV
}

func main() {
	runtime.LockOSThread()

	window := initGlfw(width, height)
	defer glfw.Terminate()
	program := initOpenGL()

	block1 := block{200, 400, 100, 0, 1, 50}
	block2 := block{400, 400, 200, -1.0 / float64(steps), int(math.Pow(float64(100), float64(digits-1))), 0}

	for !window.ShouldClose() {
		for i := 0; i < steps; i++ {
			if block1.hitWall() {
				block1.reverse()
				count++
			}

			if block1.collide(block2) {
				v1 := block1.bounce(block2)
				v2 := block2.bounce(block1)
				block1.v = v1
				block2.v = v2
				count++
			}
			block1.update()
			block2.update()
		}

		convertCords(block1, block2)
		div := math.Pow(10, float64(digits-1))
		fmt.Println(float64(count) / div)
		//fmt.Println(square)
		vao := makeVao(square)
		draw(vao, window, program, square)
	}
}

func convertCords(b1, b2 block) {
	var i float32 = 1.0
	//t
	square[0] = float32((b1.x)*stepX) - i
	square[1] = float32(float64(b1.y+b1.w)*stepY) - i
	square[2] = 0

	square[3] = float32((b1.x)*stepX) - i
	square[4] = float32(float64(b1.y)*stepY) - i
	square[5] = 0

	square[6] = float32((b1.x+float64(b1.w))*stepX) - i
	square[7] = float32(float64(b1.y)*stepY) - i
	square[8] = 0
	//t
	square[9] = float32((b1.x)*stepX) - i
	square[10] = float32(float64(b1.y+b1.w)*stepY) - i
	square[11] = 0

	square[12] = float32((b1.x+float64(b1.w))*stepX) - i
	square[13] = float32(float64(b1.y+b1.w)*stepY) - i
	square[14] = 0

	square[15] = float32((b1.x+float64(b1.w))*stepX) - i
	square[16] = float32(float64(b1.y)*stepY) - i
	square[17] = 0
	//t
	square[18] = float32((b2.x)*stepX) - i
	square[19] = float32(float64(b2.y+b2.w)*stepY) - i
	square[20] = 0

	square[21] = float32((b2.x)*stepX) - i
	square[22] = float32(float64(b2.y)*stepY) - i
	square[23] = 0

	square[24] = float32((b2.x+float64(b2.w))*stepX) - i
	square[25] = float32(float64(b2.y)*stepY) - i
	square[26] = 0
	//t
	square[27] = float32((b2.x)*stepX) - i
	square[28] = float32(float64(b2.y+b2.w)*stepY) - i
	square[29] = 0

	square[30] = float32((b2.x+float64(b2.w))*stepX) - i
	square[31] = float32(float64(b2.y+b2.w)*stepY) - i
	square[32] = 0

	square[33] = float32((b2.x+float64(b2.w))*stepX) - i
	square[34] = float32(float64(b2.y)*stepY) - i
	square[35] = 0
}
