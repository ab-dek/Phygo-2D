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

func VectorMul(v Vector, scale float32) Vector {
	return NewVector(v.X*scale, v.Y*scale)
}

func VectorNormalize(v Vector) Vector {
	if l := VectorLen(v); l > 0 {
		return VectorMul(v, 1/l)
	}
	return v
}

// Returns the len squared of a vector
func VectorLenSqr(v Vector) float32 {
	return v.X*v.X + v.Y*v.Y
}

// Returns the distance squared
func VectorDistSqr(v1, v2 Vector) float32 {
	return float32(float64((v1.X-v2.X)*(v1.X-v2.X) + (v1.Y-v2.Y)*(v1.Y-v2.Y)))
}

func VectorCrossProduct(v1, v2 Vector) float32 {
	return v1.X*v2.Y - v1.Y*v2.X
}

func VectorLerp(v1, v2 Vector, amount float32) Vector {
	return NewVector(v1.X+amount*(v2.X-v1.X), v1.Y+amount*(v2.Y-v1.Y))
}

func ClampFloat(value, min, max float32) float32 {
	if min == max { 
		return min 
	} 
	if min > max { 
		return float32(math.NaN()) 
	} 
	if value < min { 
		return min 
	} 
	if value > max {
		return max 
	}
	return value
}

func ClampInt(value, min, max int) int {
	if min == max { 
		return min 
	} 
	if min > max { 
		return 0
	} 
	if value < min { 
		return min 
	} 
	if value > max {
		return max 
	}
	return value
}

func NearlyEqual(a, b float32) bool {
	return math.Abs(float64(a-b)) < 0.000001 
}

func VectorNearlyEqual(a, b Vector) bool {
	return VectorDistSqr(a, b) < 0.000001*0.000001
}