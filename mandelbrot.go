//go:build js && wasm

package main

import (
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

// normalize and traceRay functions are defined in benchmarks_comprehensive.go
