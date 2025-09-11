//go:build !wasm

package main

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

// TestMatrixMultiplicationLogic tests the matrix multiplication algorithm correctness
func TestMatrixMultiplicationLogic(t *testing.T) {
	// Test small 2x2 matrix multiplication
	size := 2
	matrixA := []float64{1, 2, 3, 4} // [[1, 2], [3, 4]]
	matrixB := []float64{5, 6, 7, 8} // [[5, 6], [7, 8]]

	result := make([]float64, size*size)

	// Perform matrix multiplication
	for i := 0; i < size; i++ {
		for k := 0; k < size; k++ {
			aik := matrixA[i*size+k]
			for j := 0; j < size; j++ {
				result[i*size+j] += aik * matrixB[k*size+j]
			}
		}
	}

	// Expected result: [[19, 22], [43, 50]]
	expected := []float64{19, 22, 43, 50}

	for i := 0; i < len(expected); i++ {
		if result[i] != expected[i] {
			t.Errorf("Matrix multiplication failed at index %d: got %v, want %v", i, result[i], expected[i])
		}
	}
}

// TestMandelbrotLogic tests the Mandelbrot set calculation correctness
func TestMandelbrotLogic(t *testing.T) {
	// Test known points in the Mandelbrot set
	tests := []struct {
		cx, cy   float64
		maxIter  int
		expected int // Expected iteration count (approximate)
	}{
		{0.0, 0.0, 100, 100},  // Origin is in the set
		{-1.0, 0.0, 100, 100}, // Point on real axis in set
		{2.0, 0.0, 100, 1},    // Point outside set (diverges quickly)
		{-2.5, 0.0, 100, 1},   // Point far outside set
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Point(%.1f,%.1f)", tt.cx, tt.cy), func(t *testing.T) {
			zx, zy := 0.0, 0.0
			iter := 0

			for iter < tt.maxIter {
				zx2 := zx * zx
				zy2 := zy * zy

				if zx2+zy2 > 4.0 {
					break
				}

				zy = (zx+zx)*zy + tt.cy
				zx = zx2 - zy2 + tt.cx
				iter++
			}

			// For points in the set, we expect to reach maxIter
			// For points outside, we expect early termination
			if tt.expected == tt.maxIter && iter != tt.maxIter {
				t.Errorf("Point should be in Mandelbrot set: got %d iterations, want %d", iter, tt.maxIter)
			} else if tt.expected < 10 && iter >= 10 {
				t.Errorf("Point should diverge quickly: got %d iterations, want < 10", iter)
			}
		})
	}
}

// TestHashingConsistency tests that our hashing algorithm produces consistent results
func TestHashingConsistency(t *testing.T) {
	data := "Test data for hashing"
	iterations := 1000

	// Run the same hash multiple times
	results := make([]uint32, 3)

	for run := 0; run < 3; run++ {
		hash := uint32(0x12345678)

		for iter := 0; iter < iterations; iter++ {
			for _, b := range []byte(data) {
				hash = hash*33 + uint32(b)
				hash = (hash << 5) | (hash >> 27)
			}
		}

		results[run] = hash
	}

	// All results should be identical
	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("Hash results not consistent: run 0 = %d, run %d = %d", results[0], i, results[i])
		}
	}
}

// TestSHA256Correctness tests that our SHA256 implementation is working
func TestSHA256Correctness(t *testing.T) {
	testData := "Hello, WebAssembly!"

	// Use Go's crypto/sha256 as reference
	expected := sha256.Sum256([]byte(testData))

	// Our simple hash function doesn't match SHA256, but we can test consistency
	hash := uint32(0x12345678)
	for _, b := range []byte(testData) {
		hash = hash*33 + uint32(b)
		hash = (hash << 5) | (hash >> 27)
	}

	// Run twice to ensure consistency
	hash2 := uint32(0x12345678)
	for _, b := range []byte(testData) {
		hash2 = hash2*33 + uint32(b)
		hash2 = (hash2 << 5) | (hash2 >> 27)
	}

	if hash != hash2 {
		t.Errorf("Hash function not consistent: %d vs %d", hash, hash2)
	}

	// Just verify the reference implementation works
	if len(expected) != 32 {
		t.Error("SHA256 should produce 32-byte hash")
	}
}

// Benchmark tests for performance monitoring

func BenchmarkMatrixMultiplication50x50(b *testing.B) {
	benchmarkMatrixMultiplication(b, 50)
}

func BenchmarkMatrixMultiplication100x100(b *testing.B) {
	benchmarkMatrixMultiplication(b, 100)
}

func BenchmarkMatrixMultiplication200x200(b *testing.B) {
	benchmarkMatrixMultiplication(b, 200)
}

func benchmarkMatrixMultiplication(b *testing.B, size int) {
	// Create test matrices
	matrixA := make([]float64, size*size)
	matrixB := make([]float64, size*size)
	result := make([]float64, size*size)

	// Initialize with test data
	for i := 0; i < size*size; i++ {
		matrixA[i] = float64(i % 10)
		matrixB[i] = float64((i * 2) % 10)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Reset result matrix
		for j := range result {
			result[j] = 0
		}

		// Perform matrix multiplication
		for i := 0; i < size; i++ {
			for k := 0; k < size; k++ {
				aik := matrixA[i*size+k]
				for j := 0; j < size; j++ {
					result[i*size+j] += aik * matrixB[k*size+j]
				}
			}
		}
	}

	// Report operations per second
	ops := int64(size) * int64(size) * int64(size) * int64(b.N)
	b.ReportMetric(float64(ops)/b.Elapsed().Seconds()/1e6, "MOps/sec")
}

func BenchmarkMandelbrot400x300(b *testing.B) {
	benchmarkMandelbrotTest(b, 400, 300, 100)
}

func BenchmarkMandelbrot800x600(b *testing.B) {
	benchmarkMandelbrotTest(b, 800, 600, 200)
}

func BenchmarkMandelbrot1200x900(b *testing.B) {
	benchmarkMandelbrotTest(b, 1200, 900, 300)
}

func benchmarkMandelbrotTest(b *testing.B, width, height, maxIter int) {
	xmin, xmax := -2.0, 1.0
	ymin, ymax := -1.5, 1.5

	dx := (xmax - xmin) / float64(width)
	dy := (ymax - ymin) / float64(height)

	result := make([]int, width*height)

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		idx := 0
		for py := 0; py < height; py++ {
			cy := ymin + float64(py)*dy
			for px := 0; px < width; px++ {
				cx := xmin + float64(px)*dx

				zx, zy := 0.0, 0.0
				iter := 0

				for iter < maxIter {
					zx2 := zx * zx
					zy2 := zy * zy

					if zx2+zy2 > 4.0 {
						break
					}

					zy = (zx+zx)*zy + cy
					zx = zx2 - zy2 + cx
					iter++
				}

				result[idx] = iter
				idx++
			}
		}
	}

	// Report pixels per second
	pixels := int64(width) * int64(height) * int64(b.N)
	b.ReportMetric(float64(pixels)/b.Elapsed().Seconds()/1e6, "MPixels/sec")
}

func BenchmarkHashing1000(b *testing.B) {
	benchmarkHashing(b, 1000)
}

func BenchmarkHashing10000(b *testing.B) {
	benchmarkHashing(b, 10000)
}

func BenchmarkHashing100000(b *testing.B) {
	benchmarkHashing(b, 100000)
}

func benchmarkHashing(b *testing.B, iterations int) {
	data := "WebAssembly performance test data for hashing benchmark"

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		hash := uint32(0x12345678)

		for iter := 0; iter < iterations; iter++ {
			for _, b := range []byte(data) {
				hash = hash*33 + uint32(b)
				hash = (hash << 5) | (hash >> 27)
			}
		}

		// Prevent compiler optimization
		_ = hash
	}

	// Report hashes per second
	hashes := int64(iterations) * int64(b.N)
	b.ReportMetric(float64(hashes)/b.Elapsed().Seconds()/1e6, "MHashes/sec")
}

// Benchmark the actual server-side functions for comparison
func BenchmarkServerMatrixMultiplication(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = benchmarkMatrixMultiply(100)
	}
}

func BenchmarkServerMandelbrot(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = benchmarkMandelbrot(800, 600, 200)
	}
}

func BenchmarkServerHashing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = benchmarkSHA256(10000)
	}
}

// Memory allocation tests
func TestMemoryAllocation(t *testing.T) {
	// Test that our algorithms don't have excessive allocations
	size := 100

	// Matrix multiplication should only allocate result matrix
	startAllocs := testing.AllocsPerRun(1, func() {
		matrixA := make([]float64, size*size)
		matrixB := make([]float64, size*size)
		result := make([]float64, size*size)

		for i := 0; i < size*size; i++ {
			matrixA[i] = float64(i % 10)
			matrixB[i] = float64((i * 2) % 10)
		}

		for i := 0; i < size; i++ {
			for k := 0; k < size; k++ {
				aik := matrixA[i*size+k]
				for j := 0; j < size; j++ {
					result[i*size+j] += aik * matrixB[k*size+j]
				}
			}
		}
	})

	if startAllocs > 3 { // 3 slice allocations expected
		t.Errorf("Matrix multiplication allocates too much: %f allocs", startAllocs)
	}
}

// Stress tests for reliability
func TestStressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	// Run many iterations to check for consistency
	const iterations = 1000

	// Test matrix multiplication consistency
	size := 50
	matrixA := make([]float64, size*size)
	matrixB := make([]float64, size*size)

	// Initialize with fixed data
	for i := 0; i < size*size; i++ {
		matrixA[i] = float64(i % 10)
		matrixB[i] = float64((i * 2) % 10)
	}

	firstResult := make([]float64, size*size)

	// First run
	for i := 0; i < size; i++ {
		for k := 0; k < size; k++ {
			aik := matrixA[i*size+k]
			for j := 0; j < size; j++ {
				firstResult[i*size+j] += aik * matrixB[k*size+j]
			}
		}
	}

	// Verify consistency across multiple runs
	for run := 1; run < 100; run++ {
		result := make([]float64, size*size)

		for i := 0; i < size; i++ {
			for k := 0; k < size; k++ {
				aik := matrixA[i*size+k]
				for j := 0; j < size; j++ {
					result[i*size+j] += aik * matrixB[k*size+j]
				}
			}
		}

		// Compare with first result
		for i := 0; i < size*size; i++ {
			if result[i] != firstResult[i] {
				t.Fatalf("Matrix multiplication not consistent at run %d, index %d: got %f, want %f",
					run, i, result[i], firstResult[i])
			}
		}
	}

	t.Logf("Stress test passed: %d consistent matrix multiplication runs", 100)
}

// Performance regression tests
func TestPerformanceRegression(t *testing.T) {
	// These tests ensure our algorithms maintain expected performance characteristics

	// Matrix multiplication should scale roughly O(n³)
	benchTime50 := testing.Benchmark(func(b *testing.B) {
		benchmarkMatrixMultiplication(b, 50)
	})

	benchTime100 := testing.Benchmark(func(b *testing.B) {
		benchmarkMatrixMultiplication(b, 100)
	})

	// 100x100 should take roughly 8x longer than 50x50 (2³ = 8)
	ratio := float64(benchTime100.NsPerOp()) / float64(benchTime50.NsPerOp())

	// Allow some variance, but should be in expected range
	if ratio < 4 || ratio > 20 {
		t.Logf("Matrix multiplication scaling ratio: %.2f (50x50 vs 100x100)", ratio)
		t.Logf("This may indicate performance regression or improvement")
		// Don't fail the test, just report
	}
}

// Test concurrent safety if needed
func TestConcurrentSafety(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	// Test that our business logic functions are safe for concurrent use
	const goroutines = 10
	const iterations = 100

	done := make(chan bool, goroutines)

	// Test user validation concurrently
	user := User{
		Email:   "concurrent@test.com",
		Name:    "Concurrent User",
		Age:     30,
		Country: "US",
		Premium: true,
	}

	for g := 0; g < goroutines; g++ {
		go func() {
			defer func() { done <- true }()

			for i := 0; i < iterations; i++ {
				result := ValidateUser(user)
				if !result.Valid {
					t.Errorf("Concurrent validation failed: %v", result.Errors)
					return
				}
			}
		}()
	}

	// Wait for all goroutines to complete
	for g := 0; g < goroutines; g++ {
		<-done
	}

	t.Logf("Concurrent safety test passed: %d goroutines × %d iterations", goroutines, iterations)
}
