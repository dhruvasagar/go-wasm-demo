// Optimized JavaScript implementations for fair benchmark comparison

// Optimized matrix multiplication with cache-friendly access pattern
function matrixMultiplyJSOptimized(matrixA, matrixB, size) {
    const result = new Float64Array(size * size);
    
    // Use ikj loop order for better cache locality (no zero check for fair comparison)
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
function mandelbrotJSOptimized(width, height, xmin, xmax, ymin, ymax, maxIter) {
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
function sha256HashJSOptimized(data, iterations) {
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