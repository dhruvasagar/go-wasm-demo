//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"syscall/js"
)

func main() {
	// ====================================================================
	// WASM FUNCTION REGISTRATION
	// ====================================================================
	// This function registers all Go functions to be callable from JavaScript.
	// Functions are organized by category for maintainability.
	//
	// Naming Convention:
	// - Business Logic: [function]Wasm (e.g., validateUserWasm)
	// - Benchmarks: [algorithm]Wasm (e.g., mandelbrotWasm)
	// - Optimized: [algorithm]OptimizedWasm (e.g., matrixMultiplyOptimizedWasm)
	// - Concurrent: [algorithm]ConcurrentWasm (e.g., sha256HashConcurrentWasm)
	// - Legacy: Short names for compatibility (e.g., rayTracing)
	//
	// All functions follow consistent error handling patterns.
	// ====================================================================

	// ====================================================================
	// BUSINESS LOGIC FUNCTIONS
	// Shared business logic that runs identically on client and server
	// ====================================================================
	js.Global().Set("validateUserWasm", js.FuncOf(validateUserWasm))
	js.Global().Set("validateProductWasm", js.FuncOf(validateProductWasm))
	js.Global().Set("calculateOrderTotalWasm", js.FuncOf(calculateOrderTotalWasm))
	js.Global().Set("recommendProductsWasm", js.FuncOf(recommendProductsWasm))
	js.Global().Set("analyzeUserBehaviorWasm", js.FuncOf(analyzeUserBehaviorWasm))

	// ====================================================================
	// BENCHMARK FUNCTIONS - SINGLE-THREADED VERSIONS
	// Basic single-threaded implementations for performance comparison
	// ====================================================================
	js.Global().Set("mandelbrotWasm", js.FuncOf(mandelbrotWasmSingle))
	js.Global().Set("matrixMultiplyWasm", js.FuncOf(matrixMultiplyWasmSingle))
	js.Global().Set("sha256HashWasm", js.FuncOf(sha256HashWasmSingle))
	js.Global().Set("rayTracingWasm", js.FuncOf(rayTracingWasmSingle))

	// ====================================================================
	// BENCHMARK FUNCTIONS - OPTIMIZED VERSIONS
	// Highly optimized single-threaded implementations with boundary call reduction
	// ====================================================================
	js.Global().Set("mandelbrotOptimizedWasm", js.FuncOf(mandelbrotOptimizedWasm))
	js.Global().Set("matrixMultiplyOptimizedWasm", js.FuncOf(matrixMultiplyOptimizedWasm))
	js.Global().Set("sha256HashOptimizedWasm", js.FuncOf(sha256HashOptimizedWasm))
	js.Global().Set("rayTracingOptimizedWasm", js.FuncOf(rayTracingOptimizedWasm))

	// ====================================================================
	// BENCHMARK FUNCTIONS - CONCURRENT VERSIONS
	// Multi-threaded implementations using goroutines for parallel processing
	// ====================================================================
	js.Global().Set("mandelbrotConcurrentWasm", js.FuncOf(mandelbrotWasmConcurrentV2))
	js.Global().Set("matrixMultiplyConcurrentWasm", js.FuncOf(matrixMultiplyWasmConcurrentV2))
	js.Global().Set("sha256HashConcurrentWasm", js.FuncOf(sha256HashWasmConcurrentV2))
	js.Global().Set("rayTracingConcurrentWasm", js.FuncOf(rayTracingWasmConcurrentV2))

	// ====================================================================
	// LEGACY/COMPATIBILITY ALIASES
	// Shorter function names for backward compatibility and ease of use
	// ====================================================================
	// Basic ray tracing from mandelbrot.go
	js.Global().Set("rayTracing", js.FuncOf(rayTracingWasm))

	// User-friendly names for optimized versions
	js.Global().Set("mandelbrotFast", js.FuncOf(mandelbrotOptimizedWasm))
	js.Global().Set("matrixMultiplyFast", js.FuncOf(matrixMultiplyOptimizedWasm))
	js.Global().Set("sha256HashFast", js.FuncOf(sha256HashOptimizedWasm))

	// ====================================================================
	// UTILITY FUNCTIONS
	// Debugging and system information functions
	// ====================================================================
	js.Global().Set("debugConcurrency", js.FuncOf(debugConcurrencyWasm))

	// Keep the program running
	select {}
}

// WebAssembly wrapper for user validation
func validateUserWasm(this js.Value, args []js.Value) interface{} {
	// Handle edge cases and validate input
	if len(args) != 1 {
		return map[string]interface{}{
			"valid":  false,
			"errors": []string{"Invalid number of arguments - expected 1"},
		}
	}

	// Check if argument is valid
	if args[0].Type() != js.TypeString {
		return map[string]interface{}{
			"valid":  false,
			"errors": []string{"Invalid argument type - expected string"},
		}
	}

	// Parse JSON input with safety check
	userJSON := args[0].String()
	if len(userJSON) == 0 {
		return map[string]interface{}{
			"valid":  false,
			"errors": []string{"Empty JSON input"},
		}
	}

	user, err := UserFromJSON(userJSON)
	if err != nil {
		return map[string]interface{}{
			"valid":  false,
			"errors": []string{"Invalid JSON format: " + err.Error()},
		}
	}

	// Use shared business logic
	result := ValidateUser(user)

	// Convert back to JavaScript-compatible format
	// Convert errors slice to JavaScript array
	jsErrors := make([]interface{}, len(result.Errors))
	for i, err := range result.Errors {
		jsErrors[i] = err
	}

	return map[string]interface{}{
		"valid":  result.Valid,
		"errors": jsErrors,
	}
}

// WebAssembly wrapper for product validation
func validateProductWasm(this js.Value, args []js.Value) interface{} {
	// Handle edge cases and validate input
	if len(args) != 1 {
		return map[string]interface{}{
			"valid":  false,
			"errors": []string{"Invalid number of arguments - expected 1"},
		}
	}

	// Check if argument is valid
	if args[0].Type() != js.TypeString {
		return map[string]interface{}{
			"valid":  false,
			"errors": []string{"Invalid argument type - expected string"},
		}
	}

	productJSON := args[0].String()
	if len(productJSON) == 0 {
		return map[string]interface{}{
			"valid":  false,
			"errors": []string{"Empty JSON input"},
		}
	}

	product, err := ProductFromJSON(productJSON)
	if err != nil {
		return map[string]interface{}{
			"valid":  false,
			"errors": []string{"Invalid JSON format: " + err.Error()},
		}
	}

	result := ValidateProduct(product)

	// Convert errors slice to JavaScript array
	jsErrors := make([]interface{}, len(result.Errors))
	for i, err := range result.Errors {
		jsErrors[i] = err
	}

	return map[string]interface{}{
		"valid":  result.Valid,
		"errors": jsErrors,
	}
}

// WebAssembly wrapper for order total calculation
func calculateOrderTotalWasm(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return map[string]interface{}{
			"error": "Invalid number of arguments - expected order and user JSON",
		}
	}

	// Validate argument types
	if args[0].Type() != js.TypeString || args[1].Type() != js.TypeString {
		return map[string]interface{}{
			"error": "Invalid argument types - expected strings",
		}
	}

	orderJSON := args[0].String()
	userJSON := args[1].String()

	// Validate input is not empty
	if len(orderJSON) == 0 {
		return map[string]interface{}{
			"error": "Empty order JSON",
		}
	}
	if len(userJSON) == 0 {
		return map[string]interface{}{
			"error": "Empty user JSON",
		}
	}

	order, err := OrderFromJSON(orderJSON)
	if err != nil {
		return map[string]interface{}{
			"error": "Invalid order JSON: " + err.Error(),
		}
	}

	user, err := UserFromJSON(userJSON)
	if err != nil {
		return map[string]interface{}{
			"error": "Invalid user JSON: " + err.Error(),
		}
	}

	// Validate order has products
	if len(order.Products) == 0 {
		return map[string]interface{}{
			"error": "Order must contain at least one product",
		}
	}

	// Validate quantities match products
	if len(order.Products) != len(order.Quantities) {
		return map[string]interface{}{
			"error": "Product and quantity arrays must be the same length",
		}
	}

	// Use shared business logic
	CalculateOrderTotal(&order, user)

	// Return updated order with validation
	return map[string]interface{}{
		"subtotal": order.Subtotal,
		"tax":      order.Tax,
		"shipping": order.Shipping,
		"discount": order.Discount,
		"total":    order.Total,
	}
}

// WebAssembly wrapper for product recommendations
func recommendProductsWasm(this js.Value, args []js.Value) interface{} {
	if len(args) != 3 {
		return map[string]interface{}{
			"error":           "Invalid number of arguments - expected 3",
			"recommendations": []interface{}{},
		}
	}

	// Validate argument types
	for i, arg := range args {
		if arg.Type() != js.TypeString {
			return map[string]interface{}{
				"error":           fmt.Sprintf("Argument %d is not a string", i),
				"recommendations": []interface{}{},
			}
		}
	}

	userJSON := args[0].String()
	productsJSON := args[1].String()
	orderJSON := args[2].String()

	// Validate inputs are not empty
	if len(userJSON) == 0 || len(productsJSON) == 0 || len(orderJSON) == 0 {
		return map[string]interface{}{
			"error":           "One or more JSON inputs are empty",
			"recommendations": []interface{}{},
		}
	}

	user, err := UserFromJSON(userJSON)
	if err != nil {
		return map[string]interface{}{
			"error":           "Invalid user JSON: " + err.Error(),
			"recommendations": []interface{}{},
		}
	}

	var products []Product
	err = json.Unmarshal([]byte(productsJSON), &products)
	if err != nil {
		return map[string]interface{}{
			"error":           "Invalid products JSON: " + err.Error(),
			"recommendations": []interface{}{},
		}
	}

	order, err := OrderFromJSON(orderJSON)
	if err != nil {
		return map[string]interface{}{
			"error":           "Invalid order JSON: " + err.Error(),
			"recommendations": []interface{}{},
		}
	}

	// Use shared business logic
	recommendations := RecommendProducts(user, products, order)

	// Convert to JavaScript-compatible format
	result := make([]interface{}, len(recommendations))
	for i, product := range recommendations {
		result[i] = map[string]interface{}{
			"id":          product.ID,
			"name":        product.Name,
			"price":       product.Price,
			"category":    product.Category,
			"in_stock":    product.InStock,
			"rating":      product.Rating,
			"description": product.Description,
		}
	}

	return map[string]interface{}{
		"error":           "",
		"recommendations": result,
	}
}

// WebAssembly wrapper for user behavior analysis
func analyzeUserBehaviorWasm(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return map[string]interface{}{
			"error": "Invalid number of arguments - expected 2",
		}
	}

	// Validate argument types
	if args[0].Type() != js.TypeString || args[1].Type() != js.TypeString {
		return map[string]interface{}{
			"error": "Invalid argument types - expected strings",
		}
	}

	usersJSON := args[0].String()
	ordersJSON := args[1].String()

	// Validate inputs are not empty
	if len(usersJSON) == 0 {
		return map[string]interface{}{
			"error": "Empty users JSON",
		}
	}
	if len(ordersJSON) == 0 {
		return map[string]interface{}{
			"error": "Empty orders JSON",
		}
	}

	var users []User
	err := json.Unmarshal([]byte(usersJSON), &users)
	if err != nil {
		return map[string]interface{}{
			"error": "Invalid users JSON: " + err.Error(),
		}
	}

	var orders []Order
	err = json.Unmarshal([]byte(ordersJSON), &orders)
	if err != nil {
		return map[string]interface{}{
			"error": "Invalid orders JSON: " + err.Error(),
		}
	}

	// Use shared business logic
	analytics := AnalyzeUserBehavior(users, orders)

	return map[string]interface{}{
		"error":               "",
		"average_age":         analytics.AverageAge,
		"premium_percentage":  analytics.PremiumPercentage,
		"top_countries":       analytics.TopCountries,
		"total_revenue":       analytics.TotalRevenue,
		"average_order_value": analytics.AverageOrderValue,
	}
}

// ====================================================================
// UTILITY FUNCTIONS
// ====================================================================

// Debug function to check concurrency and system info
func debugConcurrencyWasm(this js.Value, args []js.Value) interface{} {
	return map[string]interface{}{
		"GOMAXPROCS":    runtime.GOMAXPROCS(0),
		"NumCPU":        runtime.NumCPU(),
		"NumGoroutines": runtime.NumGoroutine(),
		"GoVersion":     runtime.Version(),
		"GOARCH":        runtime.GOARCH,
		"GOOS":          runtime.GOOS,
	}
}

// Additional helper functions for parsing products
func ProductFromJSON(jsonStr string) (Product, error) {
	var product Product
	err := json.Unmarshal([]byte(jsonStr), &product)
	return product, err
}
