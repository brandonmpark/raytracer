package transforms

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl64"
)

func MulElems(a, b mgl.Vec3) mgl.Vec3 {
	return mgl.Vec3{a[0] * b[0], a[1] * b[1], a[2] * b[2]}
}

func Rotate(degrees float64, axis mgl.Vec3) mgl.Mat3 {
	a := axis.Normalize()
	theta := degrees * math.Pi / 180

	term1 := mgl.Ident3().Mul(math.Cos(theta))
	term2 := a.OuterProd3(a).Mul(1 - math.Cos(theta))
	term3 := mgl.Mat3{
		0, -a.Z(), a.Y(),
		a.Z(), 0, -a.X(),
		-a.Y(), a.X(), 0,
	}.Transpose().Mul(math.Sin(theta))
	return term1.Add(term2).Add(term3)
}

func Scale(sx, sy, sz float64) mgl.Mat4 {
	return mgl.Mat4{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}
}

func Translate(tx, ty, tz float64) mgl.Mat4 {
	return mgl.Mat4{
		1, 0, 0, tx,
		0, 1, 0, ty,
		0, 0, 1, tz,
		0, 0, 0, 1,
	}.Transpose()
}

func UpVec(up, dir mgl.Vec3) mgl.Vec3 {
	x := up.Cross(dir)
	y := x.Cross(dir)
	return y.Normalize()
}
