package camera

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type OrbitCamera struct {
	Yaw, Pitch float32
	Distance   float32
}

func NewOrbitCamera() *OrbitCamera {
	return &OrbitCamera{
		Yaw:      -math.Pi / 4,
		Pitch:    math.Pi / 6,
		Distance: 6.0,
	}
}

func (c *OrbitCamera) ViewMatrix() mgl32.Mat4 {
	x := c.Distance * float32(math.Cos(float64(c.Pitch))*math.Cos(float64(c.Yaw)))
	y := c.Distance * float32(math.Sin(float64(c.Pitch)))
	z := c.Distance * float32(math.Cos(float64(c.Pitch))*math.Sin(float64(c.Yaw)))

	eye := mgl32.Vec3{x, y, z}
	center := mgl32.Vec3{0, 0, 0}
	up := mgl32.Vec3{0, 1, 0}
	return mgl32.LookAtV(eye, center, up)
}

func (c *OrbitCamera) Rotate(dx, dy float32) {
	c.Yaw += dx * 0.005
	c.Pitch += dy * 0.005

	maxPitch := float32(math.Pi/2 - 0.1)
	if c.Pitch > maxPitch {
		c.Pitch = maxPitch
	}
	if c.Pitch < -maxPitch {
		c.Pitch = -maxPitch
	}
}
