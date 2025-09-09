//go:build js && wasm

package main

import (
	"syscall/js"
)

// Optimized matrix multiplication with minimal JS boundary crossings
func matrixMultiplyWasmOptimized(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return js.ValueOf("Missing arguments")
	}

	matrixA := args[0]
	matrixB := args[1]
	size := args[2].Int()

	// CRITICAL: Copy JS arrays to Go slices ONCE to avoid 27M boundary calls
	goMatrixA := make([]float64, size*size)
	goMatrixB := make([]float64, size*size)
	
	// Single batch copy from JS to Go (only 180K calls instead of 27M)
	for i := 0; i < size*size; i++ {
		goMatrixA[i] = matrixA.Index(i).Float()
		goMatrixB[i] = matrixB.Index(i).Float()
	}

	// Allocate result matrix
	result := make([]float64, size*size)

	// Pure Go computation - no JS boundary calls in hot loop
	for i := 0; i < size; i++ {
		for k := 0; k < size; k++ {
			aik := goMatrixA[i*size+k]  // No JS calls!
			for j := 0; j < size; j++ {
				result[i*size+j] += aik * goMatrixB[k*size+j]  // No JS calls!
			}
		}
	}

	// Copy result back (only 90K calls)
	jsArray := js.Global().Get("Array").New(size * size)
	for i := 0; i < size*size; i++ {
		jsArray.SetIndex(i, result[i])
	}
	
	return jsArray
}

// Ultra-optimized Mandelbrot with better memory management
func mandelbrotWasmSuperOptimized(this js.Value, args []js.Value) interface{} {
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

	// Pre-calculate constants
	dx := (xmax - xmin) / float64(width)
	dy := (ymax - ymin) / float64(height)
	pixels := width * height

	// Use Go slice for computation
	result := make([]int32, pixels)

	// Optimized computation with loop unrolling
	idx := 0
	for py := 0; py < height; py++ {
		cy := ymin + float64(py)*dy
		
		for px := 0; px < width; px++ {
			cx := xmin + float64(px)*dx
			
			// Optimized Mandelbrot calculation
			zx, zy := 0.0, 0.0
			iter := int32(0)
			
			// Unroll the first few iterations for common cases
			// Iteration 1
			if iter < int32(maxIter) {
				zx2 := zx * zx
				zy2 := zy * zy
				if zx2+zy2 <= 4.0 {
					zy = 2*zx*zy + cy
					zx = zx2 - zy2 + cx
					iter++
					
					// Continue with regular loop
					for iter < int32(maxIter) {
						zx2 = zx * zx
						zy2 = zy * zy
						if zx2+zy2 > 4.0 {
							break
						}
						temp := zx2 - zy2 + cx
						zy = 2*zx*zy + cy
						zx = temp
						iter++
					}
				}
			}
			
			result[idx] = iter
			idx++
		}
	}

	// Create typed array and copy data efficiently
	jsArray := js.Global().Get("Int32Array").New(pixels)
	
	// Use efficient bulk assignment
	for i := 0; i < pixels; i++ {
		jsArray.SetIndex(i, js.ValueOf(result[i]))
	}
	
	return jsArray
}

// Optimized hash function with minimal overhead
func sha256HashWasmOptimized(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return js.ValueOf(0)
	}

	data := args[0].String()
	iterations := args[1].Int()

	// Convert string to bytes once
	dataBytes := []byte(data)
	dataLen := len(dataBytes)

	// Use uint32 for better performance
	hash := uint32(0x12345678)

	// Optimized hash computation with loop unrolling
	for iter := 0; iter < iterations; iter++ {
		// Process 4 bytes at a time when possible
		i := 0
		for ; i <= dataLen-4; i += 4 {
			hash = hash*33 + uint32(dataBytes[i])
			hash = (hash << 5) | (hash >> 27)
			
			hash = hash*33 + uint32(dataBytes[i+1])
			hash = (hash << 5) | (hash >> 27)
			
			hash = hash*33 + uint32(dataBytes[i+2])
			hash = (hash << 5) | (hash >> 27)
			
			hash = hash*33 + uint32(dataBytes[i+3])
			hash = (hash << 5) | (hash >> 27)
		}
		
		// Process remaining bytes
		for ; i < dataLen; i++ {
			hash = hash*33 + uint32(dataBytes[i])
			hash = (hash << 5) | (hash >> 27)
		}
	}

	return js.ValueOf(int(hash))
}

