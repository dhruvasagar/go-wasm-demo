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

$ECHO_CMD "🚀 Building WebAssembly in Go: Bridging Web and Backend"
$ECHO_CMD "======================================================="

$ECHO_CMD "${BLUE}📦 Building WebAssembly module...${NC}"
GOOS=js GOARCH=wasm go build -o main.wasm main_wasm.go shared_models.go benchmarks.go benchmarks_wasm.go benchmarks_types.go benchmarks_comprehensive.go benchmarks_optimized.go benchmarks_shared.go mandelbrot.go mandelbrot_concurrent.go

if [ $? -eq 0 ]; then
    $ECHO_CMD "${GREEN}✅ WebAssembly module built successfully: main.wasm${NC}"
else
    $ECHO_CMD "${RED}❌ Failed to build WebAssembly module${NC}"
    exit 1
fi

$ECHO_CMD "${BLUE}🖥️  Building server binary...${NC}"
go build -o server main_server.go shared_models.go benchmarks.go

if [ $? -eq 0 ]; then
    $ECHO_CMD "${GREEN}✅ Server binary built successfully: server${NC}"
else
    $ECHO_CMD "${RED}❌ Failed to build server binary${NC}"
    exit 1
fi

$ECHO_CMD ""
$ECHO_CMD "${GREEN}🎉 Build completed successfully!${NC}"
$ECHO_CMD ""
$ECHO_CMD "${YELLOW}📋 Next Steps:${NC}"
$ECHO_CMD "1. Start the server:"
$ECHO_CMD "   ${CYAN}./server${NC}"
$ECHO_CMD ""
$ECHO_CMD "2. Open your browser and visit:"
$ECHO_CMD "   ${CYAN}http://localhost:8181/${NC} - WebAssembly client demo"
$ECHO_CMD "   ${CYAN}http://localhost:8181/server.html${NC} - Server API demo"
$ECHO_CMD ""
$ECHO_CMD "${YELLOW}🔧 Development Commands:${NC}"
$ECHO_CMD "• Rebuild WebAssembly:"
$ECHO_CMD "  ${CYAN}GOOS=js GOARCH=wasm go build -o main.wasm main_wasm.go shared_models.go benchmarks.go benchmarks_wasm.go benchmarks_types.go benchmarks_comprehensive.go benchmarks_optimized.go benchmarks_shared.go mandelbrot.go mandelbrot_concurrent.go${NC}"
$ECHO_CMD ""
$ECHO_CMD "• Rebuild Server:"
$ECHO_CMD "  ${CYAN}go build -o server main_server.go shared_models.go benchmarks.go${NC}"
$ECHO_CMD ""
$ECHO_CMD "• Run directly:"
$ECHO_CMD "  ${CYAN}go run main_server.go shared_models.go benchmarks.go${NC}"
$ECHO_CMD ""
$ECHO_CMD "${GREEN}🌟 Enjoy exploring the power of shared Go business logic!${NC}"