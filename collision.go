package phygo

import "math"

func CheckCollisionPolygons(polygonA, polygonB []Vector, centerA, centerB Vector) (bool, float32, Vector) {
	var depth = float32(math.MaxFloat32)
	var normal Vector

	for i := 0; i < len(polygonA); i++ {
		j := 0
		if i+1 < len(polygonA) {
			j = i + 1
		}
		vertex1 := polygonA[i]
		vertex2 := polygonA[j]

		edge := VectorSubtract(vertex2, vertex1)
		axis := NewVector(-edge.Y, edge.X)
		axis = VectorNormalize(axis)

		minA, maxA := projectVertices(polygonA, axis)
		minB, maxB := projectVertices(polygonB, axis)

		if minA >= maxB || minB >= maxA {
			return false, 0, VectorZero()
		}

		currentAxisDepth := math.Min(float64(maxB-minA), float64(maxA-minB))

		if float32(currentAxisDepth) < depth {
			depth = float32(currentAxisDepth)
			normal = axis
		}
	}

	for i := 0; i < len(polygonB); i++ {
		j := 0
		if i+1 < len(polygonB) {
			j = i + 1
		}
		vertex1 := polygonB[i]
		vertex2 := polygonB[j]

		edge := VectorSubtract(vertex2, vertex1)
		axis := NewVector(-edge.Y, edge.X)
		axis = VectorNormalize(axis)

		minA, maxA := projectVertices(polygonA, axis)
		minB, maxB := projectVertices(polygonB, axis)

		if minA >= maxB || minB >= maxA {
			return false, 0, VectorZero()
		}

		currentAxisDepth := math.Min(float64(maxB-minA), float64(maxA-minB))

		if float32(currentAxisDepth) < depth {
			depth = float32(currentAxisDepth)
			normal = axis
		}
	}

	// checking if the direction polygonA is facing polygonB is the same as the normal
	direction := VectorSubtract(centerB, centerA)
	if VectorDotProduct(direction, normal) < 0 {
		normal = VectorMul(normal, -1)
	}
	return true, depth, normal
}

func projectVertices(vertices []Vector, axis Vector) (float32, float32) {
	min := float32(math.MaxFloat32)
	max := float32(math.SmallestNonzeroFloat32)
	for _, v := range vertices {
		proj := VectorDotProduct(v, axis)
		if proj < min {
			min = proj
		}
		if proj > max {
			max = proj
		}
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

func CheckCollisionPolygonCircle(circleCenter, polygonCenter Vector, radius float32, polygon []Vector) (bool, float32, Vector) {
	var depth = float32(math.MaxFloat32)
	var normal Vector

	for i := 0; i < len(polygon); i++ {
		j := 0
		if i+1 < len(polygon) {
			j = i + 1
		}
		vertex1 := polygon[i]
		vertex2 := polygon[j]

		edge := VectorSubtract(vertex2, vertex1)
		axis := NewVector(-edge.Y, edge.X)
		axis = VectorNormalize(axis)

		minA, maxA := projectVertices(polygon, axis)
		minB, maxB := projectCircle(circleCenter, axis, radius)

		if minA >= maxB || minB >= maxA {
			return false, 0, VectorZero()
		}

		currentAxisDepth := math.Min(float64(maxB-minA), float64(maxA-minB))

		if float32(currentAxisDepth) < depth {
			depth = float32(currentAxisDepth)
			normal = axis
		}
	}

	cpIndex := findClosestPoint(circleCenter, polygon)
	axis := VectorSubtract(polygon[cpIndex], circleCenter)
	axis = VectorNormalize(axis)

	minA, maxA := projectVertices(polygon, axis)
	minB, maxB := projectCircle(circleCenter, axis, radius)

	if minA >= maxB || minB >= maxA {
		return false, 0, VectorZero()
	}

	currentAxisDepth := math.Min(float64(maxB-minA), float64(maxA-minB))

	if float32(currentAxisDepth) < depth {
		depth = float32(currentAxisDepth)
		normal = axis
	}

	// checking if the direction polygonA is facing polygonB is the same as the normal
	direction := VectorSubtract(polygonCenter, circleCenter)
	if VectorDotProduct(direction, normal) < 0 {
		normal = VectorMul(normal, -1)
	}
	return true, depth, normal
}

func projectCircle(center, axis Vector, radius float32) (float32, float32) {
	direction := VectorNormalize(axis)
	p1 := VectorAdd(center, VectorMul(direction, radius))
	p2 := VectorSubtract(center, VectorMul(direction, radius))

	proj1 := VectorDotProduct(p1, axis)
	proj2 := VectorDotProduct(p2, axis)

	if proj1 < proj2 {
		return proj1, proj2
	}
	return proj2, proj1
}

func findClosestPoint(point Vector, vertices []Vector) int {
	result := -1
	minDistance := float32(math.MaxFloat32)

	for i, v := range vertices {
		distance := VectorDistance(v, point)
		if distance < minDistance {
			minDistance = distance
			result = i
		}
	}

	return result
}

func findContactPoint(centerA, centerB Vector, radiusA float32) Vector {
	dir := VectorNormalize(VectorSubtract(centerB, centerA))
	return VectorAdd(centerA, VectorMul(dir, radiusA))
}

func findContactPoints(bodyA, bodyB Body) ([2]Vector, int) {
	var contactPoints [2]Vector
	var contactCount int = 0

	shapeA := bodyA.ShapeType
	shapeB := bodyB.ShapeType

	if shapeA == RectangleShape {
		if shapeB == RectangleShape {
		} else {
		}
	} else {
		if shapeB == RectangleShape {
		} else {
			contactPoints[0] = findContactPoint(bodyA.Position, bodyB.Position, bodyA.Radius)
			contactCount = 1
		}
	}
	return contactPoints, contactCount
}