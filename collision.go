package phygo

func CheckCollisionCircle(centerA, centerB Vector, radiusA, radiusB float32) (bool, float32, Vector) {
	normal := VectorSubtract(centerB, centerA)
	dist := VectorLen(normal)
	radii := radiusA + radiusB

	if dist >= radii {
		return false, 0, Vector{}
	}
	return true, radii - dist, VectorNormalize(normal)
}
