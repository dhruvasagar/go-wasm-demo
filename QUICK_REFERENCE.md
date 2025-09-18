# 🚀 Quick Reference Guide
## Everything You Need to Know at a Glance

---

## Essential Commands

### 🔨 Build & Run
```bash
./build.sh          # Build everything (WASM + server)
./server             # Start server on localhost:8181
./test.sh            # Run comprehensive test suite
./test.sh coverage   # Generate coverage report
./test.sh bench      # Run performance benchmarks
```

### 🧪 Testing Options
```bash
./test.sh            # Quick test run
./test.sh bench      # Include benchmarks
./test.sh coverage   # Generate coverage report  
./test.sh full       # All tests + benchmarks + coverage
go test -C src -v    # Verbose Go testing
go test -C src -race # Race condition detection
```

### 🔧 Manual Build (if needed)
```bash
# WebAssembly build
GOOS=js GOARCH=wasm go build -o main.wasm src/main_wasm.go src/shared_models.go src/benchmarks_*.go src/mandelbrot*.go

# Server build
go build -o server src/main_server.go src/shared_models.go src/benchmarks_*.go src/mandelbrot*.go

# Run server manually
./server
```

---

## 🌐 Web Interfaces

### **Main Demos**
- **http://localhost:8181/** - WebAssembly interactive demo
- **http://localhost:8181/server.html** - Server API comparison
- **http://localhost:8181/performance_benchmarks.html** - Performance tests

### **Demo Features**
- **User Validation**: Email regex, age limits, country validation
- **Order Calculator**: Tax rates, shipping, premium discounts
- **Product Recommendations**: ML-style recommendation engine
- **Performance Benchmarks**: Matrix, Mandelbrot, hashing comparisons

---

## 📁 Project Structure

### 🌟 **Core Source Files**
```
src/
├── shared_models.go         # 💎 Single source of truth - business logic
├── main_wasm.go            # 🌐 WebAssembly entry point
├── main_server.go          # 🖥️  Server entry point  
├── benchmarks_optimized.go # ⚡ High-performance algorithms
├── benchmarks_unified.go   # 🔄 Unified benchmark interface
├── mandelbrot.go           # 🎨 Mandelbrot set implementation
└── *_test.go              # 🧪 Comprehensive test suite
```

### 📚 **Documentation**
```
docs/
├── TESTING.md                    # 🧪 Testing guide
├── MOBILE_WEBASSEMBLY.md         # 📱 Mobile integration guide
├── WASM_FRAMEWORKS_2025.md       # 🚀 Latest framework updates
├── WEBASSEMBLY_IN_PRODUCTION.md  # 🏢 Production case studies
├── optimizations/                # ⚡ Performance guides
│   ├── OPTIMIZATION_GUIDE.md     # Core optimization strategies
│   └── WASM_OPTIMIZATION_RESULTS.md # Benchmark results
└── presentations/               # 🎯 Conference presentations
    ├── PRESENTATION_25MIN.md    # Standard conference slot
    └── PRESENTATION_30MIN.md    # Extended workshop format
```

### 🎨 **Web Assets**
```
assets/
├── css/         # Stylesheets for all demos
├── js/          # JavaScript utilities and benchmarks
index.html       # Main WebAssembly demo
server.html      # Server comparison interface
performance_benchmarks.html # Performance testing suite
```

---

## ⚡ Performance Highlights

### 🎯 **Boundary Call Optimization**
Our optimizations dramatically reduce expensive JavaScript ↔ WebAssembly calls:

| Algorithm | Before | After | Improvement |
|-----------|--------|-------|-------------|
| **Matrix 200x200** | 40M calls | 1,500 calls | **26,666x fewer** |
| **Mandelbrot 800x600** | 480M calls | 1,000 calls | **480,000x fewer** |
| **Hash Algorithm** | Variable | Batched | **Consistent performance** |

### 🚀 **Algorithm Performance** (vs JavaScript)
- **Hash Function**: 264 MHashes/sec (4.9x faster)
- **Matrix Multiplication**: 1,631 MOps/sec (3-5x faster)
- **Mandelbrot Rendering**: 2-7x faster depending on complexity
- **Business Logic**: 4x faster validation and calculations

### 📊 **Memory Efficiency**
- **JavaScript**: ~15MB heap for complex operations
- **WebAssembly**: ~4MB linear memory
- **Reduction**: 73% less memory usage

---

## 🛠️ Development Guidelines

### ✅ **Best Practices**
1. **Always use typed arrays** for WASM ↔ JS data transfer
2. **Batch operations** to minimize boundary crossings
3. **Pre-allocate arrays** for better performance  
4. **Use `unsafe.Slice`** for bulk data operations
5. **Test in multiple browsers** for compatibility
6. **Profile performance** to identify bottlenecks

### ❌ **Common Pitfalls**
1. **Avoid goroutines** for CPU-intensive tasks in WASM (no parallelism benefit)
2. **Don't use concurrent versions** in browser environment
3. **Minimize individual value transfers** across JS/WASM boundary
4. **Avoid frequent small allocations** in hot loops

### 🔧 **Debugging Tips**
```bash
# Enable Go race detector
go test -C src -race

# Browser debugging
# Open DevTools → Sources → main.wasm for WASM debugging
# Console.log window.validateUserWasm to test functions

# Performance profiling
# Use Performance tab in DevTools during benchmark runs
```

---

## 🚨 Common Issues & Solutions

### **Build Issues**

#### "Cannot find main.wasm"
```bash
# Solution: Run build script
./build.sh
# Or manual build:
GOOS=js GOARCH=wasm go build -o main.wasm src/main_wasm.go src/shared_models.go src/benchmarks_*.go src/mandelbrot*.go
```

#### "Server not starting"
```bash
# Check if port 8181 is available
lsof -i :8181
# Kill existing process or use different port
./server --port=8182
```

### **Runtime Issues**

#### "Invalid format" errors
Different demos expect different input formats:
- **index.html Mandelbrot**: "800,600" (comma-separated)
- **performance_benchmarks.html**: "800" (single value)
- **Matrix operations**: Always use integers (100, 200, 300)

#### "Slice bounds out of range"
```go
// ❌ Wrong: Fixed array casting
result := (*[40000]int)(unsafe.Pointer(&data[0]))[:size]

// ✅ Correct: Dynamic slice creation  
result := unsafe.Slice((*int)(unsafe.Pointer(&data[0])), size)
```

#### "Performance worse than JavaScript"
1. **Increase computational complexity**: Use larger matrices (500x500+) or more Mandelbrot iterations (400+)
2. **Check browser**: Chrome typically shows best WASM performance
3. **Enable optimizations**: Ensure using `benchmarks_optimized.go` versions
4. **Profile execution**: Use browser DevTools Performance tab

### **Testing Issues**

#### Tests failing to run
```bash
# Ensure you're in the right directory
cd /path/to/go-wasm-demo

# Run tests from project root
./test.sh

# Or from src directory
cd src && go test -v ./...
```

#### Coverage reports not generating
```bash
# Install cover tool if missing
go install golang.org/x/tools/cmd/cover@latest

# Generate coverage manually
go test -C src -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

---

## 📈 Benchmark Expectations

### **When WebAssembly Should Win**
- **Matrix Multiplication**: 300x300 and larger
- **Mandelbrot Set**: 800x600 with 200+ iterations
- **Hash Operations**: Batched operations (1000+ hashes)
- **Business Logic**: Complex validation with multiple rules

### **When JavaScript Might Win**
- **Very small operations**: 50x50 matrices or simple calculations
- **Heavy DOM manipulation**: Stick to JavaScript for UI updates
- **Simple string operations**: V8's string optimization is excellent
- **Single-operation calls**: Boundary crossing overhead dominates

### **Performance Targets**
Based on modern browsers (Chrome 90+, Firefox 88+, Safari 14+):
- **Matrix 200x200**: WASM should be 2-3x faster
- **Mandelbrot 800x600**: WASM should be 3-7x faster  
- **Hash batching**: WASM should be 4-5x faster
- **Business logic**: WASM should be 3-4x faster

---

## 🎯 Quick Demo Script

Perfect for showing off the project in 5 minutes:

```bash
# 1. Build and start (30 seconds)
./build.sh && ./server

# 2. Open browser demos (2 minutes)
# → http://localhost:8181/ - Show user validation working identically
# → Try invalid email, age > 120, invalid country
# → Show order calculator with different countries
# → Demonstrate instant offline calculations

# 3. Performance comparison (2 minutes)  
# → http://localhost:8181/performance_benchmarks.html
# → Run Matrix 300x300 - show WASM winning
# → Run Mandelbrot 800x600 - show dramatic speedup
# → Explain boundary call optimization

# 4. Show the code (30 seconds)
# → Open src/shared_models.go
# → Point out ValidateUser function used by both WASM and server
# → Same 400+ lines of business logic, zero duplication
```

---

## 🌟 Next Steps

### **For Learning**
1. **Explore the demos** - Understand shared business logic in action
2. **Read the source** - Study `src/shared_models.go` for patterns
3. **Run benchmarks** - See WebAssembly performance advantages
4. **Check mobile docs** - [Mobile WebAssembly Guide](./docs/MOBILE_WEBASSEMBLY.md)

### **For Development**
1. **Clone the repo** - Use as foundation for your projects
2. **Adapt the patterns** - Apply shared logic concept to your use case
3. **Optimize performance** - Follow [Optimization Guide](./docs/optimizations/OPTIMIZATION_GUIDE.md)
4. **Contribute back** - Share your improvements and use cases

### **For Production**
1. **Review case studies** - [Production WebAssembly](./docs/WEBASSEMBLY_IN_PRODUCTION.md)
2. **Plan architecture** - Design for code sharing from day one
3. **Consider frameworks** - [2025 WebAssembly Frameworks](./docs/WASM_FRAMEWORKS_2025.md)
4. **Test thoroughly** - Use our [Testing Guide](./docs/TESTING.md) as template

---

**🚀 Ready to revolutionize your development with shared Go business logic? Everything you need is right here!**