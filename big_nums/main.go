package main

import (
	"fmt"
	"math"
	"math/big"
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
	steps  = 100000
	digits = 7
)

type block struct {
	x    big.Float
	y    int64
	w    int64
	v    big.Float
	m    big.Int
	wall big.Float
}

func (b *block) hitWall() bool {
	res := b.x.Cmp(&b.wall)
	return res <= 0
}

func (b *block) reverse() {
	b.v.Neg(&b.v)
}

func (b *block) update() {
	b.x.Add(&b.x, &b.v)
}

func (b *block) collide(b2 block) bool {
	var op1 big.Float
	op1.SetInt64(b.w)
	op1.Add(&op1, &b.x)
	op11 := op1.Cmp(&b2.x) //-1 "<"

	var op2 big.Float
	op2.SetInt64(b2.w)
	op2.Add(&op2, &b2.x)
	op21 := op2.Cmp(&b.x)
	if op11 == -1 || op21 == -1 {
		return false
	}
	return true
}

func (b block) bounce(b2 block) big.Float {
	var sumMi big.Int
	var subMi big.Int
	var DivM big.Float

	sumM := new(big.Float).SetInt(sumMi.Add(&b.m, &b2.m))
	subM := new(big.Float).SetInt(subMi.Sub(&b.m, &b2.m))

	DivM.Quo(subM, sumM)

	part1 := new(big.Float).Mul(&DivM, &b.v)

	doubleM2 := new(big.Int).Add(&b2.m, &b2.m)
	doubleM2F := new(big.Float).SetInt(doubleM2)
	doubleM2F.Quo(doubleM2F, sumM)

	part2 := new(big.Float).Mul(doubleM2F, &b2.v)
	var newV big.Float
	newV.Add(part1, part2)

	return newV
}

func main() {
	runtime.LockOSThread()

	window := initGlfw(width, height)
	defer glfw.Terminate()
	program := initOpenGL()

	var x1 big.Float = *big.NewFloat(200)
	var v1 big.Float = *big.NewFloat(0)
	var m1 big.Int = *big.NewInt(1)
	var wall1 big.Float = *big.NewFloat(50)

	var x2 big.Float = *big.NewFloat(400)
	var v2 big.Float = *big.NewFloat(-1.0 / float64(steps))
	var m2 big.Int = *big.NewInt(int64(math.Pow(float64(100), float64(digits-1))))
	var wall2 big.Float = *big.NewFloat(0)

	block1 := block{x1, 400, 100, v1, m1, wall1}
	block2 := block{x2, 400, 200, v2, m2, wall2}

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
	b1x, _ := b1.x.Float64()
	b2x, _ := b2.x.Float64()

	//t
	square[0] = float32((b1x)*stepX) - i
	square[1] = float32(float64(b1.y+b1.w)*stepY) - i
	square[2] = 0

	square[3] = float32((b1x)*stepX) - i
	square[4] = float32(float64(b1.y)*stepY) - i
	square[5] = 0

	square[6] = float32((b1x+float64(b1.w))*stepX) - i
	square[7] = float32(float64(b1.y)*stepY) - i
	square[8] = 0
	//t
	square[9] = float32((b1x)*stepX) - i
	square[10] = float32(float64(b1.y+b1.w)*stepY) - i
	square[11] = 0

	square[12] = float32((b1x+float64(b1.w))*stepX) - i
	square[13] = float32(float64(b1.y+b1.w)*stepY) - i
	square[14] = 0

	square[15] = float32((b1x+float64(b1.w))*stepX) - i
	square[16] = float32(float64(b1.y)*stepY) - i
	square[17] = 0
	//t
	square[18] = float32((b2x)*stepX) - i
	square[19] = float32(float64(b2.y+b2.w)*stepY) - i
	square[20] = 0

	square[21] = float32((b2x)*stepX) - i
	square[22] = float32(float64(b2.y)*stepY) - i
	square[23] = 0

	square[24] = float32((b2x+float64(b2.w))*stepX) - i
	square[25] = float32(float64(b2.y)*stepY) - i
	square[26] = 0
	//t
	square[27] = float32((b2x)*stepX) - i
	square[28] = float32(float64(b2.y+b2.w)*stepY) - i
	square[29] = 0

	square[30] = float32((b2x+float64(b2.w))*stepX) - i
	square[31] = float32(float64(b2.y+b2.w)*stepY) - i
	square[32] = 0

	square[33] = float32((b2x+float64(b2.w))*stepX) - i
	square[34] = float32(float64(b2.y)*stepY) - i
	square[35] = 0
}
