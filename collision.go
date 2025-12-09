package phygo

import "math"

func CheckCollision(bodyA, bodyB *Body) (bool, float32, Vector) {
	shapeA := bodyA.ShapeType
	shapeB := bodyB.ShapeType

	if shapeA == RectangleShape {
		if shapeB == RectangleShape {
			return checkCollisionPolygons(bodyA.vertices[:], bodyB.vertices[:], bodyA.position, bodyB.position)
		} else {
			c, d, n := checkCollisionPolygonCircle(bodyB.position, bodyA.position, bodyB.radius, bodyA.vertices[:])
			n = VectorMul(n, -1)
			return c, d, n
		}
	} else {
		if shapeB == RectangleShape {
			return checkCollisionPolygonCircle(bodyA.position, bodyB.position, bodyA.radius, bodyB.vertices[:])
		} else {
			return checkCollisionCircle(bodyA.position, bodyB.position, bodyA.radius, bodyB.radius)
		}
	}
}

func checkCollisionPolygons(polygonA, polygonB []Vector, centerA, centerB Vector) (bool, float32, Vector) {
	var depth = float32(math.MaxFloat32)
	var normal Vector

	for i := 0; i < len(polygonA); i++ {
		vertex1 := polygonA[i]
		vertex2 := polygonA[(i+1)%len(polygonA)]

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
		vertex1 := polygonB[i]
		vertex2 := polygonB[(i+1)%len(polygonB)]

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

func checkCollisionCircle(centerA, centerB Vector, radiusA, radiusB float32) (bool, float32, Vector) {
	normal := VectorSubtract(centerB, centerA)
	dist := VectorLen(normal)
	radii := radiusA + radiusB

	if dist >= radii {
		return false, 0, VectorZero()
	}
	return true, radii - dist, VectorNormalize(normal)
}

func checkCollisionPolygonCircle(circleCenter, polygonCenter Vector, radius float32, polygon []Vector) (bool, float32, Vector) {
	var depth = float32(math.MaxFloat32)
	var normal Vector

	for i := 0; i < len(polygon); i++ {
		vertex1 := polygon[i]
		vertex2 := polygon[(i+1)%len(polygon)]

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

func findContactPoints(bodyA, bodyB Body) ([2]Vector, int) {
	var contactPoints [2]Vector
	var contactCount int = 0

	shapeA := bodyA.ShapeType
	shapeB := bodyB.ShapeType

	if shapeA == RectangleShape {
		if shapeB == RectangleShape {
			contactPoints, contactCount = findContactPointsPolygons(bodyA.vertices[:], bodyB.vertices[:])
		} else {
			contactPoints[0] = findContactPointCirclePolygon(bodyB.position, bodyA.vertices[:])
			contactCount = 1
		}
	} else {
		if shapeB == RectangleShape {
			contactPoints[0] = findContactPointCirclePolygon(bodyA.position, bodyB.vertices[:])
			contactCount = 1
		} else {
			contactPoints[0] = findContactPointCircles(bodyA.position, bodyB.position, bodyA.radius)
			contactCount = 1
		}
	}
	return contactPoints, contactCount
}

func findContactPointCircles(centerA, centerB Vector, radiusA float32) Vector {
	dir := VectorNormalize(VectorSubtract(centerB, centerA))
	return VectorAdd(centerA, VectorMul(dir, radiusA))
}

func findContactPointsPolygons(verticesA, verticesB []Vector) ([2]Vector, int) {
	var contactPoints [2]Vector
	contactCount := 0
	minDist := float32(math.MaxFloat32)

	for _, p := range verticesA {
		for i := range verticesB {
			va := verticesB[i]
			vb := verticesB[(i+1)%len(verticesB)]
			distSqr, contact := pointSegmentDistance(p, va, vb)

			if NearlyEqual(distSqr, minDist) {
				if !VectorNearlyEqual(contactPoints[1], contact) {
					contactCount = 2
					contactPoints[1] = contact
				}
			} else if distSqr < minDist {
				minDist = distSqr
				contactCount = 1
				contactPoints[0] = contact
			}
		}
	}

	for _, p := range verticesB {
		for i := range verticesA {
			va := verticesA[i]
			vb := verticesA[(i+1)%len(verticesA)]
			distSqr, contact := pointSegmentDistance(p, va, vb)

			if NearlyEqual(distSqr, minDist) {
				if !VectorNearlyEqual(contactPoints[1], contact) {
					contactCount = 2
					contactPoints[1] = contact
				}
			} else if distSqr < minDist {
				minDist = distSqr
				contactCount = 1
				contactPoints[0] = contact
			}
		}
	}

	return contactPoints, contactCount
}

func findContactPointCirclePolygon(circleCenter Vector, vertices []Vector) Vector {
	minDist := float32(math.MaxFloat32)
	var contactPoint Vector

	for i := range vertices {
		vertexA := vertices[i]
		vertexB := vertices[(i+1)%len(vertices)]
		distSqr, contact := pointSegmentDistance(circleCenter, vertexA, vertexB)
		if distSqr < minDist {
			minDist = distSqr
			contactPoint = contact
		}
	}
	return contactPoint
}

func pointSegmentDistance(p, a, b Vector) (float32, Vector) {
	var closestPoint Vector

	ab := VectorSubtract(b, a)
	ap := VectorSubtract(p, a)

	proj := VectorDotProduct(ap, ab)
	abLenSqr := VectorLenSqr(ab)
	d := proj / abLenSqr

	if d <= 0 {
		closestPoint = a
	} else if d >= 1 {
		closestPoint = b
	} else {
		closestPoint = VectorAdd(a, VectorMul(ab, d))
	}

	return VectorDistSqr(p, closestPoint), closestPoint
}

func CheckCollisionAABBs(a, b AABB) bool {
	if a.Max.X <= b.Min.X || b.Max.X <= a.Min.X ||
		a.Max.Y <= b.Min.Y || b.Max.Y <= a.Min.Y {
		return false
	}
	return true
}
