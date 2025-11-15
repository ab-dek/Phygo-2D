package phygo

type Manifold struct {
	BodyA        *Body
	BodyB        *Body
	Normal       Vector
	Contacts     [2]Vector
	ContactCount int
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