package phygo

import "math"

type ShapeType int

const (
	CircleShape ShapeType = iota
	RectangleShape
)

type Body struct {
	Position        Vector
	Velocity        Vector
	Rotation        float32
	AngularVelocity float32

	Density     float32
	Mass        float32
	Restitution float32
	Area        float32
	IsStatic    bool
	ShapeType   ShapeType
	// used for circle shapes
	Radius float32
	// used for rectangle shapes
	Width  float32
	Height float32
}

func CreateBodyCircle(pos Vector, isStatic bool, radius, density float32) *Body {
	newBody := &Body{
		Position:        pos,
		Velocity:        Vector{},
		AngularVelocity: 0.0,
		Restitution:     0.0,
		IsStatic:        false,
		Radius:          radius,
		ShapeType:       CircleShape,
	}
	newBody.Area = radius * radius * math.Pi
	newBody.Mass = newBody.Area * density

	return newBody
}

func CreateBodyRectangle(pos Vector, isStatic bool, width, height, density float32) *Body {
	newBody := &Body{
		Position:        pos,
		Velocity:        Vector{},
		AngularVelocity: 0.0,
		Restitution:     0.0,
		IsStatic:        false,
		Width:           width,
		Height:          height,
		ShapeType:       RectangleShape,
	}
	newBody.Area = height * width
	newBody.Mass = newBody.Area * density

	return newBody
}
