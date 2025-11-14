package phygo

import "math"

type ShapeType int

const (
	CircleShape ShapeType = iota
	RectangleShape
)

type Body struct {
	Id int

	Position        Vector
	Velocity        Vector
	Rotation        float32
	AngularVelocity float32
	Force           Vector

	Density     float32
	Mass        float32
	InvMass     float32
	Restitution float32
	Area        float32
	Inertia     float32
	InvInertia  float32

	IsStatic    bool
	ShapeType   ShapeType
	// used for circle shapes
	Radius float32
	// used for rectangle shapes
	Width  float32
	Height float32

	Vertices                [4]Vector // centered at position (0, 0)
	TransformedVertices     [4]Vector
	transformUpdateRequired bool

	aabb               AABB
	aabbUpdateRequired bool
}

func CreateBodyCircle(pos Vector, radius, density float32, restitution float32, isStatic bool) *Body {
	newBody := &Body{
		Position:    pos,
		Restitution: ClampFloat(restitution, 0, 1),
		IsStatic:    isStatic,
		ShapeType:   CircleShape,
		Radius:      radius,
	}

	newBody.Area = radius * radius * math.Pi
	newBody.Mass = newBody.Area * density
	newBody.Inertia = newBody.CalculateRotationalInertia()
	if !newBody.IsStatic {
		newBody.InvMass = 1 / newBody.Mass
		newBody.InvInertia = 1/ newBody.Inertia
	} else {
		newBody.InvMass = 0.0
		newBody.InvInertia = 0.0
	}
	newBody.transformUpdateRequired = true
	newBody.aabbUpdateRequired = true
	addBody(newBody)

	return newBody
}

func CreateBodyRectangle(pos Vector, width, height, density float32, restitution float32, isStatic bool) *Body {
	newBody := &Body{
		Position:    pos,
		Restitution: ClampFloat(restitution, 0, 1),
		IsStatic:    isStatic,
		ShapeType:   RectangleShape,
		Width:       width,
		Height:      height,
	}
	newBody.Area = height * width
	newBody.Mass = newBody.Area * density
	newBody.Inertia = newBody.CalculateRotationalInertia()
	if !newBody.IsStatic {
		newBody.InvMass = 1 / newBody.Mass
		newBody.InvInertia = 1/ newBody.Inertia
	} else {
		newBody.InvMass = 0.0
		newBody.InvInertia = 0.0
	}

	newBody.Vertices = CreateRectangleVertices(width, height)
	newBody.transformUpdateRequired = true
	newBody.aabbUpdateRequired = true
	addBody(newBody)

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

func (b *Body) TransformVertices() {
	if b.transformUpdateRequired {
		transform := NewTransform(b.Position.X, b.Position.Y, b.Rotation)

		for i := range b.Vertices {
			b.TransformedVertices[i] = VectorTransform(b.Vertices[i], transform)
		}
	}
	b.transformUpdateRequired = false
}

func (b *Body) CalculateRotationalInertia() float32 {
	if b.ShapeType == CircleShape {
		return (b.Mass*b.Radius*b.Radius)/2
	} else {
		return b.Mass/12*(b.Height*b.Height + b.Width*b.Width)
	}
}

func (b *Body) step(time float32, iteration int) {
	if b.IsStatic {
		return
	}

	time /= float32(iteration)
	
	acceleration := VectorMul(b.Force, b.InvMass)
	b.Velocity.AddValue(VectorMul(gravity, time))
	b.Velocity.AddValue(VectorMul(acceleration, time))
	b.Position.AddValue(VectorMul(b.Velocity, time))
	b.Rotation += b.AngularVelocity * time

	if b.Velocity.X != 0 || b.Velocity.Y != 0 {
		b.transformUpdateRequired = true
		b.aabbUpdateRequired = true
	}

	b.Force = VectorZero()
}

func (b *Body) Move(deltaPos Vector) {
	b.Position.AddValue(deltaPos)
}

func (b *Body) MoveTo(newPos Vector) {
	b.Position = newPos
}

func (b *Body) Rotate(amount float32) {
	b.Rotation += amount
}

func (b *Body) ApplyForce(amount Vector) {
	b.Force.AddValue(amount)
}

func (b *Body) GetAABB() AABB {
	if b.aabbUpdateRequired {
		minX := float32(math.MaxFloat32)
		minY := float32(math.MaxFloat32)
		maxX := float32(math.SmallestNonzeroFloat32)
		maxY := float32(math.SmallestNonzeroFloat32)
		if b.ShapeType == RectangleShape {
			for _, v := range b.TransformedVertices {
				if v.X < minX {
					minX = v.X
				}
				if v.X > maxX {
					maxX = v.X
				}
				if v.Y < minY {
					minY = v.Y
				}
				if v.Y > maxY {
					maxY = v.Y
				}
			}
		} else {
			minX = b.Position.X - b.Radius
			minY = b.Position.Y - b.Radius
			maxX = b.Position.X + b.Radius
			maxY = b.Position.Y + b.Radius
		}
		b.aabb = NewAABB(minX, minY, maxX, maxY)
	}
	b.aabbUpdateRequired = false
	return b.aabb
}
