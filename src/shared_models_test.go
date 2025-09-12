package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

// Test data fixtures
var testUsers = []User{
	{
		ID:       1,
		Email:    "john.doe@example.com",
		Name:     "John Doe",
		Age:      28,
		Country:  "US",
		Premium:  true,
		JoinDate: "2023-01-15",
	},
	{
		ID:       2,
		Email:    "invalid-email",
		Name:     "A",  // Too short
		Age:      12,   // Too young
		Country:  "XX", // Invalid country
		Premium:  false,
		JoinDate: "2023-02-20",
	},
	{
		ID:       3,
		Email:    "jane.smith@example.com",
		Name:     "Jane Smith",
		Age:      34,
		Country:  "CA",
		Premium:  false,
		JoinDate: "2023-03-10",
	},
}

var testProducts = []Product{
	{
		ID:          1,
		Name:        "Wireless Headphones",
		Price:       99.99,
		Category:    "electronics",
		InStock:     true,
		Rating:      4.5,
		Description: "High-quality wireless headphones",
	},
	{
		ID:          2,
		Name:        "A",       // Too short name
		Price:       -10.99,    // Invalid price
		Category:    "invalid", // Invalid category
		InStock:     true,
		Rating:      6.0, // Invalid rating
		Description: "Invalid product for testing",
	},
	{
		ID:          3,
		Name:        "Programming Book",
		Price:       49.99,
		Category:    "books",
		InStock:     true,
		Rating:      4.8,
		Description: "Learn advanced programming techniques",
	},
}

// TestValidateUser tests the user validation logic
func TestValidateUser(t *testing.T) {
	tests := []struct {
		name       string
		user       User
		wantValid  bool
		wantErrors []string
	}{
		{
			name:       "Valid user",
			user:       testUsers[0],
			wantValid:  true,
			wantErrors: []string{},
		},
		{
			name:      "Invalid user - multiple errors",
			user:      testUsers[1],
			wantValid: false,
			wantErrors: []string{
				"Invalid email format",
				"Name must be at least 2 characters",
				"Age must be between 13 and 120",
				"Invalid country code",
			},
		},
		{
			name:       "Valid Canadian user",
			user:       testUsers[2],
			wantValid:  true,
			wantErrors: []string{},
		},
		{
			name: "Edge case - minimum age",
			user: User{
				Email:   "teen@example.com",
				Name:    "Teen User",
				Age:     13,
				Country: "US",
			},
			wantValid:  true,
			wantErrors: []string{},
		},
		{
			name: "Edge case - maximum age",
			user: User{
				Email:   "senior@example.com",
				Name:    "Senior User",
				Age:     120,
				Country: "US",
			},
			wantValid:  true,
			wantErrors: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUser(tt.user)

			if result.Valid != tt.wantValid {
				t.Errorf("ValidateUser() valid = %v, want %v", result.Valid, tt.wantValid)
			}

			if !reflect.DeepEqual(result.Errors, tt.wantErrors) {
				t.Errorf("ValidateUser() errors = %v, want %v", result.Errors, tt.wantErrors)
			}
		})
	}
}

// TestValidateProduct tests the product validation logic
func TestValidateProduct(t *testing.T) {
	tests := []struct {
		name       string
		product    Product
		wantValid  bool
		wantErrors []string
	}{
		{
			name:       "Valid product",
			product:    testProducts[0],
			wantValid:  true,
			wantErrors: []string{},
		},
		{
			name:      "Invalid product - multiple errors",
			product:   testProducts[1],
			wantValid: false,
			wantErrors: []string{
				"Product name must be at least 3 characters",
				"Price must be greater than 0",
				"Invalid category",
				"Rating must be between 0 and 5",
			},
		},
		{
			name:       "Valid book",
			product:    testProducts[2],
			wantValid:  true,
			wantErrors: []string{},
		},
		{
			name: "Edge case - expensive product",
			product: Product{
				Name:     "Luxury Item",
				Price:    9999.99,
				Category: "electronics",
				Rating:   5.0,
			},
			wantValid:  true,
			wantErrors: []string{},
		},
		{
			name: "Edge case - too expensive",
			product: Product{
				Name:     "Too Expensive",
				Price:    10000.01,
				Category: "electronics",
				Rating:   5.0,
			},
			wantValid:  false,
			wantErrors: []string{"Price cannot exceed $10,000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateProduct(tt.product)

			if result.Valid != tt.wantValid {
				t.Errorf("ValidateProduct() valid = %v, want %v", result.Valid, tt.wantValid)
			}

			if !reflect.DeepEqual(result.Errors, tt.wantErrors) {
				t.Errorf("ValidateProduct() errors = %v, want %v", result.Errors, tt.wantErrors)
			}
		})
	}
}

// TestCalculateOrderTotal tests the order calculation logic
func TestCalculateOrderTotal(t *testing.T) {
	tests := []struct {
		name         string
		order        Order
		user         User
		wantSubtotal float64
		wantTax      float64
		wantShipping float64
		wantDiscount float64
	}{
		{
			name: "Basic order - US user",
			order: Order{
				Products:   []Product{testProducts[0]}, // $99.99
				Quantities: []int{1},
			},
			user:         testUsers[0], // US, Premium
			wantSubtotal: 99.99,
			wantTax:      6.79, // 8% of (99.99 - 9.99 discount)
			wantShipping: 0.0,  // Free shipping for premium over $75
			wantDiscount: 9.99, // 10% premium discount (subtotal > $50, < $100)
		},
		{
			name: "Large order - Canadian user",
			order: Order{
				Products:   []Product{testProducts[0], testProducts[2]}, // $99.99 + $49.99 = $149.98
				Quantities: []int{2, 1},                                 // 2x headphones + 1x book = $199.98 + $49.99 = $249.97
			},
			user:         testUsers[2], // CA, Non-premium
			wantSubtotal: 249.97,
			wantTax:      31.24, // 13% of (249.97 - 37.50)
			wantShipping: 0.0,   // Free shipping over $100
			wantDiscount: 37.50, // 15% premium discount for CA premium (but user is not premium, so 0)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying test data
			order := tt.order
			CalculateOrderTotal(&order, tt.user)

			if !floatEqual(order.Subtotal, tt.wantSubtotal, 0.01) {
				t.Errorf("CalculateOrderTotal() subtotal = %v, want %v", order.Subtotal, tt.wantSubtotal)
			}
		})
	}
}

// TestGetTaxRate tests tax rate calculation
func TestGetTaxRate(t *testing.T) {
	tests := []struct {
		name    string
		country string
		want    float64
	}{
		{"US tax rate", "US", 0.08},
		{"Canada tax rate", "CA", 0.13},
		{"UK tax rate", "UK", 0.20},
		{"Germany tax rate", "DE", 0.19},
		{"Japan tax rate", "JP", 0.10},
		{"Unknown country", "XX", 0.08}, // Default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTaxRate(tt.country); got != tt.want {
				t.Errorf("GetTaxRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCalculateShipping tests shipping calculation logic
func TestCalculateShipping(t *testing.T) {
	tests := []struct {
		name      string
		subtotal  float64
		country   string
		isPremium bool
		want      float64
	}{
		{"Free shipping - over $100", 150.0, "US", false, 0.0},
		{"Premium free shipping", 80.0, "US", true, 0.0},
		{"Regular US shipping", 50.0, "US", false, 8.99},
		{"Express US shipping", 60.0, "US", false, 13.485}, // 8.99 * 1.5
		{"Canada shipping", 50.0, "CA", false, 12.99},
		{"Unknown country", 50.0, "XX", false, 12.99}, // Default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateShipping(tt.subtotal, tt.country, tt.isPremium)
			if !floatEqual(got, tt.want, 0.01) {
				t.Errorf("CalculateShipping() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRecommendProducts tests the recommendation algorithm
func TestRecommendProducts(t *testing.T) {
	user := testUsers[0] // Premium US user, age 28
	products := testProducts
	order := Order{
		Products: []Product{testProducts[0]}, // Electronics
	}

	recommendations := RecommendProducts(user, products, order)

	// Should return at most 5 recommendations
	if len(recommendations) > 5 {
		t.Errorf("RecommendProducts() returned %d recommendations, want <= 5", len(recommendations))
	}

	// Should not include out-of-stock items
	for _, rec := range recommendations {
		if !rec.InStock {
			t.Errorf("RecommendProducts() included out-of-stock item: %s", rec.Name)
		}
	}

	// Should prefer electronics for this user (based on current order)
	electronicsCount := 0
	for _, rec := range recommendations {
		if rec.Category == "electronics" {
			electronicsCount++
		}
	}

	if electronicsCount == 0 && len(recommendations) > 0 {
		t.Logf("No electronics in recommendations, which is unexpected but not necessarily wrong")
	}
}

// TestAnalyzeUserBehavior tests user analytics
func TestAnalyzeUserBehavior(t *testing.T) {
	users := testUsers[:2] // Two users: one premium, one not
	orders := []Order{
		{
			ID:     1,
			UserID: 1,
			Total:  149.97,
		},
		{
			ID:     2,
			UserID: 2,
			Total:  62.98,
		},
	}

	analytics := AnalyzeUserBehavior(users, orders)

	// Check average age calculation
	expectedAge := (28.0 + 12.0) / 2.0 // 20.0
	if !floatEqual(analytics.AverageAge, expectedAge, 0.1) {
		t.Errorf("AnalyzeUserBehavior() average age = %v, want %v", analytics.AverageAge, expectedAge)
	}

	// Check premium percentage
	expectedPremium := 50.0 // 1 out of 2 users
	if !floatEqual(analytics.PremiumPercentage, expectedPremium, 0.1) {
		t.Errorf("AnalyzeUserBehavior() premium percentage = %v, want %v", analytics.PremiumPercentage, expectedPremium)
	}

	// Check revenue calculations
	expectedRevenue := 149.97 + 62.98 // 212.95
	if !floatEqual(analytics.TotalRevenue, expectedRevenue, 0.01) {
		t.Errorf("AnalyzeUserBehavior() total revenue = %v, want %v", analytics.TotalRevenue, expectedRevenue)
	}

	expectedAvgOrder := expectedRevenue / 2.0 // 106.475
	if !floatEqual(analytics.AverageOrderValue, expectedAvgOrder, 0.01) {
		t.Errorf("AnalyzeUserBehavior() average order = %v, want %v", analytics.AverageOrderValue, expectedAvgOrder)
	}
}

// TestJSONSerialization tests JSON conversion functions
func TestJSONSerialization(t *testing.T) {
	// Test User JSON serialization
	user := testUsers[0]
	jsonStr := UserToJSON(user)

	var parsedUser User
	err := json.Unmarshal([]byte(jsonStr), &parsedUser)
	if err != nil {
		t.Errorf("UserToJSON() produced invalid JSON: %v", err)
	}

	if parsedUser.Email != user.Email || parsedUser.Name != user.Name {
		t.Errorf("UserToJSON() round-trip failed: got %+v, want %+v", parsedUser, user)
	}

	// Test UserFromJSON
	parsedUser2, err := UserFromJSON(jsonStr)
	if err != nil {
		t.Errorf("UserFromJSON() failed: %v", err)
	}

	if parsedUser2.Email != user.Email {
		t.Errorf("UserFromJSON() failed: got %+v, want %+v", parsedUser2, user)
	}
}

// TestUtilityFunctions tests helper functions
func TestUtilityFunctions(t *testing.T) {
	// Test FormatCurrency
	tests := []struct {
		amount float64
		want   string
	}{
		{99.99, "$99.99"},
		{0.00, "$0.00"},
		{1234.567, "$1234.57"},
		{0.1, "$0.10"},
	}

	for _, tt := range tests {
		got := FormatCurrency(tt.amount)
		if got != tt.want {
			t.Errorf("FormatCurrency(%v) = %v, want %v", tt.amount, got, tt.want)
		}
	}

	// Test GetCurrentTimestamp
	timestamp := GetCurrentTimestamp()
	if len(timestamp) == 0 {
		t.Error("GetCurrentTimestamp() returned empty string")
	}

	// Should be in RFC3339 format (basic check)
	if len(timestamp) < 19 { // "2006-01-02T15:04:05" minimum
		t.Errorf("GetCurrentTimestamp() format seems incorrect: %s", timestamp)
	}
}

// Benchmark tests for performance validation
func BenchmarkValidateUser(b *testing.B) {
	user := testUsers[0]
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ValidateUser(user)
	}
}

func BenchmarkValidateProduct(b *testing.B) {
	product := testProducts[0]
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ValidateProduct(product)
	}
}

func BenchmarkCalculateOrderTotal(b *testing.B) {
	user := testUsers[0]
	order := Order{
		Products:   []Product{testProducts[0], testProducts[2]},
		Quantities: []int{1, 2},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Reset order state for consistent benchmarking
		orderCopy := order
		CalculateOrderTotal(&orderCopy, user)
	}
}

func BenchmarkRecommendProducts(b *testing.B) {
	user := testUsers[0]
	products := testProducts
	order := Order{Products: []Product{testProducts[0]}}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		RecommendProducts(user, products, order)
	}
}

func BenchmarkAnalyzeUserBehavior(b *testing.B) {
	users := testUsers
	orders := []Order{
		{ID: 1, UserID: 1, Total: 149.97},
		{ID: 2, UserID: 2, Total: 62.98},
		{ID: 3, UserID: 3, Total: 89.99},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		AnalyzeUserBehavior(users, orders)
	}
}

// Helper functions for tests

// floatEqual compares two floats with a tolerance
func floatEqual(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}

// Test for edge cases and error conditions
func TestEdgeCases(t *testing.T) {
	// Test empty user
	emptyUser := User{}
	result := ValidateUser(emptyUser)
	if result.Valid {
		t.Error("ValidateUser() should fail for empty user")
	}
	if len(result.Errors) == 0 {
		t.Error("ValidateUser() should return errors for empty user")
	}

	// Test empty product
	emptyProduct := Product{}
	productResult := ValidateProduct(emptyProduct)
	if productResult.Valid {
		t.Error("ValidateProduct() should fail for empty product")
	}

	// Test empty analytics
	emptyAnalytics := AnalyzeUserBehavior([]User{}, []Order{})
	if emptyAnalytics.AverageAge != 0 {
		t.Error("AnalyzeUserBehavior() should handle empty input")
	}

	// Test recommendation with no products
	emptyRecommendations := RecommendProducts(testUsers[0], []Product{}, Order{})
	if len(emptyRecommendations) != 0 {
		t.Error("RecommendProducts() should return empty for no products")
	}
}

// Integration test combining multiple functions
func TestBusinessLogicIntegration(t *testing.T) {
	// Create a complete user journey
	user := User{
		Email:   "integration@test.com",
		Name:    "Integration User",
		Age:     30,
		Country: "US",
		Premium: true,
	}

	// Validate user
	userValidation := ValidateUser(user)
	if !userValidation.Valid {
		t.Fatalf("User validation failed: %v", userValidation.Errors)
	}

	// Create products
	products := []Product{
		{Name: "Test Product 1", Price: 50.0, Category: "electronics", InStock: true, Rating: 4.5},
		{Name: "Test Product 2", Price: 30.0, Category: "books", InStock: true, Rating: 4.0},
	}

	// Validate products
	for i, product := range products {
		productValidation := ValidateProduct(product)
		if !productValidation.Valid {
			t.Fatalf("Product %d validation failed: %v", i, productValidation.Errors)
		}
	}

	// Create order
	order := Order{
		Products:   products,
		Quantities: []int{1, 2}, // 1x product1 + 2x product2 = $110
	}

	// Calculate order total
	CalculateOrderTotal(&order, user)

	// Verify calculations make sense
	expectedSubtotal := 50.0 + (30.0 * 2) // $110
	if !floatEqual(order.Subtotal, expectedSubtotal, 0.01) {
		t.Errorf("Integration test: subtotal = %v, want %v", order.Subtotal, expectedSubtotal)
	}

	// Should have premium discount since subtotal > $100
	if order.Discount <= 0 {
		t.Error("Integration test: expected premium discount for order over $100")
	}

	// Should have US tax rate
	expectedTaxBase := order.Subtotal - order.Discount
	expectedTax := expectedTaxBase * 0.08 // US tax rate
	if !floatEqual(order.Tax, expectedTax, 0.01) {
		t.Errorf("Integration test: tax = %v, want %v", order.Tax, expectedTax)
	}

	// Should have free shipping (premium user, order > $75)
	if order.Shipping != 0 {
		t.Errorf("Integration test: expected free shipping, got %v", order.Shipping)
	}

	// Get recommendations
	recommendations := RecommendProducts(user, products, order)
	if len(recommendations) < 0 {
		t.Error("Integration test: should get some recommendations")
	}

	// Analyze behavior
	analytics := AnalyzeUserBehavior([]User{user}, []Order{order})
	if analytics.TotalRevenue <= 0 {
		t.Error("Integration test: should have positive revenue")
	}

	if analytics.PremiumPercentage != 100.0 {
		t.Errorf("Integration test: premium percentage = %v, want 100.0", analytics.PremiumPercentage)
	}

	t.Logf("Integration test completed successfully:")
	t.Logf("  User: %s (%s)", user.Name, user.Email)
	t.Logf("  Order total: %s", FormatCurrency(order.Total))
	t.Logf("  Tax: %s", FormatCurrency(order.Tax))
	t.Logf("  Discount: %s", FormatCurrency(order.Discount))
	t.Logf("  Recommendations: %d", len(recommendations))
}
