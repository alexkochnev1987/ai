import React, { useState, useEffect, useRef } from "react";

export default function DashboardWithExternalWorker() {
  const [result, setResult] = useState<number | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [progress, setProgress] = useState(0);
  const workerRef = useRef<Worker | null>(null);

  useEffect(() => {
    setIsLoading(true);

    // Use external worker file
    workerRef.current = new Worker("/computation-worker.js");

    // Set up message handler
    workerRef.current.onmessage = (e) => {
      const {
        type,
        result: workerResult,
        progress: workerProgress,
        error,
      } = e.data;

      if (type === "result") {
        setResult(workerResult);
        setIsLoading(false);
        setProgress(100);
      } else if (type === "progress") {
        setProgress(Math.round(workerProgress * 100));
      } else if (type === "error") {
        console.error("Worker error:", error);
        setIsLoading(false);
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
      }
    };
  }, []);

  return (
    <div className="p-8 max-w-md mx-auto bg-white rounded-xl shadow-lg">
      <h2 className="text-2xl font-bold text-gray-800 mb-4">
        Dashboard (External Worker)
      </h2>

      {isLoading ? (
        <div className="space-y-4">
          <div className="flex items-center space-x-2">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
            <span className="text-gray-600">
              Computing heavy calculation...
            </span>
          </div>

          {/* Progress bar */}
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div
              className="bg-blue-600 h-2 rounded-full transition-all duration-300"
              style={{ width: `${progress}%` }}
            ></div>
          </div>
          <p className="text-sm text-gray-500 text-center">
            {progress}% complete
          </p>
        </div>
      ) : (
        <div className="space-y-2">
          <p className="text-gray-600">Computation Result:</p>
          <p className="text-3xl font-mono text-green-600 bg-gray-50 p-3 rounded">
            {result?.toLocaleString()}
          </p>
        </div>
      )}

      <div className="mt-4 text-sm text-gray-500">
        <p>✅ External Web Worker file</p>
        <p>✅ Progress tracking</p>
        <p>✅ Zero main thread blocking</p>
      </div>
    </div>
  );
}
