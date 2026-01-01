package main

import (
	"math"
	"math/rand"
	
	"github.com/ab-dek/Phygo-2D/phygo"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(screenWidth, screenHeight, "Phygo example")
	defer rl.CloseWindow()
	defer phygo.Close() // Clean up

	rl.SetTargetFPS(60)
	phygo.SetGravity(0, 0.3) // Initialize Gravity (x, y)

	for i := 0; i < 30; i++ {
		if i % 2 == 0 {
			// Create rectangle bodies
			r := phygo.CreateBodyRectangle(phygo.NewVector(rand.Float32()*float32(screenWidth), -rand.Float32()*float32(screenHeight-50)), 30, 30, 1, false)
			r.SetRestitution(0.3)
			r.SetDynamicFriction(0.1)
		} else {
			// Create circle bodies
			c := phygo.CreateBodyCircle(phygo.NewVector(rand.Float32()*float32(screenWidth), -rand.Float32()*float32(screenHeight-50)), float32(rand.Intn(20)+10), 1, false)
			c.SetRestitution(0.7)
			c.SetDynamicFriction(0.1)
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

		// Iterate over all bodies to draw them
		for _, b := range phygo.GetBodies() {
			if b.ShapeType == phygo.RectangleShape {
				for i := range b.GetVertices() {
					vertexA := b.GetVertices()[i]
					vertexB := b.GetVertices()[(i + 1) % len(b.GetVertices())]
					rl.DrawLineV(rl.NewVector2(vertexA.X, vertexA.Y), rl.NewVector2(vertexB.X, vertexB.Y), rl.White)
				}
			} else {
				rl.DrawCircleLines(int32(b.GetPos().X), int32(b.GetPos().Y), b.GetRadius(), rl.White)
			}
		}

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}
