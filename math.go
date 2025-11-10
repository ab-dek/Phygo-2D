package phygo

import "math"

func VectorLen(v Vector) float32 {
	return float32(math.Sqrt(float64((v.X * v.X) + (v.Y * v.Y))))
}

func VectorDotProduct(v1, v2 Vector) float32 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func VectorDistance(v1, v2 Vector) float32 {
	return float32(math.Sqrt(float64((v1.X-v2.X)*(v1.X-v2.X) + (v1.Y-v2.Y)*(v1.Y-v2.Y))))
}

func VectorScale(v Vector, scale float32) Vector {
	return NewVector(v.X*scale, v.Y*scale)
}

func VectorNormalize(v Vector) Vector {
	if l := VectorLen(v); l > 0 {
		return VectorScale(v, 1/l)
	}
	return v
}

// Returns the len squared of a vector
func VectorLenSqr(v Vector) float32 {
	return v.X*v.X + v.Y*v.Y
}

func VectorCrossProduct(v1, v2 Vector) float32 {
	return v1.X*v2.Y - v1.Y*v2.X
}