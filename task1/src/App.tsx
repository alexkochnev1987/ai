import PricingCard from "./PricingCard";
import Dashboard from "./Dashboard";
import DashboardWithExternalWorker from "./DashboardWithExternalWorker";
import PerformanceComparison from "./PerformanceComparison";

function App() {
  return (
    <div className="bg-gray-900 min-h-screen py-12">
      <div className="max-w-7xl mx-auto px-4">
        <h1 className="text-4xl font-bold text-white text-center mb-12">
          Performance Optimization Demo
        </h1>

        {/* Performance comparison: Bad vs Good */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-12">
          <PerformanceComparison />
          <Dashboard />
          <DashboardWithExternalWorker />
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8 justify-items-center">
          <PricingCard
            plan="Standard"
            price="$100"
            features={[
              "50,000 Requests",
              "4 contributors",
              "Up to 3 GB storage space",
            ]}
          />
          <PricingCard
            plan="Pro"
            price="$200"
            features={[
              "100,000 Requests",
              "7 contributors",
              "Up to 6 GB storage space",
            ]}
            isFeatured={true}
          />
          <PricingCard
            plan="Expert"
            price="$500"
            features={[
              "200,000 Requests",
              "11 contributors",
              "Up to 10 GB storage space",
            ]}
          />
        </div>
      </div>
    </div>
  );
}

export default App;
