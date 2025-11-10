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
	Width    float32
	Height   float32
	Vertices [4]Vector

	TriangleVertexIndices [6]int

	TransformedVertices     [4]Vector
	transformUpdateRequired bool
}

func CreateBodyCircle(pos Vector, isStatic bool, radius, density float32) *Body {
	newBody := &Body{
		Position:        pos,
		Velocity:        Vector{},
		AngularVelocity: 0.0,
		Restitution:     0.0,
		IsStatic:        false,
		ShapeType:       CircleShape,
		Radius:          radius,
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
		ShapeType:       RectangleShape,
		Width:           width,
		Height:          height,
	}
	newBody.Area = height * width
	newBody.Mass = newBody.Area * density
	newBody.Vertices = CreateRectangleVertices(width, height)
	newBody.TriangleVertexIndices = CreateRectangleTriangles()
	newBody.transformUpdateRequired = true

	return newBody
}

func CreateRectangleVertices(width, height float32) [4]Vector {
	left := -width / 2
	right := left + width
	bottom := -height / 2
	top := bottom + height
	return [4]Vector{
		NewVector(left, top),
		NewVector(right, top),
		NewVector(right, bottom),
		NewVector(left, bottom),
	}
}

func CreateRectangleTriangles() [6]int {
	return [6]int{0, 1, 2, 0, 2, 3}
}

func (b *Body) GetTransformedVertices() {
	if b.transformUpdateRequired {
		transform := NewTransform(b.Position.X, b.Position.Y, b.Rotation)

		for i := range b.Vertices {
			b.TransformedVertices[i] = VectorTransform(b.Vertices[i], transform)
		}
	}
	b.transformUpdateRequired = false
}

func (b *Body) Move(deltaPos Vector) {
	b.Position.AddValue(deltaPos)
	b.transformUpdateRequired = true
}

func (b *Body) MoveTo(newPos Vector) {
	b.Position = newPos
	b.transformUpdateRequired = true
}

func (b *Body) Rotate(amount float32) {
	b.Rotation += amount
	b.transformUpdateRequired = true
}