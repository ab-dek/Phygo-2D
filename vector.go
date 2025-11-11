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

func (v *Vector) AddValue(vAdd Vector) {
	v.X += vAdd.X
	v.Y += vAdd.Y
}

func VectorSubtract(v1, v2 Vector) Vector {
	return NewVector(v1.X-v2.X, v1.Y-v2.Y)
}

func (v *Vector) SubtractValue(vSub Vector) {
	v.X -= vSub.X
	v.Y -= vSub.Y
}

func VectorEquals(v1, v2 Vector) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}
