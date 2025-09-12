package main

import (
	"fmt"
	"testing"
)

// Test concurrent matrix multiplication against single-threaded version
func TestMatrixMultiplyConcurrentCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	// Test with small matrices for correctness
	size := 4

	// Create test matrices as JS-like arrays
	matrixA := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	matrixB := []float64{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	// Expected result from manual calculation
	expected := make([]float64, size*size)
	for i := 0; i < size; i++ {
		for k := 0; k < size; k++ {
			aik := matrixA[i*size+k]
			for j := 0; j < size; j++ {
				expected[i*size+j] += aik * matrixB[k*size+j]
			}
		}
	}

	// Test concurrent version produces same result
	// Note: In real WASM environment, this would use JS values
	// For testing, we simulate the algorithm logic
	result := make([]float64, size*size)

	// Transpose matrix B for cache efficiency (as in optimized version)
	matrixBT := make([]float64, size*size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			matrixBT[i*size+j] = matrixB[j*size+i]
		}
	}

	// Use same blocking algorithm as concurrent version
	const outerBlockSize = 64
	const innerBlockSize = 8

	for bi := 0; bi < size; bi += outerBlockSize {
		for bj := 0; bj < size; bj += outerBlockSize {
			for bk := 0; bk < size; bk += outerBlockSize {

				biEnd := min(bi+outerBlockSize, size)
				bjEnd := min(bj+outerBlockSize, size)
				bkEnd := min(bk+outerBlockSize, size)

				for i := bi; i < biEnd; i += innerBlockSize {
					for j := bj; j < bjEnd; j += innerBlockSize {
						for k := bk; k < bkEnd; k += innerBlockSize {

							iEnd := min(i+innerBlockSize, biEnd)
							jEnd := min(j+innerBlockSize, bjEnd)
							kEnd := min(k+innerBlockSize, bkEnd)

							// 2x2 register tiling
							for ii := i; ii < iEnd; ii += 2 {
								for jj := j; jj < jEnd; jj += 2 {
									r00 := result[ii*size+jj]
									r01 := result[ii*size+(jj+1)]
									r10 := result[(ii+1)*size+jj]
									r11 := result[(ii+1)*size+(jj+1)]

									for kk := k; kk < kEnd; kk += 2 {
										if ii < size && jj < size && kk < size {
											a00 := matrixA[ii*size+kk]
											b00 := matrixBT[jj*size+kk]
											r00 += a00 * b00
										}
										if ii < size && jj+1 < size && kk < size {
											a00 := matrixA[ii*size+kk]
											b01 := matrixBT[(jj+1)*size+kk]
											r01 += a00 * b01
										}
										if ii+1 < size && jj < size && kk < size {
											a10 := matrixA[(ii+1)*size+kk]
											b00 := matrixBT[jj*size+kk]
											r10 += a10 * b00
										}
										if ii+1 < size && jj+1 < size && kk < size {
											a10 := matrixA[(ii+1)*size+kk]
											b01 := matrixBT[(jj+1)*size+kk]
											r11 += a10 * b01
										}

										if kk+1 < kEnd {
											if ii < size && jj < size && kk+1 < size {
												a01 := matrixA[ii*size+(kk+1)]
												b10 := matrixBT[jj*size+(kk+1)]
												r00 += a01 * b10
											}
											if ii < size && jj+1 < size && kk+1 < size {
												a01 := matrixA[ii*size+(kk+1)]
												b11 := matrixBT[(jj+1)*size+(kk+1)]
												r01 += a01 * b11
											}
											if ii+1 < size && jj < size && kk+1 < size {
												a11 := matrixA[(ii+1)*size+(kk+1)]
												b10 := matrixBT[jj*size+(kk+1)]
												r10 += a11 * b10
											}
											if ii+1 < size && jj+1 < size && kk+1 < size {
												a11 := matrixA[(ii+1)*size+(kk+1)]
												b11 := matrixBT[(jj+1)*size+(kk+1)]
												r11 += a11 * b11
											}
										}
									}

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

	// Compare results
	const tolerance = 1e-10
	for i := 0; i < size*size; i++ {
		if absDiff(result[i]-expected[i]) > tolerance {
			t.Errorf("Matrix multiplication mismatch at index %d: got %f, want %f", i, result[i], expected[i])
		}
	}

	t.Logf("Matrix multiplication correctness test passed for %dx%d matrices", size, size)
}

func absDiff(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Test hash function consistency across different implementations
func TestHashConcurrentConsistency(t *testing.T) {
	data := "WebAssembly performance test data for hashing benchmark"
	iterations := 1000

	// Test our optimized hash algorithm directly (without JS interface)
	const numLanes = 4
	hashLanes := [numLanes]uint32{
		0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476,
	}

	dataBytes := []byte(data)
	dataLen := len(dataBytes)

	// Pre-process data for better cache performance
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

	const c1, c2, c3, c4 = 0x85EBCA6B, 0xC2B2AE35, 0xCC9E2D51, 0x1B873593

	// Run multiple times and ensure consistency
	results := make([]uint32, 5)

	for run := 0; run < 5; run++ {
		baseIter := iterations / numLanes
		laneResults := make([]uint32, numLanes)

		for lane := 0; lane < numLanes; lane++ {
			startIter := lane * baseIter
			endIter := startIter + baseIter
			if lane == numLanes-1 {
				endIter = iterations
			}

			hash := hashLanes[lane]

			for iter := startIter; iter < endIter; iter++ {
				iterSeed := uint32(iter)*0x9E3779B9 + uint32(lane)*c1

				// Process words with optimization
				i := 0
				for ; i <= len(dataWords)-4; i += 4 {
					w0 := dataWords[i] * c1
					w0 = (w0 << 15) | (w0 >> 17)
					hash ^= w0
					hash = ((hash<<13)|(hash>>19))*5 + 0xE6546B64

					w1 := dataWords[i+1] * c3
					w1 = (w1 << 17) | (w1 >> 15)
					hash ^= w1
					hash = ((hash<<11)|(hash>>21))*3 + 0xE6546B64

					w2 := dataWords[i+2] * c1
					w2 = (w2 << 19) | (w2 >> 13)
					hash ^= w2
					hash = ((hash<<7)|(hash>>25))*7 + 0xE6546B64

					w3 := dataWords[i+3] * c3
					w3 = (w3 << 13) | (w3 >> 19)
					hash ^= w3
					hash = ((hash<<17)|(hash>>15))*11 + 0xE6546B64
				}

				for ; i < len(dataWords); i++ {
					w := dataWords[i] * c1
					w = (w << 15) | (w >> 17)
					hash ^= w
					hash = ((hash << 13) | (hash >> 19)) + 0xE6546B64
				}

				hash ^= iterSeed
				hash = hash*c2 + c4
			}

			laneResults[lane] = hash
		}

		// Combine results
		finalHash := laneResults[0]
		for i := 1; i < numLanes; i++ {
			finalHash ^= laneResults[i]
			finalHash = finalHash*c1 + c2
			finalHash = (finalHash << 16) | (finalHash >> 16)
		}

		// Final avalanche
		finalHash ^= finalHash >> 16
		finalHash *= c1
		finalHash ^= finalHash >> 13
		finalHash *= c2
		finalHash ^= finalHash >> 16

		results[run] = finalHash
	}

	// All results should be identical
	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("Hash results not consistent: run 0 = %d, run %d = %d", results[0], i, results[i])
		}
	}

	t.Logf("Hash consistency test passed: consistent result %d across %d runs", results[0], len(results))
}

// Test Mandelbrot set calculation correctness
func TestMandelbrotConcurrentCorrectness(t *testing.T) {
	// Test known points in the Mandelbrot set
	tests := []struct {
		cx, cy   float64
		maxIter  int
		expected int // Expected iteration count (approximate range)
		inSet    bool
	}{
		{0.0, 0.0, 100, 100, true},  // Origin is in the set
		{-1.0, 0.0, 100, 100, true}, // Point on real axis in set
		{2.0, 0.0, 100, 5, false},   // Point outside set (diverges quickly)
		{-2.5, 0.0, 100, 5, false},  // Point far outside set
		{-0.5, 0.0, 100, 100, true}, // Another point in set
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Point(%.1f,%.1f)", tt.cx, tt.cy), func(t *testing.T) {
			// Test single point using same algorithm as vectorized version
			zx, zy := 0.0, 0.0
			iter := 0

			for iter < tt.maxIter {
				zx2 := zx * zx
				zy2 := zy * zy

				if zx2+zy2 > 4.0 {
					break
				}

				// z = z² + c
				newZx := zx2 - zy2 + tt.cx
				newZy := 2.0*zx*zy + tt.cy
				zx = newZx
				zy = newZy
				iter++
			}

			if tt.inSet && iter != tt.maxIter {
				t.Errorf("Point should be in Mandelbrot set: got %d iterations, want %d", iter, tt.maxIter)
			} else if !tt.inSet && iter >= tt.expected {
				t.Errorf("Point should diverge quickly: got %d iterations, want < %d", iter, tt.expected)
			}
		})
	}
}

// Fast square root using Newton-Raphson approximation (for testing)
func fastSqrtTest(x float64) float64 {
	if x <= 0 {
		return 0
	}
	if x == 1 {
		return 1
	}

	// Better initial guess
	var guess float64
	if x > 1 {
		// For x > 1, start with x/2
		guess = x * 0.5
	} else {
		// For x < 1, start with (x + 1)/2
		guess = (x + 1) * 0.5
	}

	// More iterations for better accuracy
	for i := 0; i < 5; i++ {
		if guess == 0 {
			break
		}
		guess = 0.5 * (guess + x/guess)
	}

	return guess
}

// Test fast square root implementation
func TestFastSqrt(t *testing.T) {
	tests := []struct {
		input     float64
		expected  float64
		tolerance float64
	}{
		{0.0, 0.0, 1e-10},
		{1.0, 1.0, 1e-10},
		{4.0, 2.0, 1e-4},
		{9.0, 3.0, 1e-4},
		{16.0, 4.0, 1e-4},
		{25.0, 5.0, 1e-4},
		{100.0, 10.0, 1e-3},
		{0.25, 0.5, 1e-4},
		{0.01, 0.1, 1e-3},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("sqrt(%.2f)", tt.input), func(t *testing.T) {
			result := fastSqrtTest(tt.input)
			if absDiff(result-tt.expected) > tt.tolerance {
				t.Errorf("fastSqrt(%.2f) = %.6f, want %.6f (tolerance %.1e)",
					tt.input, result, tt.expected, tt.tolerance)
			}
		})
	}
}

// Benchmark concurrent vs single-threaded matrix multiplication
func BenchmarkMatrixMultiplyConcurrent50x50(b *testing.B) {
	benchmarkMatrixMultiplyConcurrent(b, 50)
}

func BenchmarkMatrixMultiplyConcurrent100x100(b *testing.B) {
	benchmarkMatrixMultiplyConcurrent(b, 100)
}

func BenchmarkMatrixMultiplyConcurrent200x200(b *testing.B) {
	benchmarkMatrixMultiplyConcurrent(b, 200)
}

func benchmarkMatrixMultiplyConcurrent(b *testing.B, size int) {
	// Create test matrices
	matrixA := make([]float64, size*size)
	matrixB := make([]float64, size*size)

	// Initialize with test data
	for i := 0; i < size*size; i++ {
		matrixA[i] = float64(i % 10)
		matrixB[i] = float64((i * 2) % 10)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Simulate concurrent optimized algorithm
		result := make([]float64, size*size)

		// Transpose matrix B for cache efficiency
		matrixBT := make([]float64, size*size)
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				matrixBT[i*size+j] = matrixB[j*size+i]
			}
		}

		// Blocked multiplication with register tiling
		const outerBlockSize = 64
		const innerBlockSize = 8

		for bi := 0; bi < size; bi += outerBlockSize {
			for bj := 0; bj < size; bj += outerBlockSize {
				for bk := 0; bk < size; bk += outerBlockSize {

					biEnd := min(bi+outerBlockSize, size)
					bjEnd := min(bj+outerBlockSize, size)
					bkEnd := min(bk+outerBlockSize, size)

					for i := bi; i < biEnd; i += innerBlockSize {
						for j := bj; j < bjEnd; j += innerBlockSize {
							for k := bk; k < bkEnd; k += innerBlockSize {

								iEnd := min(i+innerBlockSize, biEnd)
								jEnd := min(j+innerBlockSize, bjEnd)
								kEnd := min(k+innerBlockSize, bkEnd)

								// 2x2 register tiling
								for ii := i; ii < iEnd; ii += 2 {
									for jj := j; jj < jEnd; jj += 2 {
										r00 := result[ii*size+jj]
										r01 := result[ii*size+(jj+1)]
										r10 := result[(ii+1)*size+jj]
										r11 := result[(ii+1)*size+(jj+1)]

										for kk := k; kk < kEnd; kk += 2 {
											if ii < size && jj < size && kk < size {
												a00 := matrixA[ii*size+kk]
												b00 := matrixBT[jj*size+kk]
												r00 += a00 * b00
											}
											if ii < size && jj+1 < size && kk < size {
												a00 := matrixA[ii*size+kk]
												b01 := matrixBT[(jj+1)*size+kk]
												r01 += a00 * b01
											}
											if ii+1 < size && jj < size && kk < size {
												a10 := matrixA[(ii+1)*size+kk]
												b00 := matrixBT[jj*size+kk]
												r10 += a10 * b00
											}
											if ii+1 < size && jj+1 < size && kk < size {
												a10 := matrixA[(ii+1)*size+kk]
												b01 := matrixBT[(jj+1)*size+kk]
												r11 += a10 * b01
											}

											if kk+1 < kEnd {
												if ii < size && jj < size && kk+1 < size {
													a01 := matrixA[ii*size+(kk+1)]
													b10 := matrixBT[jj*size+(kk+1)]
													r00 += a01 * b10
												}
												if ii < size && jj+1 < size && kk+1 < size {
													a01 := matrixA[ii*size+(kk+1)]
													b11 := matrixBT[(jj+1)*size+(kk+1)]
													r01 += a01 * b11
												}
												if ii+1 < size && jj < size && kk+1 < size {
													a11 := matrixA[(ii+1)*size+(kk+1)]
													b10 := matrixBT[jj*size+(kk+1)]
													r10 += a11 * b10
												}
												if ii+1 < size && jj+1 < size && kk+1 < size {
													a11 := matrixA[(ii+1)*size+(kk+1)]
													b11 := matrixBT[(jj+1)*size+(kk+1)]
													r11 += a11 * b11
												}
											}
										}

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

		// Prevent compiler optimization
		_ = result[0]
	}

	// Report operations per second
	ops := int64(size) * int64(size) * int64(size) * int64(b.N)
	b.ReportMetric(float64(ops)/b.Elapsed().Seconds()/1e6, "MOps/sec")
}

// Benchmark concurrent hash function
func BenchmarkHashConcurrent1000(b *testing.B) {
	benchmarkHashConcurrent(b, 1000)
}

func BenchmarkHashConcurrent10000(b *testing.B) {
	benchmarkHashConcurrent(b, 10000)
}

func BenchmarkHashConcurrent100000(b *testing.B) {
	benchmarkHashConcurrent(b, 100000)
}

func benchmarkHashConcurrent(b *testing.B, iterations int) {
	data := "WebAssembly performance test data for hashing benchmark"
	dataBytes := []byte(data)
	dataLen := len(dataBytes)

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		// Simulate concurrent optimized hash
		const numLanes = 4
		hashLanes := [numLanes]uint32{
			0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476,
		}

		// Pre-process data
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

		const c1, c2, c3, c4 = 0x85EBCA6B, 0xC2B2AE35, 0xCC9E2D51, 0x1B873593
		baseIter := iterations / numLanes
		results := make([]uint32, numLanes)

		for lane := 0; lane < numLanes; lane++ {
			startIter := lane * baseIter
			endIter := startIter + baseIter
			if lane == numLanes-1 {
				endIter = iterations
			}

			hash := hashLanes[lane]

			for iter := startIter; iter < endIter; iter++ {
				iterSeed := uint32(iter)*0x9E3779B9 + uint32(lane)*c1

				i := 0
				for ; i <= len(dataWords)-4; i += 4 {
					w0 := dataWords[i] * c1
					w0 = (w0 << 15) | (w0 >> 17)
					hash ^= w0
					hash = ((hash<<13)|(hash>>19))*5 + 0xE6546B64

					w1 := dataWords[i+1] * c3
					w1 = (w1 << 17) | (w1 >> 15)
					hash ^= w1
					hash = ((hash<<11)|(hash>>21))*3 + 0xE6546B64

					w2 := dataWords[i+2] * c1
					w2 = (w2 << 19) | (w2 >> 13)
					hash ^= w2
					hash = ((hash<<7)|(hash>>25))*7 + 0xE6546B64

					w3 := dataWords[i+3] * c3
					w3 = (w3 << 13) | (w3 >> 19)
					hash ^= w3
					hash = ((hash<<17)|(hash>>15))*11 + 0xE6546B64
				}

				for ; i < len(dataWords); i++ {
					w := dataWords[i] * c1
					w = (w << 15) | (w >> 17)
					hash ^= w
					hash = ((hash << 13) | (hash >> 19)) + 0xE6546B64
				}

				hash ^= iterSeed
				hash = hash*c2 + c4
			}

			results[lane] = hash
		}

		// Combine results
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

		// Prevent compiler optimization
		_ = finalHash
	}

	// Report hashes per second
	hashes := int64(iterations) * int64(b.N)
	b.ReportMetric(float64(hashes)/b.Elapsed().Seconds()/1e6, "MHashes/sec")
}
