#!/bin/bash

# Function to check if terminal supports colors
supports_color() {
    if [[ -t 1 ]] && [[ "${TERM}" != "dumb" ]] && command -v tput >/dev/null 2>&1; then
        if (( $(tput colors 2>/dev/null || echo 0) >= 8 )); then
            return 0
        fi
    fi
    return 1
}

# Set colors only if terminal supports them
if supports_color; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    BLUE='\033[0;34m'
    YELLOW='\033[1;33m'
    CYAN='\033[0;36m'
    BOLD='\033[1m'
    NC='\033[0m' # No Color
    ECHO_CMD="echo -e"
else
    RED=''
    GREEN=''
    BLUE=''
    YELLOW=''
    CYAN=''
    BOLD=''
    NC=''
    ECHO_CMD="echo"
fi

$ECHO_CMD "🧪 WebAssembly in Go: Comprehensive Testing Suite"
$ECHO_CMD "================================================="

# Track test results
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run test and track results
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    $ECHO_CMD "${BLUE}🔄 Running $test_name...${NC}"
    
    if eval "$test_command"; then
        $ECHO_CMD "${GREEN}✅ $test_name passed${NC}"
        ((TESTS_PASSED++))
    else
        $ECHO_CMD "${RED}❌ $test_name failed${NC}"
        ((TESTS_FAILED++))
    fi
    $ECHO_CMD ""
}

# Function to run benchmarks
run_benchmark() {
    local bench_name="$1"
    local bench_command="$2"
    
    $ECHO_CMD "${YELLOW}⚡ Running $bench_name...${NC}"
    eval "$bench_command"
    $ECHO_CMD ""
}

$ECHO_CMD "${BLUE}📋 Test Plan:${NC}"
$ECHO_CMD "1. Unit Tests - Business Logic"
$ECHO_CMD "2. Integration Tests - API Endpoints"
$ECHO_CMD "3. Performance Benchmarks"
$ECHO_CMD "4. Build Verification"
$ECHO_CMD "5. Code Quality Checks"
$ECHO_CMD ""

# 1. Unit Tests
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
$ECHO_CMD "${YELLOW}📝 UNIT TESTS - BUSINESS LOGIC${NC}"
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"

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
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
$ECHO_CMD "${YELLOW}🔗 INTEGRATION TESTS - API ENDPOINTS${NC}"
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"

run_test "Server API Endpoints" "go test -v -run TestServerAPIEndpoints"
run_test "Benchmark Endpoints" "go test -v -run TestBenchmarkEndpoints"
run_test "Error Handling" "go test -v -run TestErrorHandling"
run_test "CORS Headers" "go test -v -run TestCORSHeaders"
run_test "Demo Data Endpoints" "go test -v -run TestDemoDataEndpoints"
run_test "Data Consistency" "go test -v -run TestDataConsistency"

# 3. Algorithm Tests
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
$ECHO_CMD "${YELLOW}🔬 ALGORITHM CORRECTNESS TESTS${NC}"
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"

run_test "Matrix Multiplication Logic" "go test -v -run TestMatrixMultiplicationLogic"
run_test "Mandelbrot Set Logic" "go test -v -run TestMandelbrotLogic"
run_test "Hash Function Consistency" "go test -v -run TestHashingConsistency"
run_test "Memory Allocation" "go test -v -run TestMemoryAllocation"

# 4. Performance Tests (if not in short mode)
if [[ "$1" != "short" ]]; then
    $ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
    $ECHO_CMD "${YELLOW}🚀 PERFORMANCE TESTS${NC}"
    $ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
    
    run_test "Server Performance" "go test -v -run TestServerPerformance"
    run_test "Stress Testing" "go test -v -run TestStressTest"
    run_test "Performance Regression" "go test -v -run TestPerformanceRegression"
    run_test "Concurrent Safety" "go test -v -run TestConcurrentSafety"
fi

# 5. Business Logic Integration
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
$ECHO_CMD "${YELLOW}🎯 BUSINESS LOGIC INTEGRATION${NC}"
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"

run_test "End-to-End Business Logic" "go test -v -run TestBusinessLogicIntegration"

# 6. Build Verification
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
$ECHO_CMD "${YELLOW}🔧 BUILD VERIFICATION${NC}"
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"

run_test "WebAssembly Build" "GOOS=js GOARCH=wasm go build -o test_main.wasm main_wasm.go shared_models.go benchmarks_wasm.go benchmarks_types.go benchmarks_comprehensive.go benchmarks_optimized.go benchmarks_shared.go mandelbrot.go mandelbrot_concurrent.go"
run_test "Server Build" "go build -o test_server main_server.go shared_models.go"
run_test "Test Compilation" "go test -c -o test_binary"

# Clean up build artifacts
rm -f test_main.wasm test_server test_binary 2>/dev/null

# 7. Performance Benchmarks
if [[ "$1" == "bench" ]] || [[ "$1" == "full" ]]; then
    $ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
    $ECHO_CMD "${YELLOW}⚡ PERFORMANCE BENCHMARKS${NC}"
    $ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
    
    run_benchmark "Business Logic Benchmarks" "go test -bench=BenchmarkValidate -benchmem"
    run_benchmark "Algorithm Benchmarks" "go test -bench=BenchmarkMatrix -benchmem -benchtime=3s"
    run_benchmark "Mandelbrot Benchmarks" "go test -bench=BenchmarkMandelbrot -benchmem -benchtime=2s"
    run_benchmark "Hash Benchmarks" "go test -bench=BenchmarkHashing -benchmem"
    run_benchmark "Server Function Benchmarks" "go test -bench=BenchmarkServer -benchmem"
fi

# 8. Code Quality Checks
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
$ECHO_CMD "${YELLOW}✨ CODE QUALITY CHECKS${NC}"
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"

# Check for gofmt
if command -v gofmt >/dev/null 2>&1; then
    if gofmt -l . | grep -q .; then
        $ECHO_CMD "${RED}❌ Code formatting issues found:${NC}"
        gofmt -l .
        ((TESTS_FAILED++))
    else
        $ECHO_CMD "${GREEN}✅ Code formatting check passed${NC}"
        ((TESTS_PASSED++))
    fi
else
    $ECHO_CMD "${YELLOW}⚠️  gofmt not available, skipping format check${NC}"
fi

# Check for go vet
if command -v go >/dev/null 2>&1; then
    if go vet ./...; then
        $ECHO_CMD "${GREEN}✅ Go vet check passed${NC}"
        ((TESTS_PASSED++))
    else
        $ECHO_CMD "${RED}❌ Go vet found issues${NC}"
        ((TESTS_FAILED++))
    fi
fi

# Test Coverage Report
if [[ "$1" == "coverage" ]] || [[ "$1" == "full" ]]; then
    $ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
    $ECHO_CMD "${YELLOW}📊 TEST COVERAGE REPORT${NC}"
    $ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
    
    $ECHO_CMD "${BLUE}Generating coverage report...${NC}"
    go test -coverprofile=coverage.out -covermode=count
    
    if command -v go >/dev/null 2>&1; then
        $ECHO_CMD "${BLUE}Coverage by function:${NC}"
        go tool cover -func=coverage.out
        
        $ECHO_CMD ""
        $ECHO_CMD "${BLUE}Generating HTML coverage report...${NC}"
        go tool cover -html=coverage.out -o coverage.html
        $ECHO_CMD "${GREEN}Coverage report saved to coverage.html${NC}"
    fi
    
    # Clean up
    rm -f coverage.out 2>/dev/null
fi

# Summary
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"
$ECHO_CMD "${YELLOW}📋 TEST SUMMARY${NC}"
$ECHO_CMD "${YELLOW}═══════════════════════════════════════${NC}"

TOTAL_TESTS=$((TESTS_PASSED + TESTS_FAILED))
PASS_RATE=0

if [[ $TOTAL_TESTS -gt 0 ]]; then
    PASS_RATE=$(( (TESTS_PASSED * 100) / TOTAL_TESTS ))
fi

$ECHO_CMD "Total Tests: $TOTAL_TESTS"
$ECHO_CMD "${GREEN}Passed: $TESTS_PASSED${NC}"
$ECHO_CMD "${RED}Failed: $TESTS_FAILED${NC}"
$ECHO_CMD "Pass Rate: ${PASS_RATE}%"

if [[ $TESTS_FAILED -eq 0 ]]; then
    $ECHO_CMD ""
    $ECHO_CMD "${GREEN}🎉 All tests passed! The project is fully functional.${NC}"
    $ECHO_CMD "${GREEN}✨ Ready for demonstration and production use!${NC}"
    exit 0
else
    $ECHO_CMD ""
    $ECHO_CMD "${RED}⚠️  Some tests failed. Please review the output above.${NC}"
    exit 1
fi
