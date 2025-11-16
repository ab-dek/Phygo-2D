package phygo

type Manifold struct {
	BodyA        *Body
	BodyB        *Body
	Normal       Vector
	Depth 		 float32
	Contacts     [2]Vector
	ContactCount int
}

func createManifold(bodyA *Body, bodyB *Body, normal Vector, depth float32, contacts [2]Vector, contactCount int) {
	newManifold := &Manifold{
		BodyA:        bodyA,
		BodyB:        bodyB,
		Normal:       normal,
		Depth: 		  depth,
		Contacts:     contacts,
		ContactCount: contactCount,
	}
	manifolds[manifoldCount] = newManifold
	manifoldCount++
}