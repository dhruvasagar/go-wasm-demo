//go:build js && wasm

package main

import (
	"math"
	"syscall/js"
)

// Ultra-optimized Mandelbrot calculation
func mandelbrotWasmUltra(this js.Value, args []js.Value) interface{} {
	if len(args) < 6 {
		return js.ValueOf("Missing arguments")
	}

	width := args[0].Int()
	height := args[1].Int()
	xmin := args[2].Float()
	xmax := args[3].Float()
	ymin := args[4].Float()
	ymax := args[5].Float()
	maxIter := 100
	if len(args) > 6 {
		maxIter = args[6].Int()
	}

	// Pre-compute constants
	dx := (xmax - xmin) / float64(width)
	dy := (ymax - ymin) / float64(height)

	// Allocate result array once
	jsResult := js.Global().Get("Array").New(width * height)

	// Highly optimized nested loops with minimal function calls
	idx := 0
	for py := 0; py < height; py++ {
		cy := ymin + float64(py)*dy

		for px := 0; px < width; px++ {
			cx := xmin + float64(px)*dx

			// Inline Mandelbrot calculation - no function calls
			zx := 0.0
			zy := 0.0
			iter := 0

			// Manual loop unrolling for first few iterations (common case)
			for iter < maxIter {
				zx2 := zx * zx
				zy2 := zy * zy

				if zx2+zy2 > 4.0 {
					break
				}

				newZy := zx*zy + zx*zy + cy // 2*zx*zy + cy optimized
				zx = zx2 - zy2 + cx
				zy = newZy
				iter++
			}

			// Direct assignment to JS array
			jsResult.SetIndex(idx, js.ValueOf(iter))
			idx++
		}
	}

	return jsResult
}

// Even more optimized version with better memory patterns
func mandelbrotWasmOptimized(this js.Value, args []js.Value) interface{} {
	if len(args) < 6 {
		return js.ValueOf("Missing arguments")
	}

	width := args[0].Int()
	height := args[1].Int()
	xmin := args[2].Float()
	xmax := args[3].Float()
	ymin := args[4].Float()
	ymax := args[5].Float()
	maxIter := 100
	if len(args) > 6 {
		maxIter = args[6].Int()
	}

	// Use native Go slice for computation
	result := make([]int, width*height)

	dx := (xmax - xmin) / float64(width)
	dy := (ymax - ymin) / float64(height)

	// Super optimized computation
	idx := 0
	for py := 0; py < height; py++ {
		cy := ymin + float64(py)*dy

		for px := 0; px < width; px++ {
			cx := xmin + float64(px)*dx

			// Ultra-fast mandelbrot calculation
			var zx, zy, zx2, zy2 float64
			iter := 0

			// Optimized inner loop
			for iter < maxIter {
				zx2 = zx * zx
				zy2 = zy * zy

				if zx2+zy2 > 4.0 {
					break
				}

				zy = (zx+zx)*zy + cy // 2*zx*zy + cy
				zx = zx2 - zy2 + cx
				iter++
			}

			result[idx] = iter
			idx++
		}
	}

	// Convert to JS array efficiently
	jsArray := js.Global().Get("Array").New(len(result))
	for i := 0; i < len(result); i++ {
		jsArray.SetIndex(i, js.ValueOf(result[i]))
	}

	return jsArray
}

// Raw performance version
func mandelbrotWasmRaw(this js.Value, args []js.Value) interface{} {
	width := args[0].Int()
	height := args[1].Int()
	xmin := args[2].Float()
	xmax := args[3].Float()
	ymin := args[4].Float()
	ymax := args[5].Float()
	maxIter := args[6].Int()

	// Direct computation without extra allocations
	jsArray := js.Global().Get("Array").New(width * height)

	dx := (xmax - xmin) / float64(width)
	dy := (ymax - ymin) / float64(height)

	for py := 0; py < height; py++ {
		cy := ymin + float64(py)*dy
		for px := 0; px < width; px++ {
			cx := xmin + float64(px)*dx

			zx, zy := 0.0, 0.0
			iter := 0

			for iter < maxIter && zx*zx+zy*zy <= 4.0 {
				zx, zy = zx*zx-zy*zy+cx, 2*zx*zy+cy
				iter++
			}

			jsArray.SetIndex(py*width+px, js.ValueOf(iter))
		}
	}

	return jsArray
}

// Matrix multiplication for additional benchmarking
func matrixMultiplyWasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return js.ValueOf("Missing arguments")
	}

	matrixA := args[0]
	matrixB := args[1]
	size := args[2].Int()

	result := make([]float64, size*size)

	// Optimized matrix multiplication
	for i := 0; i < size; i++ {
		for k := 0; k < size; k++ {
			aik := matrixA.Index(i*size + k).Float()
			for j := 0; j < size; j++ {
				result[i*size+j] += aik * matrixB.Index(k*size+j).Float()
			}
		}
	}

	jsArray := js.Global().Get("Array").New(len(result))
	for i, val := range result {
		jsArray.SetIndex(i, js.ValueOf(val))
	}

	return jsArray
}

// Simple hash computation
func sha256HashWasm(this js.Value, args []js.Value) interface{} {
	data := args[0].String()
	iterations := args[1].Int()

	hash := uint32(0x12345678)

	for iter := 0; iter < iterations; iter++ {
		for _, b := range []byte(data) {
			hash = hash*33 + uint32(b)
			hash = (hash << 5) | (hash >> 27)
		}
	}

	return js.ValueOf(int(hash))
}

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
