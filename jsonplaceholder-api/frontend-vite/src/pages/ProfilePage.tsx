import React, { useState, useEffect } from "react";
import { useAuth } from "../contexts/AuthContext";
import { authApi } from "../services/api";
import { Card } from "../components/ui/Card";
import { Button } from "../components/ui/Button";
import type { UserResponse } from "../types/api";

export const ProfilePage: React.FC = () => {
  const { user: contextUser, isAuthenticated } = useAuth();
  const [user, setUser] = useState<UserResponse | null>(contextUser);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const fetchProfile = async () => {
    if (!isAuthenticated) return;

    setLoading(true);
    setError("");

    try {
      const response = await authApi.me();
      setUser(response.data);
    } catch (err) {
      console.error("Error fetching profile:", err);
      setError("Failed to load profile");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isAuthenticated && !user) {
      fetchProfile();
    }
  }, [isAuthenticated, user]);

  if (!isAuthenticated) {
    return (
      <div className="text-center">
        <div className="bg-yellow-50 border border-yellow-200 text-yellow-600 px-4 py-3 rounded-md mb-4">
          You must be logged in to view your profile.
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="text-center">
        <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-md mb-4">
          {error || "Profile not found"}
        </div>
        <Button onClick={fetchProfile}>Retry</Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">My Profile</h1>
          <p className="text-gray-600">Your account information</p>
        </div>
        <Button onClick={fetchProfile} loading={loading}>
          Refresh Profile
        </Button>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Personal Information */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Personal Information
          </h2>
          <div className="space-y-3">
            <div>
              <span className="font-medium text-gray-700">Full Name:</span>
              <p className="text-gray-900">{user.name}</p>
            </div>
            <div>
              <span className="font-medium text-gray-700">Username:</span>
              <p className="text-gray-900">@{user.username}</p>
            </div>
            <div>
              <span className="font-medium text-gray-700">Email:</span>
              <p className="text-gray-900">{user.email}</p>
            </div>
            <div>
              <span className="font-medium text-gray-700">User ID:</span>
              <p className="text-gray-900">{user.id}</p>
            </div>
            {user.phone && (
              <div>
                <span className="font-medium text-gray-700">Phone:</span>
                <p className="text-gray-900">{user.phone}</p>
              </div>
            )}
            {user.website && (
              <div>
                <span className="font-medium text-gray-700">Website:</span>
                <a
                  href={
                    user.website.startsWith("http")
                      ? user.website
                      : `https://${user.website}`
                  }
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-blue-600 hover:text-blue-800"
                >
                  {user.website}
                </a>
              </div>
            )}
          </div>
        </Card>

        {/* Address Information */}
        {user.address && (
          <Card>
            <h2 className="text-xl font-semibold text-gray-900 mb-4">
              Address
            </h2>
            <div className="space-y-3">
              <div>
                <span className="font-medium text-gray-700">Street:</span>
                <p className="text-gray-900">
                  {user.address.street} {user.address.suite}
                </p>
              </div>
              <div>
                <span className="font-medium text-gray-700">City:</span>
                <p className="text-gray-900">{user.address.city}</p>
              </div>
              <div>
                <span className="font-medium text-gray-700">Zipcode:</span>
                <p className="text-gray-900">{user.address.zipcode}</p>
              </div>
              {user.address.geo &&
                (user.address.geo.lat || user.address.geo.lng) && (
                  <div>
                    <span className="font-medium text-gray-700">
                      Coordinates:
                    </span>
                    <p className="text-gray-900">
                      {user.address.geo.lat}, {user.address.geo.lng}
                    </p>
                  </div>
                )}
            </div>
          </Card>
        )}

        {/* Company Information */}
        {user.company && (
          <Card className="md:col-span-2">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">
              Company
            </h2>
            <div className="space-y-3">
              <div>
                <span className="font-medium text-gray-700">Company Name:</span>
                <p className="text-gray-900">{user.company.name}</p>
              </div>
              {user.company.catchPhrase && (
                <div>
                  <span className="font-medium text-gray-700">
                    Catch Phrase:
                  </span>
                  <p className="text-gray-900 italic">
                    "{user.company.catchPhrase}"
                  </p>
                </div>
              )}
              {user.company.bs && (
                <div>
                  <span className="font-medium text-gray-700">Business:</span>
                  <p className="text-gray-900">{user.company.bs}</p>
                </div>
              )}
            </div>
          </Card>
        )}
      </div>
    </div>
  );
};
