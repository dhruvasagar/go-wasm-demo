//go:build js && wasm

package main

import (
	"runtime"
	"sync"
	"syscall/js"
)

// Concurrent Mandelbrot implementation using goroutines
func mandelbrotWasmConcurrent(this js.Value, args []js.Value) interface{} {
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

	// Determine number of workers
	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers < 1 {
		numWorkers = 4 // Default for WASM
	}
	if numWorkers > height {
		numWorkers = height // Don't create more workers than rows
	}

	// Channel for work distribution
	rowChan := make(chan int, height)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go mandelbrotWorker(rowChan, &wg, result, width, height, dx, dy, xmin, ymin, maxIter)
	}

	// Send work to workers
	for y := 0; y < height; y++ {
		rowChan <- y
	}
	close(rowChan)

	// Wait for all workers to complete
	wg.Wait()

	// Create typed array and copy data efficiently
	jsArray := js.Global().Get("Int32Array").New(pixels)
	
	// Use efficient bulk assignment
	for i := 0; i < pixels; i++ {
		jsArray.SetIndex(i, js.ValueOf(result[i]))
	}
	
	return jsArray
}

// Worker function that processes rows of the Mandelbrot set
func mandelbrotWorker(rowChan <-chan int, wg *sync.WaitGroup, result []int32, width, height int, dx, dy, xmin, ymin float64, maxIter int) {
	defer wg.Done()

	for py := range rowChan {
		cy := ymin + float64(py)*dy
		rowOffset := py * width
		
		for px := 0; px < width; px++ {
			cx := xmin + float64(px)*dx
			
			// Optimized Mandelbrot calculation
			zx, zy := 0.0, 0.0
			iter := int32(0)
			
			// Hot loop - optimized for performance
			for iter < int32(maxIter) {
				zx2 := zx * zx
				zy2 := zy * zy
				if zx2+zy2 > 4.0 {
					break
				}
				zy = 2*zx*zy + cy
				zx = zx2 - zy2 + cx
				iter++
			}
			
			result[rowOffset+px] = iter
		}
	}
}

// Advanced concurrent version with work stealing for better load balancing
func mandelbrotWasmWorkStealing(this js.Value, args []js.Value) interface{} {
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

	result := make([]int32, pixels)

	// Work stealing approach: divide image into chunks
	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers < 1 {
		numWorkers = 4
	}

	// Create smaller chunks for better load balancing
	chunkSize := height / (numWorkers * 4) // More chunks than workers
	if chunkSize < 1 {
		chunkSize = 1
	}

	workChan := make(chan mandelbrotChunk, numWorkers*4)
	var wg sync.WaitGroup

	// Create work chunks
	for y := 0; y < height; y += chunkSize {
		endY := y + chunkSize
		if endY > height {
			endY = height
		}
		workChan <- mandelbrotChunk{startY: y, endY: endY}
	}
	close(workChan)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go mandelbrotChunkWorker(workChan, &wg, result, width, dx, dy, xmin, ymin, maxIter)
	}

	wg.Wait()

	// Return result
	jsArray := js.Global().Get("Int32Array").New(pixels)
	for i := 0; i < pixels; i++ {
		jsArray.SetIndex(i, js.ValueOf(result[i]))
	}
	
	return jsArray
}

// Chunk worker for work-stealing approach
func mandelbrotChunkWorker(workChan <-chan mandelbrotChunk, wg *sync.WaitGroup, result []int32, width int, dx, dy, xmin, ymin float64, maxIter int) {
	defer wg.Done()

	for chunk := range workChan {
		for py := chunk.startY; py < chunk.endY; py++ {
			cy := ymin + float64(py)*dy
			rowOffset := py * width
			
			for px := 0; px < width; px++ {
				cx := xmin + float64(px)*dx
				
				// Optimized Mandelbrot calculation with loop unrolling
				zx, zy := 0.0, 0.0
				iter := int32(0)
				
				// Unroll first iteration
				zx2 := zx * zx  // 0
				zy2 := zy * zy  // 0
				if zx2+zy2 <= 4.0 {
					zy = 2*zx*zy + cy  // cy
					zx = zx2 - zy2 + cx  // cx
					iter = 1
					
					// Continue with main loop
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
				
				result[rowOffset+px] = iter
			}
		}
	}
}