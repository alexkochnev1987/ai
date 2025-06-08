import React, { useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { usersApi } from "../services/api";
import { Card } from "../components/ui/Card";
import { Input } from "../components/ui/Input";
import { Button } from "../components/ui/Button";
import type { UpdateUserRequest, UserResponse } from "../types/api";

export const EditUserPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [user, setUser] = useState<UserResponse | null>(null);
  const [formData, setFormData] = useState<UpdateUserRequest>({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [apiError, setApiError] = useState("");

  const navigate = useNavigate();

  useEffect(() => {
    const fetchUser = async () => {
      if (!id) return;

      setLoading(true);
      setError("");

      try {
        const response = await usersApi.getUser(parseInt(id));
        const userData = response.data;
        setUser(userData);

        // Initialize form with current user data
        setFormData({
          name: userData.name,
          username: userData.username,
          email: userData.email,
          phone: userData.phone,
          website: userData.website,
          address: userData.address,
          company: userData.company,
        });
      } catch (err) {
        console.error("Error fetching user:", err);
        setError("Failed to load user");
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [id]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) return;

    setApiError("");
    setSaving(true);

    try {
      const response = await usersApi.updateUser(user.id, formData);
      navigate(`/users/${response.data.id}`);
    } catch (error) {
      console.error("Update user error:", error);
      setApiError("Failed to update user. Please try again.");
    } finally {
      setSaving(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    if (name.startsWith("address.")) {
      const addressField = name.split(".")[1];
      if (addressField === "lat" || addressField === "lng") {
        setFormData((prev) => ({
          ...prev,
          address: {
            ...prev.address!,
            geo: {
              ...prev.address!.geo,
              [addressField]: value,
            },
          },
        }));
      } else {
        setFormData((prev) => ({
          ...prev,
          address: {
            ...prev.address!,
            [addressField]: value,
          },
        }));
      }
    } else if (name.startsWith("company.")) {
      const companyField = name.split(".")[1];
      setFormData((prev) => ({
        ...prev,
        company: {
          ...prev.company!,
          [companyField]: value,
        },
      }));
    } else {
      setFormData((prev) => ({ ...prev, [name]: value }));
    }
  };

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
          {error || "User not found"}
        </div>
        <Link to="/users">
          <Button>Back to Users</Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Edit User</h1>
          <p className="text-gray-600">Update {user.name}'s information</p>
        </div>
        <div className="flex space-x-2">
          <Link to={`/users/${user.id}`}>
            <Button variant="outline">View User</Button>
          </Link>
          <Link to="/users">
            <Button variant="outline">Back to Users</Button>
          </Link>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {apiError && (
          <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-md">
            {apiError}
          </div>
        )}

        {/* Personal Information */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Personal Information
          </h2>
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <Input
              label="Full Name"
              type="text"
              name="name"
              value={formData.name || ""}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Username"
              type="text"
              name="username"
              value={formData.username || ""}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Email"
              type="email"
              name="email"
              value={formData.email || ""}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Phone"
              type="tel"
              name="phone"
              value={formData.phone || ""}
              onChange={handleChange}
              fullWidth
            />

            <div className="sm:col-span-2">
              <Input
                label="Website"
                type="url"
                name="website"
                value={formData.website || ""}
                onChange={handleChange}
                fullWidth
              />
            </div>
          </div>
        </Card>

        {/* Address Information */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Address Information
          </h2>
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <Input
              label="Street"
              type="text"
              name="address.street"
              value={formData.address?.street || ""}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Suite"
              type="text"
              name="address.suite"
              value={formData.address?.suite || ""}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="City"
              type="text"
              name="address.city"
              value={formData.address?.city || ""}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Zipcode"
              type="text"
              name="address.zipcode"
              value={formData.address?.zipcode || ""}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Latitude"
              type="text"
              name="address.lat"
              value={formData.address?.geo?.lat || ""}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Longitude"
              type="text"
              name="address.lng"
              value={formData.address?.geo?.lng || ""}
              onChange={handleChange}
              fullWidth
            />
          </div>
        </Card>

        {/* Company Information */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Company Information
          </h2>
          <div className="grid grid-cols-1 gap-6">
            <Input
              label="Company Name"
              type="text"
              name="company.name"
              value={formData.company?.name || ""}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Catch Phrase"
              type="text"
              name="company.catchPhrase"
              value={formData.company?.catchPhrase || ""}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Business"
              type="text"
              name="company.bs"
              value={formData.company?.bs || ""}
              onChange={handleChange}
              fullWidth
            />
          </div>
        </Card>

        <div className="flex justify-end space-x-4">
          <Link to={`/users/${user.id}`}>
            <Button variant="outline">Cancel</Button>
          </Link>
          <Button type="submit" loading={saving}>
            Update User
          </Button>
        </div>
      </form>
    </div>
  );
};
