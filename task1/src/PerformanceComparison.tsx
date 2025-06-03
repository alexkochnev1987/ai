import React, { useState } from "react";

export default function PerformanceComparison() {
  const [result, setResult] = useState<number | null>(null);
  const [isBlocking, setIsBlocking] = useState(false);
  const [clickCount, setClickCount] = useState(0);

  // BAD: Blocking computation on main thread
  const runBlockingComputation = () => {
    setIsBlocking(true);
    setResult(null);

    // This will freeze the UI!
    setTimeout(() => {
      let total = 0;
      for (let i = 0; i < 1e8; i++) {
        total += i;
      }
      setResult(total);
      setIsBlocking(false);
    }, 10); // Small delay to show loading state
  };

  // Test UI responsiveness
  const handleTestClick = () => {
    setClickCount((prev) => prev + 1);
  };

  return (
    <div className="p-8 max-w-md mx-auto bg-red-50 border-2 border-red-200 rounded-xl">
      <h2 className="text-2xl font-bold text-red-800 mb-4">
        ⚠️ BAD Example (Main Thread)
      </h2>

      <div className="space-y-4">
        <button
          onClick={runBlockingComputation}
          disabled={isBlocking}
          className="w-full py-2 px-4 bg-red-600 hover:bg-red-700 disabled:bg-gray-400 text-white rounded font-semibold"
        >
          {isBlocking
            ? "Computing... (UI FROZEN!)"
            : "Run Blocking Computation"}
        </button>

        {/* Test button to show UI freezing */}
        <div className="border-t pt-4">
          <p className="text-sm text-gray-600 mb-2">Test UI responsiveness:</p>
          <button
            onClick={handleTestClick}
            className="py-2 px-4 bg-blue-500 hover:bg-blue-600 text-white rounded mr-2"
          >
            Click Me! ({clickCount})
          </button>
          <span className="text-sm text-gray-600">
            {isBlocking
              ? "🔴 Try clicking - UI is frozen!"
              : "🟢 UI responsive"}
          </span>
        </div>

        {result && (
          <div className="bg-white p-3 rounded border">
            <p className="text-sm text-gray-600">Result:</p>
            <p className="font-mono text-lg text-red-600">
              {result.toLocaleString()}
            </p>
          </div>
        )}
      </div>

      <div className="mt-4 text-xs text-red-600 bg-red-100 p-2 rounded">
        <p>❌ Blocks main thread for 3+ seconds</p>
        <p>❌ UI becomes completely unresponsive</p>
        <p>❌ High Total Blocking Time</p>
      </div>
    </div>
  );
}
