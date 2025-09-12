//go:build !wasm

package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
	}

	// Setup HTTP server with timeouts
	server := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Serve static files
	http.HandleFunc("/", serveStaticFile)

	// API endpoints using shared business logic
	http.HandleFunc("/api/validate-user", handleValidateUser)
	http.HandleFunc("/api/validate-product", handleValidateProduct)
	http.HandleFunc("/api/calculate-order", handleCalculateOrder)
	http.HandleFunc("/api/recommend-products", handleRecommendProducts)
	http.HandleFunc("/api/analyze-behavior", handleAnalyzeBehavior)

	// Demo data endpoints
	http.HandleFunc("/api/demo-users", handleDemoUsers)
	http.HandleFunc("/api/demo-products", handleDemoProducts)
	http.HandleFunc("/api/demo-orders", handleDemoOrders)

	// Performance benchmark endpoints
	http.HandleFunc("/api/benchmark/matrix", handleMatrixBenchmark)
	http.HandleFunc("/api/benchmark/mandelbrot", handleMandelbrotBenchmark)
	http.HandleFunc("/api/benchmark/hash", handleHashBenchmark)

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("🚀 Server starting on http://localhost:%s\n", port)
		fmt.Println("📊 Visit /server.html for server-side demo")
		fmt.Println("🌐 Visit / for WebAssembly demo")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-c
	fmt.Println("\n🛑 Shutting down server gracefully...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		fmt.Println("✅ Server shut down successfully")
	}
}

func serveStaticFile(w http.ResponseWriter, r *http.Request) {
	// Handle root path
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "./index.html")
		return
	}

	// Serve other static files
	http.ServeFile(w, r, "."+r.URL.Path)
}

// API endpoint for user validation using shared business logic
func handleValidateUser(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Add content length check to prevent memory issues
	if r.ContentLength > 1024*1024 { // 1MB limit
		http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Use shared business logic - identical to WebAssembly version
	result := ValidateUser(user)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// API endpoint for product validation using shared business logic
func handleValidateProduct(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Use shared business logic - identical to WebAssembly version
	result := ValidateProduct(product)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// API endpoint for order calculation using shared business logic
func handleCalculateOrder(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Add content length check to prevent memory issues
	if r.ContentLength > 1024*1024 { // 1MB limit
		http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	var requestData struct {
		Order Order `json:"order"`
		User  User  `json:"user"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate order has products
	if len(requestData.Order.Products) == 0 {
		http.Error(w, "Order must contain at least one product", http.StatusBadRequest)
		return
	}

	// Validate quantities match products
	if len(requestData.Order.Products) != len(requestData.Order.Quantities) {
		http.Error(w, "Product and quantity arrays must be the same length", http.StatusBadRequest)
		return
	}

	// Validate user data is present
	if requestData.User.Country == "" {
		http.Error(w, "User country is required", http.StatusBadRequest)
		return
	}

	// Use shared business logic - identical to WebAssembly version
	CalculateOrderTotal(&requestData.Order, requestData.User)

	response := map[string]interface{}{
		"subtotal": requestData.Order.Subtotal,
		"tax":      requestData.Order.Tax,
		"shipping": requestData.Order.Shipping,
		"discount": requestData.Order.Discount,
		"total":    requestData.Order.Total,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// API endpoint for product recommendations using shared business logic
func handleRecommendProducts(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		User     User      `json:"user"`
		Products []Product `json:"products"`
		Order    Order     `json:"order"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Use shared business logic - identical to WebAssembly version
	recommendations := RecommendProducts(requestData.User, requestData.Products, requestData.Order)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}

// API endpoint for user behavior analysis using shared business logic
func handleAnalyzeBehavior(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Users  []User  `json:"users"`
		Orders []Order `json:"orders"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Use shared business logic - identical to WebAssembly version
	analytics := AnalyzeUserBehavior(requestData.Users, requestData.Orders)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// Demo data endpoints
func handleDemoUsers(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	users := generateDemoUsers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func handleDemoProducts(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	products := generateDemoProducts()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func handleDemoOrders(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	orders := generateDemoOrders()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// Performance benchmark endpoints
func handleMatrixBenchmark(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	size := 100
	if sizeParam := r.URL.Query().Get("size"); sizeParam != "" {
		if parsedSize, err := strconv.Atoi(sizeParam); err == nil {
			size = parsedSize
		}
	}

	result := benchmarkMatrixMultiply(size)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleMandelbrotBenchmark(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	width, height := 400, 300
	iterations := 100

	if w := r.URL.Query().Get("width"); w != "" {
		if parsed, err := strconv.Atoi(w); err == nil {
			width = parsed
		}
	}
	if h := r.URL.Query().Get("height"); h != "" {
		if parsed, err := strconv.Atoi(h); err == nil {
			height = parsed
		}
	}
	if iter := r.URL.Query().Get("iterations"); iter != "" {
		if parsed, err := strconv.Atoi(iter); err == nil {
			iterations = parsed
		}
	}

	result := benchmarkMandelbrot(width, height, iterations)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleHashBenchmark(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	count := 10000
	if countParam := r.URL.Query().Get("count"); countParam != "" {
		if parsed, err := strconv.Atoi(countParam); err == nil {
			count = parsed
		}
	}

	result := benchmarkSHA256(count)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func enableCORS(w http.ResponseWriter) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
}

// Demo data generators
func generateDemoUsers() []User {
	return []User{
		{ID: 1, Email: "john.doe@example.com", Name: "John Doe", Age: 28, Country: "US", Premium: true, JoinDate: "2023-01-15"},
		{ID: 2, Email: "jane.smith@example.com", Name: "Jane Smith", Age: 34, Country: "CA", Premium: false, JoinDate: "2023-02-20"},
		{ID: 3, Email: "alice.johnson@example.com", Name: "Alice Johnson", Age: 22, Country: "UK", Premium: true, JoinDate: "2023-03-10"},
		{ID: 4, Email: "bob.wilson@example.com", Name: "Bob Wilson", Age: 45, Country: "AU", Premium: false, JoinDate: "2023-01-30"},
		{ID: 5, Email: "carol.brown@example.com", Name: "Carol Brown", Age: 31, Country: "DE", Premium: true, JoinDate: "2023-04-05"},
	}
}

// Server-side benchmark implementations using the same algorithms
func benchmarkMatrixMultiply(size int) map[string]interface{} {
	start := time.Now()

	// Create test matrices
	matrixA := make([]float64, size*size)
	matrixB := make([]float64, size*size)
	result := make([]float64, size*size)

	// Initialize with test data
	for i := 0; i < size*size; i++ {
		matrixA[i] = float64(i % 10)
		matrixB[i] = float64((i * 2) % 10)
	}

	// Matrix multiplication
	for i := 0; i < size; i++ {
		for k := 0; k < size; k++ {
			aik := matrixA[i*size+k]
			for j := 0; j < size; j++ {
				result[i*size+j] += aik * matrixB[k*size+j]
			}
		}
	}

	duration := time.Since(start)

	return map[string]interface{}{
		"operation":   "Matrix Multiplication",
		"size":        fmt.Sprintf("%dx%d", size, size),
		"duration_ms": float64(duration.Nanoseconds()) / 1000000,
		"operations":  size * size * size,
		"result_hash": int(result[0] + result[size-1] + result[len(result)-1]),
	}
}

func benchmarkMandelbrot(width, height, iterations int) map[string]interface{} {
	start := time.Now()

	xmin, xmax := -2.0, 1.0
	ymin, ymax := -1.5, 1.5

	result := make([]int, width*height)
	dx := (xmax - xmin) / float64(width)
	dy := (ymax - ymin) / float64(height)

	idx := 0
	for py := 0; py < height; py++ {
		cy := ymin + float64(py)*dy
		for px := 0; px < width; px++ {
			cx := xmin + float64(px)*dx

			zx, zy := 0.0, 0.0
			iter := 0

			for iter < iterations {
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

	duration := time.Since(start)

	return map[string]interface{}{
		"operation":   "Mandelbrot Set",
		"size":        fmt.Sprintf("%dx%d", width, height),
		"iterations":  iterations,
		"duration_ms": float64(duration.Nanoseconds()) / 1000000,
		"pixels":      width * height,
		"result_hash": result[0] + result[len(result)/2] + result[len(result)-1],
	}
}

func benchmarkSHA256(count int) map[string]interface{} {
	start := time.Now()

	data := "WebAssembly performance test data for hashing benchmark"
	hash := 0

	for i := 0; i < count; i++ {
		hasher := sha256.New()
		hasher.Write([]byte(fmt.Sprintf("%s-%d", data, i)))
		sum := hasher.Sum(nil)
		hash += int(sum[0])
	}

	duration := time.Since(start)

	return map[string]interface{}{
		"operation":   "SHA256 Hashing",
		"count":       count,
		"duration_ms": float64(duration.Nanoseconds()) / 1000000,
		"result_hash": hash,
	}
}

func generateDemoProducts() []Product {
	return []Product{
		{ID: 1, Name: "Wireless Headphones", Price: 99.99, Category: "electronics", InStock: true, Rating: 4.5, Description: "High-quality wireless headphones with noise cancellation"},
		{ID: 2, Name: "Cotton T-Shirt", Price: 24.99, Category: "clothing", InStock: true, Rating: 4.2, Description: "Comfortable 100% cotton t-shirt"},
		{ID: 3, Name: "Programming Book", Price: 49.99, Category: "books", InStock: true, Rating: 4.8, Description: "Learn advanced programming techniques"},
		{ID: 4, Name: "Coffee Mug", Price: 12.99, Category: "home", InStock: true, Rating: 4.0, Description: "Ceramic coffee mug with handle"},
		{ID: 5, Name: "Running Shoes", Price: 129.99, Category: "sports", InStock: true, Rating: 4.6, Description: "Lightweight running shoes for athletes"},
		{ID: 6, Name: "Smartphone", Price: 699.99, Category: "electronics", InStock: false, Rating: 4.7, Description: "Latest smartphone with advanced features"},
		{ID: 7, Name: "Jeans", Price: 79.99, Category: "clothing", InStock: true, Rating: 4.3, Description: "Classic blue jeans"},
		{ID: 8, Name: "Cookbook", Price: 29.99, Category: "books", InStock: true, Rating: 4.4, Description: "Delicious recipes for home cooking"},
	}
}

func generateDemoOrders() []Order {
	products := generateDemoProducts()
	return []Order{
		{
			ID:         1,
			UserID:     1,
			Products:   products[0:2],
			Quantities: []int{1, 2},
			Subtotal:   149.97,
			Tax:        12.00,
			Shipping:   0.00,
			Total:      161.97,
			Discount:   0.00,
			OrderDate:  "2023-05-01",
			Status:     "delivered",
		},
		{
			ID:         2,
			UserID:     2,
			Products:   products[2:4],
			Quantities: []int{1, 1},
			Subtotal:   62.98,
			Tax:        8.19,
			Shipping:   12.99,
			Total:      84.16,
			Discount:   0.00,
			OrderDate:  "2023-05-03",
			Status:     "shipped",
		},
	}
}
