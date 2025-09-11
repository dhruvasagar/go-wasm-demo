//go:build js && wasm

package main

import (
	"sync"
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

	// CRITICAL: Return result efficiently - Create typed array and copy data
	resultTyped := js.Global().Get("Float64Array").New(totalElements)

	// Use efficient bulk copy through array buffer when possible
	arrayBuffer := resultTyped.Get("buffer")
	uint8View := js.Global().Get("Uint8Array").New(arrayBuffer)

	// Copy bytes to Uint8Array view of the Float64Array buffer
	js.CopyBytesToJS(
		uint8View,
		unsafe.Slice((*byte)(unsafe.Pointer(&result[0])), totalElements*8),
	)

	return resultTyped
}

// sha256HashOptimizedWasm - ZERO boundary calls during computation
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

	const numLanes = 4
	hashLanes := [numLanes]uint32{
		0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476,
	}

	// Pre-process data into 32-bit words (pure Go)
	dataWords := make([]uint32, (dataLen+3)/4)
	for i := 0; i < len(dataWords); i++ {
		wordIdx := i * 4
		if wordIdx+3 < dataLen {
			dataWords[i] = uint32(dataBytes[wordIdx]) |
				uint32(dataBytes[wordIdx+1])<<8 |
				uint32(dataBytes[wordIdx+2])<<16 |
				uint32(dataBytes[wordIdx+3])<<24
		} else {
			for b := 0; b < 4 && wordIdx+b < dataLen; b++ {
				dataWords[i] |= uint32(dataBytes[wordIdx+b]) << (b * 8)
			}
		}
	}

	// Advanced mixing constants
	const c1, c2, c3, c4 = 0x85EBCA6B, 0xC2B2AE35, 0xCC9E2D51, 0x1B873593
	baseIter := iterations / numLanes

	var wg sync.WaitGroup
	results := make([]uint32, numLanes)

	// Pure Go goroutines - no boundary calls
	for lane := 0; lane < numLanes; lane++ {
		wg.Add(1)
		go func(laneID int) {
			defer wg.Done()

			startIter := laneID * baseIter
			endIter := startIter + baseIter
			if laneID == numLanes-1 {
				endIter = iterations
			}

			hash := hashLanes[laneID]

			// Heavily optimized pure Go computation
			for iter := startIter; iter < endIter; iter++ {
				iterSeed := uint32(iter)*0x9E3779B9 + uint32(laneID)*c1

				// Unrolled word processing for maximum performance
				i := 0
				for ; i <= len(dataWords)-8; i += 8 {
					// Process 8 words with optimized mixing
					w0 := dataWords[i] * c1
					w0 = (w0 << 15) | (w0 >> 17)
					w0 *= c2
					hash ^= w0
					hash = ((hash<<13)|(hash>>19))*5 + 0xE6546B64

					w1 := dataWords[i+1] * c3
					w1 = (w1 << 17) | (w1 >> 15)
					w1 *= c4
					hash ^= w1
					hash = ((hash<<11)|(hash>>21))*3 + 0xE6546B64

					w2 := dataWords[i+2] * c1
					w2 = (w2 << 19) | (w2 >> 13)
					w2 *= c2
					hash ^= w2
					hash = ((hash<<7)|(hash>>25))*7 + 0xE6546B64

					w3 := dataWords[i+3] * c3
					w3 = (w3 << 13) | (w3 >> 19)
					w3 *= c4
					hash ^= w3
					hash = ((hash<<17)|(hash>>15))*11 + 0xE6546B64

					w4 := dataWords[i+4] * c1
					w4 = (w4 << 21) | (w4 >> 11)
					w4 *= c2
					hash ^= w4
					hash = ((hash<<9)|(hash>>23))*13 + 0xE6546B64

					w5 := dataWords[i+5] * c3
					w5 = (w5 << 23) | (w5 >> 9)
					w5 *= c4
					hash ^= w5
					hash = ((hash<<15)|(hash>>17))*17 + 0xE6546B64

					w6 := dataWords[i+6] * c1
					w6 = (w6 << 11) | (w6 >> 21)
					w6 *= c2
					hash ^= w6
					hash = ((hash<<19)|(hash>>13))*19 + 0xE6546B64

					w7 := dataWords[i+7] * c3
					w7 = (w7 << 7) | (w7 >> 25)
					w7 *= c4
					hash ^= w7
					hash = ((hash<<23)|(hash>>9))*23 + 0xE6546B64
				}

				// Handle remaining words
				for ; i < len(dataWords); i++ {
					w := dataWords[i] * c1
					w = (w << 15) | (w >> 17)
					hash ^= w * c2
					hash = ((hash << 13) | (hash >> 19)) + 0xE6546B64
				}

				hash ^= iterSeed
				hash = hash*c1 + c2
			}

			results[laneID] = hash
		}(lane)
	}

	wg.Wait()

	// Pure Go result combination
	finalHash := results[0]
	for i := 1; i < numLanes; i++ {
		finalHash ^= results[i]
		finalHash = finalHash*c1 + c2
		finalHash = (finalHash << 16) | (finalHash >> 16)
	}

	// Final avalanche
	finalHash ^= finalHash >> 16
	finalHash *= c1
	finalHash ^= finalHash >> 13
	finalHash *= c2
	finalHash ^= finalHash >> 16

	// Single boundary call for result
	return js.ValueOf(int(finalHash))
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

	// CRITICAL: Return result efficiently - Create typed array and copy data
	resultTyped := js.Global().Get("Int32Array").New(pixels)

	// Use efficient bulk copy through array buffer
	arrayBuffer := resultTyped.Get("buffer")
	uint8View := js.Global().Get("Uint8Array").New(arrayBuffer)

	// Copy bytes to Uint8Array view of the Int32Array buffer
	js.CopyBytesToJS(
		uint8View,
		unsafe.Slice((*byte)(unsafe.Pointer(&result[0])), pixels*4),
	)

	return resultTyped
}

// rayTracingOptimizedWasm - ULTRA-SIMPLE ZERO-FUNCTION-CALL VERSION
func rayTracingOptimizedWasm(this js.Value, args []js.Value) interface{} {
	width := args[0].Int()
	height := args[1].Int()
	samples := args[2].Int()

	result := make([]float64, width*height*3)

	// Sphere properties (consistent with JavaScript)
	const sphereX, sphereY, sphereZ = 0.0, 0.0, -5.0
	const sphereRadius2 = 1.0
	const lightX, lightY, lightZ = -0.57735027, -0.57735027, -0.57735027

	// ULTRA-SIMPLE LOOP - NO TILING, NO BATCHING, NO FUNCTION CALLS
	for y := 0; y < height; y++ {
		ny := (float64(y)/float64(height))*2.0 - 1.0

		for x := 0; x < width; x++ {
			nx := (float64(x)/float64(width))*2.0 - 1.0

			var colorR, colorG, colorB float64

			for s := 0; s < samples; s++ {
				// FULLY INLINED: Ray direction normalization
				rayLenSq := nx*nx + ny*ny + 1.0

				// Inlined fast square root
				var rayLen float64
				if rayLenSq <= 0 {
					rayLen = 0
				} else if rayLenSq == 1 {
					rayLen = 1
				} else {
					guess := rayLenSq * 0.5
					guess = 0.5 * (guess + rayLenSq/guess)
					rayLen = 0.5 * (guess + rayLenSq/guess)
				}

				invRayLen := 1.0 / rayLen
				dirX := nx * invRayLen
				dirY := ny * invRayLen
				dirZ := -1.0 * invRayLen

				// FULLY INLINED: Ray-sphere intersection
				ocX := 0.0 - sphereX
				ocY := 0.0 - sphereY
				ocZ := 0.0 - sphereZ

				rayA := dirX*dirX + dirY*dirY + dirZ*dirZ
				rayB := 2.0 * (ocX*dirX + ocY*dirY + ocZ*dirZ)
				rayC := ocX*ocX + ocY*ocY + ocZ*ocZ - sphereRadius2

				discriminant := rayB*rayB - 4.0*rayA*rayC

				if discriminant < 0 {
					// Background color
					colorR += 0.2
					colorG += 0.2
					colorB += 0.8
				} else {
					// Hit the sphere - inlined square root
					var sqrtDisc float64
					if discriminant <= 0 {
						sqrtDisc = 0
					} else if discriminant == 1 {
						sqrtDisc = 1
					} else {
						guess := discriminant * 0.5
						guess = 0.5 * (guess + discriminant/guess)
						sqrtDisc = 0.5 * (guess + discriminant/guess)
					}

					t := (-rayB - sqrtDisc) / (2.0 * rayA)
					if t < 0 {
						t = (-rayB + sqrtDisc) / (2.0 * rayA)
					}

					if t < 0 {
						// Behind camera
						colorR += 0.2
						colorG += 0.2
						colorB += 0.8
					} else {
						// FULLY INLINED: Calculate intersection, normal, lighting
						ix := 0.0 + t*dirX
						iy := 0.0 + t*dirY
						iz := 0.0 + t*dirZ

						normalX := ix - sphereX
						normalY := iy - sphereY
						normalZ := iz - sphereZ

						// Inlined max(0, dot)
						dot := normalX*lightX + normalY*lightY + normalZ*lightZ
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
			idx := (y*width + x) * 3
			result[idx] = colorR * invSamples
			result[idx+1] = colorG * invSamples
			result[idx+2] = colorB * invSamples
		}
	}

	// Efficient bulk copy
	resultTyped := js.Global().Get("Float64Array").New(len(result))
	arrayBuffer := resultTyped.Get("buffer")
	uint8View := js.Global().Get("Uint8Array").New(arrayBuffer)

	js.CopyBytesToJS(
		uint8View,
		unsafe.Slice((*byte)(unsafe.Pointer(&result[0])), len(result)*8),
	)

	return resultTyped
}

// ============================================================================
// UTILITY FUNCTIONS - PURE GO (No Boundary Calls)
// ============================================================================

func traceRayOptimized(ox, oy, oz, dx, dy, dz float64) [3]float64 {
	// Sphere at (0, 0, -5) with radius 1
	const sphereX, sphereY, sphereZ = 0.0, 0.0, -5.0
	const sphereRadius2 = 1.0

	// Optimized ray-sphere intersection
	ocx := ox - sphereX
	ocy := oy - sphereY
	ocz := oz - sphereZ

	a := dx*dx + dy*dy + dz*dz
	b := 2.0 * (ocx*dx + ocy*dy + ocz*dz)
	c := ocx*ocx + ocy*ocy + ocz*ocz - sphereRadius2

	discriminant := b*b - 4.0*a*c
	if discriminant < 0 {
		return [3]float64{0.2, 0.2, 0.8} // Background
	}

	sqrtDisc := fastSqrtOptimized(discriminant)
	t1 := (-b - sqrtDisc) / (2.0 * a)
	t2 := (-b + sqrtDisc) / (2.0 * a)

	t := t1
	if t < 0 {
		t = t2
	}
	if t < 0 {
		return [3]float64{0.2, 0.2, 0.8} // Behind camera
	}

	// Intersection point and normal
	ix := ox + t*dx
	iy := oy + t*dy
	iz := oz + t*dz

	nx := ix - sphereX
	ny := iy - sphereY
	nz := iz - sphereZ

	// Optimized lighting
	const lightX, lightY, lightZ = -0.57735027, -0.57735027, -0.57735027
	dot := nx*lightX + ny*lightY + nz*lightZ
	intensity := maxFloat(0.0, dot)

	baseColor := 0.2 + 0.8*intensity
	return [3]float64{
		baseColor * 1.0,
		baseColor * 0.7,
		baseColor * 0.3,
	}
}

func normalizeOptimized(x, y, z float64) [3]float64 {
	lenSq := x*x + y*y + z*z
	if lenSq == 0 {
		return [3]float64{0, 0, 0}
	}

	invLen := 1.0 / fastSqrtOptimized(lenSq)
	return [3]float64{x * invLen, y * invLen, z * invLen}
}

func fastSqrtOptimized(x float64) float64 {
	if x <= 0 {
		return 0
	}
	if x == 1 {
		return 1
	}

	// Better initial guess
	var guess float64
	if x >= 1 {
		guess = x * 0.5
	} else {
		guess = (x + 1) * 0.5
	}

	// 4 iterations for optimal accuracy/performance balance
	for i := 0; i < 4; i++ {
		if guess == 0 {
			break
		}
		guess = 0.5 * (guess + x/guess)
	}

	return guess
}

func minInt(a, b int) int {
	if a < b {
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
