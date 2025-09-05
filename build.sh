#!/bin/bash

echo "🚀 Building WebAssembly in Go: Bridging Web and Backend"
echo "======================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}📦 Building WebAssembly module...${NC}"
GOOS=js GOARCH=wasm go build -o main.wasm main_wasm.go shared_models.go benchmarks.go mandelbrot.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ WebAssembly module built successfully: main.wasm${NC}"
else
    echo -e "${RED}❌ Failed to build WebAssembly module${NC}"
    exit 1
fi

echo -e "${BLUE}🖥️  Building server binary...${NC}"
go build -o server main_server.go shared_models.go benchmarks.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Server binary built successfully: server${NC}"
else
    echo -e "${RED}❌ Failed to build server binary${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 Build completed successfully!${NC}"
echo ""
echo -e "${YELLOW}📋 Next Steps:${NC}"
echo "1. Start the server:"
echo "   ${BLUE}./server${NC}"
echo ""
echo "2. Open your browser and visit:"
echo "   ${BLUE}http://localhost:8080/${NC} - WebAssembly client demo"
echo "   ${BLUE}http://localhost:8080/server.html${NC} - Server API demo"
echo ""
echo -e "${YELLOW}🔧 Development Commands:${NC}"
echo "• Rebuild WebAssembly: ${BLUE}GOOS=js GOARCH=wasm go build -o main.wasm main_wasm.go shared_models.go benchmarks.go mandelbrot.go${NC}"
echo "• Rebuild Server: ${BLUE}go build -o server main_server.go shared_models.go benchmarks.go${NC}"
echo "• Run directly: ${BLUE}go run main_server.go shared_models.go benchmarks.go${NC}"
echo ""
echo -e "${GREEN}🌟 Enjoy exploring the power of shared Go business logic!${NC}"