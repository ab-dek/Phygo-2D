package phygo

type Manifold struct {
	BodyA        *Body
	BodyB        *Body
	Normal       Vector
	Contacts     [2]Vector
	ContactCount int
}
