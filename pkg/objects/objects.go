package objects

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl64"
)

type Shape interface {
	Normal(point mgl.Vec3) mgl.Vec3
	Intersect(origin, ray mgl.Vec3) float64
	Ambient() mgl.Vec3
	Diffuse() mgl.Vec3
	Specular() mgl.Vec3
	Emission() mgl.Vec3
	Shininess() float64
}

type Light struct {
	Pos, Color    mgl.Vec3
	IsDirectional bool
}

type Triangle struct {
	v1, v2, v3                           mgl.Vec3
	ambient, diffuse, specular, emission mgl.Vec3
	shininess                            float64
	transform                            mgl.Mat4
	norm                                 mgl.Vec3
}

func NewTriangle(V1, V2, V3, Ambient, Diffuse, Specular, Emission mgl.Vec3, Shininess float64, Transform mgl.Mat4) Triangle {
	return Triangle{
		V1, V2, V3, Ambient, Diffuse, Specular, Emission, Shininess, Transform, V3.Sub(V1).Cross(V2.Sub(V1)).Normalize(),
	}
}

func (object Triangle) Normal(_ mgl.Vec3) mgl.Vec3 {
	transform_invT := object.transform.Transpose().Inv()
	return transform_invT.Mul4x1(object.norm.Vec4(0)).Vec3().Normalize()
}

func (object Triangle) Intersect(origin, ray mgl.Vec3) float64 {
	ray = object.transform.Inv().Mul4x1(ray.Vec4(0)).Vec3()
	eye := object.transform.Inv().Mul4x1(origin.Vec4(1)).Vec3()

	if ray.Dot(object.norm) == 0 {
		return -1
	}
	t := (object.v1.Dot(object.norm) - eye.Dot(object.norm)) / ray.Dot(object.norm)
	if t < 0 {
		return -1
	}
	p := eye.Add(ray.Mul(t))

	alpha := object.norm.Dot(object.v2.Sub(p).Cross(object.v3.Sub(p))) / object.norm.Dot(object.v2.Sub(object.v1).Cross(object.v3.Sub(object.v1)))
	beta := object.norm.Dot(object.v3.Sub(p).Cross(object.v1.Sub(p))) / object.norm.Dot(object.v2.Sub(object.v1).Cross(object.v3.Sub(object.v1)))
	if alpha >= 0 && alpha <= 1 && beta >= 0 && beta <= 1 && alpha+beta <= 1 {
		return t
	}
	return -1
}

func (object Triangle) Ambient() mgl.Vec3 {
	return object.ambient
}

func (object Triangle) Diffuse() mgl.Vec3 {
	return object.diffuse
}

func (object Triangle) Specular() mgl.Vec3 {
	return object.specular
}

func (object Triangle) Emission() mgl.Vec3 {
	return object.emission
}

func (object Triangle) Shininess() float64 {
	return object.shininess
}

type Sphere struct {
	pos                                  mgl.Vec3
	radius                               float64
	ambient, diffuse, specular, emission mgl.Vec3
	shininess                            float64
	transform                            mgl.Mat4
}

func NewSphere(pos mgl.Vec3, radius float64, ambient, diffuse, specular, emission mgl.Vec3, shininess float64, transform mgl.Mat4) Sphere {
	return Sphere{
		pos, radius, ambient, diffuse, specular, emission, shininess, transform,
	}
}

func (object Sphere) Normal(point mgl.Vec3) mgl.Vec3 {
	p_transform := object.transform.Inv().Mul4x1(point.Vec4(1)).Vec3()
	norm_transform := p_transform.Sub(object.pos).Normalize()
	transform_invT := object.transform.Transpose().Inv()
	return transform_invT.Mul4x1(norm_transform.Vec4(0)).Vec3().Normalize()
}

func (object Sphere) Intersect(origin, ray mgl.Vec3) float64 {
	ray = object.transform.Inv().Mul4x1(ray.Vec4(0)).Vec3()
	eye := object.transform.Inv().Mul4x1(origin.Vec4(1)).Vec3()

	a := ray.Dot(ray)
	b := ray.Dot(eye.Sub(object.pos)) * 2
	c := eye.Sub(object.pos).Dot(eye.Sub(object.pos)) - object.radius*object.radius

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return -1
	}

	t1 := (-b + math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b - math.Sqrt(discriminant)) / (2 * a)
	switch {
	case t1 > 0 && t2 > 0:
		return math.Min(t1, t2)
	case t1 > 0:
		return t1
	case t2 > 0:
		return t2
	default:
		return -1
	}
}

func (object Sphere) Ambient() mgl.Vec3 {
	return object.ambient
}

func (object Sphere) Diffuse() mgl.Vec3 {
	return object.diffuse
}

func (object Sphere) Specular() mgl.Vec3 {
	return object.specular
}

func (object Sphere) Emission() mgl.Vec3 {
	return object.emission
}

func (object Sphere) Shininess() float64 {
	return object.shininess
}
