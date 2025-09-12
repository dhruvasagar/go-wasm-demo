# 🧪 Testing Guide: Ensuring Quality and Reliability

This project includes a comprehensive testing suite to ensure the shared business logic works correctly across both WebAssembly and server environments. The testing demonstrates proper practices for Go projects and validates the core narrative of identical business logic execution.

## 🎯 **Testing Philosophy**

### **Shared Logic, Shared Tests**
- **Same business functions** tested in both environments
- **Identical results** verified across WebAssembly and HTTP APIs  
- **Performance benchmarks** ensure WebAssembly delivers promised speedups
- **Integration tests** validate end-to-end functionality

### **Comprehensive Coverage**
- ✅ **Unit Tests** - Individual business logic functions
- ✅ **Integration Tests** - API endpoints and workflows
- ✅ **Performance Tests** - Benchmarks and regression detection
- ✅ **Algorithm Tests** - Mathematical correctness verification
- ✅ **Stress Tests** - Reliability under load
- ✅ **Build Tests** - Compilation verification

## 🚀 **Quick Start Testing**

### **Run All Tests**
```bash
# Run comprehensive test suite
./test.sh

# Run with benchmarks
./test.sh bench

# Run with coverage report
./test.sh coverage

# Run full suite (all options)
./test.sh full
```

### **Run Specific Test Categories**
```bash
# Unit tests only
go test -C src -v -run TestValidate

# Integration tests only  
go test -C src -v -run TestServer

# Performance benchmarks
go test -C src -bench=. -benchmem

# Algorithm correctness
go test -C src -v -run TestMatrix
```

## 📊 **Test Categories**

### **1. Business Logic Unit Tests** (`src/shared_models_test.go`)

#### **User Validation Tests**
```go
func TestValidateUser(t *testing.T) {
    // Tests email validation with regex
    // Age boundary conditions (13-120)
    // Country code validation
    // Name length requirements
    // Multiple error conditions
}
```

**What's Tested:**
- ✅ Valid user data passes validation
- ✅ Invalid email formats are rejected  
- ✅ Age limits (13-120) are enforced
- ✅ Country codes must be from approved list
- ✅ Names must be at least 2 characters
- ✅ Multiple validation errors are collected

#### **Product Validation Tests**
```go
func TestValidateProduct(t *testing.T) {
    // Tests product name length
    // Price validation (positive, under $10K)
    // Category validation against approved list
    // Rating validation (0-5 range)
}
```

**What's Tested:**
- ✅ Product names must be 3+ characters
- ✅ Prices must be positive and under $10,000
- ✅ Categories must be from approved list
- ✅ Ratings must be 0-5 range
- ✅ Stock status validation

#### **Order Calculation Tests**
```go
func TestCalculateOrderTotal(t *testing.T) {
    // Complex tax rate calculations by country
    // Premium user discounts
    // Shipping cost calculations
    // Multi-item order handling
}
```

**What's Tested:**
- ✅ Subtotal calculations with quantities
- ✅ Country-specific tax rates (8%-20%)
- ✅ Premium discounts (10-15% based on order size)
- ✅ Shipping rules (free over $100, premium exceptions)
- ✅ Total calculation accuracy

#### **Advanced Business Logic Tests**
```go
func TestRecommendProducts(t *testing.T) {
    // Recommendation algorithm based on:
    // - User demographics
    // - Purchase history  
    // - Category preferences
    // - Price sensitivity
}

func TestAnalyzeUserBehavior(t *testing.T) {
    // Analytics calculations:
    // - Average age computation
    // - Premium percentage
    // - Revenue analysis
    // - Top countries ranking
}
```

### **2. Integration Tests** (`src/src/integration_test.go`)

#### **HTTP API Endpoint Tests**
Tests that server endpoints use identical business logic:

```go
func TestServerAPIEndpoints(t *testing.T) {
    // POST /api/validate-user
    // POST /api/validate-product  
    // POST /api/calculate-order
    // POST /api/recommend-products
    // POST /api/analyze-behavior
}
```

**What's Verified:**
- ✅ API endpoints return correct HTTP status codes
- ✅ JSON request/response parsing works correctly
- ✅ Business logic produces expected results
- ✅ Error conditions are handled properly
- ✅ CORS headers are set for browser access

#### **Performance Benchmark Endpoints**
```go
func TestBenchmarkEndpoints(t *testing.T) {
    // GET /api/benchmark/matrix?size=100
    // GET /api/benchmark/mandelbrot?width=800&height=600
    // GET /api/benchmark/hash?count=10000
}
```

#### **Data Consistency Tests**
Critical tests that verify identical results:

```go
func TestDataConsistency(t *testing.T) {
    // Compare direct function calls vs API calls
    // Ensure identical validation results
    // Verify order calculations match exactly
    // Check floating-point precision
}
```

### **3. Algorithm Correctness Tests** (`src/src/benchmarks_test.go`)

#### **Matrix Multiplication Verification**
```go
func TestMatrixMultiplicationLogic(t *testing.T) {
    // Test known 2x2 matrix multiplication:
    // [[1,2],[3,4]] × [[5,6],[7,8]] = [[19,22],[43,50]]
}
```

#### **Mandelbrot Set Mathematical Verification**
```go
func TestMandelbrotLogic(t *testing.T) {
    // Test known points:
    // (0,0) should be in the set (reaches max iterations)
    // (2,0) should diverge quickly (< 10 iterations)  
    // (-1,0) should be in the set
    // (-2.5,0) should diverge immediately
}
```

#### **Hash Function Consistency**
```go
func TestHashingConsistency(t *testing.T) {
    // Verify same input produces same hash
    // Test multiple runs for consistency
    // Validate hash function determinism
}
```

### **4. Performance Benchmarks**

#### **Business Logic Performance**
```bash
BenchmarkValidateUser-8         3000000    450 ns/op    96 B/op    2 allocs/op
BenchmarkCalculateOrderTotal-8   500000   2847 ns/op   640 B/op    8 allocs/op
BenchmarkRecommendProducts-8     100000  12456 ns/op  2048 B/op   45 allocs/op
```

#### **Algorithm Performance**
```bash
BenchmarkMatrixMultiplication100x100-8    50    21.2 ms/op    800000 B/op    3 allocs/op
BenchmarkMandelbrot800x600-8               10   120.5 ms/op    1920000 B/op    1 allocs/op  
BenchmarkHashing10000-8                  1000     1.8 ms/op        0 B/op    0 allocs/op
```

**Performance Metrics Reported:**
- **Operations per second** (MOps/sec for matrix, MPixels/sec for Mandelbrot)
- **Memory allocations** (B/op and allocs/op)
- **Execution time** (ns/op, ms/op)
- **Throughput rates** for comparison with JavaScript

### **5. Stress and Reliability Tests**

#### **Concurrent Safety Testing**
```go
func TestConcurrentSafety(t *testing.T) {
    // Run 10 goroutines × 100 iterations
    // Verify business logic is thread-safe
    // Check for race conditions
}
```

#### **Stress Testing**
```go
func TestStressTest(t *testing.T) {
    // 1000+ iterations of matrix multiplication
    // Verify results remain consistent
    // Check for memory leaks or performance degradation
}
```

#### **Performance Regression Detection**
```go
func TestPerformanceRegression(t *testing.T) {
    // Compare 50x50 vs 100x100 matrix performance
    // Verify expected O(n³) scaling
    // Alert on unexpected performance changes
}
```

### **6. Build and Quality Tests**

#### **Build Verification**
```bash
# WebAssembly compilation
GOOS=js GOARCH=wasm go build -o main.wasm src/main_wasm.go src/shared_models.go src/benchmarks_*.go src/mandelbrot*.go

# Server compilation  
go build -o server src/main_server.go src/shared_models.go

# Test compilation
go test -C src -c
```

#### **Code Quality Checks**
```bash
# Format checking
gofmt -l .

# Vet analysis  
go vet ./...

# Test coverage
go test -C src -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📈 **Test Coverage Goals**

### **Current Coverage Targets**
- **Business Logic Functions**: 95%+ coverage
- **API Endpoints**: 90%+ coverage  
- **Error Handling**: 85%+ coverage
- **Algorithm Implementations**: 100% coverage

### **Coverage Report Example**
```
src/shared_models.go:45:    ValidateUser           100.0%
src/shared_models.go:78:    ValidateProduct        100.0%
src/shared_models.go:112:   CalculateOrderTotal     95.2%
src/shared_models.go:145:   GetTaxRate             100.0%
src/shared_models.go:167:   CalculateShipping       92.3%
src/shared_models.go:203:   RecommendProducts       88.7%
src/shared_models.go:267:   AnalyzeUserBehavior     91.4%
Total Coverage:         94.1%
```

## 🚀 **Running Tests in CI/CD**

### **GitHub Actions Example**
```yaml
name: Test Suite
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: '1.19'
    - name: Run Tests
      run: ./test.sh full
    - name: Upload Coverage
      uses: codecov/codecov-action@v3
```

### **Docker Testing Environment**
```dockerfile
FROM golang:1.19-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN ./test.sh full
```

## 🔧 **Test Development Guidelines**

### **Writing New Tests**

#### **1. Follow Naming Conventions**
```go
func TestFunctionName(t *testing.T)           // Unit tests
func BenchmarkFunctionName(b *testing.B)     // Benchmarks  
func ExampleFunctionName()                   // Examples
```

#### **2. Use Table-Driven Tests**
```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
        wantErr  bool
    }{
        {"valid input", validInput, expectedOutput, false},
        {"invalid input", invalidInput, nil, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionUnderTest(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(got, tt.expected) {
                t.Errorf("got %v, want %v", got, tt.expected)
            }
        })
    }
}
```

#### **3. Test Both Success and Failure Cases**
```go
func TestBusinessLogic(t *testing.T) {
    // Test successful path
    t.Run("success case", func(t *testing.T) { /* ... */ })
    
    // Test error conditions
    t.Run("validation errors", func(t *testing.T) { /* ... */ })
    t.Run("edge cases", func(t *testing.T) { /* ... */ })
    t.Run("boundary conditions", func(t *testing.T) { /* ... */ })
}
```

### **Benchmark Writing Guidelines**

#### **Proper Benchmark Structure**
```go
func BenchmarkFunction(b *testing.B) {
    // Setup (not timed)
    input := setupBenchmarkData()
    
    b.ResetTimer()      // Start timing here
    b.ReportAllocs()    // Report memory allocations
    
    for i := 0; i < b.N; i++ {
        FunctionUnderTest(input)
    }
    
    // Report custom metrics
    ops := int64(b.N) * operationsPerIteration
    b.ReportMetric(float64(ops)/b.Elapsed().Seconds(), "ops/sec")
}
```

## 📊 **Interpreting Test Results**

### **Understanding Performance Benchmarks**
```bash
BenchmarkMatrixMultiplication100x100-8    50    21.2 ms/op    800000 B/op    3 allocs/op
│                                          │     │            │             │
│                                          │     │            │             └─ Allocations per operation
│                                          │     │            └─ Bytes allocated per operation  
│                                          │     └─ Nanoseconds per operation
│                                          └─ Number of iterations run
```

### **Performance Comparison Guidelines**
- **WebAssembly vs JavaScript**: Expected 2-10x speedup
- **Memory usage**: WASM should use less memory
- **Allocation rate**: Lower allocation rates are better
- **Throughput**: Higher ops/sec indicates better performance

## 🎯 **Continuous Quality Assurance**

### **Pre-commit Hooks**
```bash
#!/bin/sh
# Run tests before every commit
./test.sh short
if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi
```

### **Performance Monitoring**
- Track benchmark results over time
- Alert on performance regressions > 20%
- Monitor memory usage trends
- Validate WebAssembly speedup claims

### **Test Maintenance**
- Review test coverage monthly
- Update tests when business logic changes
- Add regression tests for bugs found
- Keep performance benchmarks current

---

## 🌟 **Testing Success Metrics**

### **Quality Gates**
- ✅ **All unit tests pass** - Business logic correctness
- ✅ **API tests pass** - Integration functionality  
- ✅ **Performance benchmarks meet targets** - WebAssembly speedup validated
- ✅ **Code coverage > 90%** - Comprehensive testing
- ✅ **No race conditions** - Concurrent safety verified
- ✅ **Builds succeed** - Both WASM and server compile

### **Confidence Indicators**
- **Identical results** between WebAssembly and server
- **Performance improvements** documented and measured
- **Error handling** comprehensive and tested
- **Edge cases** identified and covered
- **Regression protection** in place

---

**🎉 A comprehensive test suite ensures your shared business logic works correctly across all environments, giving you confidence to demonstrate the full power of WebAssembly in Go!**