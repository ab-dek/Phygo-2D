package main

import (
	"math"
	"math/rand"
	
	"github.com/ab-dek/phygo2d/phygo"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "Phygo example")
	defer rl.CloseWindow()
	defer phygo.Close()

	rl.SetTargetFPS(60)

	for i := 0; i < 30; i++ {
		if i % 2 == 0 {
			phygo.CreateBodyRectangle(phygo.NewVector(rand.Float32()*float32(screenWidth), -rand.Float32()*float32(screenHeight-50)), 30, 30, 0.5, false)
		} else {
			phygo.CreateBodyCircle(phygo.NewVector(rand.Float32()*float32(screenWidth), -rand.Float32()*float32(screenHeight-50)), 20, 0.1, false)
		}
	}
	
	// ground body
	phygo.CreateBodyRectangle(phygo.NewVector(float32(screenWidth)/2, float32(screenHeight)-25), float32(screenWidth), 50, 1, true)

	// walls
	phygo.CreateBodyRectangle(phygo.NewVector(0, float32(screenHeight/2)), 20, float32(screenHeight)-100, 1, true)
	phygo.CreateBodyRectangle(phygo.NewVector(float32(screenWidth), float32(screenHeight/2)), 20, float32(screenHeight)-100, 1, true)

	// slants
	slant1 := phygo.CreateBodyRectangle(phygo.NewVector(float32(screenWidth*2/3), float32(screenHeight/3)), 10, float32(screenHeight)-150, 1, true)
	slant1.RotateTo(60*math.Pi/180)
	slant2 := phygo.CreateBodyRectangle(phygo.NewVector(float32(screenWidth/3), float32(screenHeight*2/3)), 10, float32(screenHeight)-150, 1, true)
	slant2.RotateTo(-70*math.Pi/180)

	for !rl.WindowShouldClose() {
		phygo.UpdatePhysics(rl.GetFrameTime())

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		for _, b := range phygo.GetBodies() {
			if b.ShapeType == phygo.RectangleShape {
				for i := range b.TransformedVertices {
					j := 0
					vertexA := b.TransformedVertices[i]
					if i+1 < 4 {
						j = i + 1
					}
					vertexB := b.TransformedVertices[j]
					rl.DrawLineV(rl.NewVector2(vertexA.X, vertexA.Y), rl.NewVector2(vertexB.X, vertexB.Y), rl.White)
				}
			} else {
				rl.DrawCircleLines(int32(b.Position.X), int32(b.Position.Y), b.Radius, rl.White)
			}
		}

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}
