package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Shared data models - used identically on both server and client
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Country  string `json:"country"`
	Premium  bool   `json:"premium"`
	JoinDate string `json:"join_date"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	InStock     bool    `json:"in_stock"`
	Rating      float64 `json:"rating"`
	Description string  `json:"description"`
}

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Products   []Product `json:"products"`
	Quantities []int     `json:"quantities"`
	Subtotal   float64   `json:"subtotal"`
	Tax        float64   `json:"tax"`
	Shipping   float64   `json:"shipping"`
	Total      float64   `json:"total"`
	Discount   float64   `json:"discount"`
	OrderDate  string    `json:"order_date"`
	Status     string    `json:"status"`
}

type ValidationResult struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors"`
}

// Shared business logic - identical implementation on server and client
func ValidateUser(user User) ValidationResult {
	result := ValidationResult{Valid: true, Errors: []string{}}

	// Email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		result.Valid = false
		result.Errors = append(result.Errors, "Invalid email format")
	}

	// Name validation
	if len(strings.TrimSpace(user.Name)) < 2 {
		result.Valid = false
		result.Errors = append(result.Errors, "Name must be at least 2 characters")
	}

	// Age validation
	if user.Age < 13 || user.Age > 120 {
		result.Valid = false
		result.Errors = append(result.Errors, "Age must be between 13 and 120")
	}

	// Country validation
	validCountries := []string{"US", "CA", "UK", "DE", "FR", "JP", "AU", "IN", "BR", "MX"}
	isValidCountry := false
	for _, country := range validCountries {
		if user.Country == country {
			isValidCountry = true
			break
		}
	}
	if !isValidCountry {
		result.Valid = false
		result.Errors = append(result.Errors, "Invalid country code")
	}

	return result
}

func ValidateProduct(product Product) ValidationResult {
	result := ValidationResult{Valid: true, Errors: []string{}}

	// Name validation
	if len(strings.TrimSpace(product.Name)) < 3 {
		result.Valid = false
		result.Errors = append(result.Errors, "Product name must be at least 3 characters")
	}

	// Price validation
	if product.Price <= 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "Price must be greater than 0")
	}

	if product.Price > 10000 {
		result.Valid = false
		result.Errors = append(result.Errors, "Price cannot exceed $10,000")
	}

	// Category validation
	validCategories := []string{"electronics", "clothing", "books", "home", "sports", "toys", "beauty"}
	isValidCategory := false
	for _, cat := range validCategories {
		if strings.ToLower(product.Category) == cat {
			isValidCategory = true
			break
		}
	}
	if !isValidCategory {
		result.Valid = false
		result.Errors = append(result.Errors, "Invalid category")
	}

	// Rating validation
	if product.Rating < 0 || product.Rating > 5 {
		result.Valid = false
		result.Errors = append(result.Errors, "Rating must be between 0 and 5")
	}

	return result
}

func CalculateOrderTotal(order *Order, user User) {
	// Calculate subtotal
	order.Subtotal = 0
	for i, product := range order.Products {
		if i < len(order.Quantities) {
			order.Subtotal += product.Price * float64(order.Quantities[i])
		}
	}

	// Apply premium discount
	order.Discount = 0
	if user.Premium {
		if order.Subtotal > 100 {
			order.Discount = order.Subtotal * 0.15 // 15% premium discount
		} else if order.Subtotal > 50 {
			order.Discount = order.Subtotal * 0.10 // 10% premium discount
		}
	}

	// Calculate tax (varies by country)
	taxRate := GetTaxRate(user.Country)
	order.Tax = (order.Subtotal - order.Discount) * taxRate

	// Calculate shipping
	order.Shipping = CalculateShipping(order.Subtotal, user.Country, user.Premium)

	// Calculate total
	order.Total = order.Subtotal - order.Discount + order.Tax + order.Shipping
}

func GetTaxRate(country string) float64 {
	taxRates := map[string]float64{
		"US": 0.08, // 8%
		"CA": 0.13, // 13% GST+PST
		"UK": 0.20, // 20% VAT
		"DE": 0.19, // 19% VAT
		"FR": 0.20, // 20% VAT
		"JP": 0.10, // 10% consumption tax
		"AU": 0.10, // 10% GST
		"IN": 0.18, // 18% GST
		"BR": 0.17, // 17% ICMS
		"MX": 0.16, // 16% IVA
	}

	if rate, exists := taxRates[country]; exists {
		return rate
	}
	return 0.08 // Default 8%
}

func CalculateShipping(subtotal float64, country string, isPremium bool) float64 {
	if isPremium && subtotal > 75 {
		return 0 // Free shipping for premium users over $75
	}

	// Base shipping rates by country
	shippingRates := map[string]float64{
		"US": 8.99,
		"CA": 12.99,
		"UK": 15.99,
		"DE": 14.99,
		"FR": 14.99,
		"JP": 18.99,
		"AU": 19.99,
		"IN": 9.99,
		"BR": 16.99,
		"MX": 13.99,
	}

	baseRate := 12.99 // Default
	if rate, exists := shippingRates[country]; exists {
		baseRate = rate
	}

	// Free shipping threshold
	if subtotal > 100 {
		return 0
	}

	// Express shipping for orders over $50
	if subtotal > 50 {
		return baseRate * 1.5
	}

	return baseRate
}

// Advanced business logic - recommendation algorithm
func RecommendProducts(user User, allProducts []Product, currentOrder Order) []Product {
	recommendations := []Product{}
	userCategory := inferUserPreference(user, currentOrder)

	// Score-based recommendation
	productScores := make(map[int]float64)

	for _, product := range allProducts {
		if !product.InStock {
			continue
		}

		score := 0.0

		// Category preference
		if strings.ToLower(product.Category) == userCategory {
			score += 3.0
		}

		// Price preference based on user's current order
		avgOrderPrice := getAverageProductPrice(currentOrder)
		priceDiff := abs(product.Price - avgOrderPrice)
		if priceDiff < avgOrderPrice*0.3 { // Within 30% of average
			score += 2.0
		}

		// Rating boost
		score += product.Rating * 0.5

		// Premium user gets higher-end recommendations
		if user.Premium && product.Price > avgOrderPrice*1.2 {
			score += 1.0
		}

		// Age-based preferences
		if user.Age < 25 && (product.Category == "electronics" || product.Category == "toys") {
			score += 1.0
		} else if user.Age > 40 && (product.Category == "home" || product.Category == "books") {
			score += 1.0
		}

		productScores[product.ID] = score
	}

	// Sort by score and return top 5
	for len(recommendations) < 5 {
		bestID := -1
		bestScore := -1.0

		for id, score := range productScores {
			if score > bestScore {
				bestScore = score
				bestID = id
			}
		}

		if bestID == -1 {
			break
		}

		// Find and add the product
		for _, product := range allProducts {
			if product.ID == bestID {
				recommendations = append(recommendations, product)
				break
			}
		}

		delete(productScores, bestID)
	}

	return recommendations
}

func inferUserPreference(user User, order Order) string {
	if len(order.Products) == 0 {
		// Default preferences by age
		if user.Age < 25 {
			return "electronics"
		} else if user.Age < 40 {
			return "clothing"
		} else {
			return "home"
		}
	}

	// Find most common category in current order
	categoryCount := make(map[string]int)
	for _, product := range order.Products {
		categoryCount[strings.ToLower(product.Category)]++
	}

	mostCommon := ""
	maxCount := 0
	for category, count := range categoryCount {
		if count > maxCount {
			maxCount = count
			mostCommon = category
		}
	}

	return mostCommon
}

func getAverageProductPrice(order Order) float64 {
	if len(order.Products) == 0 {
		return 50.0 // Default
	}

	total := 0.0
	for _, product := range order.Products {
		total += product.Price
	}

	return total / float64(len(order.Products))
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Data processing and analytics - same algorithms on server and client
func AnalyzeUserBehavior(users []User, orders []Order) UserAnalytics {
	analytics := UserAnalytics{}

	if len(users) == 0 {
		return analytics
	}

	// Calculate demographics
	ageSum := 0
	countryCount := make(map[string]int)
	premiumCount := 0

	for _, user := range users {
		ageSum += user.Age
		countryCount[user.Country]++
		if user.Premium {
			premiumCount++
		}
	}

	analytics.AverageAge = float64(ageSum) / float64(len(users))
	analytics.PremiumPercentage = (float64(premiumCount) / float64(len(users))) * 100
	analytics.TopCountries = getTopCountries(countryCount, 3)

	// Analyze orders
	if len(orders) > 0 {
		totalRevenue := 0.0
		totalOrders := len(orders)

		for _, order := range orders {
			totalRevenue += order.Total
		}

		analytics.TotalRevenue = totalRevenue
		analytics.AverageOrderValue = totalRevenue / float64(totalOrders)
	}

	return analytics
}

type UserAnalytics struct {
	AverageAge        float64  `json:"average_age"`
	PremiumPercentage float64  `json:"premium_percentage"`
	TopCountries      []string `json:"top_countries"`
	TotalRevenue      float64  `json:"total_revenue"`
	AverageOrderValue float64  `json:"average_order_value"`
}

func getTopCountries(countryCount map[string]int, limit int) []string {
	type countryPair struct {
		country string
		count   int
	}

	pairs := []countryPair{}
	for country, count := range countryCount {
		pairs = append(pairs, countryPair{country, count})
	}

	// Simple bubble sort
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].count > pairs[i].count {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	result := []string{}
	for i := 0; i < limit && i < len(pairs); i++ {
		result = append(result, pairs[i].country)
	}

	return result
}

// JSON serialization helpers - identical on both sides
func UserToJSON(user User) string {
	data, _ := json.Marshal(user)
	return string(data)
}

func UserFromJSON(jsonStr string) (User, error) {
	var user User
	err := json.Unmarshal([]byte(jsonStr), &user)
	return user, err
}

func OrderToJSON(order Order) string {
	data, _ := json.Marshal(order)
	return string(data)
}

func OrderFromJSON(jsonStr string) (Order, error) {
	var order Order
	err := json.Unmarshal([]byte(jsonStr), &order)
	return order, err
}

// Utility functions
func FormatCurrency(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

func GetCurrentTimestamp() string {
	return time.Now().Format("2006-01-02T15:04:05Z")
}
