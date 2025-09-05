#!/bin/bash

echo "🧪 WebAssembly in Go: Comprehensive Testing Suite"
echo "================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track test results
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run test and track results
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    echo -e "${BLUE}🔄 Running $test_name...${NC}"
    
    if eval "$test_command"; then
        echo -e "${GREEN}✅ $test_name passed${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}❌ $test_name failed${NC}"
        ((TESTS_FAILED++))
    fi
    echo ""
}

# Function to run benchmarks
run_benchmark() {
    local bench_name="$1"
    local bench_command="$2"
    
    echo -e "${YELLOW}⚡ Running $bench_name...${NC}"
    eval "$bench_command"
    echo ""
}

echo -e "${BLUE}📋 Test Plan:${NC}"
echo "1. Unit Tests - Business Logic"
echo "2. Integration Tests - API Endpoints"
echo "3. Performance Benchmarks"
echo "4. Build Verification"
echo "5. Code Quality Checks"
echo ""

# 1. Unit Tests
echo -e "${YELLOW}═══════════════════════════════════════${NC}"
echo -e "${YELLOW}📝 UNIT TESTS - BUSINESS LOGIC${NC}"
echo -e "${YELLOW}═══════════════════════════════════════${NC}"

run_test "User Validation Logic" "go test -v -run TestValidateUser"
run_test "Product Validation Logic" "go test -v -run TestValidateProduct" 
run_test "Order Calculation Logic" "go test -v -run TestCalculateOrderTotal"
run_test "Tax Rate Calculation" "go test -v -run TestGetTaxRate"
run_test "Shipping Calculation" "go test -v -run TestCalculateShipping"
run_test "Recommendation Algorithm" "go test -v -run TestRecommendProducts"
run_test "User Analytics" "go test -v -run TestAnalyzeUserBehavior"
run_test "JSON Serialization" "go test -v -run TestJSONSerialization"
run_test "Utility Functions" "go test -v -run TestUtilityFunctions"
run_test "Edge Cases" "go test -v -run TestEdgeCases"

# 2. Integration Tests
echo -e "${YELLOW}═══════════════════════════════════════${NC}"
echo -e "${YELLOW}🔗 INTEGRATION TESTS - API ENDPOINTS${NC}"
echo -e "${YELLOW}═══════════════════════════════════════${NC}"

run_test "Server API Endpoints" "go test -v -run TestServerAPIEndpoints"
run_test "Benchmark Endpoints" "go test -v -run TestBenchmarkEndpoints"
run_test "Error Handling" "go test -v -run TestErrorHandling"
run_test "CORS Headers" "go test -v -run TestCORSHeaders"
run_test "Demo Data Endpoints" "go test -v -run TestDemoDataEndpoints"
run_test "Data Consistency" "go test -v -run TestDataConsistency"

# 3. Algorithm Tests
echo -e "${YELLOW}═══════════════════════════════════════${NC}"
echo -e "${YELLOW}🔬 ALGORITHM CORRECTNESS TESTS${NC}"
echo -e "${YELLOW}═══════════════════════════════════════${NC}"

run_test "Matrix Multiplication Logic" "go test -v -run TestMatrixMultiplicationLogic"
run_test "Mandelbrot Set Logic" "go test -v -run TestMandelbrotLogic"
run_test "Hash Function Consistency" "go test -v -run TestHashingConsistency"
run_test "Memory Allocation" "go test -v -run TestMemoryAllocation"

# 4. Performance Tests (if not in short mode)
if [[ "$1" != "short" ]]; then
    echo -e "${YELLOW}═══════════════════════════════════════${NC}"
    echo -e "${YELLOW}🚀 PERFORMANCE TESTS${NC}"
    echo -e "${YELLOW}═══════════════════════════════════════${NC}"
    
    run_test "Server Performance" "go test -v -run TestServerPerformance"
    run_test "Stress Testing" "go test -v -run TestStressTest"
    run_test "Performance Regression" "go test -v -run TestPerformanceRegression"
    run_test "Concurrent Safety" "go test -v -run TestConcurrentSafety"
fi

# 5. Business Logic Integration
echo -e "${YELLOW}═══════════════════════════════════════${NC}"
echo -e "${YELLOW}🎯 BUSINESS LOGIC INTEGRATION${NC}"
echo -e "${YELLOW}═══════════════════════════════════════${NC}"

run_test "End-to-End Business Logic" "go test -v -run TestBusinessLogicIntegration"

# 6. Build Verification
echo -e "${YELLOW}═══════════════════════════════════════${NC}"
echo -e "${YELLOW}🔧 BUILD VERIFICATION${NC}"
echo -e "${YELLOW}═══════════════════════════════════════${NC}"

run_test "WebAssembly Build" "GOOS=js GOARCH=wasm go build -o test_main.wasm main_wasm.go shared_models.go benchmarks.go mandelbrot.go"
run_test "Server Build" "go build -o test_server main_server.go shared_models.go benchmarks.go"
run_test "Test Compilation" "go test -c -o test_binary"

# Clean up build artifacts
rm -f test_main.wasm test_server test_binary 2>/dev/null

# 7. Performance Benchmarks
if [[ "$1" == "bench" ]] || [[ "$1" == "full" ]]; then
    echo -e "${YELLOW}═══════════════════════════════════════${NC}"
    echo -e "${YELLOW}⚡ PERFORMANCE BENCHMARKS${NC}"
    echo -e "${YELLOW}═══════════════════════════════════════${NC}"
    
    run_benchmark "Business Logic Benchmarks" "go test -bench=BenchmarkValidate -benchmem"
    run_benchmark "Algorithm Benchmarks" "go test -bench=BenchmarkMatrix -benchmem -benchtime=3s"
    run_benchmark "Mandelbrot Benchmarks" "go test -bench=BenchmarkMandelbrot -benchmem -benchtime=2s"
    run_benchmark "Hash Benchmarks" "go test -bench=BenchmarkHashing -benchmem"
    run_benchmark "Server Function Benchmarks" "go test -bench=BenchmarkServer -benchmem"
fi

# 8. Code Quality Checks
echo -e "${YELLOW}═══════════════════════════════════════${NC}"
echo -e "${YELLOW}✨ CODE QUALITY CHECKS${NC}"
echo -e "${YELLOW}═══════════════════════════════════════${NC}"

# Check for gofmt
if command -v gofmt >/dev/null 2>&1; then
    if gofmt -l . | grep -q .; then
        echo -e "${RED}❌ Code formatting issues found:${NC}"
        gofmt -l .
        ((TESTS_FAILED++))
    else
        echo -e "${GREEN}✅ Code formatting check passed${NC}"
        ((TESTS_PASSED++))
    fi
else
    echo -e "${YELLOW}⚠️  gofmt not available, skipping format check${NC}"
fi

# Check for go vet
if command -v go >/dev/null 2>&1; then
    if go vet ./...; then
        echo -e "${GREEN}✅ Go vet check passed${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}❌ Go vet found issues${NC}"
        ((TESTS_FAILED++))
    fi
fi

# Test Coverage Report
if [[ "$1" == "coverage" ]] || [[ "$1" == "full" ]]; then
    echo -e "${YELLOW}═══════════════════════════════════════${NC}"
    echo -e "${YELLOW}📊 TEST COVERAGE REPORT${NC}"
    echo -e "${YELLOW}═══════════════════════════════════════${NC}"
    
    echo -e "${BLUE}Generating coverage report...${NC}"
    go test -coverprofile=coverage.out -covermode=count
    
    if command -v go >/dev/null 2>&1; then
        echo -e "${BLUE}Coverage by function:${NC}"
        go tool cover -func=coverage.out
        
        echo ""
        echo -e "${BLUE}Generating HTML coverage report...${NC}"
        go tool cover -html=coverage.out -o coverage.html
        echo -e "${GREEN}Coverage report saved to coverage.html${NC}"
    fi
    
    # Clean up
    rm -f coverage.out 2>/dev/null
fi

# Summary
echo -e "${YELLOW}═══════════════════════════════════════${NC}"
echo -e "${YELLOW}📋 TEST SUMMARY${NC}"
echo -e "${YELLOW}═══════════════════════════════════════${NC}"

TOTAL_TESTS=$((TESTS_PASSED + TESTS_FAILED))
PASS_RATE=0

if [[ $TOTAL_TESTS -gt 0 ]]; then
    PASS_RATE=$(( (TESTS_PASSED * 100) / TOTAL_TESTS ))
fi

echo -e "Total Tests: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Failed: $TESTS_FAILED${NC}"
echo -e "Pass Rate: ${PASS_RATE}%"

if [[ $TESTS_FAILED -eq 0 ]]; then
    echo ""
    echo -e "${GREEN}🎉 All tests passed! The project is fully functional.${NC}"
    echo -e "${GREEN}✨ Ready for demonstration and production use!${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}⚠️  Some tests failed. Please review the output above.${NC}"
    exit 1
fi