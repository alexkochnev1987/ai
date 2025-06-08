import React, { useState, useEffect } from "react";
import { healthCheck, getApiInfo } from "../services/api";
import { Card } from "../components/ui/Card";
import { Button } from "../components/ui/Button";
import type { HealthResponse, APIResponse } from "../types/api";

export const HealthPage: React.FC = () => {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [apiInfo, setApiInfo] = useState<APIResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const fetchData = async () => {
    setLoading(true);
    setError("");

    try {
      const [healthResponse, infoResponse] = await Promise.all([
        healthCheck(),
        getApiInfo(),
      ]);

      setHealth(healthResponse);
      setApiInfo(infoResponse);
    } catch (err) {
      console.error("Error fetching health data:", err);
      setError("Failed to fetch API status");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            API Health Status
          </h1>
          <p className="text-gray-600">
            Monitor the status and information of the JSONPlaceholder API
          </p>
        </div>
        <Button onClick={fetchData} loading={loading}>
          Refresh Status
        </Button>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-md">
          {error}
        </div>
      )}

      <div className="grid gap-6 md:grid-cols-2">
        {/* Health Status */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Health Check
          </h2>
          {health ? (
            <div className="space-y-3">
              <div className="flex items-center">
                <span className="font-medium text-gray-700">Status:</span>
                <span
                  className={`ml-2 px-2 py-1 rounded-full text-sm font-medium ${
                    health.status === "ok"
                      ? "bg-green-100 text-green-800"
                      : "bg-red-100 text-red-800"
                  }`}
                >
                  {health.status.toUpperCase()}
                </span>
              </div>
              <div>
                <span className="font-medium text-gray-700">Message:</span>
                <p className="text-gray-900">{health.message}</p>
              </div>
              <div>
                <span className="font-medium text-gray-700">Version:</span>
                <p className="text-gray-900">{health.version}</p>
              </div>
            </div>
          ) : (
            <p className="text-gray-500">Health data not available</p>
          )}
        </Card>

        {/* API Information */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            API Information
          </h2>
          {apiInfo ? (
            <div className="space-y-3">
              <div>
                <span className="font-medium text-gray-700">Message:</span>
                <p className="text-gray-900">{apiInfo.message}</p>
              </div>
              <div>
                <span className="font-medium text-gray-700">Success:</span>
                <span
                  className={`ml-2 px-2 py-1 rounded-full text-sm font-medium ${
                    apiInfo.success
                      ? "bg-green-100 text-green-800"
                      : "bg-red-100 text-red-800"
                  }`}
                >
                  {apiInfo.success ? "TRUE" : "FALSE"}
                </span>
              </div>
            </div>
          ) : (
            <p className="text-gray-500">API info not available</p>
          )}
        </Card>
      </div>

      {/* API Endpoints */}
      {apiInfo?.data?.endpoints && (
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Available Endpoints
          </h2>
          <div className="space-y-6">
            {/* Authentication Endpoints */}
            {apiInfo.data.endpoints.authentication && (
              <div>
                <h3 className="text-lg font-medium text-gray-900 mb-3">
                  Authentication
                </h3>
                <div className="space-y-2">
                  {Object.entries(apiInfo.data.endpoints.authentication).map(
                    ([endpoint, description]) => (
                      <div
                        key={endpoint}
                        className="flex justify-between items-center p-3 bg-gray-50 rounded-md"
                      >
                        <code className="text-sm font-mono text-blue-600">
                          {endpoint}
                        </code>
                        <span className="text-sm text-gray-600">
                          {description as string}
                        </span>
                      </div>
                    )
                  )}
                </div>
              </div>
            )}

            {/* Users Endpoints */}
            {apiInfo.data.endpoints.users && (
              <div>
                <h3 className="text-lg font-medium text-gray-900 mb-3">
                  Users
                </h3>
                <div className="space-y-2">
                  {Object.entries(apiInfo.data.endpoints.users).map(
                    ([endpoint, description]) => (
                      <div
                        key={endpoint}
                        className="flex justify-between items-center p-3 bg-gray-50 rounded-md"
                      >
                        <code className="text-sm font-mono text-blue-600">
                          {endpoint}
                        </code>
                        <span className="text-sm text-gray-600">
                          {description as string}
                        </span>
                      </div>
                    )
                  )}
                </div>
              </div>
            )}
          </div>
        </Card>
      )}

      {/* Additional API Data */}
      {apiInfo?.data?.version && (
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Additional Information
          </h2>
          <div className="space-y-3">
            <div>
              <span className="font-medium text-gray-700">API Version:</span>
              <p className="text-gray-900">{apiInfo.data.version}</p>
            </div>
            {apiInfo.data.documentation && (
              <div>
                <span className="font-medium text-gray-700">
                  Documentation:
                </span>
                <p className="text-gray-900">{apiInfo.data.documentation}</p>
              </div>
            )}
          </div>
        </Card>
      )}
    </div>
  );
};
