// Initialize WebAssembly using shared function
window.initWasm().then(() => {
    console.log("✅ WebAssembly initialized for main page!");
}).catch((err) => {
    console.error("❌ Failed to initialize WebAssembly:", err);
});

// Demo data
const demoProducts = [
    {"id": 1, "name": "Wireless Headphones", "price": 99.99, "category": "electronics", "in_stock": true, "rating": 4.5, "description": "High-quality wireless headphones"},
    {"id": 2, "name": "Cotton T-Shirt", "price": 24.99, "category": "clothing", "in_stock": true, "rating": 4.2, "description": "Comfortable 100% cotton t-shirt"},
    {"id": 3, "name": "Programming Book", "price": 49.99, "category": "books", "in_stock": true, "rating": 4.8, "description": "Learn advanced programming techniques"},
    {"id": 4, "name": "Coffee Mug", "price": 12.99, "category": "home", "in_stock": true, "rating": 4.0, "description": "Ceramic coffee mug"},
    {"id": 5, "name": "Running Shoes", "price": 129.99, "category": "sports", "in_stock": true, "rating": 4.6, "description": "Lightweight running shoes"}
];

// User validation functions
function validateUserWasmButton() {
    if (!window.isWasmReady()) {
	document.getElementById('userResults').className = 'results error';
	document.getElementById('userResults').textContent = 'WebAssembly not ready yet. Please wait...';
	return;
    }

    const user = {
	email: document.getElementById('userEmail').value,
	name: document.getElementById('userName').value,
	age: parseInt(document.getElementById('userAge').value) || 0,
	country: document.getElementById('userCountry').value,
	premium: document.getElementById('userPremium').checked
    };

    try {
	const result = window.validateUserWasm(JSON.stringify(user));
	const parsedResult = (typeof result === 'string') ? JSON.parse(result) : result;
	displayResult('userResults', parsedResult, '🌐 WebAssembly Client-Side Validation');
    } catch (error) {
	displayError('userResults', error);
    }
}

function validateUserServer() {
    const user = {
	email: document.getElementById('userEmail').value,
	name: document.getElementById('userName').value,
	age: parseInt(document.getElementById('userAge').value) || 0,
	country: document.getElementById('userCountry').value,
	premium: document.getElementById('userPremium').checked
    };

    fetch('/api/validate-user', {
	method: 'POST',
	headers: { 'Content-Type': 'application/json' },
	body: JSON.stringify(user)
    })
    .then(response => response.json())
    .then(result => displayResult('userResults', result, '🖥️ Server-Side API Validation'))
    .catch(error => displayError('userResults', error));
}

// Product validation functions
function validateProductWasmButton() {
    if (!window.isWasmReady()) {
	document.getElementById('productResults').className = 'results error';
	document.getElementById('productResults').textContent = 'WebAssembly not ready yet. Please wait...';
	return;
    }

    const product = {
	name: document.getElementById('productName').value,
	price: parseFloat(document.getElementById('productPrice').value) || 0,
	category: document.getElementById('productCategory').value,
	rating: parseFloat(document.getElementById('productRating').value) || 0,
	in_stock: document.getElementById('productInStock').checked
    };

    try {
	const result = window.validateProductWasm(JSON.stringify(product));
	const parsedResult = (typeof result === 'string') ? JSON.parse(result) : result;
	displayResult('productResults', parsedResult, '🌐 WebAssembly Client-Side Validation');
    } catch (error) {
	displayError('productResults', error);
    }
}

function validateProductServer() {
    const product = {
	name: document.getElementById('productName').value,
	price: parseFloat(document.getElementById('productPrice').value) || 0,
	category: document.getElementById('productCategory').value,
	rating: parseFloat(document.getElementById('productRating').value) || 0,
	in_stock: document.getElementById('productInStock').checked
    };

    fetch('/api/validate-product', {
	method: 'POST',
	headers: { 'Content-Type': 'application/json' },
	body: JSON.stringify(product)
    })
    .then(response => response.json())
    .then(result => displayResult('productResults', result, '🖥️ Server-Side API Validation'))
    .catch(error => displayError('productResults', error));
}

// Order calculation functions
function calculateOrderWasmButton() {
    if (!window.isWasmReady()) {
	document.getElementById('orderResults').className = 'results error';
	document.getElementById('orderResults').textContent = 'WebAssembly not ready yet. Please wait...';
	return;
    }

    try {
	const products = JSON.parse(document.getElementById('orderProducts').value);
	const quantities = JSON.parse(document.getElementById('orderQuantities').value);
	const user = getCurrentUser();

	const order = {
	    products: products,
	    quantities: quantities
	};

	const result = window.calculateOrderTotalWasm(JSON.stringify(order), JSON.stringify(user));
	const parsedResult = (typeof result === 'string') ? JSON.parse(result) : result;
	displayOrderResult('orderResults', parsedResult, '🌐 WebAssembly Client-Side Calculation');
    } catch (error) {
	displayError('orderResults', error);
    }
}

function calculateOrderServer() {
    try {
	const products = JSON.parse(document.getElementById('orderProducts').value);
	const quantities = JSON.parse(document.getElementById('orderQuantities').value);
	const user = getCurrentUser();

	const requestData = {
	    order: {
		products: products,
		quantities: quantities
	    },
	    user: user
	};

	fetch('/api/calculate-order', {
	    method: 'POST',
	    headers: { 'Content-Type': 'application/json' },
	    body: JSON.stringify(requestData)
	})
	.then(response => response.json())
	.then(result => displayOrderResult('orderResults', result, '🖥️ Server-Side API Calculation'))
	.catch(error => displayError('orderResults', error));
    } catch (error) {
	displayError('orderResults', error);
    }
}

// Recommendation functions
function getRecommendationsWasmButton() {
    if (!window.isWasmReady()) {
	document.getElementById('recommendationResults').className = 'results error';
	document.getElementById('recommendationResults').textContent = 'WebAssembly not ready yet. Please wait...';
	return;
    }

    try {
	const user = getCurrentUser();
	const products = demoProducts;
	const order = { products: JSON.parse(document.getElementById('orderProducts').value) };

	const result = window.recommendProductsWasm(
	    JSON.stringify(user),
	    JSON.stringify(products),
	    JSON.stringify(order)
	);

	const parsedResult = (typeof result === 'string') ? JSON.parse(result) : result;
	
	// WASM function returns an object with error and recommendations fields
	if (parsedResult.error && parsedResult.error !== "") {
	    displayError('recommendationResults', new Error(parsedResult.error));
	} else {
	    // Extract the recommendations array from the WASM response
	    const recommendations = parsedResult.recommendations || parsedResult;
	    displayRecommendations('recommendationResults', recommendations, '🌐 WebAssembly Client-Side Recommendations');
	}
    } catch (error) {
	displayError('recommendationResults', error);
    }
}

function getRecommendationsServer() {
    try {
	const user = getCurrentUser();
	const products = demoProducts;
	const order = { products: JSON.parse(document.getElementById('orderProducts').value) };

	const requestData = {
	    user: user,
	    products: products,
	    order: order
	};

	fetch('/api/recommend-products', {
	    method: 'POST',
	    headers: { 'Content-Type': 'application/json' },
	    body: JSON.stringify(requestData)
	})
	.then(response => response.json())
	.then(result => displayRecommendations('recommendationResults', result, '🖥️ Server-Side API Recommendations'))
	.catch(error => displayError('recommendationResults', error));
    } catch (error) {
	displayError('recommendationResults', error);
    }
}

// Benchmark functions
function benchmarkMatrix() {
    const size = parseInt(document.getElementById('matrixSize').value);
    document.getElementById('matrixResults').className = 'results info';
    document.getElementById('matrixResults').textContent = 'Running matrix multiplication benchmark...\n';

    // Create test matrices using typed arrays for fair comparison
    const matrixA = new Float64Array(size * size);
    const matrixB = new Float64Array(size * size);
    for (let i = 0; i < size * size; i++) {
	matrixA[i] = i % 10;
	matrixB[i] = (i * 2) % 10;
    }

    // JavaScript benchmark (5 runs for average)
    let jsTotalTime = 0;
    for (let run = 0; run < 5; run++) {
	const jsStart = performance.now();
	const jsResult = matrixMultiplyJSOptimized(matrixA, matrixB, size);
	jsTotalTime += performance.now() - jsStart;
    }
    const jsDuration = jsTotalTime / 5;

    // WebAssembly benchmark (5 runs for average)
    if (wasmReady) {
	let wasmTotalTime = 0;
	for (let run = 0; run < 5; run++) {
	    const wasmStart = performance.now();
	    const wasmResult = matrixMultiplyWasm(matrixA, matrixB, size);
	    wasmTotalTime += performance.now() - wasmStart;
	}
	const wasmDuration = wasmTotalTime / 5;

	const speedup = (jsDuration / wasmDuration).toFixed(1);
	document.getElementById('matrixResults').textContent =
	    `Matrix Multiplication (${size}x${size})\n` +
	    `JavaScript: ${jsDuration.toFixed(2)}ms (average of 5 runs)\n` +
	    `WebAssembly: ${wasmDuration.toFixed(2)}ms (average of 5 runs)\n` +
	    `Result: ${speedup}x ${speedup >= 1 ? 'faster' : 'slower'} with WebAssembly`;

	displayPerformanceComparison('matrixComparison', jsDuration, wasmDuration);
    }
}

function benchmarkMandelbrot() {
    const [width, height] = document.getElementById('mandelbrotSize').value.split(',').map(Number);
    const iterations = parseInt(document.getElementById('mandelbrotIterations').value);

    document.getElementById('mandelbrotResults').className = 'results info';
    document.getElementById('mandelbrotResults').textContent = 'Running Mandelbrot set benchmark...\n';

    // JavaScript benchmark (5 runs for average)
    let jsTotalTime = 0;
    for (let run = 0; run < 5; run++) {
	const jsStart = performance.now();
	const jsResult = mandelbrotJSOptimized(width, height, -2, 1, -1.5, 1.5, iterations);
	jsTotalTime += performance.now() - jsStart;
    }
    const jsDuration = jsTotalTime / 5;

    // WebAssembly benchmark (5 runs for average)
    if (wasmReady) {
	let wasmTotalTime = 0;
	for (let run = 0; run < 5; run++) {
	    const wasmStart = performance.now();
	    const wasmResult = mandelbrotWasm(width, height, -2, 1, -1.5, 1.5, iterations);
	    wasmTotalTime += performance.now() - wasmStart;
	}
	const wasmDuration = wasmTotalTime / 5;

	const speedup = (jsDuration / wasmDuration).toFixed(1);
	document.getElementById('mandelbrotResults').textContent =
	    `Mandelbrot Set (${width}x${height}, ${iterations} iterations)\n` +
	    `JavaScript: ${jsDuration.toFixed(2)}ms (average of 5 runs)\n` +
	    `WebAssembly: ${wasmDuration.toFixed(2)}ms (average of 5 runs)\n` +
	    `Result: ${speedup}x ${speedup >= 1 ? 'faster' : 'slower'} with WebAssembly`;

	displayPerformanceComparison('mandelbrotComparison', jsDuration, wasmDuration);
    }
}

function benchmarkHash() {
    const count = parseInt(document.getElementById('hashCount').value);
    document.getElementById('hashResults').className = 'results info';
    document.getElementById('hashResults').textContent = 'Running simple hash benchmark...\n';

    const data = "WebAssembly performance test data";

    // JavaScript benchmark (5 runs for average)
    let jsTotalTime = 0;
    for (let run = 0; run < 5; run++) {
	const jsStart = performance.now();
	const jsHash = sha256HashJSOptimized(data, count);
	jsTotalTime += performance.now() - jsStart;
    }
    const jsDuration = jsTotalTime / 5;

    // WebAssembly benchmark (5 runs for average)
    if (wasmReady) {
	let wasmTotalTime = 0;
	for (let run = 0; run < 5; run++) {
	    const wasmStart = performance.now();
	    const wasmHash = sha256HashWasm(data, count);
	    wasmTotalTime += performance.now() - wasmStart;
	}
	const wasmDuration = wasmTotalTime / 5;

	const speedup = (jsDuration / wasmDuration).toFixed(1);
	document.getElementById('hashResults').textContent =
	    `Simple Hash Function (${count.toLocaleString()} iterations)\n` +
	    `JavaScript: ${jsDuration.toFixed(2)}ms (average of 5 runs)\n` +
	    `WebAssembly: ${wasmDuration.toFixed(2)}ms (average of 5 runs)\n` +
	    `Result: ${speedup}x ${speedup >= 1 ? 'faster' : 'slower'} with WebAssembly`;

	displayPerformanceComparison('hashComparison', jsDuration, wasmDuration);
    }
}



// Helper functions
function getCurrentUser() {
    return {
	email: document.getElementById('userEmail').value || 'demo@example.com',
	name: document.getElementById('userName').value || 'Demo User',
	age: parseInt(document.getElementById('userAge').value) || 30,
	country: document.getElementById('userCountry').value || 'US',
	premium: document.getElementById('userPremium').checked
    };
}

function displayResult(elementId, result, title) {
    const element = document.getElementById(elementId);
    const isValid = result.valid;

    element.className = `results ${isValid ? 'success' : 'error'}`;
    element.textContent = `${title}\n\n` +
	`Status: ${isValid ? '✅ Valid' : '❌ Invalid'}\n` +
	(result.errors && result.errors.length > 0 ?
	    `Errors:\n${result.errors.map(e => `• ${e}`).join('\n')}` :
	    'All validation rules passed!');
}

function displayOrderResult(elementId, result, title) {
    const element = document.getElementById(elementId);
    element.className = 'results success';
    element.textContent = `${title}\n\n` +
	`Subtotal: $${result.subtotal.toFixed(2)}\n` +
	`Tax: $${result.tax.toFixed(2)}\n` +
	`Shipping: $${result.shipping.toFixed(2)}\n` +
	`Discount: $${result.discount.toFixed(2)}\n` +
	`Total: $${result.total.toFixed(2)}`;
}

function displayRecommendations(elementId, result, title) {
    const element = document.getElementById(elementId);
    element.className = 'results info';

    if (Array.isArray(result) && result.length > 0) {
	element.textContent = `${title}\n\nRecommended Products:\n` +
	    result.map(p => `• ${p.name} - $${p.price} (${p.category})`).join('\n');
    } else {
	element.textContent = `${title}\n\nNo recommendations available.`;
    }
}

function displayError(elementId, error) {
    const element = document.getElementById(elementId);
    element.className = 'results error';
    element.textContent = `Error: ${error.message || error}`;
}

function displayPerformanceComparison(elementId, jsDuration, wasmDuration) {
    const maxDuration = Math.max(jsDuration, wasmDuration);
    const jsPercent = (jsDuration / maxDuration) * 100;
    const wasmPercent = (wasmDuration / maxDuration) * 100;

    const speedupRatio = jsDuration / wasmDuration;
    const isWasmFaster = wasmDuration < jsDuration;
    const winner = isWasmFaster ? 'WebAssembly' : 'JavaScript';
    const speedupText = isWasmFaster ?
	`${speedupRatio.toFixed(1)}x faster with WebAssembly!` :
	`${(1/speedupRatio).toFixed(1)}x faster with JavaScript (WebAssembly is slower)`;
    const speedupColor = isWasmFaster ? '#27ae60' : '#e74c3c';

    document.getElementById(elementId).innerHTML = `
	<div>JavaScript: ${jsDuration.toFixed(2)}ms</div>
	<div class="performance-bar">
	    <div class="performance-fill js" style="width: ${jsPercent}%"></div>
	    <div class="performance-label">JavaScript</div>
	</div>
	<div>WebAssembly: ${wasmDuration.toFixed(2)}ms</div>
	<div class="performance-bar">
	    <div class="performance-fill wasm" style="width: ${wasmPercent}%"></div>
	    <div class="performance-label">WebAssembly</div>
	</div>
	<div style="text-align: center; margin-top: 10px; font-weight: 600; color: ${speedupColor};">
	    Winner: ${winner} - ${speedupText}
	</div>
    `;
}

// Comprehensive benchmark functions for 3-way comparison
function benchmarkMatrixComprehensive() {
    if (!window.isWasmReady()) {
	document.getElementById('matrixResults').className = 'results error';
	document.getElementById('matrixResults').textContent = 'WebAssembly not ready yet. Please wait...';
	return;
    }

    const sizeStr = document.getElementById('matrixSize').value;
    const size = parseInt(sizeStr.split('x')[0]);

    document.getElementById('matrixResults').className = 'results info';
    document.getElementById('matrixResults').textContent = 'Running comprehensive matrix benchmark...\n';
    document.getElementById('matrixComparison').innerHTML = '';

    // Generate random matrices
    const matrixA = new Array(size * size).fill(0).map(() => Math.random());
    const matrixB = new Array(size * size).fill(0).map(() => Math.random());

    setTimeout(() => {
	// JavaScript benchmark
	const jsStart = performance.now();
	matrixMultiplyJS(matrixA, matrixB, size);
	const jsTime = performance.now() - jsStart;

	// Single-threaded WASM benchmark
	const singleStart = performance.now();
	window.matrixMultiplyWasm(matrixA, matrixB, size);
	const singleTime = performance.now() - singleStart;

	// Concurrent WASM benchmark
	const concurrentStart = performance.now();
	window.matrixMultiplyConcurrentWasm(matrixA, matrixB, size);
	const concurrentTime = performance.now() - concurrentStart;

	// Display results
	const speedupSingle = (jsTime / singleTime).toFixed(2);
	const speedupConcurrent = (jsTime / concurrentTime).toFixed(2);
	const concurrentVsSingle = (singleTime / concurrentTime).toFixed(2);

	document.getElementById('matrixResults').innerHTML = `
	    <strong>Matrix Multiplication (${sizeStr})</strong><br>
	    <span style="color: #f39c12;">JavaScript:</span> ${jsTime.toFixed(2)}ms<br>
	    <span style="color: #3498db;">Single-Thread WASM:</span> ${singleTime.toFixed(2)}ms (${speedupSingle}x vs JS)<br>
	    <span style="color: #27ae60;">Concurrent WASM:</span> ${concurrentTime.toFixed(2)}ms (${speedupConcurrent}x vs JS, ${concurrentVsSingle}x vs Single)
	`;

	displayThreeWayComparison('matrixComparison', jsTime, singleTime, concurrentTime);
    }, 10);
}

function benchmarkMandelbrotComprehensive() {
    if (!window.isWasmReady()) {
	document.getElementById('mandelbrotResults').className = 'results error';
	document.getElementById('mandelbrotResults').textContent = 'WebAssembly not ready yet. Please wait...';
	return;
    }

    const sizeStr = document.getElementById('mandelbrotSize').value;
    const sizeParts = sizeStr.split(',');
    if (sizeParts.length !== 2) {
	document.getElementById('mandelbrotResults').className = 'results error';
	document.getElementById('mandelbrotResults').textContent = 'Invalid size format. Please use format like "800,600"';
	return;
    }

    const width = parseInt(sizeParts[0].replace(/\D/g, ''));
    const height = parseInt(sizeParts[1].replace(/\D/g, ''));

    if (isNaN(width) || isNaN(height) || width <= 0 || height <= 0) {
	document.getElementById('mandelbrotResults').className = 'results error';
	document.getElementById('mandelbrotResults').textContent = 'Invalid dimensions. Please enter positive numbers.';
	return;
    }

    const maxIterStr = document.getElementById('mandelbrotIterations').value;
    const maxIter = parseInt(maxIterStr.replace(/\D/g, ''));

    if (isNaN(maxIter) || maxIter <= 0) {
	document.getElementById('mandelbrotResults').className = 'results error';
	document.getElementById('mandelbrotResults').textContent = 'Invalid iteration count. Please enter a positive number.';
	return;
    }

    document.getElementById('mandelbrotResults').className = 'results info';
    document.getElementById('mandelbrotResults').textContent = 'Running comprehensive Mandelbrot benchmark...\n';
    document.getElementById('mandelbrotComparison').innerHTML = '';

    const xmin = -2.5, xmax = 1.5, ymin = -1.5, ymax = 1.5;

    setTimeout(() => {
	// JavaScript benchmark
	const jsStart = performance.now();
	mandelbrotJSOptimized(width, height, xmin, xmax, ymin, ymax, maxIter);
	const jsTime = performance.now() - jsStart;

	// Single-threaded WASM benchmark
	const singleStart = performance.now();
	window.mandelbrotWasm(width, height, xmin, xmax, ymin, ymax, maxIter);
	const singleTime = performance.now() - singleStart;

	// Concurrent WASM benchmark
	const concurrentStart = performance.now();
	window.mandelbrotConcurrentWasm(width, height, xmin, xmax, ymin, ymax, maxIter);
	const concurrentTime = performance.now() - concurrentStart;

	// Display results
	const speedupSingle = (jsTime / singleTime).toFixed(2);
	const speedupConcurrent = (jsTime / concurrentTime).toFixed(2);
	const concurrentVsSingle = (singleTime / concurrentTime).toFixed(2);

	document.getElementById('mandelbrotResults').innerHTML = `
	    <strong>Mandelbrot Set (${sizeStr}, ${maxIter} iterations)</strong><br>
	    <span style="color: #f39c12;">JavaScript:</span> ${jsTime.toFixed(2)}ms<br>
	    <span style="color: #3498db;">Single-Thread WASM:</span> ${singleTime.toFixed(2)}ms (${speedupSingle}x vs JS)<br>
	    <span style="color: #27ae60;">Concurrent WASM:</span> ${concurrentTime.toFixed(2)}ms (${speedupConcurrent}x vs JS, ${concurrentVsSingle}x vs Single)
	`;

	displayThreeWayComparison('mandelbrotComparison', jsTime, singleTime, concurrentTime);
    }, 10);
}

function benchmarkHashComprehensive() {
    if (!window.isWasmReady()) {
	document.getElementById('hashResults').className = 'results error';
	document.getElementById('hashResults').textContent = 'WebAssembly not ready yet. Please wait...';
	return;
    }

    const iterations = parseInt(document.getElementById('hashCount').value);
    const data = 'The quick brown fox jumps over the lazy dog. '.repeat(10);

    document.getElementById('hashResults').className = 'results info';
    document.getElementById('hashResults').textContent = 'Running comprehensive hash benchmark...\n';
    document.getElementById('hashComparison').innerHTML = '';

    setTimeout(() => {
	// JavaScript benchmark
	const jsStart = performance.now();
	hashJS(data, iterations);
	const jsTime = performance.now() - jsStart;

	// Single-threaded WASM benchmark
	const singleStart = performance.now();
	window.sha256HashWasm(data, iterations);
	const singleTime = performance.now() - singleStart;

	// Concurrent WASM benchmark
	const concurrentStart = performance.now();
	window.sha256HashConcurrentWasm(data, iterations);
	const concurrentTime = performance.now() - concurrentStart;

	// Display results
	const speedupSingle = (jsTime / singleTime).toFixed(2);
	const speedupConcurrent = (jsTime / concurrentTime).toFixed(2);
	const concurrentVsSingle = (singleTime / concurrentTime).toFixed(2);

	document.getElementById('hashResults').innerHTML = `
	    <strong>Hash Computation (${iterations.toLocaleString()} iterations)</strong><br>
	    <span style="color: #f39c12;">JavaScript:</span> ${jsTime.toFixed(2)}ms<br>
	    <span style="color: #3498db;">Single-Thread WASM:</span> ${singleTime.toFixed(2)}ms (${speedupSingle}x vs JS)<br>
	    <span style="color: #27ae60;">Concurrent WASM:</span> ${concurrentTime.toFixed(2)}ms (${speedupConcurrent}x vs JS, ${concurrentVsSingle}x vs Single)
	`;

	displayThreeWayComparison('hashComparison', jsTime, singleTime, concurrentTime);
    }, 10);
}

function benchmarkRayTracingComprehensive() {
    if (!window.isWasmReady()) {
	document.getElementById('rayTracingResults').className = 'results error';
	document.getElementById('rayTracingResults').textContent = 'WebAssembly not ready yet. Please wait...';
	return;
    }

    const sizeStr = document.getElementById('rayTracingSize').value;
    const sizeParts = sizeStr.split(',');
    if (sizeParts.length !== 2) {
	document.getElementById('rayTracingResults').className = 'results error';
	document.getElementById('rayTracingResults').textContent = 'Invalid size format. Please use format like "400,300"';
	return;
    }

    const width = parseInt(sizeParts[0].replace(/\D/g, ''));
    const height = parseInt(sizeParts[1].replace(/\D/g, ''));

    if (isNaN(width) || isNaN(height) || width <= 0 || height <= 0) {
	document.getElementById('rayTracingResults').className = 'results error';
	document.getElementById('rayTracingResults').textContent = 'Invalid dimensions. Please enter positive numbers.';
	return;
    }

    const samplesStr = document.getElementById('rayTracingSamples').value;
    const samples = parseInt(samplesStr.replace(/\D/g, ''));

    if (isNaN(samples) || samples <= 0) {
	document.getElementById('rayTracingResults').className = 'results error';
	document.getElementById('rayTracingResults').textContent = 'Invalid sample count. Please enter a positive number.';
	return;
    }

    document.getElementById('rayTracingResults').className = 'results info';
    document.getElementById('rayTracingResults').textContent = 'Running comprehensive ray tracing benchmark...\n';
    document.getElementById('rayTracingComparison').innerHTML = '';

    setTimeout(() => {
	// JavaScript benchmark
	const jsStart = performance.now();
	rayTracingJSOptimized(width, height, samples);
	const jsTime = performance.now() - jsStart;

	// Single-threaded WASM benchmark
	const singleStart = performance.now();
	window.rayTracingWasm(width, height, samples);
	const singleTime = performance.now() - singleStart;

	// Concurrent WASM benchmark
	const concurrentStart = performance.now();
	window.rayTracingConcurrentWasm(width, height, samples);
	const concurrentTime = performance.now() - concurrentStart;

	// Display results
	const speedupSingle = (jsTime / singleTime).toFixed(2);
	const speedupConcurrent = (jsTime / concurrentTime).toFixed(2);
	const concurrentVsSingle = (singleTime / concurrentTime).toFixed(2);

	document.getElementById('rayTracingResults').innerHTML = `
	    <strong>Ray Tracing (${width}x${height}, ${samples} samples)</strong><br>
	    <span style="color: #f39c12;">JavaScript:</span> ${jsTime.toFixed(2)}ms<br>
	    <span style="color: #3498db;">Single-Thread WASM:</span> ${singleTime.toFixed(2)}ms (${speedupSingle}x vs JS)<br>
	    <span style="color: #27ae60;">Concurrent WASM:</span> ${concurrentTime.toFixed(2)}ms (${speedupConcurrent}x vs JS, ${concurrentVsSingle}x vs Single)
	`;

	displayThreeWayComparison('rayTracingComparison', jsTime, singleTime, concurrentTime);
    }, 10);
}

// displayThreeWayComparison is now in shared-benchmarks.js

// JavaScript implementations are now in shared-benchmarks.js

// Fill demo data on load
window.addEventListener('load', () => {
    document.getElementById('userEmail').value = 'john.doe@example.com';
    document.getElementById('userName').value = 'John Doe';
    document.getElementById('userAge').value = '28';
    document.getElementById('userCountry').value = 'US';
    document.getElementById('userPremium').checked = true;

    document.getElementById('productName').value = 'Wireless Headphones';
    document.getElementById('productPrice').value = '99.99';
    document.getElementById('productCategory').value = 'electronics';
    document.getElementById('productRating').value = '4.5';
});
