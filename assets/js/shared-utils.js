// ============================================================================
// SHARED UTILITY FUNCTIONS
// Common JavaScript utilities used across multiple pages
// ============================================================================

// ============================================================================
// WEBASSEMBLY INITIALIZATION (Centralized)
// ============================================================================

// Global WebAssembly state
window.wasmReady = false;
window.initWasm = function() {
    return window.initializeWasm();
};

window.isWasmReady = function() {
    return window.wasmReady;
};

// ============================================================================
// COMMON UI UTILITIES
// ============================================================================

// Standardized result display function
function displayResult(elementId, result, title, elapsed = null, category = null) {
    const element = document.getElementById(elementId);
    if (!element) {
        console.error(`Element with ID ${elementId} not found`);
        return;
    }

    let content = `<h4>${title}</h4>`;
    
    if (result.valid !== undefined) {
        // Validation result
        content += `<div class="validation-status ${result.valid ? 'valid' : 'invalid'}">`;
        content += `<strong>Status:</strong> ${result.valid ? '✅ Valid' : '❌ Invalid'}<br>`;
        
        if (result.errors && result.errors.length > 0) {
            content += `<strong>Errors:</strong><ul>`;
            result.errors.forEach(error => {
                content += `<li>${error}</li>`;
            });
            content += `</ul>`;
        }
        content += `</div>`;
    } else {
        // Other results (calculations, recommendations, etc.)
        content += `<pre>${JSON.stringify(result, null, 2)}</pre>`;
    }
    
    if (elapsed !== null) {
        content += `<div class="performance-info">⏱️ Execution time: ${elapsed.toFixed(2)}ms</div>`;
    }
    
    element.className = 'results success';
    element.innerHTML = content;
    
    // Update performance tracking if category provided
    if (category && window.performanceData && elapsed !== null) {
        updatePerformanceComparison(category, 'wasm', elapsed);
    }
}

// Standardized error display function
function displayError(elementId, error, title = 'Error') {
    const element = document.getElementById(elementId);
    if (!element) {
        console.error(`Element with ID ${elementId} not found`);
        return;
    }
    
    element.className = 'results error';
    element.innerHTML = `<h4>${title}</h4><div class="error-message">${error.message || error}</div>`;
}

// ============================================================================
// PERFORMANCE TRACKING UTILITIES
// ============================================================================

// Initialize performance data structure if it doesn't exist
function initializePerformanceData() {
    if (!window.performanceData) {
        window.performanceData = {};
    }
}

// Update performance comparison data
function updatePerformanceComparison(category, type, elapsed) {
    initializePerformanceData();
    
    if (!window.performanceData[category]) {
        window.performanceData[category] = {};
    }
    
    window.performanceData[category][type] = elapsed;
    
    // Auto-update comparison display if both WASM and server times are available
    updatePerformanceDisplay(category);
}

// Display performance comparison
function updatePerformanceDisplay(category) {
    const data = window.performanceData[category];
    if (!data || !data.wasm || !data.server) return;
    
    const comparisonElement = document.getElementById(`${category}Comparison`);
    if (!comparisonElement) return;
    
    const wasmTime = data.wasm;
    const serverTime = data.server;
    const faster = wasmTime < serverTime ? 'WebAssembly' : 'Server';
    const ratio = wasmTime < serverTime ? (serverTime / wasmTime).toFixed(1) : (wasmTime / serverTime).toFixed(1);
    
    comparisonElement.innerHTML = `
        <div class="performance-comparison">
            <strong>🏁 Performance Comparison:</strong><br>
            WebAssembly: ${wasmTime.toFixed(2)}ms<br>
            Server API: ${serverTime.toFixed(2)}ms<br>
            <span class="winner">${faster} is ${ratio}x faster!</span>
        </div>
    `;
}

// ============================================================================
// FORM UTILITIES
// ============================================================================

// Get form values as object
function getFormData(formPrefix) {
    const data = {};
    const inputs = document.querySelectorAll(`[id^="${formPrefix}"]`);
    
    inputs.forEach(input => {
        const key = input.id.replace(formPrefix, '').toLowerCase();
        if (input.type === 'checkbox') {
            data[key] = input.checked;
        } else if (input.type === 'number') {
            data[key] = parseInt(input.value) || 0;
        } else {
            data[key] = input.value;
        }
    });
    
    return data;
}

// Clear form values
function clearForm(formPrefix) {
    const inputs = document.querySelectorAll(`[id^="${formPrefix}"]`);
    inputs.forEach(input => {
        if (input.type === 'checkbox') {
            input.checked = false;
        } else {
            input.value = '';
        }
    });
}

// Fill form with demo data
function fillDemoData(formPrefix, demoData) {
    Object.keys(demoData).forEach(key => {
        const input = document.getElementById(formPrefix + key.charAt(0).toUpperCase() + key.slice(1));
        if (input) {
            if (input.type === 'checkbox') {
                input.checked = demoData[key];
            } else {
                input.value = demoData[key];
            }
        }
    });
}

// ============================================================================
// API UTILITIES
// ============================================================================

// Make API request with error handling
async function makeApiRequest(url, data, method = 'POST') {
    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('API request failed:', error);
        throw error;
    }
}

// ============================================================================
// LOADING OVERLAY UTILITIES
// ============================================================================

function showLoadingOverlay(text = 'Loading...', details = '') {
    const overlay = document.getElementById('loadingOverlay');
    if (!overlay) return;
    
    const loadingText = document.getElementById('loadingText');
    const loadingDetails = document.getElementById('loadingDetails');
    
    if (loadingText) loadingText.textContent = text;
    if (loadingDetails) loadingDetails.textContent = details;
    
    overlay.classList.remove('hide');
    overlay.classList.add('show');
    overlay.style.display = 'flex';
}

function hideLoadingOverlay() {
    const overlay = document.getElementById('loadingOverlay');
    if (!overlay) return;
    
    overlay.classList.remove('show');
    overlay.classList.add('hide');
    
    setTimeout(() => {
        overlay.style.display = 'none';
        overlay.classList.remove('hide');
    }, 300);
}

// ============================================================================
// INITIALIZATION
// ============================================================================

// Initialize shared utilities when DOM is ready
document.addEventListener('DOMContentLoaded', function() {
    initializePerformanceData();
    console.log('✅ Shared utilities initialized');
});