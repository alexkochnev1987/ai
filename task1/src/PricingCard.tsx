import React from "react";

interface PricingCardProps {
  plan: string;
  price: string;
  features: string[];
  isFeatured?: boolean;
}

const PricingCard: React.FC<PricingCardProps> = ({
  plan,
  price,
  features,
  isFeatured = false,
}) => {
  return (
    <div className="w-full max-w-sm mx-auto">
      <div
        className={`
          relative p-8 rounded-lg shadow-lg transition-all duration-300 
          hover:shadow-xl hover:-translate-y-2 focus-within:ring-4 focus-within:ring-blue-500 focus-within:ring-opacity-50
          ${
            isFeatured
              ? "bg-slate-700 text-white border-2 border-blue-500 scale-105"
              : "bg-white text-gray-800 border border-gray-200"
          }
        `}
        tabIndex={0}
      >
        {isFeatured && (
          <div className="absolute -top-4 left-1/2 transform -translate-x-1/2">
            <span className="bg-blue-500 text-white px-4 py-2 rounded-full text-sm font-semibold">
              Most Popular
            </span>
          </div>
        )}

        <div className="text-center">
          <h3
            className={`text-xl font-semibold mb-4 ${
              isFeatured ? "text-white" : "text-gray-800"
            }`}
          >
            {plan}
          </h3>

          <div className="mb-6">
            <span
              className={`text-5xl font-bold ${
                isFeatured ? "text-white" : "text-gray-900"
              }`}
            >
              {price}
            </span>
          </div>

          <div className="space-y-4 mb-8">
            {features.map((feature, index) => (
              <div
                key={index}
                className={`text-sm ${
                  isFeatured ? "text-gray-200" : "text-gray-600"
                } border-b ${
                  isFeatured ? "border-gray-600" : "border-gray-200"
                } pb-2`}
              >
                {feature}
              </div>
            ))}
          </div>

          <button
            className={`
              w-full py-3 px-6 rounded-lg font-semibold text-sm uppercase tracking-wide
              transition-all duration-200 focus:outline-none focus:ring-4 focus:ring-opacity-50
              ${
                isFeatured
                  ? "bg-blue-600 hover:bg-blue-700 text-white focus:ring-blue-500"
                  : "bg-gray-200 hover:bg-gray-300 text-gray-800 focus:ring-gray-500"
              }
            `}
          >
            Subscribe
          </button>
        </div>
      </div>
    </div>
  );
};

export default PricingCard;
