package phygo

type AABB struct {
	Min Vector
	Max Vector
}

func NewAABB(minX, minY, maxX, maxY float32) AABB {
	return AABB{
		NewVector(minX, minY),
		NewVector(maxX, maxY),
	}
}