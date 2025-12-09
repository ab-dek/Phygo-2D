package phygo

import "math"

// constants
const (
	minIterations = 1
	maxIterations = 64

	minFriction    = 0
	maxFriction    = 1
	minRestitution = 0
	maxRestitution = 1

	ppu = 50 // pixels per unit

	maxBodies   = 300
	maxManifold = 1000
)

// globals
var (
	bodies        [maxBodies]*Body
	bodyCount     = 0 // number of bodies
	gravity       = NewVector(0, 1)
	manifolds     [maxManifold]*Manifold
	manifoldCount = 0

	iterations = 32 // number of steps per frame
)

func SetIteration(i int) {
	iterations = ClampInt(i, minIterations, maxIterations)
}

func SetGravity(x, y float32) {
	gravity.X = x
	gravity.Y = y
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

func addBody(b *Body) {
	bodies[bodyCount] = b
	b.Id = getId()
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

func getId() int {
	index := -1
	for i := 0; i < maxBodies; i++ {
		currentID := i

		for j := 0; j < bodyCount; j++ {
			if bodies[j].Id == currentID {
				currentID++
				break
			}
		}

		if currentID == i {
			index = i
			break
		}
	}
	return index
}

func UpdatePhysics(time float32) {
	for i := 0; i < iterations; i++ {
		step(time, iterations)
	}
}

func step(time float32, iteration int) {
	// movement step
	for _, b := range bodies[:bodyCount] {
		b.step(time, iteration)
		b.IsOnGround = false
		b.transformVertices()
		b.updateAABB()
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
		for j := i + 1; j < bodyCount; j++ {
			bodyB := bodies[j]

			if bodyA.IsStatic && bodyB.IsStatic {
				continue
			}

			if !CheckCollisionAABBs(bodyA.aabb, bodyB.aabb) {
				continue
			}

			if ok, depth, normal := CheckCollision(bodyA, bodyB); ok {
				cntPoints, cntCount := findContactPoints(*bodyA, *bodyB)
				createManifold(bodyA, bodyB, normal, depth, cntPoints, cntCount)
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
	contactPoints := manifold.Contacts
	contactCount := manifold.ContactCount
	depth := manifold.Depth

	if !bodyA.IsOnGround {
		bodyA.IsOnGround = manifold.Normal.Y > 0
	}

	if !bodyB.IsOnGround {
		bodyB.IsOnGround = manifold.Normal.Y < 0
	}

	// separating overlapping bodies
	if bodyA.IsStatic {
		bodyB.move(VectorMul(normal, depth))
	} else if bodyB.IsStatic {
		bodyA.move(VectorMul(normal, -depth))
	} else {
		bodyA.move(VectorMul(normal, -depth/2))
		bodyB.move(VectorMul(normal, depth/2))
	}

	// applying impulse
	e := (bodyA.restitution + bodyB.restitution) / 2

	var impulseList [2]Vector
	var raList [2]Vector
	var rbList [2]Vector

	var jList [2]float32

	for i, p := range contactPoints[:contactCount] {
		ra := VectorSubtract(p, bodyA.position)
		rb := VectorSubtract(p, bodyB.position)

		raList[i] = ra
		rbList[i] = rb

		raPerp := NewVector(-ra.Y, ra.X)
		rbPerp := NewVector(-rb.Y, rb.X)

		angularVelocityA := VectorMul(raPerp, bodyA.AngularVelocity)
		angularVelocityB := VectorMul(rbPerp, bodyB.AngularVelocity)

		relativeVelocity := VectorSubtract(VectorAdd(bodyB.Velocity, angularVelocityB), VectorAdd(bodyA.Velocity, angularVelocityA))
		rvProj := VectorDotProduct(relativeVelocity, normal)

		if rvProj > 0 {
			continue
		}

		raPerpDotN := VectorDotProduct(raPerp, normal)
		rbPerpDotN := VectorDotProduct(rbPerp, normal)

		denom := bodyA.invMass + bodyB.invMass + raPerpDotN*raPerpDotN*bodyA.invInertia + rbPerpDotN*rbPerpDotN*bodyB.invInertia
		j := -(1 + e) * rvProj
		j /= denom
		j /= float32(contactCount)

		jList[i] = j

		impulseList[i] = VectorMul(normal, j)
	}

	for i, imp := range impulseList[:contactCount] {
		ra := raList[i]
		rb := rbList[i]

		bodyA.Velocity.AddValue(VectorMul(imp, -bodyA.invMass))
		if !bodyA.RotationDisabled {
			bodyA.AngularVelocity += -VectorCrossProduct(ra, imp) * bodyA.invInertia
		}
		bodyB.Velocity.AddValue(VectorMul(imp, bodyB.invMass))
		if !bodyB.RotationDisabled {
			bodyB.AngularVelocity += VectorCrossProduct(rb, imp) * bodyB.invInertia
		}
	}

	// applying friction
	var frictionImpulseList [2]Vector
	sf := (bodyA.staticFriction + bodyB.staticFriction) / 2
	df := (bodyA.dynamicFriction + bodyB.dynamicFriction) / 2

	for i, p := range contactPoints[:contactCount] {
		ra := VectorSubtract(p, bodyA.position)
		rb := VectorSubtract(p, bodyB.position)

		raList[i] = ra
		rbList[i] = rb

		raPerp := NewVector(-ra.Y, ra.X)
		rbPerp := NewVector(-rb.Y, rb.X)

		angularVelocityA := VectorMul(raPerp, bodyA.AngularVelocity)
		angularVelocityB := VectorMul(rbPerp, bodyB.AngularVelocity)

		relativeVelocity := VectorSubtract(VectorAdd(bodyB.Velocity, angularVelocityB), VectorAdd(bodyA.Velocity, angularVelocityA))
		tangent := VectorSubtract(relativeVelocity, VectorMul(normal, VectorDotProduct(relativeVelocity, normal)))

		if VectorNearlyEqual(tangent, VectorZero()) {
			continue
		} else {
			tangent = VectorNormalize(tangent)
		}

		raPerpDotT := VectorDotProduct(raPerp, tangent)
		rbPerpDotT := VectorDotProduct(rbPerp, tangent)

		denom := bodyA.invMass + bodyB.invMass + raPerpDotT*raPerpDotT*bodyA.invInertia + rbPerpDotT*rbPerpDotT*bodyB.invInertia
		jt := -VectorDotProduct(relativeVelocity, tangent)
		jt /= denom
		jt /= float32(contactCount)

		var friction Vector
		j := jList[i]

		if float32(math.Abs(float64(jt))) <= j*sf {
			friction = VectorMul(tangent, jt)
		} else {
			friction = VectorMul(tangent, -j*df)
		}

		frictionImpulseList[i] = friction
	}

	for i, f := range frictionImpulseList[:contactCount] {
		ra := raList[i]
		rb := rbList[i]

		bodyA.Velocity.AddValue(VectorMul(f, -bodyA.invMass))
		if !bodyA.RotationDisabled {
			bodyA.AngularVelocity += -VectorCrossProduct(ra, f) * bodyA.invInertia
		}
		bodyB.Velocity.AddValue(VectorMul(f, bodyB.invMass))
		if !bodyB.RotationDisabled {
			bodyB.AngularVelocity += VectorCrossProduct(rb, f) * bodyB.invInertia
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
