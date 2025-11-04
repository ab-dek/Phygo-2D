package phygo

type ShapeType int

const (
	CircleShape = iota
	RectangleShape
)

type Body struct {
	Position        Vector
	Velocity        Vector
	Rotation        float32
	AngularVelocity float32

	Density     float32
	Mass        float32
	Restitution float32
	Area        float32
	IsStatic    bool
	ShapeType   ShapeType
	// used for circle shapes
	Radius float32
	// used for rectangle shapes
	Width  float32
	Height float32
}
