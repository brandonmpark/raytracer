package display

import (
	"math"
	"runtime"
	"sync"

	"github.com/brandonmpark/raytracer/internal/reader"
	"github.com/brandonmpark/raytracer/pkg/objects"
	"github.com/brandonmpark/raytracer/pkg/transforms"
	mgl "github.com/go-gl/mathgl/mgl64"
)

func ray(i, j float64, scene reader.Scene) mgl.Vec3 {
	fovx := 2 * math.Atan(scene.FovY/2) * float64(scene.Width) / float64(scene.Height)
	w := scene.Center.Sub(scene.Eye).Normalize()
	u := scene.Up.Cross(w).Normalize()
	v := w.Cross(u)

	alpha := math.Tan(fovx/2) * (i - float64(scene.Width)/2) / (float64(scene.Width) / 2)
	beta := math.Tan(scene.FovY/2) * -(j - float64(scene.Height)/2) / (float64(scene.Height) / 2)

	return u.Mul(alpha).Add(v.Mul(beta)).Add(w).Normalize()
}

func trace(origin, dir mgl.Vec3, scene reader.Scene) (*objects.Shape, float64) {
	closest_i := -1
	min_t := -1.
	for i, shape := range scene.Shapes {
		t := shape.Intersect(origin, dir)
		if t > 0 && (closest_i == -1 || t < min_t) {
			closest_i = i
			min_t = t
		}
	}
	if closest_i == -1 {
		return nil, -1
	}
	return &scene.Shapes[closest_i], min_t
}

func isLit(point mgl.Vec3, light objects.Light, scene reader.Scene) bool {
	if light.IsDirectional {
		_, t_shape := trace(point.Add(light.Pos.Mul(0.01)), light.Pos, scene)
		return t_shape == -1
	}
	r := light.Pos.Sub(point)
	t_light := r.Len()
	r = r.Normalize()
	_, t_shape := trace(point.Add(r.Mul(0.01)), r, scene)
	return t_shape == -1 || t_shape >= t_light
}

func calcColor(eye, point mgl.Vec3, shape objects.Shape, depth int, scene reader.Scene) mgl.Vec3 {
	color := shape.Ambient().Add(shape.Emission())
	n := shape.Normal(point)
	e := eye.Sub(point).Normalize()

	for _, light := range scene.Lights {
		if !isLit(point, light, scene) {
			continue
		}
		l := light.Pos.Normalize()
		if !light.IsDirectional {
			l = light.Pos.Sub(point).Normalize()
		}
		h := l.Add(e).Normalize()

		lambert := transforms.MulElems(shape.Diffuse(), light.Color).Mul(math.Max(n.Dot(l), 0))
		phong := transforms.MulElems(shape.Specular(), light.Color).Mul(math.Pow(math.Max(n.Dot(h), 0), shape.Shininess()))
		if !light.IsDirectional {
			d := light.Pos.Sub(point).Len()
			atten := scene.Attenuation[0] + scene.Attenuation[1]*d + scene.Attenuation[2]*d*d
			lambert = lambert.Mul(1 / atten)
			phong = phong.Mul(1 / atten)
		}
		color = color.Add(lambert).Add(phong)
	}

	r := n.Mul(2 * n.Dot(e)).Sub(e).Normalize()
	closest, t := trace(point.Add(r.Mul(0.001)), r, scene)
	if t != -1 && depth > 0 {
		color = color.Add(transforms.MulElems(shape.Specular(), calcColor(point, point.Add(r.Mul(t)), *closest, depth-1, scene)))
	}
	color[0] = math.Min(color[0], 1)
	color[1] = math.Min(color[1], 1)
	color[2] = math.Min(color[2], 1)
	return color
}

func Draw(scene reader.Scene) [][]mgl.Vec3 {
	screen := make([][]mgl.Vec3, scene.Width)
	for i := 0; i < scene.Width; i++ {
		screen[i] = make([]mgl.Vec3, scene.Height)
	}

	var wg sync.WaitGroup
	maxRoutines := runtime.GOMAXPROCS(0)
	guard := make(chan struct{}, maxRoutines)

	for i := 0; i < scene.Width; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			guard <- struct{}{}
			for j := 0; j < scene.Height; j++ {
				r := ray(float64(x)+0.5, float64(j)+0.5, scene)
				closest, t := trace(scene.Eye, r, scene)

				if t != -1 {
					intersection := scene.Eye.Add(r.Mul(t))
					screen[x][j] = calcColor(scene.Eye, intersection, *closest, scene.MaxDepth, scene)
				}
			}
			<-guard
		}(i)
	}

	wg.Wait()
	return screen
}
