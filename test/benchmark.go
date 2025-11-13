package main

import (
	"fmt"
	"math/rand"
	"phygo/phygo"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "Phygo example")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for i := 0; i < 97; i++ {
		if i % 2 == 0 {
			phygo.CreateBodyRectangle(phygo.NewVector(rand.Float32()*float32(screenWidth), rand.Float32()*float32(screenHeight-50)), 40, 40, 0.1, 0.1, false)
		} else {
			phygo.CreateBodyCircle(phygo.NewVector(rand.Float32()*float32(screenWidth), rand.Float32()*float32(screenHeight-50)), 30, 0.5, 0.1, false)
		}
	}
	
	// ground body
	phygo.CreateBodyRectangle(phygo.NewVector(float32(screenWidth)/2, float32(screenHeight)-25), float32(screenWidth), 50, 10, 1, true)

	// walls
	phygo.CreateBodyRectangle(phygo.NewVector(0, float32(screenHeight/2)), 20, float32(screenHeight)-100, 10, 1, true)
	phygo.CreateBodyRectangle(phygo.NewVector(float32(screenWidth), float32(screenHeight/2)), 20, float32(screenHeight)-100, 10, 1, true)

	var stepTime time.Duration
	var bodyCount int
	ticker := time.NewTicker(time.Second*1)
	defer ticker.Stop()

	for !rl.WindowShouldClose() {

		start := time.Now()
		phygo.UpdatePhysics(rl.GetFrameTime())
		elapsed := time.Since(start)

		select{
		case <-ticker.C:
			stepTime = elapsed
			bodyCount = phygo.GetBodiesCount()
		default:
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawFPS(10, 10)
		rl.DrawText(fmt.Sprintf("step time: %s", stepTime), 10, 30, 20, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("body count: %d", bodyCount), 10, 50, 20, rl.DarkGray)

		for _, b := range phygo.GetBodies() {
			for i := range b.TransformedVertices {
				j := 0
				vertexA := b.TransformedVertices[i]
				if i+1 < 4 {
					j = i + 1
				}
				vertexB := b.TransformedVertices[j]
				if b.ShapeType == phygo.RectangleShape {
					rl.DrawLineV(rl.NewVector2(vertexA.X, vertexA.Y), rl.NewVector2(vertexB.X, vertexB.Y), rl.Black)
				} else {
					rl.DrawCircleLines(int32(b.Position.X), int32(b.Position.Y), b.Radius, rl.Black)
				}
			}
		}
		rl.EndDrawing()
	}
}
