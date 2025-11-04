package phygo

type Vector struct {
	X float32
	Y float32
}

func NewVector(x, y float32) Vector {
	return Vector{x, y}
}

// Vector with components value 0.0
func VectorZero() Vector {
	return NewVector(0.0, 0.0)
}

func VectorAdd(v1, v2 Vector) Vector {
	return NewVector(v1.X+v2.X, v1.Y+v2.Y)
}

func VectorAddValue(v Vector, add float32) Vector {
	return NewVector(v.X+add, v.Y+add)
}

func VectorSubtract(v1, v2 Vector) Vector {
	return NewVector(v1.X-v2.X, v1.Y-v2.Y)
}

func VectorSubtractValue(v Vector, sub float32) Vector {
	return NewVector(v.X-sub, v.Y-sub)
}