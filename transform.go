package phygo

import "math"

type transform struct {
	PosX float32
	PosY float32
	Sin  float32
	Cos  float32
}

func NewTransform(x, y, angle float32) transform {
	return transform{
		PosX: x,
		PosY: y,
		Sin:  float32(math.Sin(float64(angle))),
		Cos:  float32(math.Cos(float64(angle))),
	}
}

func TransformZero() transform {
	return transform{0, 0, 0, 0}
}

func VectorTransform(v Vector, t transform) Vector {
	// applying rotation and translation transformations
	return NewVector(t.Cos*v.X-t.Sin*v.Y+t.PosX,
		t.Sin*v.X+t.Cos*v.Y+t.PosY)
}
