# 🚀 Quick Reference Guide

## Essential Commands

### Build & Run
```bash
./build.sh          # Build everything
./server           # Start server on localhost:8181
```

### Testing
```bash
./test.sh          # Run all tests
go test -C src -bench .   # Run benchmarks
```

## Key Files

### 🌐 Web Interfaces
- **http://localhost:8181/** - WebAssembly demo (index.html)
- **http://localhost:8181/server.html** - Server API demo
- **http://localhost:8181/performance_benchmarks.html** - Performance tests

### 📝 Core Code Files
- `src/shared_models.go` - Business logic used by both WASM and server
- `src/benchmarks_optimized.go` - High-performance WASM implementations
- `src/main_wasm.go` - WebAssembly entry point
- `src/main_server.go` - Server implementation

### 📚 Documentation
- `README.md` - Project overview
- `docs/optimizations/` - Performance optimization guides
- `docs/summaries/FINAL_OPTIMIZATION_SUMMARY.md` - Key optimizations

## Performance Highlights

### Boundary Call Optimization
- Matrix 200x200: **26,666x fewer** boundary calls
- Mandelbrot 800x600: **480,000x fewer** boundary calls

### Algorithm Performance
- Hash function: **264 MHashes/sec** (4.9x faster)
- Matrix multiplication: **1631 MOps/sec**
- All benchmarks use optimal memory patterns

## Development Tips

1. **Always use typed arrays** for WASM data transfer
2. **Avoid goroutines** for CPU tasks in WASM (no parallelism)
3. **Use `unsafe.Slice`** for bulk data operations
4. **Test in multiple browsers** for compatibility

## Quick Fixes

### "Invalid format" error
- Mandelbrot in index.html uses comma: "800,600"
- Mandelbrot in performance uses single value: "800"

### Slice bounds panic
- Use `unsafe.Slice` instead of fixed-size array casting
- Ensure typed arrays for `js.CopyBytesToJS`

### Performance issues
- Use optimized single-threaded versions
- Avoid concurrent versions in WASM
- Enable browser developer tools for profiling