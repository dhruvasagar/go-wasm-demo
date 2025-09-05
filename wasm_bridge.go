//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"
)

// WebAssembly bridge functions - expose shared logic to JavaScript

func validateUserWasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return js.ValueOf("Missing user data")
	}

	userJSON := args[0].String()
	user, err := UserFromJSON(userJSON)
	if err != nil {
		return js.ValueOf("Invalid user JSON: " + err.Error())
	}

	result := ValidateUser(user)
	resultJSON, _ := json.Marshal(result)
	return js.ValueOf(string(resultJSON))
}

func validateProductWasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return js.ValueOf("Missing product data")
	}

	productJSON := args[0].String()
	var product Product
	err := json.Unmarshal([]byte(productJSON), &product)
	if err != nil {
		return js.ValueOf("Invalid product JSON: " + err.Error())
	}

	result := ValidateProduct(product)
	resultJSON, _ := json.Marshal(result)
	return js.ValueOf(string(resultJSON))
}

func calculateOrderTotalWasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return js.ValueOf("Missing order or user data")
	}

	orderJSON := args[0].String()
	userJSON := args[1].String()

	order, err := OrderFromJSON(orderJSON)
	if err != nil {
		return js.ValueOf("Invalid order JSON: " + err.Error())
	}

	user, err := UserFromJSON(userJSON)
	if err != nil {
		return js.ValueOf("Invalid user JSON: " + err.Error())
	}

	CalculateOrderTotal(&order, user)

	resultJSON, _ := json.Marshal(order)
	return js.ValueOf(string(resultJSON))
}

func recommendProductsWasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return js.ValueOf("Missing required data")
	}

	userJSON := args[0].String()
	allProductsJSON := args[1].String()
	currentOrderJSON := args[2].String()

	user, err := UserFromJSON(userJSON)
	if err != nil {
		return js.ValueOf("Invalid user JSON: " + err.Error())
	}

	var allProducts []Product
	err = json.Unmarshal([]byte(allProductsJSON), &allProducts)
	if err != nil {
		return js.ValueOf("Invalid products JSON: " + err.Error())
	}

	currentOrder, err := OrderFromJSON(currentOrderJSON)
	if err != nil {
		return js.ValueOf("Invalid order JSON: " + err.Error())
	}

	recommendations := RecommendProducts(user, allProducts, currentOrder)

	resultJSON, _ := json.Marshal(recommendations)
	return js.ValueOf(string(resultJSON))
}

func analyzeUsersWasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return js.ValueOf("Missing users or orders data")
	}

	usersJSON := args[0].String()
	ordersJSON := args[1].String()

	var users []User
	err := json.Unmarshal([]byte(usersJSON), &users)
	if err != nil {
		return js.ValueOf("Invalid users JSON: " + err.Error())
	}

	var orders []Order
	err = json.Unmarshal([]byte(ordersJSON), &orders)
	if err != nil {
		return js.ValueOf("Invalid orders JSON: " + err.Error())
	}

	analytics := AnalyzeUserBehavior(users, orders)

	resultJSON, _ := json.Marshal(analytics)
	return js.ValueOf(string(resultJSON))
}

func getTaxRateWasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return js.ValueOf(0.08) // Default rate
	}

	country := args[0].String()
	rate := GetTaxRate(country)
	return js.ValueOf(rate)
}

func calculateShippingWasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return js.ValueOf(12.99) // Default shipping
	}

	subtotal := args[0].Float()
	country := args[1].String()
	isPremium := args[2].Bool()

	shipping := CalculateShipping(subtotal, country, isPremium)
	return js.ValueOf(shipping)
}

func formatCurrencyWasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return js.ValueOf("$0.00")
	}

	amount := args[0].Float()
	formatted := FormatCurrency(amount)
	return js.ValueOf(formatted)
}

func getCurrentTimestampWasm(this js.Value, args []js.Value) interface{} {
	timestamp := GetCurrentTimestamp()
	return js.ValueOf(timestamp)
}

// Performance comparison - same Mandelbrot algorithm as before
func mandelbrotWasmBridge(this js.Value, args []js.Value) interface{} {
	return mandelbrotWasmUltra(this, args)
}
