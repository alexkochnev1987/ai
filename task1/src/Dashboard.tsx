import { useState, useEffect, useRef } from "react";

export default function Dashboard() {
  const [result, setResult] = useState<number | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const workerRef = useRef<Worker | null>(null);

  useEffect(() => {
    // Initialize Web Worker
    setIsLoading(true);

    // Create worker from inline script to avoid separate file complexity
    const workerScript = `
      self.onmessage = function(e) {
        const { iterations } = e.data;
        let total = 0;
        
        // Perform heavy computation
        for (let i = 0; i < iterations; i++) {
          total += i;
          
          // Periodically report progress to avoid blocking
          if (i % 10000000 === 0) {
            self.postMessage({ type: 'progress', progress: i / iterations });
          }
        }
        
        // Send final result
        self.postMessage({ type: 'result', result: total });
      };
    `;

    const blob = new Blob([workerScript], { type: "application/javascript" });
    const workerUrl = URL.createObjectURL(blob);

    workerRef.current = new Worker(workerUrl);

    // Set up message handler
    workerRef.current.onmessage = (e) => {
      const { type, result: workerResult, progress } = e.data;

      if (type === "result") {
        setResult(workerResult);
        setIsLoading(false);
      } else if (type === "progress") {
        // Could update a progress bar here
        console.log(`Computation progress: ${Math.round(progress * 100)}%`);
      }
    };

    // Handle worker errors
    workerRef.current.onerror = (error) => {
      console.error("Worker error:", error);
      setIsLoading(false);
    };

    // Start the computation
    workerRef.current.postMessage({ iterations: 1e8 });

    // Cleanup function
    return () => {
      if (workerRef.current) {
        workerRef.current.terminate();
        URL.revokeObjectURL(workerUrl);
      }
    };
  }, []);

  return (
    <div className="p-8 max-w-md mx-auto bg-white rounded-xl shadow-lg">
      <h2 className="text-2xl font-bold text-gray-800 mb-4">Dashboard</h2>

      {isLoading ? (
        <div className="flex items-center space-x-2">
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
          <span className="text-gray-600">Computing heavy calculation...</span>
        </div>
      ) : (
        <div className="space-y-2">
          <p className="text-gray-600">Computation Result:</p>
          <p className="text-3xl font-mono text-blue-600 bg-gray-50 p-3 rounded">
            {result?.toLocaleString()}
          </p>
        </div>
      )}

      <div className="mt-4 text-sm text-gray-500">
        <p>✅ Non-blocking computation using Web Worker</p>
        <p>✅ UI remains responsive during calculation</p>
      </div>
    </div>
  );
}
