//go:build js && wasm

package main

import (
	"math"
	"syscall/js"
)

// Ray tracing
func rayTracingWasm(this js.Value, args []js.Value) interface{} {
	width := args[0].Int()
	height := args[1].Int()
	samples := args[2].Int()

	result := make([]float64, width*height*3)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			nx := (float64(x)/float64(width))*2.0 - 1.0
			ny := (float64(y)/float64(height))*2.0 - 1.0

			var r, g, b float64

			for s := 0; s < samples; s++ {
				rayDir := normalize(nx, ny, -1.0)
				color := traceRay(0.0, 0.0, 0.0, rayDir[0], rayDir[1], rayDir[2])
				r += color[0]
				g += color[1]
				b += color[2]
			}

			r /= float64(samples)
			g /= float64(samples)
			b /= float64(samples)

			idx := (y*width + x) * 3
			result[idx] = r
			result[idx+1] = g
			result[idx+2] = b
		}
	}

	jsArray := js.Global().Get("Array").New(len(result))
	for i, val := range result {
		jsArray.SetIndex(i, js.ValueOf(val))
	}

	return jsArray
}

func normalize(x, y, z float64) []float64 {
	length := math.Sqrt(x*x + y*y + z*z)
	return []float64{x / length, y / length, z / length}
}

func traceRay(originX, originY, originZ, dirX, dirY, dirZ float64) []float64 {
	sphereX, sphereY, sphereZ := 0.0, 0.0, -3.0
	radius := 1.0

	oc := []float64{originX - sphereX, originY - sphereY, originZ - sphereZ}
	a := dirX*dirX + dirY*dirY + dirZ*dirZ
	b := 2.0 * (oc[0]*dirX + oc[1]*dirY + oc[2]*dirZ)
	c := oc[0]*oc[0] + oc[1]*oc[1] + oc[2]*oc[2] - radius*radius

	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return []float64{0.2, 0.3, 0.8}
	} else {
		t := (-b - math.Sqrt(discriminant)) / (2.0 * a)

		hitX := originX + t*dirX
		hitY := originY + t*dirY
		hitZ := originZ + t*dirZ

		normalX := (hitX - sphereX) / radius
		normalY := (hitY - sphereY) / radius
		normalZ := (hitZ - sphereZ) / radius

		lightDirX, lightDirY, lightDirZ := -1.0, -1.0, -1.0
		lightIntensity := math.Max(0, normalX*lightDirX+normalY*lightDirY+normalZ*lightDirZ)

		return []float64{
			0.8 * lightIntensity,
			0.2 * lightIntensity,
			0.2 * lightIntensity,
		}
	}
}
