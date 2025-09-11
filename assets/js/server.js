// Server Dashboard JavaScript

// Test business logic APIs
async function testUserValidation() {
    const testUser = {
        email: "john.doe@example.com",
        name: "John Doe",
        age: 28,
        country: "US",
        premium: true
    };
    
    document.getElementById('userValidationResults').textContent = 'Testing user validation API...\n';
    
    try {
        const response = await fetch('/api/validate-user', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(testUser)
        });
        const result = await response.json();
        
        document.getElementById('userValidationResults').textContent = 
            `✅ User Validation API Response:\n\n` +
            `Status: ${result.valid ? 'Valid' : 'Invalid'}\n` +
            `Input: ${JSON.stringify(testUser, null, 2)}\n` +
            `Result: ${JSON.stringify(result, null, 2)}`;
    } catch (error) {
        document.getElementById('userValidationResults').textContent = `❌ Error: ${error.message}`;
    }
}

async function testProductValidation() {
    const testProduct = {
        name: "Wireless Headphones",
        price: 99.99,
        category: "electronics",
        rating: 4.5,
        in_stock: true
    };
    
    document.getElementById('productValidationResults').textContent = 'Testing product validation API...\n';
    
    try {
        const response = await fetch('/api/validate-product', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(testProduct)
        });
        const result = await response.json();
        
        document.getElementById('productValidationResults').textContent = 
            `✅ Product Validation API Response:\n\n` +
            `Status: ${result.valid ? 'Valid' : 'Invalid'}\n` +
            `Input: ${JSON.stringify(testProduct, null, 2)}\n` +
            `Result: ${JSON.stringify(result, null, 2)}`;
    } catch (error) {
        document.getElementById('productValidationResults').textContent = `❌ Error: ${error.message}`;
    }
}

async function testOrderCalculation() {
    const testData = {
        order: {
            products: [
                { id: 1, name: "Wireless Headphones", price: 99.99, category: "electronics" },
                { id: 2, name: "Cotton T-Shirt", price: 24.99, category: "clothing" }
            ],
            quantities: [1, 2]
        },
        user: {
            email: "john.doe@example.com",
            name: "John Doe",
            age: 28,
            country: "US",
            premium: true
        }
    };
    
    document.getElementById('orderCalculationResults').textContent = 'Testing order calculation API...\n';
    
    try {
        const response = await fetch('/api/calculate-order', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(testData)
        });
        const result = await response.json();
        
        document.getElementById('orderCalculationResults').textContent = 
            `✅ Order Calculation API Response:\n\n` +
            `Subtotal: $${result.subtotal.toFixed(2)}\n` +
            `Tax: $${result.tax.toFixed(2)}\n` +
            `Shipping: $${result.shipping.toFixed(2)}\n` +
            `Discount: $${result.discount.toFixed(2)}\n` +
            `Total: $${result.total.toFixed(2)}`;
    } catch (error) {
        document.getElementById('orderCalculationResults').textContent = `❌ Error: ${error.message}`;
    }
}

async function testRecommendations() {
    const testData = {
        user: {
            email: "john.doe@example.com",
            name: "John Doe",
            age: 28,
            country: "US",
            premium: true
        },
        products: [
            { id: 1, name: "Wireless Headphones", price: 99.99, category: "electronics", in_stock: true, rating: 4.5 },
            { id: 2, name: "Cotton T-Shirt", price: 24.99, category: "clothing", in_stock: true, rating: 4.2 },
            { id: 3, name: "Programming Book", price: 49.99, category: "books", in_stock: true, rating: 4.8 },
            { id: 4, name: "Coffee Mug", price: 12.99, category: "home", in_stock: true, rating: 4.0 },
            { id: 5, name: "Running Shoes", price: 129.99, category: "sports", in_stock: true, rating: 4.6 }
        ],
        order: {
            products: [
                { id: 1, name: "Wireless Headphones", price: 99.99, category: "electronics" }
            ]
        }
    };
    
    document.getElementById('recommendationsResults').textContent = 'Testing recommendations API...\n';
    
    try {
        const response = await fetch('/api/recommend-products', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(testData)
        });
        const result = await response.json();
        
        document.getElementById('recommendationsResults').textContent = 
            `✅ Recommendations API Response:\n\n` +
            `Recommended Products:\n` +
            result.map(p => `• ${p.name} - $${p.price} (${p.category})`).join('\n');
    } catch (error) {
        document.getElementById('recommendationsResults').textContent = `❌ Error: ${error.message}`;
    }
}

// Performance benchmarks
async function benchmarkServerMatrix() {
    const size = parseInt(document.getElementById('serverMatrixSize').value);
    document.getElementById('serverMatrixResults').textContent = 'Running server-side matrix benchmark...\n';
    
    try {
        const response = await fetch(`/api/benchmark/matrix?size=${size}`);
        const result = await response.json();
        
        document.getElementById('serverMatrixResults').textContent = 
            `✅ Server Matrix Multiplication Benchmark:\n\n` +
            `Operation: ${result.operation}\n` +
            `Size: ${result.size}\n` +
            `Duration: ${result.duration_ms.toFixed(2)}ms\n` +
            `Operations: ${result.operations.toLocaleString()}\n` +
            `Throughput: ${(result.operations / result.duration_ms * 1000).toFixed(0)} ops/sec`;
    } catch (error) {
        document.getElementById('serverMatrixResults').textContent = `❌ Error: ${error.message}`;
    }
}

async function benchmarkServerMandelbrot() {
    const [width, height] = document.getElementById('serverMandelbrotSize').value.split(',').map(Number);
    const iterations = parseInt(document.getElementById('serverMandelbrotIterations').value);
    
    document.getElementById('serverMandelbrotResults').textContent = 'Running server-side Mandelbrot benchmark...\n';
    
    try {
        const response = await fetch(`/api/benchmark/mandelbrot?width=${width}&height=${height}&iterations=${iterations}`);
        const result = await response.json();
        
        document.getElementById('serverMandelbrotResults').textContent = 
            `✅ Server Mandelbrot Set Benchmark:\n\n` +
            `Operation: ${result.operation}\n` +
            `Size: ${result.size}\n` +
            `Iterations: ${result.iterations}\n` +
            `Duration: ${result.duration_ms.toFixed(2)}ms\n` +
            `Pixels: ${result.pixels.toLocaleString()}\n` +
            `Throughput: ${(result.pixels / result.duration_ms * 1000).toFixed(0)} pixels/sec`;
    } catch (error) {
        document.getElementById('serverMandelbrotResults').textContent = `❌ Error: ${error.message}`;
    }
}

async function benchmarkServerHash() {
    const count = parseInt(document.getElementById('serverHashCount').value);
    document.getElementById('serverHashResults').textContent = 'Running server-side hash benchmark...\n';
    
    try {
        const response = await fetch(`/api/benchmark/hash?count=${count}`);
        const result = await response.json();
        
        document.getElementById('serverHashResults').textContent = 
            `✅ Server SHA256 Hash Benchmark:\n\n` +
            `Operation: ${result.operation}\n` +
            `Count: ${result.count.toLocaleString()}\n` +
            `Duration: ${result.duration_ms.toFixed(2)}ms\n` +
            `Throughput: ${(result.count / result.duration_ms * 1000).toFixed(0)} hashes/sec`;
    } catch (error) {
        document.getElementById('serverHashResults').textContent = `❌ Error: ${error.message}`;
    }
}