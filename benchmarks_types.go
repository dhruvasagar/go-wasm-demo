//go:build js && wasm

package main

// Common types for benchmarks

// Matrix work chunk for concurrent processing
type matrixWorkChunk struct {
	startRow, endRow int
}

// Mandelbrot work chunk
type mandelbrotChunk struct {
	startY, endY int
}

// Tile for ray tracing
type tile struct {
	startX, endX, startY, endY int
}

// Pixel work for ray tracing
type pixelWork struct {
	x, y int
}

// Vector work for vectorized processing
type vectorWork struct {
	startIdx, endIdx int
}

// Matrix block for block-based multiplication
type matrixBlock struct {
	startI, endI, startK, endK int
}

// Hash result with worker ID
type hashResult struct {
	worker int
	hash   uint32
}

// Tile work for enhanced ray tracing
type tileWork struct {
	tileX, tileY int
}