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

	mass            float32
	invMass         float32
	restitution     float32
	area            float32
	inertia         float32
	invInertia      float32
	staticFriction  float32
	dynamicFriction float32

	IsStatic         bool
	RotationDisabled bool
	IsOnGround       bool
	ShapeType        ShapeType
	// used for circle shapes
	Radius float32
	// used for rectangle shapes
	Width  float32
	Height float32

	vertices                [4]Vector // centered at position (0, 0)
	TransformedVertices     [4]Vector
	transformUpdateRequired bool

	aabb               AABB
	aabbUpdateRequired bool
}

func CreateBodyCircle(pos Vector, radius, density float32, isStatic bool) *Body {
	newBody := &Body{
		Position:         pos,
		restitution:      0.5,
		staticFriction:   0.6,
		dynamicFriction:  0.3,
		IsStatic:         isStatic,
		IsOnGround:       false,
		RotationDisabled: false,
		ShapeType:        CircleShape,
		Radius:           radius,
	}

	newBody.area = radius * radius * math.Pi
	newBody.mass = newBody.area * density
	newBody.inertia = (newBody.mass * radius * radius) / 2
	if !newBody.IsStatic {
		newBody.invMass = 1 / newBody.mass
		newBody.invInertia = 1 / newBody.inertia
	} else {
		newBody.invMass = 0.0
		newBody.invInertia = 0.0
	}
	newBody.transformUpdateRequired = true
	newBody.aabbUpdateRequired = true
	addBody(newBody)

	return newBody
}

func CreateBodyRectangle(pos Vector, width, height, density float32, isStatic bool) *Body {
	newBody := &Body{
		Position:         pos,
		restitution:      0.1,
		staticFriction:   0.6,
		dynamicFriction:  0.3,
		IsStatic:         isStatic,
		IsOnGround:       false,
		RotationDisabled: false,
		ShapeType:        RectangleShape,
		Width:            width,
		Height:           height,
	}
	newBody.area = height * width
	newBody.mass = newBody.area * density
	newBody.inertia = newBody.mass / 12 * (height*height + width*width)
	if !newBody.IsStatic {
		newBody.invMass = 1 / newBody.mass
		newBody.invInertia = 1 / newBody.inertia
	} else {
		newBody.invMass = 0.0
		newBody.invInertia = 0.0
	}

	newBody.vertices = createRectangleVertices(width, height)
	newBody.transformUpdateRequired = true
	newBody.aabbUpdateRequired = true
	addBody(newBody)

	return newBody
}

func createRectangleVertices(width, height float32) [4]Vector {
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

func (b *Body) SetRestitution(restitution float32) {
	b.restitution = ClampFloat(restitution, minRestitution, maxRestitution)
}

func (b *Body) SetStaticFriction(sFriction float32) {
	b.staticFriction = ClampFloat(sFriction, minFriction, maxFriction)
}

func (b *Body) SetDynamicFriction(dFriction float32) {
	b.dynamicFriction = ClampFloat(dFriction, minFriction, maxFriction)
}

func (b *Body) TransformVertices() {
	if b.transformUpdateRequired {
		transform := NewTransform(b.Position.X, b.Position.Y, b.Rotation)

		for i := range b.vertices {
			b.TransformedVertices[i] = VectorTransform(b.vertices[i], transform)
		}
	}
	b.transformUpdateRequired = false
}

func (b *Body) step(time float32, iteration int) {
	if b.IsStatic {
		return
	}

	time /= float32(iteration)

	b.IsOnGround = false

	acceleration := VectorMul(b.Force, b.invMass)
	b.Velocity.AddValue(VectorMul(acceleration, time))
	b.Velocity.AddValue(VectorMul(gravity, time))
	b.Position.AddValue(VectorMul(b.Velocity, time))
	if !b.RotationDisabled {
		b.Rotation += b.AngularVelocity * time
	}

	if !VectorNearlyEqual(b.Velocity, VectorZero()) || NearlyEqual(b.Rotation, 0.0) {
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

func (b *Body) RotateTo(amount float32) {
	b.Rotation = amount
}

func (b *Body) ApplyForce(amount Vector) {
	b.Force = amount
}

func (b *Body) getAABB() AABB {
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
		b.aabb = newAABB(minX, minY, maxX, maxY)
	}
	b.aabbUpdateRequired = false
	return b.aabb
}
