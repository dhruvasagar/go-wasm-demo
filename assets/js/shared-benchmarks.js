// Shared benchmark functions for both index.html and performance_benchmarks.html
// This file consolidates common benchmark logic to avoid duplication

// ============================================================================
// WASM INITIALIZATION AND STATUS
// ============================================================================
let wasmReady = false;

function initializeWasm() {
    const go = new Go();
    return WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
        .then((result) => {
            go.run(result.instance);
            wasmReady = true;
            console.log("✅ WebAssembly module loaded and ready!");
            return true;
        })
        .catch((err) => {
            console.error("❌ Failed to load WebAssembly:", err);
            throw err;
        });
}

// ============================================================================
// SHARED JAVASCRIPT IMPLEMENTATIONS (FOR FAIR COMPARISON)
// ============================================================================

// Optimized matrix multiplication with cache-friendly access pattern
function matrixMultiplyJSShared(matrixA, matrixB, size) {
    const result = new Float64Array(size * size);
    
    // Use ikj loop order for better cache locality
    for (let i = 0; i < size; i++) {
        for (let k = 0; k < size; k++) {
            const aik = matrixA[i * size + k];
            for (let j = 0; j < size; j++) {
                result[i * size + j] += aik * matrixB[k * size + j];
            }
        }
    }
    
    return result;
}

// Optimized Mandelbrot with typed arrays
function mandelbrotJSShared(width, height, xmin, xmax, ymin, ymax, maxIter) {
    const result = new Int32Array(width * height);
    const dx = (xmax - xmin) / width;
    const dy = (ymax - ymin) / height;
    
    let idx = 0;
    for (let py = 0; py < height; py++) {
        const cy = ymin + py * dy;
        
        for (let px = 0; px < width; px++) {
            const cx = xmin + px * dx;
            
            let zx = 0, zy = 0;
            let iter = 0;
            
            // Optimized inner loop
            while (iter < maxIter) {
                const zx2 = zx * zx;
                const zy2 = zy * zy;
                
                if (zx2 + zy2 > 4) {
                    break;
                }
                
                const temp = zx2 - zy2 + cx;
                zy = 2 * zx * zy + cy;
                zx = temp;
                iter++;
            }
            
            result[idx++] = iter;
        }
    }
    
    return result;
}

// Optimized hash function with loop unrolling
function hashJSShared(data, iterations) {
    const dataBytes = new TextEncoder().encode(data);
    const dataLen = dataBytes.length;
    let hash = 0x12345678 >>> 0; // Ensure unsigned 32-bit
    
    for (let iter = 0; iter < iterations; iter++) {
        // Process 4 bytes at a time when possible
        let i = 0;
        for (; i <= dataLen - 4; i += 4) {
            hash = ((hash * 33) + dataBytes[i]) >>> 0;
            hash = ((hash << 5) | (hash >>> 27)) >>> 0;
            
            hash = ((hash * 33) + dataBytes[i + 1]) >>> 0;
            hash = ((hash << 5) | (hash >>> 27)) >>> 0;
            
            hash = ((hash * 33) + dataBytes[i + 2]) >>> 0;
            hash = ((hash << 5) | (hash >>> 27)) >>> 0;
            
            hash = ((hash * 33) + dataBytes[i + 3]) >>> 0;
            hash = ((hash << 5) | (hash >>> 27)) >>> 0;
        }
        
        // Process remaining bytes
        for (; i < dataLen; i++) {
            hash = ((hash * 33) + dataBytes[i]) >>> 0;
            hash = ((hash << 5) | (hash >>> 27)) >>> 0;
        }
    }
    
    return hash | 0; // Convert to signed 32-bit
}

// Fixed and optimized Ray Tracing implementation matching WASM complexity
function rayTracingJSShared(width, height, samples) {
    console.log(`[DEBUG] rayTracingJSShared called with ${width}x${height}, ${samples} samples`);
    const result = new Float64Array(width * height * 3);

    // Sphere properties (same as WASM implementation)
    const sphereX = 0.0, sphereY = 0.0, sphereZ = -5.0;
    const sphereRadius2 = 1.0;
    
    // Light direction (same as WASM implementation)
    const lightX = -0.57735027, lightY = -0.57735027, lightZ = -0.57735027;

    for (let y = 0; y < height; y++) {
        const ny = (y / height) * 2.0 - 1.0;
        
        for (let x = 0; x < width; x++) {
            const nx = (x / width) * 2.0 - 1.0;
            
            let colorR = 0, colorG = 0, colorB = 0;

            // Sample accumulation (same complexity as WASM)
            for (let s = 0; s < samples; s++) {
                // Ray direction normalization
                const rayLen = Math.sqrt(nx * nx + ny * ny + 1.0);
                const invRayLen = 1.0 / rayLen;
                const dirX = nx * invRayLen;
                const dirY = ny * invRayLen;
                const dirZ = -1.0 * invRayLen;

                // Ray-sphere intersection (same algorithm as WASM)
                const ocX = 0.0 - sphereX;
                const ocY = 0.0 - sphereY;
                const ocZ = 0.0 - sphereZ;
                
                const rayA = dirX * dirX + dirY * dirY + dirZ * dirZ;
                const rayB = 2.0 * (ocX * dirX + ocY * dirY + ocZ * dirZ);
                const rayC = ocX * ocX + ocY * ocY + ocZ * ocZ - sphereRadius2;
                
                const discriminant = rayB * rayB - 4.0 * rayA * rayC;
                
                if (discriminant < 0) {
                    // Background color
                    colorR += 0.2;
                    colorG += 0.2;
                    colorB += 0.8;
                } else {
                    // Hit the sphere
                    const sqrtDisc = Math.sqrt(discriminant);
                    let t = (-rayB - sqrtDisc) / (2.0 * rayA);
                    if (t < 0) {
                        t = (-rayB + sqrtDisc) / (2.0 * rayA);
                    }
                    
                    if (t < 0) {
                        // Behind camera
                        colorR += 0.2;
                        colorG += 0.2;
                        colorB += 0.8;
                    } else {
                        // Calculate intersection point and normal
                        const ix = 0.0 + t * dirX;
                        const iy = 0.0 + t * dirY;
                        const iz = 0.0 + t * dirZ;
                        
                        const normalX = ix - sphereX;
                        const normalY = iy - sphereY;
                        const normalZ = iz - sphereZ;
                        
                        // Lighting calculation (same as WASM)
                        const dot = normalX * lightX + normalY * lightY + normalZ * lightZ;
                        const intensity = Math.max(0.0, dot);
                        
                        const baseColor = 0.2 + 0.8 * intensity;
                        colorR += baseColor * 1.0;
                        colorG += baseColor * 0.7;
                        colorB += baseColor * 0.3;
                    }
                }
            }

            const invSamples = 1.0 / samples;
            const idx = (y * width + x) * 3;
            result[idx] = colorR * invSamples;
            result[idx + 1] = colorG * invSamples;
            result[idx + 2] = colorB * invSamples;
        }
    }

    return result;
}

// ============================================================================
// SHARED RESULT DISPLAY FUNCTIONS
// ============================================================================

// Display three-way performance comparison
function displayThreeWayComparison(elementId, jsTime, singleTime, concurrentTime) {
    const maxTime = Math.max(jsTime, singleTime, concurrentTime);
    const jsPercent = (jsTime / maxTime) * 100;
    const singlePercent = (singleTime / maxTime) * 100;
    const concurrentPercent = (concurrentTime / maxTime) * 100;

    document.getElementById(elementId).innerHTML = `
        <div class="performance-bar" style="margin: 5px 0;">
            <div class="performance-fill" style="width: ${jsPercent}%; background: #f39c12;"></div>
            <div class="performance-label" style="color: #333;">JavaScript</div>
        </div>
        <div class="performance-bar" style="margin: 5px 0;">
            <div class="performance-fill" style="width: ${singlePercent}%; background: #3498db;"></div>
            <div class="performance-label" style="color: #333;">Single WASM</div>
        </div>
        <div class="performance-bar" style="margin: 5px 0;">
            <div class="performance-fill" style="width: ${concurrentPercent}%; background: #27ae60;"></div>
            <div class="performance-label" style="color: #333;">Concurrent WASM</div>
        </div>
    `;
}

// ============================================================================
// EXPORTED UTILITY FUNCTIONS
// ============================================================================

// Check if WebAssembly is ready
window.isWasmReady = () => wasmReady;

// Initialize WebAssembly (call this from your main scripts)
window.initWasm = initializeWasm;

// Export shared display functions
window.displayThreeWayComparison = displayThreeWayComparison;

// Export individual JS implementations for compatibility
window.matrixMultiplyJSOptimized = matrixMultiplyJSShared;
window.mandelbrotJSOptimized = mandelbrotJSShared;
window.sha256HashJSOptimized = hashJSShared;
window.rayTracingJSOptimized = rayTracingJSShared;

// Export simpler names for convenience
window.matrixMultiplyJS = matrixMultiplyJSShared;
window.mandelbrotJS = mandelbrotJSShared;
window.hashJS = hashJSShared;
window.rayTracingJS = rayTracingJSShared;