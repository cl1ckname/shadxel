package camera

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type OrbitCamera struct {
	Distance    float32
	Orientation mgl32.Quat
	Target      mgl32.Vec3
}

func NewOrbitCamera() *OrbitCamera {
	initialYaw := mgl32.QuatRotate(mgl32.DegToRad(-45), mgl32.Vec3{0, 1, 0})
	initialPitch := mgl32.QuatRotate(mgl32.DegToRad(30), mgl32.Vec3{1, 0, 0})
	return &OrbitCamera{
		Orientation: initialYaw.Mul(initialPitch),
		Distance:    6.0,
		Target:      mgl32.Vec3{0, 0, 0},
	}
}

func (c *OrbitCamera) ViewMatrix() mgl32.Mat4 {
	// Default forward direction (looking down -Z)
	forward := c.Orientation.Rotate(mgl32.Vec3{0, 0, -1})
	up := c.Orientation.Rotate(mgl32.Vec3{0, 1, 0})
	eye := c.Target.Sub(forward.Mul(c.Distance))

	return mgl32.LookAtV(eye, c.Target, up)
}

func (c *OrbitCamera) Rotate(dx, dy float32) {
	// Sensitivity
	const sensitivity = 0.005

	// Convert mouse movement into angles
	yaw := dx * sensitivity
	pitch := dy * sensitivity

	// Build rotation quaternions
	yawRot := mgl32.QuatRotate(yaw, mgl32.Vec3{0, -1, 0})     // Y-axis world up
	pitchRot := mgl32.QuatRotate(pitch, mgl32.Vec3{-1, 0, 0}) // X-axis local

	// Apply rotation
	c.Orientation = yawRot.Mul(c.Orientation)   // yaw rotates globally
	c.Orientation = c.Orientation.Mul(pitchRot) // pitch rotates locally
}

const twoPi = float32(2 * math.Pi)

func wrapAngle(angle float32) float32 {
	if angle > twoPi || angle < -twoPi {
		angle = float32(math.Mod(float64(angle), float64(twoPi)))
	}
	return angle
}
