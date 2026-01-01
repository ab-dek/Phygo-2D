package main

import (
	"github.com/ab-dek/Phygo-2D/phygo"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "Phygo example")
	defer rl.CloseWindow()
	defer phygo.Close() // Clean up

	phygo.SetGravity(0, 2) // Initialize Gravity (x, y)

	rl.SetTargetFPS(60)

	// Create Rectange Body
	player := phygo.CreateBodyRectangle(phygo.NewVector(float32(screenWidth)/2, 0), 45, 45, 1, false)
	player.RotationDisabled = true
	player.SetDynamicFriction(0.8)
	
	// ground body(Static)
	phygo.CreateBodyRectangle(phygo.NewVector(float32(screenWidth)/2, float32(screenHeight)-25), float32(screenWidth), 50, 1, true)
	
	// walls(Static)
	phygo.CreateBodyRectangle(phygo.NewVector(0, float32(screenHeight/2)), 20, float32(screenHeight)-100, 1, true)
	phygo.CreateBodyRectangle(phygo.NewVector(float32(screenWidth), float32(screenHeight/2)), 20, float32(screenHeight)-100, 1, true)
	
	// platforms(Static)
	phygo.CreateBodyRectangle(phygo.NewVector(float32(screenWidth*2/3), float32(screenHeight/3)), float32(screenWidth)/3, 10, 1, true)
	phygo.CreateBodyRectangle(phygo.NewVector(float32(screenWidth/3), float32(screenHeight*2/3)), float32(screenWidth)/3, 10, 1, true)
	
	for !rl.WindowShouldClose() {
		phygo.UpdatePhysics(rl.GetFrameTime())

		if rl.IsKeyDown(rl.KeyLeft) {
			player.Velocity.X = -0.15
		} else if rl.IsKeyDown(rl.KeyRight) {
			player.Velocity.X = 0.15
		}
		if rl.IsKeyPressed(rl.KeyUp) && player.IsOnGround {
			player.Velocity.Y = -0.55
		}

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
