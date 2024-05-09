package reader

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/brandonmpark/raytracer/pkg/objects"
	"github.com/brandonmpark/raytracer/pkg/transforms"
	mgl "github.com/go-gl/mathgl/mgl64"
)

type Scene struct {
	Width           int
	Height          int
	MaxDepth        int
	OutputFile      string
	Attenuation     mgl.Vec3
	Eye, Up, Center mgl.Vec3
	FovY            float64
	Shapes          []objects.Shape
	Lights          []objects.Light
}

func top(stack []mgl.Mat4) mgl.Mat4 {
	return stack[len(stack)-1]
}

func push(stack *[]mgl.Mat4, m mgl.Mat4) {
	*stack = append(*stack, m)
}

func pop(stack *[]mgl.Mat4) {
	*stack = (*stack)[:len(*stack)-1]
}

func rightMultiply(m mgl.Mat4, stack *[]mgl.Mat4) {
	curr := top(*stack)
	pop(stack)
	push(stack, curr.Mul4(m))
}

func ReadFile(filename string) Scene {
	scene := Scene{OutputFile: "output.png", Attenuation: mgl.Vec3{1, 0, 0}, MaxDepth: 5}
	vertices := []mgl.Vec3{}
	ambient := mgl.Vec3{0.2, 0.2, 0.2}
	var diffuse, specular, emission mgl.Vec3
	var shininess float64

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return scene
	}
	defer file.Close()

	stack := []mgl.Mat4{mgl.Ident4()}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		str := strings.TrimSpace(scanner.Text())
		if str == "" || strings.HasPrefix(str, "#") {
			continue
		}

		parts := strings.Fields(str)
		cmd := parts[0]
		args := make([]float64, len(parts)-1)
		if cmd != "output" {
			for i, arg := range parts[1:] {
				val, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					fmt.Printf("Failed to parse argument %v: %v", arg, err)
					return scene
				}
				args[i] = val
			}
		} else {
			scene.OutputFile = parts[1]
			continue
		}

		switch cmd {
		case "size":
			scene.Width = int(args[0])
			scene.Height = int(args[1])
		case "maxdepth":
			scene.MaxDepth = int(args[0])
		case "camera":
			scene.Eye = mgl.Vec3{args[0], args[1], args[2]}
			scene.Center = mgl.Vec3{args[3], args[4], args[5]}
			scene.Up = transforms.UpVec(mgl.Vec3{args[6], args[7], args[8]}, scene.Center.Sub(scene.Eye))
			scene.FovY = args[9] * math.Pi / 180
		case "attenuation":
			scene.Attenuation = mgl.Vec3{args[0], args[1], args[2]}
		case "maxverts":
			// do nothing
		case "sphere":
			sphere := objects.NewSphere(
				mgl.Vec3{args[0], args[1], args[2]},
				args[3],
				ambient, diffuse, specular, emission,
				shininess,
				top(stack),
			)
			scene.Shapes = append(scene.Shapes, sphere)
		case "vertex":
			vertices = append(vertices, mgl.Vec3{args[0], args[1], args[2]})
		case "tri":
			triangle := objects.NewTriangle(
				vertices[int(args[0])],
				vertices[int(args[2])],
				vertices[int(args[1])],
				ambient, diffuse, specular, emission,
				shininess,
				top(stack),
			)
			scene.Shapes = append(scene.Shapes, triangle)
		case "translate":
			rightMultiply(transforms.Translate(args[0], args[1], args[2]), &stack)
		case "scale":
			rightMultiply(transforms.Scale(args[0], args[1], args[2]), &stack)
		case "rotate":
			rightMultiply(transforms.Rotate(args[3], mgl.Vec3{args[0], args[1], args[2]}).Mat4(), &stack)
		case "pushTransform":
			push(&stack, top(stack))
		case "popTransform":
			pop(&stack)
		case "directional":
			light := objects.Light{
				Pos:           mgl.Vec3{args[0], args[1], args[2]},
				Color:         mgl.Vec3{args[3], args[4], args[5]},
				IsDirectional: true,
			}
			scene.Lights = append(scene.Lights, light)
		case "point":
			light := objects.Light{
				Pos:           mgl.Vec3{args[0], args[1], args[2]},
				Color:         mgl.Vec3{args[3], args[4], args[5]},
				IsDirectional: false,
			}
			scene.Lights = append(scene.Lights, light)
		case "ambient":
			ambient = mgl.Vec3{args[0], args[1], args[2]}
		case "diffuse":
			diffuse = mgl.Vec3{args[0], args[1], args[2]}
		case "specular":
			specular = mgl.Vec3{args[0], args[1], args[2]}
		case "emission":
			emission = mgl.Vec3{args[0], args[1], args[2]}
		case "shininess":
			shininess = args[0]
		default:
			fmt.Printf("Unknown command: %v\n", cmd)
		}
	}
	return scene
}
