package phygo

import "math"

type ShapeType int

const (
	CircleShape ShapeType = iota
	RectangleShape
)

type Body struct {
	Id int

	position        Vector
	Velocity        Vector
	Rotation        float32
	AngularVelocity float32
	Force           Vector

	mass, invMass, restitution      float32
	area, inertia, invInertia       float32
	staticFriction, dynamicFriction float32

	IsStatic         bool
	RotationDisabled bool
	IsOnGround       bool
	UseGravity       bool
	ShapeType        ShapeType
	// used for circle shapes
	radius float32
	// used for rectangle shapes
	width  float32
	height float32

	verticesAtOrigin        [4]Vector // centered at position (0, 0)
	vertices                [4]Vector
	transformUpdateRequired bool

	aabb               AABB
	aabbUpdateRequired bool
}

func CreateBodyCircle(pos Vector, radius, density float32, isStatic bool) *Body {
	radius /= ppu

	newBody := &Body{
		position:         NewVector(pos.X/ppu, pos.Y/ppu),
		restitution:      0.0,
		staticFriction:   0.6,
		dynamicFriction:  0.3,
		IsStatic:         isStatic,
		IsOnGround:       false,
		RotationDisabled: false,
		UseGravity:       true,
		ShapeType:        CircleShape,
		radius:           radius,
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
	width /= ppu
	height /= ppu

	newBody := &Body{
		position:         NewVector(pos.X/ppu, pos.Y/ppu),
		restitution:      0.0,
		staticFriction:   0.6,
		dynamicFriction:  0.3,
		IsStatic:         isStatic,
		IsOnGround:       false,
		RotationDisabled: false,
		UseGravity:       true,
		ShapeType:        RectangleShape,
		width:            width,
		height:           height,
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

	newBody.verticesAtOrigin = createRectangleVertices(width, height)
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

func (b *Body) transformVertices() {
	if b.transformUpdateRequired {
		transform := NewTransform(b.position.X, b.position.Y, b.Rotation)

		for i := range b.verticesAtOrigin {
			b.vertices[i] = VectorTransform(b.verticesAtOrigin[i], transform)
		}
	}
	b.transformUpdateRequired = false
}

func (b *Body) step(time float32, iteration int) {
	if b.IsStatic {
		return
	}

	time /= float32(iteration)

	acceleration := VectorMul(b.Force, b.invMass)
	b.Velocity.AddValue(VectorMul(acceleration, time))
	if b.UseGravity {
		b.Velocity.AddValue(VectorMul(gravity, time))
	}
	b.position.AddValue(VectorMul(b.Velocity, ppu*time))
	if !b.RotationDisabled {
		b.Rotation += b.AngularVelocity * ppu * time
	}

	if !VectorNearlyEqual(b.Velocity, VectorZero()) || !NearlyEqual(b.Rotation, 0.0) {
		b.transformUpdateRequired = true
		b.aabbUpdateRequired = true
	}

	b.Force = VectorZero()
}

func (b *Body) Move(deltaPos Vector) {
	b.position.AddValue(VectorMul(deltaPos, 1/float32(ppu)))
	b.transformUpdateRequired = true
}

func (b *Body) move(deltaPos Vector) {
	b.position.AddValue(deltaPos)
	b.transformUpdateRequired = true
	b.aabbUpdateRequired = true
}

func (b *Body) MoveTo(newPos Vector) {
	b.position = VectorMul(newPos, 1/float32(ppu))
	b.transformUpdateRequired = true
	b.aabbUpdateRequired = true
}

func (b *Body) Rotate(amount float32) {
	b.Rotation += amount
	b.transformUpdateRequired = true
	b.aabbUpdateRequired = true
}

func (b *Body) RotateTo(amount float32) {
	b.Rotation = amount
	b.transformUpdateRequired = true
	b.aabbUpdateRequired = true
}

func (b *Body) ApplyForce(amount Vector) {
	b.Force = amount
}

func (b *Body) updateAABB() {
	if b.aabbUpdateRequired {
		minX := float32(math.MaxFloat32)
		minY := float32(math.MaxFloat32)
		maxX := float32(math.SmallestNonzeroFloat32)
		maxY := float32(math.SmallestNonzeroFloat32)
		if b.ShapeType == RectangleShape {
			for _, v := range b.vertices {
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
			minX = b.position.X - b.radius
			minY = b.position.Y - b.radius
			maxX = b.position.X + b.radius
			maxY = b.position.Y + b.radius
		}
		b.aabb = newAABB(minX, minY, maxX, maxY)
	}
	b.aabbUpdateRequired = false
}

func (b Body) GetAABB() AABB {
	b.transformVertices()

	minX := float32(math.MaxFloat32)
	minY := float32(math.MaxFloat32)
	maxX := float32(math.SmallestNonzeroFloat32)
	maxY := float32(math.SmallestNonzeroFloat32)
	if b.ShapeType == RectangleShape {
		for _, v := range b.vertices {
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
		minX = b.position.X - b.radius
		minY = b.position.Y - b.radius
		maxX = b.position.X + b.radius
		maxY = b.position.Y + b.radius
	}
	return newAABB(minX*ppu, minY*ppu, maxX*ppu, maxY*ppu)
}

func (b *Body) GetPos() Vector {
	return VectorMul(b.position, ppu)
}

func (b *Body) GetVertices() [4]Vector {
	b.transformVertices()

	var verts [4]Vector
	for i, v := range b.vertices {
		verts[i] = VectorMul(v, ppu)
	}
	return verts
}

func (b *Body) GetRadius() float32 {
	return b.radius * ppu
}

func (b *Body) GetWidth() float32 {
	return b.width * ppu
}

func (b *Body) GetHeight() float32 {
	return b.height * ppu
}
