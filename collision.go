package phygo

import "math"

func CheckCollisionPolygons(verticesA, verticesB []Vector, centerA, centerB Vector) (bool, float32, Vector) {
	var depth = float32(math.MaxFloat32)
	var normal Vector

	for i := 0; i < len(verticesA); i++ {
		j := 0
		if i+1 < len(verticesA) {
			j = i+1
		}
		vertex1 := verticesA[i]
		vertex2 := verticesA[j]

		edge := VectorSubtract(vertex2, vertex1)
		axis := NewVector(-edge.Y, edge.X)

		minA, maxA := projectVertices(verticesA, axis)
		minB, maxB := projectVertices(verticesB, axis)

		if minA >= maxB || minB >= maxA {
			return false, 0, VectorZero()
		}

		currentAxisDepth := math.Min(float64(maxB-minA), float64(maxA-minB))

		if float32(currentAxisDepth) < depth {
			depth = float32(currentAxisDepth)
			normal = axis
		}
	}

	for i := 0; i < len(verticesB); i++ {
		j := 0
		if i+1 < len(verticesB) {
			j = i+1
		}
		vertex1 := verticesB[i]
		vertex2 := verticesB[j]

		edge := VectorSubtract(vertex2, vertex1)
		axis := NewVector(-edge.Y, edge.X)

		minA, maxA := projectVertices(verticesA, axis)
		minB, maxB := projectVertices(verticesB, axis)

		if minA >= maxB || minB >= maxA {
			return false, 0, VectorZero()
		}

		currentAxisDepth := math.Min(float64(maxB-minA), float64(maxA-minB))

		if float32(currentAxisDepth) < depth {
			depth = float32(currentAxisDepth)
			normal = axis
		}
	}

	depth /= VectorLen(normal)
	normal = VectorNormalize(normal)

	// checking if the direction polygonA is facing polygonB is the same as the normal
	direction := VectorSubtract(centerB, centerA)
	if VectorDotProduct(direction, normal) < 0 {
		normal = VectorScale(normal, -1)
	}
	return true, depth, normal
}

func projectVertices(vertices []Vector, axis Vector) (float32, float32) {
	min := float32(math.MaxFloat32)
	max := float32(math.SmallestNonzeroFloat32)
	for _, v := range vertices {
		proj := VectorDotProduct(v, axis)
		if proj < min { min = proj }
		if proj > max { max = proj }
	}
	return min, max
}

func CheckCollisionCircle(centerA, centerB Vector, radiusA, radiusB float32) (bool, float32, Vector) {
	normal := VectorSubtract(centerB, centerA)
	dist := VectorLen(normal)
	radii := radiusA + radiusB

	if dist >= radii {
		return false, 0, VectorZero()
	}
	return true, radii - dist, VectorNormalize(normal)
}
