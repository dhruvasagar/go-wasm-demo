//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"
)

func main() {
	// Register shared business logic functions for WebAssembly
	js.Global().Set("validateUserWasm", js.FuncOf(validateUserWasm))
	js.Global().Set("validateProductWasm", js.FuncOf(validateProductWasm))
	js.Global().Set("calculateOrderTotalWasm", js.FuncOf(calculateOrderTotalWasm))
	js.Global().Set("recommendProductsWasm", js.FuncOf(recommendProductsWasm))
	js.Global().Set("analyzeUserBehaviorWasm", js.FuncOf(analyzeUserBehaviorWasm))
	
	// Performance benchmark functions
	js.Global().Set("mandelbrotWasm", js.FuncOf(mandelbrotWasmUltra))
	js.Global().Set("matrixMultiplyWasm", js.FuncOf(matrixMultiplyWasm))
	js.Global().Set("sha256HashWasm", js.FuncOf(sha256HashWasm))
	js.Global().Set("rayTracingWasm", js.FuncOf(rayTracingWasm))
	
	// Keep the program running
	select {}
}

// WebAssembly wrapper for user validation
func validateUserWasm(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return map[string]interface{}{
			"valid":  false,
			"errors": []string{"Invalid number of arguments"},
		}
	}
	
	// Parse JSON input
	userJSON := args[0].String()
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
	return map[string]interface{}{
		"valid":  result.Valid,
		"errors": result.Errors,
	}
}

// WebAssembly wrapper for product validation
func validateProductWasm(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return map[string]interface{}{
			"valid":  false,
			"errors": []string{"Invalid number of arguments"},
		}
	}
	
	productJSON := args[0].String()
	product, err := ProductFromJSON(productJSON)
	if err != nil {
		return map[string]interface{}{
			"valid":  false,
			"errors": []string{"Invalid JSON format: " + err.Error()},
		}
	}
	
	result := ValidateProduct(product)
	
	return map[string]interface{}{
		"valid":  result.Valid,
		"errors": result.Errors,
	}
}

// WebAssembly wrapper for order total calculation
func calculateOrderTotalWasm(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return map[string]interface{}{
			"error": "Invalid number of arguments - expected order and user JSON",
		}
	}
	
	orderJSON := args[0].String()
	userJSON := args[1].String()
	
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
	
	// Use shared business logic
	CalculateOrderTotal(&order, user)
	
	// Return updated order
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
		return []interface{}{}
	}
	
	userJSON := args[0].String()
	productsJSON := args[1].String()
	orderJSON := args[2].String()
	
	user, err := UserFromJSON(userJSON)
	if err != nil {
		return []interface{}{}
	}
	
	var products []Product
	err = json.Unmarshal([]byte(productsJSON), &products)
	if err != nil {
		return []interface{}{}
	}
	
	order, err := OrderFromJSON(orderJSON)
	if err != nil {
		return []interface{}{}
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
	
	return result
}

// WebAssembly wrapper for user behavior analysis
func analyzeUserBehaviorWasm(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return map[string]interface{}{
			"error": "Invalid number of arguments",
		}
	}
	
	usersJSON := args[0].String()
	ordersJSON := args[1].String()
	
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
		"average_age":         analytics.AverageAge,
		"premium_percentage":  analytics.PremiumPercentage,
		"top_countries":       analytics.TopCountries,
		"total_revenue":       analytics.TotalRevenue,
		"average_order_value": analytics.AverageOrderValue,
	}
}

// Additional helper functions for parsing products
func ProductFromJSON(jsonStr string) (Product, error) {
	var product Product
	err := json.Unmarshal([]byte(jsonStr), &product)
	return product, err
}