// Performance Benchmarks JavaScript
const benchmarkResults = {};

// Loading overlay functions
function showLoadingOverlay(text = 'Loading...', details = '') {
    const overlay = document.getElementById('loadingOverlay');
    const loadingText = document.getElementById('loadingText');
    const loadingDetails = document.getElementById('loadingDetails');
    
    loadingText.textContent = text;
    loadingDetails.textContent = details;
    
    overlay.classList.remove('hide');
    overlay.classList.add('show');
}

function hideLoadingOverlay() {
    const overlay = document.getElementById('loadingOverlay');
    overlay.classList.remove('show');
    overlay.classList.add('hide');
    
    setTimeout(() => {
        overlay.style.display = 'none';
        overlay.classList.remove('hide');
    }, 300);
}

function updateLoadingProgress(completed, total, benchmarkName = '', testName = '') {
    const progressPercent = Math.round((completed / total) * 100);
    const progressFill = document.getElementById('loadingProgressFill');
    const progressPercentSpan = document.getElementById('progressPercent');
    const completedTestsSpan = document.getElementById('completedTests');
    const totalTestsSpan = document.getElementById('totalTests');
    const benchmarkStatus = document.getElementById('benchmarkStatus');
    const currentBenchmark = document.getElementById('currentBenchmark');
    const currentTest = document.getElementById('currentTest');
    const loadingStats = document.getElementById('loadingStats');
    
    // Update progress bar
    progressFill.style.width = progressPercent + '%';
    progressPercentSpan.textContent = progressPercent + '%';
    
    // Update stats
    completedTestsSpan.textContent = completed;
    totalTestsSpan.textContent = total;
    
    // Show stats and benchmark info after first test
    if (completed > 0) {
        loadingStats.style.display = 'grid';
        benchmarkStatus.style.display = 'block';
        
        if (benchmarkName) {
            currentBenchmark.textContent = benchmarkName;
        }
        if (testName) {
            currentTest.textContent = `Running: ${testName}`;
        }
    }
    
    // Update loading text based on progress
    const loadingText = document.getElementById('loadingText');
    const loadingDetails = document.getElementById('loadingDetails');
    
    if (progressPercent === 0) {
        loadingText.textContent = 'Starting Benchmarks...';
        loadingDetails.textContent = 'Preparing performance tests';
    } else if (progressPercent < 25) {
        loadingText.textContent = 'Running Initial Tests...';
        loadingDetails.textContent = 'Warming up and measuring baseline performance';
    } else if (progressPercent < 75) {
        loadingText.textContent = 'Processing Benchmarks...';
        loadingDetails.textContent = 'Comparing JavaScript vs WebAssembly performance';
    } else if (progressPercent < 100) {
        loadingText.textContent = 'Finalizing Results...';
        loadingDetails.textContent = 'Completing final tests and calculations';
    } else {
        loadingText.textContent = 'Tests Complete!';
        loadingDetails.textContent = 'Generating performance summary and charts';
    }
}

// Show initial loading overlay for WebAssembly
document.addEventListener('DOMContentLoaded', function() {
    showLoadingOverlay('Loading WebAssembly...', 'Downloading and initializing Go WASM module');
    updateLoadingProgress(0, 100);
    
    // Simulate progress for WASM loading
    let wasmProgress = 0;
    const wasmProgressInterval = setInterval(() => {
        wasmProgress += Math.random() * 15;
        if (wasmProgress > 90) wasmProgress = 90;
        updateLoadingProgress(wasmProgress, 100, 'WebAssembly', 'Loading module...');
    }, 200);
    
    // Initialize WebAssembly using shared function
    window.initWasm().then(() => {
        clearInterval(wasmProgressInterval);
        updateLoadingProgress(100, 100, 'WebAssembly', 'Module loaded successfully!');
        
        document.getElementById('status').className = 'status ready';
        document.getElementById('status').textContent = '✅ WebAssembly module loaded and ready!';
        document.getElementById('runBenchmark').disabled = false;
        
        // Hide overlay after a brief delay
        setTimeout(() => {
            hideLoadingOverlay();
        }, 1500);
        
    }).catch((err) => {
        clearInterval(wasmProgressInterval);
        console.error("Failed to load WebAssembly:", err);
        document.getElementById('status').className = 'status error';
        document.getElementById('status').textContent = '❌ Failed to load WebAssembly: ' + err.message;
        
        // Update overlay to show error
        document.getElementById('loadingText').textContent = 'Error Loading WebAssembly';
        document.getElementById('loadingDetails').textContent = err.message;
        
        setTimeout(() => {
            hideLoadingOverlay();
        }, 3000);
    });
});

// Benchmark functions
const benchmarks = [
    {
        name: 'Matrix Multiplication',
        icon: '🔢',
        tests: [
            { name: 'JavaScript', fn: 'matrixMultiplyJSOptimized' },
            { name: 'Single-Thread WASM', fn: 'matrixMultiplyWasm' },
            { name: 'Concurrent WASM', fn: 'matrixMultiplyConcurrentWasm' }
        ],
        setup: () => {
            const size = parseInt(document.getElementById('matrixSize').value);
            const matrixA = new Array(size * size);
            const matrixB = new Array(size * size);
            for (let i = 0; i < size * size; i++) {
                matrixA[i] = Math.random();
                matrixB[i] = Math.random();
            }
            return { matrixA, matrixB, size };
        }
    },
    {
        name: 'Mandelbrot Set',
        icon: '🌀',
        tests: [
            { name: 'JavaScript', fn: 'mandelbrotJSOptimized' },
            { name: 'Single-Thread WASM', fn: 'mandelbrotWasm' },
            { name: 'Concurrent WASM', fn: 'mandelbrotConcurrentWasm' }
        ],
        setup: () => {
            const sizeStr = document.getElementById('mandelbrotSize').value;
            const width = parseInt(sizeStr.replace(/\D/g, ''));
            
            if (isNaN(width) || width <= 0) {
                throw new Error('Invalid dimensions');
            }
            
            // Calculate height based on 4:3 aspect ratio
            const height = Math.floor(width * 0.75);
            
            return { 
                width, 
                height,
                xmin: -2.5, 
                xmax: 1.5, 
                ymin: -1.5, 
                ymax: 1.5,
                maxIter: 150
            };
        }
    },
    {
        name: 'Cryptographic Hash',
        icon: '🔐',
        tests: [
            { name: 'JavaScript', fn: 'sha256HashJSOptimized' },
            { name: 'Single-Thread WASM', fn: 'sha256HashWasm' },
            { name: 'Concurrent WASM', fn: 'sha256HashConcurrentWasm' }
        ],
        setup: () => {
            const iterations = parseInt(document.getElementById('hashIterations').value);
            const data = 'The quick brown fox jumps over the lazy dog. '.repeat(10);
            return { data, iterations };
        }
    },
    {
        name: 'Ray Tracing',
        icon: '🎨',
        tests: [
            { name: 'JavaScript', fn: 'rayTracingJSOptimized' },
            { name: 'Single-Thread WASM', fn: 'rayTracingWasm' },
            { name: 'Concurrent WASM', fn: 'rayTracingConcurrentWasm' }
        ],
        setup: () => {
            return { width: 200, height: 150, samples: 10 };
        }
    }
];

async function runAllBenchmarks() {
    if (!window.isWasmReady()) {
        alert('WebAssembly module not loaded yet. Please wait.');
        return;
    }

    // Show loading overlay
    showLoadingOverlay('Starting Benchmarks...', 'Preparing performance tests');

    // Clear previous results
    document.getElementById('benchmarkResults').innerHTML = '';
    document.getElementById('summary').style.display = 'none';
    // Clear the benchmarkResults object
    Object.keys(benchmarkResults).forEach(key => delete benchmarkResults[key]);

    // Hide normal progress bar (we'll use overlay progress)
    const progressBar = document.getElementById('progressBar');
    progressBar.style.display = 'none';
    
    const iterations = parseInt(document.getElementById('iterations').value);
    const totalTests = benchmarks.length * 3 * iterations;
    let completedTests = 0;
    
    // Initialize progress
    updateLoadingProgress(0, totalTests);

    function updateProgress(benchmarkName = '', testName = '') {
        completedTests++;
        updateLoadingProgress(completedTests, totalTests, benchmarkName, testName);
    }

    // Run each benchmark
    for (const benchmark of benchmarks) {
        const card = createBenchmarkCard(benchmark.name, benchmark.icon);
        document.getElementById('benchmarkResults').appendChild(card);
        
        const params = benchmark.setup();
        const results = [];

        // Run each test variant
        for (const test of benchmark.tests) {
            const times = [];
            const fn = window[test.fn];
            
            if (!fn) {
                console.error(`Function ${test.fn} not found`);
                results.push({ name: test.name, avg: 0, times: [] });
                // Still update progress for missing functions
                for (let i = 0; i < iterations; i++) {
                    updateProgress(benchmark.name, test.name);
                }
                continue;
            }

            // Warm-up run
            try {
                if (benchmark.name === 'Matrix Multiplication') {
                    fn(params.matrixA, params.matrixB, params.size);
                } else if (benchmark.name === 'Mandelbrot Set') {
                    fn(params.width, params.height, params.xmin, params.xmax, params.ymin, params.ymax, params.maxIter);
                } else if (benchmark.name === 'Cryptographic Hash') {
                    fn(params.data, params.iterations);
                } else if (benchmark.name === 'Ray Tracing') {
                    fn(params.width, params.height, params.samples);
                }
            } catch (e) {
                console.error('Warm-up error:', e);
            }

            // Timed runs
            for (let i = 0; i < iterations; i++) {
                try {
                    const start = performance.now();
                    
                    if (benchmark.name === 'Matrix Multiplication') {
                        fn(params.matrixA, params.matrixB, params.size);
                    } else if (benchmark.name === 'Mandelbrot Set') {
                        fn(params.width, params.height, params.xmin, params.xmax, params.ymin, params.ymax, params.maxIter);
                    } else if (benchmark.name === 'Cryptographic Hash') {
                        fn(params.data, params.iterations);
                    } else if (benchmark.name === 'Ray Tracing') {
                        fn(params.width, params.height, params.samples);
                    }
                    
                    const time = performance.now() - start;
                    times.push(time);
                    updateProgress(benchmark.name, test.name);
                    
                    // Add small delay to make progress visible
                    await new Promise(resolve => setTimeout(resolve, 10));
                } catch (e) {
                    console.error('Test error:', e);
                    updateProgress(benchmark.name, test.name);
                }
            }

            const avg = times.reduce((a, b) => a + b, 0) / times.length;
            results.push({ name: test.name, avg, times });
        }

        // Update card with results
        updateBenchmarkCard(card, results);
        benchmarkResults[benchmark.name] = results;
    }

    // Complete progress and show summary
    updateLoadingProgress(totalTests, totalTests, 'Complete!', 'Generating summary...');
    
    // Small delay before showing summary
    await new Promise(resolve => setTimeout(resolve, 500));
    
    showSummary();
    
    // Hide loading overlay
    setTimeout(() => {
        hideLoadingOverlay();
    }, 1000);
}

function createBenchmarkCard(name, icon) {
    const card = document.createElement('div');
    card.className = 'benchmark-card';
    card.innerHTML = `
        <h2>${icon} ${name}</h2>
        <div class="benchmark-results">
            <div class="result-row">
                <div class="result-label">Implementation</div>
                <div class="result-time">Avg Time</div>
                <div class="result-speedup">vs JavaScript</div>
            </div>
        </div>
    `;
    return card;
}

function updateBenchmarkCard(card, results) {
    const resultsDiv = card.querySelector('.benchmark-results');
    const jsTime = results.find(r => r.name === 'JavaScript').avg;

    results.forEach(result => {
        const speedup = jsTime / result.avg;
        const row = document.createElement('div');
        row.className = 'result-row';
        
        let speedupClass = '';
        let speedupText = '';
        
        if (result.name === 'JavaScript') {
            speedupText = '(baseline)';
        } else {
            speedupClass = speedup > 1 ? 'speedup-positive' : 'speedup-negative';
            speedupText = speedup > 1 ? 
                `${speedup.toFixed(2)}x faster` : 
                `${(1/speedup).toFixed(2)}x slower`;
        }

        row.innerHTML = `
            <div class="result-label">${result.name}</div>
            <div class="result-time">${result.avg.toFixed(1)}ms</div>
            <div class="result-speedup ${speedupClass}">${speedupText}</div>
        `;
        resultsDiv.appendChild(row);
    });
}

function showSummary() {
    const summaryDiv = document.getElementById('summary');
    const summaryGrid = document.getElementById('summaryGrid');
    summaryDiv.style.display = 'block';
    summaryGrid.innerHTML = '';

    // Calculate overall statistics
    const stats = {
        singleThreadSpeedup: [],
        concurrentSpeedup: [],
        concurrentVsSingle: []
    };

    Object.values(benchmarkResults).forEach(results => {
        // Ensure results is an array
        if (!Array.isArray(results)) {
            console.error('Invalid results format:', results);
            return;
        }
        
        const jsResult = results.find(r => r.name === 'JavaScript');
        const singleResult = results.find(r => r.name === 'Single-Thread WASM');
        const concurrentResult = results.find(r => r.name === 'Concurrent WASM');
        
        if (jsResult && singleResult && concurrentResult) {
            const js = jsResult.avg;
            const single = singleResult.avg;
            const concurrent = concurrentResult.avg;

            stats.singleThreadSpeedup.push(js / single);
            stats.concurrentSpeedup.push(js / concurrent);
            stats.concurrentVsSingle.push(single / concurrent);
        }
    });

    // Create summary cards
    const summaryData = [
        {
            title: 'Avg Single-Thread Speedup',
            value: average(stats.singleThreadSpeedup).toFixed(2) + 'x',
            color: '#2196F3'
        },
        {
            title: 'Avg Concurrent Speedup',
            value: average(stats.concurrentSpeedup).toFixed(2) + 'x',
            color: '#4CAF50'
        },
        {
            title: 'Concurrent vs Single-Thread',
            value: average(stats.concurrentVsSingle).toFixed(2) + 'x',
            color: '#FF9800'
        },
        {
            title: 'Best Improvement',
            value: (stats.concurrentSpeedup.length > 0 ? Math.max(...stats.concurrentSpeedup).toFixed(2) : '0.00') + 'x',
            color: '#9C27B0'
        }
    ];

    summaryData.forEach(data => {
        const card = document.createElement('div');
        card.className = 'summary-card';
        card.innerHTML = `
            <h3>${data.title}</h3>
            <div class="summary-value" style="color: ${data.color}">${data.value}</div>
        `;
        summaryGrid.appendChild(card);
    });
}

function average(arr) {
    if (!arr || arr.length === 0) return 0;
    return arr.reduce((a, b) => a + b, 0) / arr.length;
}

// JavaScript implementations are now in shared-benchmarks.js