package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
	"time"
)

const (
	UP    int = 0
	DOWN      = 1
	LEFT      = 2
	RIGHT     = 3
)

const (
	WIDTH  float32 = 800
	HEIGHT         = 450
)

const spd = 20

type Chain struct {
	head rl.Rectangle
	tail *Chain
	dir  int
}

func (c *Chain) MoveAndCollapse(dir int, head *rl.Rectangle, clicked bool) bool {
	if clicked {
		c.dir = dir
	}
	if c.dir == UP {
		c.head.Y -= spd
	}
	if c.dir == DOWN {
		c.head.Y += spd
	}
	if c.dir == LEFT {
		c.head.X -= spd
	}
	if c.dir == RIGHT {
		c.head.X += spd
	}
	var collision bool
	if &c.head != head {
		collision = rl.CheckCollisionRecs(c.head, *head)
		if collision {
			return true
		}
	} else {
		if c.head.X+c.head.Width > WIDTH || c.head.Y+c.head.Height > HEIGHT || c.head.X < 0 || c.head.Y < 0 {
			return true
		}
		head = &c.head
	}
	if c.tail != nil {
		collision = c.tail.MoveAndCollapse(c.dir, head, false)
		if collision {
			return true
		}
	}
	if !clicked {
		c.dir = dir
	}
	return false
}

func (c *Chain) AddTail() {
	if c.tail != nil {
		c.tail.AddTail()
	} else {
		bChain := c.head
		switch c.dir {
		case UP:
			bChain.Y += 20
		case DOWN:
			bChain.Y -= 20
		case LEFT:
			bChain.X += 20
		case RIGHT:
			bChain.X -= 20
		}
		c.tail = &Chain{
			head: bChain,
			tail: nil,
			dir:  c.dir,
		}
	}
}

func (c *Chain) IsEatFood(food rl.Rectangle) bool {
	collision := rl.CheckCollisionRecs(c.head, food)
	return collision
}

func renderChain(chain *Chain) {
	rl.DrawRectangleRec(chain.head, rl.Green)
	if chain.tail != nil {
		renderChain(chain.tail)
	}
}

func randFood(curr rl.Vector2) rl.Vector2 {
	x := float32(rand.Intn(int(WIDTH)-50) / 20 * 20)
	y := float32(rand.Intn(int(HEIGHT)-50) / 20 * 20)
	if x == curr.X && y == curr.Y {
		return randFood(curr)
	}
	return rl.Vector2{
		X: x,
		Y: y,
	}
}

func main() {
	rl.InitWindow(int32(WIDTH), HEIGHT, "snake")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	chainBlock := rl.NewRectangle(
		200+1,
		200+1,
		18,
		18,
	)
	player := &Chain{
		head: chainBlock,
		tail: nil,
		dir:  UP,
	}

	for i := 0; i < 5; i++ {
		player.AddTail()
	}

	dir := UP
	clicked := false
	foodPos := randFood(rl.Vector2{})
	food := rl.NewRectangle(foodPos.X, foodPos.Y, 20, 20)
	level := time.Duration(100)
	currLevel := time.Duration(50)

	startTick := time.Now()

	for !rl.WindowShouldClose() {

		// Controls
		if rl.IsKeyDown(rl.KeyLeft) && dir != RIGHT && !clicked {
			clicked = true
			dir = LEFT
		}
		if rl.IsKeyDown(rl.KeyRight) && dir != LEFT && !clicked {
			clicked = true
			dir = RIGHT
		}
		if rl.IsKeyDown(rl.KeyUp) && dir != DOWN && !clicked {
			clicked = true
			dir = UP
		}
		if rl.IsKeyDown(rl.KeyDown) && dir != UP && !clicked {
			clicked = true
			dir = DOWN
		}

		// Updates
		if time.Since(startTick) > time.Second/level*(level-currLevel) {
			collision := player.MoveAndCollapse(dir, &player.head, clicked)
			if collision {
				break
			}
			if player.IsEatFood(food) {
				player.AddTail()
				foodPos = randFood(foodPos)
				food.X = foodPos.X
				food.Y = foodPos.Y
				currLevel++
			}
			clicked = false
			startTick = time.Now()
		}

		// Renders
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawRectangleRec(food, rl.Red)
		renderChain(player)

		rl.EndDrawing()
	}
}
