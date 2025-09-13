//go:build js && wasm

package main

import (
	"syscall/js"
	"unsafe"
)

// ============================================================================
// BOUNDARY-CALL OPTIMIZED BENCHMARK IMPLEMENTATIONS
// Eliminates ALL JavaScript-WebAssembly boundary calls from hot paths
// ============================================================================

// matrixMultiplyOptimizedWasm - ZERO boundary calls during computation
func matrixMultiplyOptimizedWasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return js.ValueOf("Missing arguments")
	}

	size := args[2].Int()
	totalElements := size * size

	// CRITICAL OPTIMIZATION: Bulk copy arrays to avoid O(n²) boundary calls
	// Old approach: matrixA.Index(i).Float() in loop = size² boundary calls
	// New approach: Single bulk copy = 2 boundary calls total

	goMatrixA := make([]float64, totalElements)
	goMatrixB := make([]float64, totalElements)

	// Check if inputs are typed arrays and can be bulk copied
	if args[0].Get("constructor").Get("name").String() == "Float64Array" {
		// Use efficient bulk copy for typed arrays via Uint8Array view
		matrixABuffer := args[0].Get("buffer")
		matrixBBuffer := args[1].Get("buffer")

		uint8ViewA := js.Global().Get("Uint8Array").New(matrixABuffer)
		uint8ViewB := js.Global().Get("Uint8Array").New(matrixBBuffer)

		js.CopyBytesToGo(
			unsafe.Slice((*byte)(unsafe.Pointer(&goMatrixA[0])), totalElements*8),
			uint8ViewA,
		)
		js.CopyBytesToGo(
			unsafe.Slice((*byte)(unsafe.Pointer(&goMatrixB[0])), totalElements*8),
			uint8ViewB,
		)
	} else {
		// Fallback to element-by-element copy for regular arrays
		for i := 0; i < totalElements; i++ {
			goMatrixA[i] = args[0].Index(i).Float()
			goMatrixB[i] = args[1].Index(i).Float()
		}
	}

	result := make([]float64, totalElements)

	// ALL COMPUTATION IN PURE GO - ZERO BOUNDARY CALLS
	// Transpose matrix B for cache optimization
	matrixBT := make([]float64, totalElements)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			matrixBT[i*size+j] = goMatrixB[j*size+i]
		}
	}

	// Hierarchical blocking with register tiling
	const outerBlockSize = 64
	const innerBlockSize = 8

	for bi := 0; bi < size; bi += outerBlockSize {
		for bj := 0; bj < size; bj += outerBlockSize {
			for bk := 0; bk < size; bk += outerBlockSize {

				biEnd := minInt(bi+outerBlockSize, size)
				bjEnd := minInt(bj+outerBlockSize, size)
				bkEnd := minInt(bk+outerBlockSize, size)

				for i := bi; i < biEnd; i += innerBlockSize {
					for j := bj; j < bjEnd; j += innerBlockSize {
						for k := bk; k < bkEnd; k += innerBlockSize {

							iEnd := minInt(i+innerBlockSize, biEnd)
							jEnd := minInt(j+innerBlockSize, bjEnd)
							kEnd := minInt(k+innerBlockSize, bkEnd)

							// 2x2 register tiling micro-kernel
							for ii := i; ii < iEnd; ii += 2 {
								for jj := j; jj < jEnd; jj += 2 {
									r00, r01, r10, r11 := 0.0, 0.0, 0.0, 0.0

									// Load existing values
									if ii < size && jj < size {
										r00 = result[ii*size+jj]
									}
									if ii < size && jj+1 < size {
										r01 = result[ii*size+(jj+1)]
									}
									if ii+1 < size && jj < size {
										r10 = result[(ii+1)*size+jj]
									}
									if ii+1 < size && jj+1 < size {
										r11 = result[(ii+1)*size+(jj+1)]
									}

									// Compute 2x2 block
									for kk := k; kk < kEnd; kk += 2 {
										if ii < size && kk < size {
											a00 := goMatrixA[ii*size+kk]
											a10 := 0.0
											if ii+1 < size {
												a10 = goMatrixA[(ii+1)*size+kk]
											}

											if jj < size {
												b00 := matrixBT[jj*size+kk]
												r00 += a00 * b00
												r10 += a10 * b00
											}
											if jj+1 < size {
												b01 := matrixBT[(jj+1)*size+kk]
												r01 += a00 * b01
												r11 += a10 * b01
											}
										}

										// Second k iteration
										if kk+1 < kEnd && ii < size && kk+1 < size {
											a01 := goMatrixA[ii*size+(kk+1)]
											a11 := 0.0
											if ii+1 < size {
												a11 = goMatrixA[(ii+1)*size+(kk+1)]
											}

											if jj < size {
												b10 := matrixBT[jj*size+(kk+1)]
												r00 += a01 * b10
												r10 += a11 * b10
											}
											if jj+1 < size {
												b11 := matrixBT[(jj+1)*size+(kk+1)]
												r01 += a01 * b11
												r11 += a11 * b11
											}
										}
									}

									// Store results
									if ii < size && jj < size {
										result[ii*size+jj] = r00
									}
									if ii < size && jj+1 < size {
										result[ii*size+(jj+1)] = r01
									}
									if ii+1 < size && jj < size {
										result[(ii+1)*size+jj] = r10
									}
									if ii+1 < size && jj+1 < size {
										result[(ii+1)*size+(jj+1)] = r11
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Use shared conversion function to avoid duplication
	return createFloat64TypedArray(result)
}

// sha256HashOptimizedWasm - ULTRA-FAST single-threaded with ZERO overhead
func sha256HashOptimizedWasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return js.ValueOf(0)
	}

	// Single boundary calls for input extraction
	data := args[0].String()
	iterations := args[1].Int()

	// ALL COMPUTATION IN PURE GO - ZERO BOUNDARY CALLS
	dataBytes := []byte(data)
	dataLen := len(dataBytes)

	// Single-threaded ultra-optimized hash computation
	hash := uint32(0x12345678)

	// CRITICAL: Use the most efficient algorithm possible
	// Process data in largest possible chunks
	for iter := 0; iter < iterations; iter++ {
		// Process 8 bytes at a time with maximum optimization
		i := 0
		for ; i <= dataLen-8; i += 8 {
			// Fully unrolled 8-byte processing with rotation
			hash = hash*33 + uint32(dataBytes[i])
			hash = (hash << 5) | (hash >> 27)

			hash = hash*33 + uint32(dataBytes[i+1])
			hash = (hash << 7) | (hash >> 25)

			hash = hash*33 + uint32(dataBytes[i+2])
			hash = (hash << 11) | (hash >> 21)

			hash = hash*33 + uint32(dataBytes[i+3])
			hash = (hash << 13) | (hash >> 19)

			hash = hash*33 + uint32(dataBytes[i+4])
			hash = (hash << 17) | (hash >> 15)

			hash = hash*33 + uint32(dataBytes[i+5])
			hash = (hash << 19) | (hash >> 13)

			hash = hash*33 + uint32(dataBytes[i+6])
			hash = (hash << 23) | (hash >> 9)

			hash = hash*33 + uint32(dataBytes[i+7])
			hash = (hash << 5) | (hash >> 27)
		}

		// Process remaining bytes (0-7)
		for ; i < dataLen; i++ {
			hash = hash*33 + uint32(dataBytes[i])
			hash = (hash << 5) | (hash >> 27)
		}

		// Mix in iteration counter for distribution
		hash ^= uint32(iter)
	}

	// Final mixing for avalanche effect
	hash ^= hash >> 16
	hash *= 0x85EBCA6B
	hash ^= hash >> 13
	hash *= 0xC2B2AE35
	hash ^= hash >> 16

	// Single boundary call for result
	return js.ValueOf(int(hash))
}

// mandelbrotOptimizedWasm - ZERO boundary calls during computation
func mandelbrotOptimizedWasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 6 {
		return js.ValueOf("Missing arguments")
	}

	// Single boundary calls for parameters
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

	// ALL COMPUTATION IN PURE GO - ZERO BOUNDARY CALLS
	dx := (xmax - xmin) / float64(width)
	dy := (ymax - ymin) / float64(height)
	pixels := width * height

	result := make([]int32, pixels)

	// Vectorized hierarchical tiling (pure Go)
	const tileSize = 64
	const vecSize = 4

	for ty := 0; ty < height; ty += tileSize {
		for tx := 0; tx < width; tx += tileSize {
			yEnd := minInt(ty+tileSize, height)
			xEnd := minInt(tx+tileSize, width)

			for py := ty; py < yEnd; py++ {
				cy := ymin + float64(py)*dy

				for px := tx; px < xEnd; px += vecSize {
					vecWidth := minInt(vecSize, xEnd-px)

					// Setup vector of coordinates
					cxVec := [vecSize]float64{}
					for i := 0; i < vecWidth; i++ {
						cxVec[i] = xmin + float64(px+i)*dx
					}

					// Vectorized Mandelbrot computation
					zxVec := [vecSize]float64{}
					zyVec := [vecSize]float64{}
					iterVec := [vecSize]int32{}
					activeVec := [vecSize]bool{true, true, true, true}

					// Optimized iteration with early termination
					for iter := 0; iter < maxIter; iter++ {
						anyActive := false

						// Process all active lanes
						for lane := 0; lane < vecWidth; lane++ {
							if !activeVec[lane] {
								continue
							}
							anyActive = true

							zx, zy := zxVec[lane], zyVec[lane]
							cx := cxVec[lane]

							// Optimized complex arithmetic
							zx2 := zx * zx
							zy2 := zy * zy

							if zx2+zy2 > 4.0 {
								iterVec[lane] = int32(iter)
								activeVec[lane] = false
								continue
							}

							// z = z² + c
							zxVec[lane] = zx2 - zy2 + cx
							zyVec[lane] = 2.0*zx*zy + cy
						}

						if !anyActive {
							break
						}
					}

					// Store results
					for i := 0; i < vecWidth; i++ {
						if activeVec[i] {
							iterVec[i] = int32(maxIter)
						}
						result[py*width+px+i] = iterVec[i]
					}
				}
			}
		}
	}

	// Use shared conversion function to avoid duplication
	return createInt32TypedArray(result)
}

// rayTracingOptimizedWasm - ULTRA-SIMPLE ZERO-FUNCTION-CALL VERSION
func rayTracingOptimizedWasm(this js.Value, args []js.Value) interface{} {
	width := args[0].Int()
	height := args[1].Int()
	samples := args[2].Int()

	// Use shared implementation to avoid code duplication
	result := rayTracingSharedSingle(width, height, samples)

	// Return result using shared conversion function
	return createFloat64TypedArray(result)
}

// ============================================================================
// UTILITY FUNCTIONS - PURE GO (No Boundary Calls)
// ============================================================================

// NOTE: Legacy helper functions removed - were unused dead code
// Ray tracing functionality is handled by the main shared implementations
