//go:build js && wasm

package main

import (
	"math"
	"runtime"
	"sync"
	"syscall/js"
	"unsafe"
)

// ============================================================================
// SINGLE-THREADED IMPLEMENTATIONS (for comparison)
// ============================================================================

// Single-threaded matrix multiplication
func matrixMultiplyWasmSingle(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return js.ValueOf("Missing arguments")
	}

	matrixA := args[0]
	matrixB := args[1]
	size := args[2].Int()

	// Copy JS arrays to Go slices once
	goMatrixA := make([]float64, size*size)
	goMatrixB := make([]float64, size*size)

	for i := 0; i < size*size; i++ {
		goMatrixA[i] = matrixA.Index(i).Float()
		goMatrixB[i] = matrixB.Index(i).Float()
	}

	result := make([]float64, size*size)

	// Single-threaded computation
	for i := 0; i < size; i++ {
		for k := 0; k < size; k++ {
			aik := goMatrixA[i*size+k]
			for j := 0; j < size; j++ {
				result[i*size+j] += aik * goMatrixB[k*size+j]
			}
		}
	}

	// Convert result back to JavaScript
	jsArray := js.Global().Get("Array").New(size * size)
	for i := 0; i < size*size; i++ {
		jsArray.SetIndex(i, result[i])
	}

	return jsArray
}

// Single-threaded Mandelbrot
func mandelbrotWasmSingle(this js.Value, args []js.Value) interface{} {
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

	dx := (xmax - xmin) / float64(width)
	dy := (ymax - ymin) / float64(height)
	pixels := width * height
	result := make([]int32, pixels)

	idx := 0
	for py := 0; py < height; py++ {
		cy := ymin + float64(py)*dy

		for px := 0; px < width; px++ {
			cx := xmin + float64(px)*dx

			zx, zy := 0.0, 0.0
			iter := int32(0)

			for iter < int32(maxIter) {
				zx2 := zx * zx
				zy2 := zy * zy
				if zx2+zy2 > 4.0 {
					break
				}
				temp := zx2 - zy2 + cx
				zy = 2*zx*zy + cy
				zx = temp
				iter++
			}

			result[idx] = iter
			idx++
		}
	}

	jsArray := js.Global().Get("Int32Array").New(pixels)
	for i := 0; i < pixels; i++ {
		jsArray.SetIndex(i, js.ValueOf(result[i]))
	}

	return jsArray
}

// Single-threaded hash computation
func sha256HashWasmSingle(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return js.ValueOf(0)
	}

	data := args[0].String()
	iterations := args[1].Int()

	dataBytes := []byte(data)
	dataLen := len(dataBytes)
	hash := uint32(0x12345678)

	for iter := 0; iter < iterations; iter++ {
		for i := 0; i < dataLen; i++ {
			hash = hash*33 + uint32(dataBytes[i])
			hash = (hash << 5) | (hash >> 27)
		}
	}

	return js.ValueOf(int(hash))
}

// Single-threaded ray tracing - OPTIMIZED with inlined calculations
func rayTracingWasmSingle(this js.Value, args []js.Value) interface{} {
	width := args[0].Int()
	height := args[1].Int()
	samples := args[2].Int()

	result := make([]float64, width*height*3)

	// Sphere properties (same as JavaScript)
	const sphereX, sphereY, sphereZ = 0.0, 0.0, -5.0
	const sphereRadius2 = 1.0

	// Light direction (same as JavaScript)
	const lightX, lightY, lightZ = -0.57735027, -0.57735027, -0.57735027

	for y := 0; y < height; y++ {
		ny := (float64(y)/float64(height))*2.0 - 1.0

		for x := 0; x < width; x++ {
			nx := (float64(x)/float64(width))*2.0 - 1.0

			var colorR, colorG, colorB float64

			for s := 0; s < samples; s++ {
				// FULLY INLINED: Ray direction normalization (no function calls)
				rayLenSq := nx*nx + ny*ny + 1.0

				// Inlined fast square root (Newton-Raphson, 2 iterations)
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
						// FULLY INLINED: Calculate intersection point, normal, and lighting
						ix := 0.0 + t*dirX
						iy := 0.0 + t*dirY
						iz := 0.0 + t*dirZ

						normalX := ix - sphereX
						normalY := iy - sphereY
						normalZ := iz - sphereZ

						// Inlined max(0, dot) - no function call
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

// Fast square root approximation (avoids math.Sqrt overhead)
func fastSqrt(x float64) float64 {
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

	// 3 iterations is usually enough for good precision
	for i := 0; i < 3; i++ {
		if guess == 0 {
			break
		}
		guess = 0.5 * (guess + x/guess)
	}

	return guess
}

// ============================================================================
// ENHANCED CONCURRENT IMPLEMENTATIONS
// ============================================================================

// Enhanced concurrent matrix multiplication with adaptive chunking
func matrixMultiplyWasmConcurrentV2(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return js.ValueOf("Missing arguments")
	}

	matrixA := args[0]
	matrixB := args[1]
	size := args[2].Int()

	// Copy matrices once
	goMatrixA := make([]float64, size*size)
	goMatrixB := make([]float64, size*size)

	for i := 0; i < size*size; i++ {
		goMatrixA[i] = matrixA.Index(i).Float()
		goMatrixB[i] = matrixB.Index(i).Float()
	}

	result := make([]float64, size*size)

	// Adaptive worker count based on problem size
	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers < 1 {
		numWorkers = 4
	}

	// Adjust for small matrices
	if size < 100 {
		numWorkers = 2
	} else if size > 1000 {
		numWorkers = min(numWorkers, 8) // Cap workers for very large matrices
	}

	// Dynamic chunk size based on matrix size and worker count
	chunkSize := size / (numWorkers * 2)
	if chunkSize < 1 {
		chunkSize = 1
	}
	if chunkSize > 64 {
		chunkSize = 64 // Cache-friendly chunk size
	}

	workChan := make(chan matrixWorkChunk, numWorkers*2)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go matrixMultiplyChunkWorker(workChan, &wg, goMatrixA, goMatrixB, result, size)
	}

	// Generate work chunks
	go func() {
		defer close(workChan)
		for row := 0; row < size; row += chunkSize {
			endRow := min(row+chunkSize, size)
			workChan <- matrixWorkChunk{startRow: row, endRow: endRow}
		}
	}()

	wg.Wait()

	// Convert result
	jsArray := js.Global().Get("Array").New(size * size)
	for i := 0; i < size*size; i++ {
		jsArray.SetIndex(i, result[i])
	}

	return jsArray
}

func matrixMultiplyChunkWorker(workChan <-chan matrixWorkChunk, wg *sync.WaitGroup, matrixA, matrixB, result []float64, size int) {
	defer wg.Done()

	for chunk := range workChan {
		// Process chunk of rows
		for i := chunk.startRow; i < chunk.endRow; i++ {
			rowOffset := i * size

			// Cache-optimized computation
			for k := 0; k < size; k++ {
				aik := matrixA[rowOffset+k]
				bRowOffset := k * size

				// Vectorizable inner loop
				for j := 0; j < size; j++ {
					result[rowOffset+j] += aik * matrixB[bRowOffset+j]
				}
			}
		}
	}
}

// Enhanced concurrent Mandelbrot with adaptive load balancing
func mandelbrotWasmConcurrentV2(this js.Value, args []js.Value) interface{} {
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

	dx := (xmax - xmin) / float64(width)
	dy := (ymax - ymin) / float64(height)
	pixels := width * height
	result := make([]int32, pixels)

	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers < 1 {
		numWorkers = 4
	}

	// Adaptive chunk size for better load balancing
	totalChunks := numWorkers * 8 // More chunks than workers for better load distribution
	chunkHeight := height / totalChunks
	if chunkHeight < 1 {
		chunkHeight = 1
		totalChunks = height
	}

	workChan := make(chan mandelbrotChunk, totalChunks)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go mandelbrotChunkWorkerV2(workChan, &wg, result, width, dx, dy, xmin, ymin, maxIter)
	}

	// Generate work chunks
	go func() {
		defer close(workChan)
		for y := 0; y < height; y += chunkHeight {
			endY := min(y+chunkHeight, height)
			workChan <- mandelbrotChunk{startY: y, endY: endY}
		}
	}()

	wg.Wait()

	jsArray := js.Global().Get("Int32Array").New(pixels)
	for i := 0; i < pixels; i++ {
		jsArray.SetIndex(i, js.ValueOf(result[i]))
	}

	return jsArray
}

func mandelbrotChunkWorkerV2(workChan <-chan mandelbrotChunk, wg *sync.WaitGroup, result []int32, width int, dx, dy, xmin, ymin float64, maxIter int) {
	defer wg.Done()

	for chunk := range workChan {
		for py := chunk.startY; py < chunk.endY; py++ {
			cy := ymin + float64(py)*dy
			rowOffset := py * width

			for px := 0; px < width; px++ {
				cx := xmin + float64(px)*dx

				// Optimized Mandelbrot with early escape
				zx, zy := 0.0, 0.0
				iter := int32(0)

				// Unrolled first few iterations for better performance
				for iter < int32(maxIter) {
					zx2 := zx * zx
					zy2 := zy * zy
					if zx2+zy2 > 4.0 {
						break
					}

					// Compute next iteration
					temp := zx2 - zy2 + cx
					zy = 2*zx*zy + cy
					zx = temp
					iter++
				}

				result[rowOffset+px] = iter
			}
		}
	}
}

// Enhanced concurrent hash computation with better work distribution
func sha256HashWasmConcurrentV2(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return js.ValueOf(0)
	}

	data := args[0].String()
	iterations := args[1].Int()

	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers < 1 {
		numWorkers = 4
	}

	// Better work distribution - avoid integer division issues
	baseIterations := iterations / numWorkers
	remainder := iterations % numWorkers

	resultChan := make(chan uint32, numWorkers)
	var wg sync.WaitGroup

	// Start workers with balanced workload
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		workerIterations := baseIterations
		if i < remainder {
			workerIterations++ // Distribute remainder evenly
		}

		go hashWorkerV2(data, workerIterations, i, resultChan, &wg)
	}

	wg.Wait()
	close(resultChan)

	// Combine results with better mixing
	finalHash := uint32(0x9E3779B9) // Better initial value
	for workerResult := range resultChan {
		finalHash ^= workerResult
		finalHash = finalHash*0x85EBCA6B + 0xC2B2AE35     // Better mixing constants
		finalHash = (finalHash << 13) | (finalHash >> 19) // Rotate for better distribution
	}

	return js.ValueOf(int(finalHash))
}

func hashWorkerV2(data string, iterations, workerId int, resultChan chan<- uint32, wg *sync.WaitGroup) {
	defer wg.Done()

	dataBytes := []byte(data)
	dataLen := len(dataBytes)
	hash := uint32(0x12345678) + uint32(workerId)*0x1000 // Different starting point per worker

	for iter := 0; iter < iterations; iter++ {
		// Process 8 bytes at a time when possible (loop unrolling)
		i := 0
		for ; i <= dataLen-8; i += 8 {
			for j := 0; j < 8; j++ {
				hash = hash*33 + uint32(dataBytes[i+j])
				hash = (hash << 5) | (hash >> 27)
			}
		}

		// Process remaining bytes
		for ; i < dataLen; i++ {
			hash = hash*33 + uint32(dataBytes[i])
			hash = (hash << 5) | (hash >> 27)
		}

		// Add iteration mixing to improve distribution
		hash ^= uint32(iter)
	}

	resultChan <- hash
}

// Enhanced concurrent ray tracing with tile-based rendering
func rayTracingWasmConcurrentV2(this js.Value, args []js.Value) interface{} {
	width := args[0].Int()
	height := args[1].Int()
	samples := args[2].Int()

	result := make([]float64, width*height*3)

	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers < 1 {
		numWorkers = 4
	}

	// Tile-based rendering for better cache performance
	tileSize := 32 // 32x32 tiles
	if tileSize > width {
		tileSize = width
	}
	if tileSize > height {
		tileSize = height
	}

	tileChan := make(chan tile, numWorkers*4)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go rayTracingTileWorker(tileChan, &wg, result, width, height, samples)
	}

	// Generate tiles
	go func() {
		defer close(tileChan)
		for y := 0; y < height; y += tileSize {
			endY := min(y+tileSize, height)
			for x := 0; x < width; x += tileSize {
				endX := min(x+tileSize, width)
				tileChan <- tile{startX: x, endX: endX, startY: y, endY: endY}
			}
		}
	}()

	wg.Wait()

	// FIXED: Use efficient bulk copy instead of O(n) boundary calls
	resultTyped := js.Global().Get("Float64Array").New(len(result))
	arrayBuffer := resultTyped.Get("buffer")
	uint8View := js.Global().Get("Uint8Array").New(arrayBuffer)

	// Copy bytes to Uint8Array view of the Float64Array buffer
	js.CopyBytesToJS(
		uint8View,
		unsafe.Slice((*byte)(unsafe.Pointer(&result[0])), len(result)*8),
	)

	return resultTyped
}

func rayTracingTileWorker(tileChan chan tile, wg *sync.WaitGroup, result []float64, width, height, samples int) {
	defer wg.Done()

	// Sphere properties (same as JavaScript and single-threaded)
	const sphereX, sphereY, sphereZ = 0.0, 0.0, -5.0
	const sphereRadius2 = 1.0
	const lightX, lightY, lightZ = -0.57735027, -0.57735027, -0.57735027

	for t := range tileChan {
		for y := t.startY; y < t.endY; y++ {
			ny := (float64(y)/float64(height))*2.0 - 1.0

			for x := t.startX; x < t.endX; x++ {
				nx := (float64(x)/float64(width))*2.0 - 1.0

				var colorR, colorG, colorB float64

				// FULLY INLINED: Zero function calls in hot path
				for s := 0; s < samples; s++ {
					// Fully inlined ray direction normalization
					rayLenSq := nx*nx + ny*ny + 1.0

					// Inlined square root (Newton-Raphson, 2 iterations)
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

					// Fully inlined ray-sphere intersection
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
							// Fully inlined intersection point, normal, and lighting
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
	}
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Helper functions for ray tracing
func normalize(x, y, z float64) []float64 {
	length := math.Sqrt(x*x + y*y + z*z)
	if length == 0 {
		return []float64{0, 0, 0}
	}
	invLength := 1.0 / length
	return []float64{x * invLength, y * invLength, z * invLength}
}

// Simple ray tracing for a sphere
func traceRay(originX, originY, originZ, dirX, dirY, dirZ float64) []float64 {
	// Sphere at center
	sphereX, sphereY, sphereZ := 0.0, 0.0, -3.0
	radius := 1.0

	// Ray-sphere intersection
	ocX := originX - sphereX
	ocY := originY - sphereY
	ocZ := originZ - sphereZ

	a := dirX*dirX + dirY*dirY + dirZ*dirZ
	b := 2.0 * (ocX*dirX + ocY*dirY + ocZ*dirZ)
	c := ocX*ocX + ocY*ocY + ocZ*ocZ - radius*radius

	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		// No intersection - return sky color
		return []float64{0.2, 0.3, 0.8}
	}

	// Calculate intersection point
	t := (-b - math.Sqrt(discriminant)) / (2.0 * a)

	// Calculate hit point and normal
	hitX := originX + t*dirX
	hitY := originY + t*dirY
	hitZ := originZ + t*dirZ

	normalX := (hitX - sphereX) / radius
	normalY := (hitY - sphereY) / radius
	normalZ := (hitZ - sphereZ) / radius

	// Simple lighting calculation
	lightDirX, lightDirY, lightDirZ := -0.57735, -0.57735, -0.57735
	lightIntensity := math.Max(0, normalX*lightDirX+normalY*lightDirY+normalZ*lightDirZ)

	// Basic shading
	intensity := 0.1 + 0.9*lightIntensity

	return []float64{intensity, intensity * 0.8, intensity * 0.6}
}
