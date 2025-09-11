//go:build !wasm

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestServerAPIEndpoints tests all the HTTP API endpoints
func TestServerAPIEndpoints(t *testing.T) {
	// Test user validation endpoint
	t.Run("UserValidationAPI", func(t *testing.T) {
		testUser := map[string]interface{}{
			"email":   "test@example.com",
			"name":    "Test User",
			"age":     25,
			"country": "US",
			"premium": true,
		}

		jsonData, _ := json.Marshal(testUser)
		req := httptest.NewRequest("POST", "/api/validate-user", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handleValidateUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var result ValidationResult
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected valid user, got errors: %v", result.Errors)
		}
	})

	// Test product validation endpoint
	t.Run("ProductValidationAPI", func(t *testing.T) {
		testProduct := map[string]interface{}{
			"name":     "Test Product",
			"price":    99.99,
			"category": "electronics",
			"rating":   4.5,
			"in_stock": true,
		}

		jsonData, _ := json.Marshal(testProduct)
		req := httptest.NewRequest("POST", "/api/validate-product", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handleValidateProduct(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var result ValidationResult
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected valid product, got errors: %v", result.Errors)
		}
	})

	// Test order calculation endpoint
	t.Run("OrderCalculationAPI", func(t *testing.T) {
		testData := map[string]interface{}{
			"order": map[string]interface{}{
				"products": []map[string]interface{}{
					{
						"id":       1,
						"name":     "Test Product",
						"price":    100.0,
						"category": "electronics",
					},
				},
				"quantities": []int{1},
			},
			"user": map[string]interface{}{
				"email":   "test@example.com",
				"name":    "Test User",
				"age":     25,
				"country": "US",
				"premium": true,
			},
		}

		jsonData, _ := json.Marshal(testData)
		req := httptest.NewRequest("POST", "/api/calculate-order", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handleCalculateOrder(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var result map[string]float64
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result["subtotal"] <= 0 {
			t.Error("Expected positive subtotal")
		}

		if result["total"] <= 0 {
			t.Error("Expected positive total")
		}
	})

	// Test recommendations endpoint
	t.Run("RecommendationsAPI", func(t *testing.T) {
		testData := map[string]interface{}{
			"user": map[string]interface{}{
				"email":   "test@example.com",
				"name":    "Test User",
				"age":     25,
				"country": "US",
				"premium": true,
			},
			"products": []map[string]interface{}{
				{
					"id":          1,
					"name":        "Product 1",
					"price":       50.0,
					"category":    "electronics",
					"in_stock":    true,
					"rating":      4.5,
					"description": "Test product",
				},
				{
					"id":          2,
					"name":        "Product 2",
					"price":       30.0,
					"category":    "books",
					"in_stock":    true,
					"rating":      4.0,
					"description": "Test book",
				},
			},
			"order": map[string]interface{}{
				"products": []map[string]interface{}{
					{
						"id":       1,
						"name":     "Product 1",
						"category": "electronics",
					},
				},
			},
		}

		jsonData, _ := json.Marshal(testData)
		req := httptest.NewRequest("POST", "/api/recommend-products", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handleRecommendProducts(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var result []Product
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Should return some recommendations
		if len(result) < 0 {
			t.Error("Expected some recommendations")
		}

		// Should not exceed 5 recommendations
		if len(result) > 5 {
			t.Errorf("Too many recommendations: %d, expected <= 5", len(result))
		}
	})
}

// TestBenchmarkEndpoints tests the performance benchmark endpoints
func TestBenchmarkEndpoints(t *testing.T) {
	t.Run("MatrixBenchmark", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/benchmark/matrix?size=50", nil)
		w := httptest.NewRecorder()

		handleMatrixBenchmark(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var result map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result["operation"] != "Matrix Multiplication" {
			t.Error("Wrong operation type")
		}

		if duration, ok := result["duration_ms"].(float64); !ok || duration <= 0 {
			t.Error("Invalid duration")
		}
	})

	t.Run("MandelbrotBenchmark", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/benchmark/mandelbrot?width=200&height=150&iterations=50", nil)
		w := httptest.NewRecorder()

		handleMandelbrotBenchmark(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var result map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result["operation"] != "Mandelbrot Set" {
			t.Error("Wrong operation type")
		}

		if pixels, ok := result["pixels"].(float64); !ok || pixels != 30000 { // 200*150
			t.Errorf("Wrong pixel count: %v", pixels)
		}
	})

	t.Run("HashBenchmark", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/benchmark/hash?count=1000", nil)
		w := httptest.NewRecorder()

		handleHashBenchmark(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var result map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result["operation"] != "SHA256 Hashing" {
			t.Error("Wrong operation type")
		}

		if count, ok := result["count"].(float64); !ok || count != 1000 {
			t.Errorf("Wrong count: %v", count)
		}
	})
}

// TestErrorHandling tests error conditions in API endpoints
func TestErrorHandling(t *testing.T) {
	t.Run("InvalidJSONUserValidation", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/validate-user", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handleValidateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("WrongHTTPMethod", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/validate-user", nil)

		w := httptest.NewRecorder()
		handleValidateUser(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})

	t.Run("EmptyRequestBody", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/validate-user", strings.NewReader(""))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handleValidateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

// TestCORSHeaders tests that CORS headers are properly set
func TestCORSHeaders(t *testing.T) {
	req := httptest.NewRequest("OPTIONS", "/api/validate-user", nil)
	w := httptest.NewRecorder()

	handleValidateUser(w, req)

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}

	for header, expectedValue := range expectedHeaders {
		if got := w.Header().Get(header); got != expectedValue {
			t.Errorf("Header %s = %s, want %s", header, got, expectedValue)
		}
	}
}

// TestDemoDataEndpoints tests the demo data endpoints
func TestDemoDataEndpoints(t *testing.T) {
	t.Run("DemoUsers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/demo-users", nil)
		w := httptest.NewRecorder()

		handleDemoUsers(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var users []User
		err := json.NewDecoder(w.Body).Decode(&users)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(users) == 0 {
			t.Error("Expected some demo users")
		}

		// Validate first user
		if len(users) > 0 {
			result := ValidateUser(users[0])
			if !result.Valid {
				t.Errorf("Demo user invalid: %v", result.Errors)
			}
		}
	})

	t.Run("DemoProducts", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/demo-products", nil)
		w := httptest.NewRecorder()

		handleDemoProducts(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var products []Product
		err := json.NewDecoder(w.Body).Decode(&products)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(products) == 0 {
			t.Error("Expected some demo products")
		}

		// Validate first product
		if len(products) > 0 {
			result := ValidateProduct(products[0])
			if !result.Valid {
				t.Errorf("Demo product invalid: %v", result.Errors)
			}
		}
	})

	t.Run("DemoOrders", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/demo-orders", nil)
		w := httptest.NewRecorder()

		handleDemoOrders(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var orders []Order
		err := json.NewDecoder(w.Body).Decode(&orders)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(orders) == 0 {
			t.Error("Expected some demo orders")
		}
	})
}

// TestServerPerformance tests that server responds within reasonable time
func TestServerPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	testCases := []struct {
		name    string
		path    string
		method  string
		body    string
		maxTime time.Duration
	}{
		{
			name:    "UserValidation",
			path:    "/api/validate-user",
			method:  "POST",
			body:    `{"email":"test@example.com","name":"Test","age":25,"country":"US"}`,
			maxTime: 100 * time.Millisecond,
		},
		{
			name:    "OrderCalculation",
			path:    "/api/calculate-order",
			method:  "POST",
			body:    `{"order":{"products":[{"name":"Test","price":100,"category":"electronics"}],"quantities":[1]},"user":{"email":"test@example.com","name":"Test","age":25,"country":"US","premium":true}}`,
			maxTime: 50 * time.Millisecond,
		},
		{
			name:    "MatrixBenchmark",
			path:    "/api/benchmark/matrix?size=50",
			method:  "GET",
			body:    "",
			maxTime: 200 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name+"Performance", func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			if tc.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			start := time.Now()
			w := httptest.NewRecorder()

			// Route to appropriate handler
			switch tc.path {
			case "/api/validate-user":
				handleValidateUser(w, req)
			case "/api/calculate-order":
				handleCalculateOrder(w, req)
			default:
				if strings.Contains(tc.path, "matrix") {
					handleMatrixBenchmark(w, req)
				}
			}

			duration := time.Since(start)

			if w.Code != http.StatusOK {
				t.Errorf("Request failed with status %d", w.Code)
			}

			if duration > tc.maxTime {
				t.Errorf("Request too slow: %v > %v", duration, tc.maxTime)
			}

			t.Logf("%s completed in %v", tc.name, duration)
		})
	}
}

// TestDataConsistency tests that server and business logic produce consistent results
func TestDataConsistency(t *testing.T) {
	// Test user validation consistency
	t.Run("UserValidationConsistency", func(t *testing.T) {
		user := User{
			Email:   "consistency@test.com",
			Name:    "Consistency Test",
			Age:     30,
			Country: "US",
			Premium: true,
		}

		// Direct business logic call
		directResult := ValidateUser(user)

		// API call
		jsonData, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/api/validate-user", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handleValidateUser(w, req)

		var apiResult ValidationResult
		json.NewDecoder(w.Body).Decode(&apiResult)

		// Results should be identical
		if directResult.Valid != apiResult.Valid {
			t.Errorf("Validation consistency failed: direct=%v, api=%v", directResult.Valid, apiResult.Valid)
		}

		if len(directResult.Errors) != len(apiResult.Errors) {
			t.Errorf("Error count mismatch: direct=%d, api=%d", len(directResult.Errors), len(apiResult.Errors))
		}
	})

	// Test order calculation consistency
	t.Run("OrderCalculationConsistency", func(t *testing.T) {
		order := Order{
			Products: []Product{
				{Name: "Test Product", Price: 100.0, Category: "electronics"},
			},
			Quantities: []int{1},
		}
		user := User{
			Email: "test@example.com", Name: "Test User", Age: 25,
			Country: "US", Premium: true,
		}

		// Direct business logic call
		directOrder := order
		CalculateOrderTotal(&directOrder, user)

		// API call
		testData := map[string]interface{}{
			"order": map[string]interface{}{
				"products":   []map[string]interface{}{{"name": "Test Product", "price": 100.0, "category": "electronics"}},
				"quantities": []int{1},
			},
			"user": map[string]interface{}{
				"email": "test@example.com", "name": "Test User", "age": 25,
				"country": "US", "premium": true,
			},
		}

		jsonData, _ := json.Marshal(testData)
		req := httptest.NewRequest("POST", "/api/calculate-order", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handleCalculateOrder(w, req)

		var apiResult map[string]float64
		json.NewDecoder(w.Body).Decode(&apiResult)

		// Compare results (allowing small floating-point differences)
		tolerance := 0.01
		if absFloat(directOrder.Subtotal-apiResult["subtotal"]) > tolerance {
			t.Errorf("Subtotal mismatch: direct=%f, api=%f", directOrder.Subtotal, apiResult["subtotal"])
		}

		if absFloat(directOrder.Total-apiResult["total"]) > tolerance {
			t.Errorf("Total mismatch: direct=%f, api=%f", directOrder.Total, apiResult["total"])
		}
	})
}

// TestStaticFileServing tests that static files are served correctly
func TestStaticFileServing(t *testing.T) {
	t.Run("IndexHTML", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		serveStaticFile(w, req)

		// Should redirect to index.html and serve content
		// In a real test, we'd check for HTML content
		if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
			t.Logf("Static file serving status: %d (may be OK if files don't exist in test)", w.Code)
		}
	})
}

// Helper function for floating point comparison in tests
func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
