package phygo

// globals
var (
	bodies        [100]*Body
	bodyCount     = 0 // number of bodies
	gravity       = NewVector(0, 900)
	manifolds     [1000]*Manifold
	manifoldCount = 0

	iterations = 32 // number of steps per frame
)

// constants
const (
	minIterations = 1
	maxIterations = 64
)

func addBody(b *Body) {
	bodies[bodyCount] = b
	b.Id = bodyCount
	bodyCount++
}

func RemoveBody(b *Body) {
	index := -1
	for i := 0; i < bodyCount; i++ {
		if bodies[i].Id == b.Id {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}

	bodies[index] = nil

	for i := index; i+1 < bodyCount; i++ {
		bodies[i] = bodies[i+1]
	}
	bodyCount--
}

func GetBody(index int) (bool, *Body) {
	if index < 0 || index >= bodyCount {
		return false, nil
	}

	return true, bodies[index]
}

func GetBodiesCount() int {
	return bodyCount
}

func GetBodies() []*Body {
	return bodies[:bodyCount]
}

func createManifold(bodyA *Body, bodyB *Body, normal Vector, contacts [2]Vector, contactCount int) {
	newManifold := &Manifold{
		BodyA:        bodyA,
		BodyB:        bodyB,
		Normal:       normal,
		Contacts:     contacts,
		ContactCount: contactCount,
	}
	manifolds[manifoldCount] = newManifold
	manifoldCount++
}

func UpdatePhysics(time float32) {
	iteration := ClampInt(iterations, minIterations, maxIterations)
	for i := 0; i < iteration; i++ {
		step(time, iteration)
	}
}

func step(time float32, iteration int) {
	// movement step
	for _, b := range bodies[:bodyCount] {
		b.step(time, iteration)
		b.TransformVertices()
	}

	// clearing the previous step manifold list
	for i := manifoldCount - 1; i >= 0; i-- {
		if manifold := manifolds[i]; manifold != nil {
			manifold = nil
		}
	}
	manifoldCount = 0

	//collision step
	for i := 0; i < bodyCount-1; i++ {
		bodyA := bodies[i]
		aabbA := bodyA.getAABB()
		for j := i + 1; j < bodyCount; j++ {
			bodyB := bodies[j]
			aabbB := bodyB.getAABB()

			if bodyA.IsStatic && bodyB.IsStatic {
				continue
			}

			if !CheckCollisionAABBs(aabbA, aabbB) {
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

				cntPoints, cntCount := findContactPoints(*bodyA, *bodyB)
				createManifold(bodyA, bodyB, normal, cntPoints, cntCount)

			}
		}
	}

	for _, m := range manifolds[:manifoldCount] {
		resolveCollision(m)
	}
}

func resolveCollision(manifold *Manifold) {
	bodyA := manifold.BodyA
	bodyB := manifold.BodyB
	normal := manifold.Normal

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

func Close() {
	for i := manifoldCount - 1; i >= 0; i-- {
		if manifold := manifolds[i]; manifold != nil {
			manifold = nil
		}
	}

	for i := bodyCount - 1; i >= 0; i-- {
		RemoveBody(bodies[i])
	}
}