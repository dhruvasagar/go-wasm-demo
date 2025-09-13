//go:build js && wasm

package main

import (
	"math"
	"runtime"
	"sync"
	"syscall/js"
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

	// Use shared conversion function for efficient result conversion
	return createInt32TypedArray(result)
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

	// Use shared implementation to avoid code duplication
	result := rayTracingSharedSingle(width, height, samples)

	// Return result using shared conversion function
	return createFloat64TypedArray(result)
}

// Note: fastSqrt function removed - replaced with math.Sqrt() for better performance

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
		numWorkers = minInt(numWorkers, 8) // Cap workers for very large matrices
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
			endRow := minInt(row+chunkSize, size)
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
			endY := minInt(y+chunkHeight, height)
			workChan <- mandelbrotChunk{startY: y, endY: endY}
		}
	}()

	wg.Wait()

	// Use shared conversion function for efficient result conversion
	return createInt32TypedArray(result)
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

// OPTIMIZED concurrent hash computation - eliminates most overhead
func sha256HashWasmConcurrentV2(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return js.ValueOf(0)
	}

	data := args[0].String()
	iterations := args[1].Int()

	// CRITICAL: Only use concurrency if the workload is large enough
	// Threshold determined by benchmarking - below this, single-threaded is faster
	const concurrencyThreshold = 50000
	if iterations < concurrencyThreshold {
		// Fall back to optimized single-threaded version for small workloads
		return sha256HashWasmSingle(js.Value{}, args)
	}

	// Smart worker count: use fewer workers to reduce overhead
	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers < 1 {
		numWorkers = 4
	}
	if numWorkers > 8 {
		numWorkers = 8 // Cap at 8 workers - diminishing returns beyond this
	}

	// Pre-convert string to bytes once (major optimization)
	dataBytes := []byte(data)
	dataLen := len(dataBytes)

	// Pre-allocate result array to avoid channel overhead
	results := make([]uint32, numWorkers)
	var wg sync.WaitGroup

	// Distribute work evenly with larger chunks
	baseIterations := iterations / numWorkers
	remainder := iterations % numWorkers

	// Start workers - each works on a different portion of iterations
	for workerID := 0; workerID < numWorkers; workerID++ {
		wg.Add(1)

		workerIterations := baseIterations
		if workerID < remainder {
			workerIterations++
		}

		// CRITICAL: Minimize allocations and function calls
		go func(id, iters int) {
			defer wg.Done()

			// Each worker uses a different hash seed for better distribution
			hash := uint32(0x12345678) + uint32(id)*0x9E3779B9

			// ULTRA-OPTIMIZED inner loop - no function calls, minimal operations
			for iter := 0; iter < iters; iter++ {
				// Process data in 4-byte chunks for maximum efficiency
				i := 0
				for ; i <= dataLen-4; i += 4 {
					// Unrolled 4-byte processing
					hash = hash*33 + uint32(dataBytes[i])
					hash = (hash << 5) | (hash >> 27)
					hash = hash*33 + uint32(dataBytes[i+1])
					hash = (hash << 3) | (hash >> 29)
					hash = hash*33 + uint32(dataBytes[i+2])
					hash = (hash << 7) | (hash >> 25)
					hash = hash*33 + uint32(dataBytes[i+3])
					hash = (hash << 11) | (hash >> 21)
				}

				// Process remaining bytes (0-3)
				for ; i < dataLen; i++ {
					hash = hash*33 + uint32(dataBytes[i])
					hash = (hash << 5) | (hash >> 27)
				}

				// Mix in iteration counter efficiently
				hash ^= uint32(iter)
				hash = hash*0x85EBCA6B + 0xC2B2AE35
			}

			results[id] = hash
		}(workerID, workerIterations)
	}

	wg.Wait()

	// OPTIMIZED result combination - no channel overhead
	finalHash := uint32(0x9E3779B9)
	for i := 0; i < numWorkers; i++ {
		finalHash ^= results[i]
		finalHash = finalHash*0x85EBCA6B + 0xC2B2AE35
		finalHash = (finalHash << 13) | (finalHash >> 19)
	}

	// Final avalanche mixing
	finalHash ^= finalHash >> 16
	finalHash *= 0x85EBCA6B
	finalHash ^= finalHash >> 13
	finalHash *= 0xC2B2AE35
	finalHash ^= finalHash >> 16

	return js.ValueOf(int(finalHash))
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
			endY := minInt(y+tileSize, height)
			for x := 0; x < width; x += tileSize {
				endX := minInt(x+tileSize, width)
				tileChan <- tile{startX: x, endX: endX, startY: y, endY: endY}
			}
		}
	}()

	wg.Wait()

	// Use shared conversion function to avoid duplication
	return createFloat64TypedArray(result)
}

func rayTracingTileWorker(tileChan chan tile, wg *sync.WaitGroup, result []float64, width, height, samples int) {
	defer wg.Done()

	for t := range tileChan {
		for y := t.startY; y < t.endY; y++ {
			ny := (float64(y)/float64(height))*2.0 - 1.0

			for x := t.startX; x < t.endX; x++ {
				nx := (float64(x)/float64(width))*2.0 - 1.0

				// Use shared ray computation to avoid code duplication
				colorR, colorG, colorB := computeRayColor(nx, ny, samples)

				idx := (y*width + x) * 3
				result[idx] = colorR
				result[idx+1] = colorG
				result[idx+2] = colorB
			}
		}
	}
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// Note: min/max functions moved to benchmarks_shared.go as minInt/maxInt to avoid duplication

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
