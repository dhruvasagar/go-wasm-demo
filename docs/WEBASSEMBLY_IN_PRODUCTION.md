## Major Companies Using WebAssembly in Production

### **Design & Creative Tools**

**1. Figma** 🎨
- Because their product is written in C++, which can easily be compiled into WebAssembly, Figma is a perfect demonstration of this new format's power. If you haven't used Figma before, it's a browser-based interface design tool with a powerful 2D WebGL rendering engine that supports very large documents.
- Impressively, Figma reports that Wasm has cut load times by a factor of 3x for large documents on the web.
- **Impact**: Figma uses WebAssembly to run their C++ rendering engine directly in the browser, enabling desktop-quality design tools on the web

**2. Adobe Creative Suite** 🖼️
- One notable Wasm success story is Adobe's web-based versions of its popular Acrobat, Photoshop, Lightroom, and Express applications.
- Adobe Premiere Rush uses WebAssembly to enable powerful video editing directly in the browser.
- **Impact**: Bringing traditionally desktop-only creative applications to the web with near-native performance

**3. AutoCAD** 🏗️
- The AutoCAD web app uses emscripten to port pieces from the > 35 years old native application for AutoCAD, to the web! This is quite notable, as it proves that WebAssembly can bring these large C/C++ codebases using Emscripten to the web, to run large computationally intensive desktop applications on the web!
- As a product created in 1982, AutoCad is older than the Web and monstrous in size - it has a whopping 15M lines of C++ code. It's been wanting to move to the web for quite some time, but rewriting everything with Javascript is just impractical: lots of work and a much slower result product.
- **Impact**: Successfully ported a massive 35+ year old CAD application to run in browsers

### **Gaming & Entertainment**

**4. Unity** 🎮
- One of the most popular real world examples of WebAssembly is in the Unity WebGL exporter. It allows Unity games to run directly in the browser with native-like performance.
- Wasm is faster, smaller and more memory-efficient than asm.js, which are all pain points of the Unity WebGL export. Wasm may not solve all the existing problems, but it certainly improves the platform in all areas.
- **Impact**: Enables complex 3D games to run smoothly in browsers without plugins

### **E-commerce & Business Platforms**

**5. Shopify** 🛒
- Shopify is actively using Wasm outside the browser with Shopfiy Functions, helping developers customize the backend of their online stores.
- At Shopify, we're keeping the flexibility of untrusted Partner code, but executing it on our own infrastructure with WebAssembly.
- Shopify switched to Wasmtime from another WebAssembly engine in July 2021. With the switch, Shopify saw an average execution performance improvement of ~50%.
- **Impact**: Safely runs third-party merchant customizations in their backend infrastructure

### **Cloud & Infrastructure**

**6. Fastly** ☁️
- Fastly, a cloud computing company, integrates WebAssembly into its Compute@Edge platform to execute ultra-fast serverless functions. The Challenge: Businesses needed a way to run custom logic at the edge of the network with minimal latency. How WebAssembly Helped: WebAssembly provided a lightweight, high-speed execution environment for edge computing workloads. The Results: Reduced latency, improved performance, and a more efficient cloud computing infrastructure.
- Fastly also saw a ~50% improvement in execution time. In addition, Fastly saw a 72% to 163% increase in requests-per-second it could serve. Fastly has since served trillions of requests using Wasmtime.
- **Impact**: Powers edge computing with significantly improved performance

### **Other Notable Companies**

**7. Amazon Prime Video** 📺
- WebAssembly is now actively used in production at Amazon Prime, AutoCad, Midokura (a Sony subsidiary), and beyond.

**8. Google** 🔍
- Although Google can't comment on future product plans, Steiner reiterates Wasm's importance to its overall technical strategy: "Google is heavily invested in the Wasm ecosystem, and it plays a central role in our products, existing and new."

**9. DFINITY (Internet Computer)** 🌐
- DFINITY launched the Internet Computer blockchain using Wasmtime in May 2021. Since then, the Internet Computer has executed 1 quintillion (10^18) instructions for over 150,000 smart contracts without any production issues.

## Why These Companies Choose WebAssembly

The common threads across these implementations are:

1. **Performance**: Near-native speed for computationally intensive tasks
2. **Security**: Sandboxed execution for running untrusted code
3. **Portability**: Code reuse across different platforms and environments
4. **Legacy Code Migration**: Bringing existing C/C++ applications to the web
5. **Consistent Behavior**: Predictable performance across different environments

Companies like Figma and Adobe have demonstrated its value for high-performance browser computing. But I do think that the primary use case for WebAssembly will be on the cloud.

These examples show that WebAssembly has moved well beyond experimental status and is now a crucial technology for companies needing high-performance, secure, and portable code execution both in browsers and server environments.
