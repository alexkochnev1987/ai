// Heavy computation worker
self.onmessage = function (e) {
  const { iterations } = e.data;
  let total = 0;

  console.log(`Worker: Starting computation of ${iterations} iterations`);

  // Perform heavy computation
  for (let i = 0; i < iterations; i++) {
    total += i;

    // Report progress every 10M iterations to avoid blocking
    if (i % 10000000 === 0 && i > 0) {
      self.postMessage({
        type: "progress",
        progress: i / iterations,
        currentIteration: i,
      });
    }
  }

  console.log("Worker: Computation completed");

  // Send final result
  self.postMessage({
    type: "result",
    result: total,
  });
};

self.onerror = function (error) {
  console.error("Worker error:", error);
  self.postMessage({
    type: "error",
    error: error.message,
  });
};
