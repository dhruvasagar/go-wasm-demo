//go:build js && wasm

package main

import (
	"math"
	"syscall/js"
	"unsafe"
)

// ============================================================================
// SHARED UTILITY FUNCTIONS
// Consolidates common functions used across different benchmark files
// ============================================================================

// Math utility functions
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// NOTE: fastSqrt function removed - replaced with math.Sqrt() for better performance
// The custom Newton-Raphson implementation was slower than the standard library

// ============================================================================
// SHARED RAY TRACING FUNCTIONS
// Common ray tracing implementations used by multiple benchmarks
// ============================================================================

// Shared constants for ray tracing
const (
	SphereX       = 0.0
	SphereY       = 0.0
	SphereZ       = -5.0
	SphereRadius2 = 1.0
	LightX        = -0.57735027
	LightY        = -0.57735027
	LightZ        = -0.57735027
	BackgroundR   = 0.2
	BackgroundG   = 0.2
	BackgroundB   = 0.8
)

// Shared ray tracing core computation (fully inlined for performance)
func computeRayColor(nx, ny float64, samples int) (float64, float64, float64) {
	var colorR, colorG, colorB float64

	for s := 0; s < samples; s++ {
		// PERFORMANCE FIX: Use math.Sqrt instead of custom Newton-Raphson
		// Ray direction normalization
		rayLenSq := nx*nx + ny*ny + 1.0
		rayLen := math.Sqrt(rayLenSq)

		invRayLen := 1.0 / rayLen
		dirX := nx * invRayLen
		dirY := ny * invRayLen
		dirZ := -1.0 * invRayLen

		// FULLY INLINED: Ray-sphere intersection
		ocX := 0.0 - SphereX
		ocY := 0.0 - SphereY
		ocZ := 0.0 - SphereZ

		rayA := dirX*dirX + dirY*dirY + dirZ*dirZ
		rayB := 2.0 * (ocX*dirX + ocY*dirY + ocZ*dirZ)
		rayC := ocX*ocX + ocY*ocY + ocZ*ocZ - SphereRadius2

		discriminant := rayB*rayB - 4.0*rayA*rayC

		if discriminant < 0 {
			// Background color
			colorR += BackgroundR
			colorG += BackgroundG
			colorB += BackgroundB
		} else {
			// PERFORMANCE FIX: Use math.Sqrt instead of custom Newton-Raphson
			sqrtDisc := math.Sqrt(discriminant)

			t := (-rayB - sqrtDisc) / (2.0 * rayA)
			if t < 0 {
				t = (-rayB + sqrtDisc) / (2.0 * rayA)
			}

			if t < 0 {
				// Behind camera
				colorR += BackgroundR
				colorG += BackgroundG
				colorB += BackgroundB
			} else {
				// FULLY INLINED: Calculate intersection point, normal, and lighting
				ix := 0.0 + t*dirX
				iy := 0.0 + t*dirY
				iz := 0.0 + t*dirZ

				normalX := ix - SphereX
				normalY := iy - SphereY
				normalZ := iz - SphereZ

				// Inlined max(0, dot)
				dot := normalX*LightX + normalY*LightY + normalZ*LightZ
				var intensity float64
				if dot > 0.0 {
					intensity = dot
				} else {
					intensity = 0.0
				}

				baseColor := 0.2 + 0.8*intensity
				colorR += baseColor * 1.0
				colorG += baseColor * 0.7
				colorB += baseColor * 0.3
			}
		}
	}

	invSamples := 1.0 / float64(samples)
	return colorR * invSamples, colorG * invSamples, colorB * invSamples
}

// ============================================================================
// SHARED RESULT CONVERSION FUNCTIONS
// Common functions for converting Go results to JavaScript efficiently
// ============================================================================

// Convert Float64 slice to JavaScript typed array with bulk copy
func createFloat64TypedArray(data []float64) js.Value {
	resultTyped := js.Global().Get("Float64Array").New(len(data))
	arrayBuffer := resultTyped.Get("buffer")
	uint8View := js.Global().Get("Uint8Array").New(arrayBuffer)

	js.CopyBytesToJS(
		uint8View,
		unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), len(data)*8),
	)

	return resultTyped
}

// Convert Int32 slice to JavaScript typed array with bulk copy
func createInt32TypedArray(data []int32) js.Value {
	resultTyped := js.Global().Get("Int32Array").New(len(data))
	arrayBuffer := resultTyped.Get("buffer")
	uint8View := js.Global().Get("Uint8Array").New(arrayBuffer)

	js.CopyBytesToJS(
		uint8View,
		unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), len(data)*4),
	)

	return resultTyped
}

// ============================================================================
// SHARED RAY TRACING IMPLEMENTATIONS
// Consolidated ray tracing functions to avoid duplication
// ============================================================================

// Single-threaded ray tracing implementation used by both files
func rayTracingSharedSingle(width, height, samples int) []float64 {
	result := make([]float64, width*height*3)

	for y := 0; y < height; y++ {
		ny := (float64(y)/float64(height))*2.0 - 1.0

		for x := 0; x < width; x++ {
			nx := (float64(x)/float64(width))*2.0 - 1.0

			colorR, colorG, colorB := computeRayColor(nx, ny, samples)

			idx := (y*width + x) * 3
			result[idx] = colorR
			result[idx+1] = colorG
			result[idx+2] = colorB
		}
	}

	return result
}
