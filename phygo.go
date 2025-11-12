package phygo

// globals
var (
	bodies  [10]*Body
	count   = 0 // number of bodies
	gravity = NewVector(0, 9.81)
)

func addBody(b *Body) {
	bodies[count] = b
	b.Id = count
	count++
}

func RemoveBody(b *Body) {
	index := -1
	for i := 0; i < count; i++ {
		if bodies[i].Id == b.Id {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}

	bodies[index] = nil

	for i := index; i+1 < count; i++ {
		bodies[i] = bodies[i+1]
	}
	count--
}

func GetBody(index int) (bool, *Body) {
	if index < 0 || index >= count {
		return false, nil
	}

	return true, bodies[index]
}

func GetBodiesCount() int {
	return count
}

func GetBodies() []*Body {
	return bodies[:count]
}

func Step(time float32) {
	// movement step
	for _, b := range bodies[:count] {
		b.step(time)
		b.TransformVertices()
	}

	//collision step
	for i := 0; i < count-1; i++ {
		bodyA := bodies[i]
		for j := i + 1; j < count; j++ {
			bodyB := bodies[j]

			if bodyA.IsStatic && bodyB.IsStatic {
				continue
			}

			if ok, depth, normal := collide(*bodyA, *bodyB); ok {
				if bodyA.IsStatic {
					bodyB.Move(VectorMul(normal, depth))
				} else if bodyB.IsStatic {
					bodyA.Move(VectorMul(normal, -depth))
				} else {
					bodyA.Move(VectorMul(normal, -depth/2))
					bodyB.Move(VectorMul(normal, depth/2))
				}

				resolveCollision(bodyA, bodyB, normal)
			}
		}
	}
}

func resolveCollision(bodyA, bodyB *Body, normal Vector) {
	relativeVelocity := VectorSubtract(bodyB.Velocity, bodyA.Velocity)

	if VectorDotProduct(relativeVelocity, normal) > 0 {
		return
	}

	e := min(bodyA.Restitution, bodyB.Restitution)
	j := -(1 + e) * VectorDotProduct(relativeVelocity, normal)
	j /= (bodyA.InvMass) + (bodyB.InvMass)

	bodyA.Velocity.SubtractValue(VectorMul(normal, j*bodyA.InvMass))
	bodyB.Velocity.AddValue(VectorMul(normal, j*bodyB.InvMass))
}

func collide(bodyA, bodyB Body) (bool, float32, Vector) {
	shapeA := bodyA.ShapeType
	shapeB := bodyB.ShapeType

	if shapeA == RectangleShape {
		if shapeB == RectangleShape {
			return CheckCollisionPolygons(bodyA.TransformedVertices[:], bodyB.TransformedVertices[:], bodyA.Position, bodyB.Position)
		} else {
			c, d, n := CheckCollisionPolygonCircle(bodyB.Position, bodyA.Position, bodyB.Radius, bodyA.TransformedVertices[:])
			n = VectorMul(n, -1)
			return c, d, n
		}
	} else {
		if shapeB == RectangleShape {
			return CheckCollisionPolygonCircle(bodyA.Position, bodyB.Position, bodyA.Radius, bodyB.TransformedVertices[:])
		} else {
			return CheckCollisionCircle(bodyA.Position, bodyB.Position, bodyA.Radius, bodyB.Radius)
		}
	}
}
