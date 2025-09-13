//go:build js && wasm

package main

import (
	"syscall/js"
)

// ============================================================================
// UNIFIED BENCHMARK CONFIGURATION
// Consolidates all benchmark variants into a single, configurable interface
// ============================================================================

// BenchmarkConfig defines the optimization level and parameters for benchmarks
type BenchmarkConfig struct {
	OptimizationLevel string // "single", "optimized", "concurrent"
	UseBulkCopy       bool
	UseTypedArrays    bool
	Workers           int
}

// Predefined configurations for common use cases
var (
	SingleThreadedConfig = BenchmarkConfig{
		OptimizationLevel: "single",
		UseBulkCopy:       false,
		UseTypedArrays:    false,
		Workers:           1,
	}

	OptimizedConfig = BenchmarkConfig{
		OptimizationLevel: "optimized",
		UseBulkCopy:       true,
		UseTypedArrays:    true,
		Workers:           1,
	}

	ConcurrentConfig = BenchmarkConfig{
		OptimizationLevel: "concurrent",
		UseBulkCopy:       true,
		UseTypedArrays:    true,
		Workers:           4,
	}
)

// ============================================================================
// UNIFIED MATRIX MULTIPLICATION
// Single function that handles all matrix multiplication variants
// ============================================================================

func createUnifiedMatrixMultiplyWasm(config BenchmarkConfig) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 3 {
			return js.ValueOf("Missing arguments: expected matrixA, matrixB, size")
		}

		switch config.OptimizationLevel {
		case "single":
			return matrixMultiplyWasmSingle(this, args)
		case "optimized":
			return matrixMultiplyOptimizedWasm(this, args)
		case "concurrent":
			return matrixMultiplyWasmConcurrentV2(this, args)
		default:
			return js.ValueOf("Invalid optimization level: " + config.OptimizationLevel)
		}
	})
}

// ============================================================================
// UNIFIED MANDELBROT GENERATION
// Single function that handles all Mandelbrot variants
// ============================================================================

func createUnifiedMandelbrotWasm(config BenchmarkConfig) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 6 {
			return js.ValueOf("Missing arguments: expected width, height, xmin, xmax, ymin, ymax, [iterations]")
		}

		switch config.OptimizationLevel {
		case "single":
			return mandelbrotWasmSingle(this, args)
		case "optimized":
			return mandelbrotOptimizedWasm(this, args)
		case "concurrent":
			return mandelbrotWasmConcurrentV2(this, args)
		default:
			return js.ValueOf("Invalid optimization level: " + config.OptimizationLevel)
		}
	})
}

// ============================================================================
// UNIFIED HASH COMPUTATION
// Single function that handles all hash variants
// ============================================================================

func createUnifiedHashWasm(config BenchmarkConfig) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 2 {
			return js.ValueOf("Missing arguments: expected data, iterations")
		}

		switch config.OptimizationLevel {
		case "single":
			return sha256HashWasmSingle(this, args)
		case "optimized":
			return sha256HashOptimizedWasm(this, args)
		case "concurrent":
			return sha256HashWasmConcurrentV2(this, args)
		default:
			return js.ValueOf("Invalid optimization level: " + config.OptimizationLevel)
		}
	})
}

// ============================================================================
// UNIFIED RAY TRACING
// Single function that handles all ray tracing variants
// ============================================================================

func createUnifiedRayTracingWasm(config BenchmarkConfig) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 3 {
			return js.ValueOf("Missing arguments: expected width, height, samples")
		}

		switch config.OptimizationLevel {
		case "single":
			return rayTracingWasmSingle(this, args)
		case "optimized":
			return rayTracingOptimizedWasm(this, args)
		case "concurrent":
			return rayTracingWasmConcurrentV2(this, args)
		default:
			return js.ValueOf("Invalid optimization level: " + config.OptimizationLevel)
		}
	})
}

// ============================================================================
// BENCHMARK FACTORY FUNCTIONS
// Creates benchmark functions with specific configurations
// ============================================================================

// Creates a complete set of benchmark functions for a given configuration
func createBenchmarkSuite(prefix string, config BenchmarkConfig) map[string]js.Func {
	return map[string]js.Func{
		prefix + "MatrixMultiply": createUnifiedMatrixMultiplyWasm(config),
		prefix + "Mandelbrot":     createUnifiedMandelbrotWasm(config),
		prefix + "Hash":           createUnifiedHashWasm(config),
		prefix + "RayTracing":     createUnifiedRayTracingWasm(config),
	}
}

// Registers all benchmark variants using the unified interface
func registerUnifiedBenchmarks() {
	// Register single-threaded benchmarks
	singleSuite := createBenchmarkSuite("single", SingleThreadedConfig)
	for name, fn := range singleSuite {
		js.Global().Set(name+"Wasm", fn)
	}

	// Register optimized benchmarks
	optimizedSuite := createBenchmarkSuite("optimized", OptimizedConfig)
	for name, fn := range optimizedSuite {
		js.Global().Set(name+"WasmFast", fn)
	}

	// Register concurrent benchmarks
	concurrentSuite := createBenchmarkSuite("concurrent", ConcurrentConfig)
	for name, fn := range concurrentSuite {
		js.Global().Set(name+"WasmConcurrent", fn)
	}
}

// ============================================================================
// CONFIGURATION UTILITIES
// Helper functions to create custom benchmark configurations
// ============================================================================

// Creates a custom benchmark configuration
func createCustomConfig(optimizationLevel string, workers int, useBulkCopy bool) BenchmarkConfig {
	return BenchmarkConfig{
		OptimizationLevel: optimizationLevel,
		UseBulkCopy:       useBulkCopy,
		UseTypedArrays:    useBulkCopy, // Usually correlated
		Workers:           workers,
	}
}

// Validates a benchmark configuration
func (config BenchmarkConfig) isValid() bool {
	validLevels := map[string]bool{
		"single":     true,
		"optimized":  true,
		"concurrent": true,
	}

	return validLevels[config.OptimizationLevel] &&
		config.Workers > 0 &&
		config.Workers <= 16 // Reasonable limit
}
