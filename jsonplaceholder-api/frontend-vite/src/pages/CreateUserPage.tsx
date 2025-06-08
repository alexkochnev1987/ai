import React, { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { usersApi } from "../services/api";
import { Card } from "../components/ui/Card";
import { Input } from "../components/ui/Input";
import { Button } from "../components/ui/Button";
import type { CreateUserRequest } from "../types/api";

export const CreateUserPage: React.FC = () => {
  const [formData, setFormData] = useState<CreateUserRequest>({
    name: "",
    username: "",
    email: "",
    password: "",
    phone: "",
    website: "",
    address: {
      street: "",
      suite: "",
      city: "",
      zipcode: "",
      geo: {
        lat: "",
        lng: "",
      },
    },
    company: {
      name: "",
      catchPhrase: "",
      bs: "",
    },
  });
  const [errors, setErrors] = useState<Partial<CreateUserRequest>>({});
  const [loading, setLoading] = useState(false);
  const [apiError, setApiError] = useState("");

  const navigate = useNavigate();

  const validateForm = (): boolean => {
    const newErrors: Partial<CreateUserRequest> = {};

    if (!formData.name) {
      newErrors.name = "Name is required";
    } else if (formData.name.length < 2) {
      newErrors.name = "Name must be at least 2 characters";
    }

    if (!formData.username) {
      newErrors.username = "Username is required";
    } else if (formData.username.length < 3) {
      newErrors.username = "Username must be at least 3 characters";
    }

    if (!formData.email) {
      newErrors.email = "Email is required";
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      newErrors.email = "Email is invalid";
    }

    if (!formData.password) {
      newErrors.password = "Password is required";
    } else if (formData.password.length < 6) {
      newErrors.password = "Password must be at least 6 characters";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setApiError("");

    if (!validateForm()) return;

    setLoading(true);
    try {
      const response = await usersApi.createUser(formData);
      navigate(`/users/${response.data.id}`);
    } catch (error) {
      console.error("Create user error:", error);
      setApiError("Failed to create user. Please try again.");
    } finally {
      setLoading(false);
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
            ...prev.address,
            geo: {
              ...prev.address.geo,
              [addressField]: value,
            },
          },
        }));
      } else {
        setFormData((prev) => ({
          ...prev,
          address: {
            ...prev.address,
            [addressField]: value,
          },
        }));
      }
    } else if (name.startsWith("company.")) {
      const companyField = name.split(".")[1];
      setFormData((prev) => ({
        ...prev,
        company: {
          ...prev.company,
          [companyField]: value,
        },
      }));
    } else {
      setFormData((prev) => ({ ...prev, [name]: value }));
    }

    // Clear error when user starts typing
    if (errors[name as keyof CreateUserRequest]) {
      setErrors((prev) => ({ ...prev, [name]: "" }));
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Create New User</h1>
          <p className="text-gray-600">Add a new user to the system</p>
        </div>
        <Link to="/users">
          <Button variant="outline">Back to Users</Button>
        </Link>
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
              label="Full Name *"
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              error={errors.name}
              fullWidth
              required
            />

            <Input
              label="Username *"
              type="text"
              name="username"
              value={formData.username}
              onChange={handleChange}
              error={errors.username}
              fullWidth
              required
            />

            <Input
              label="Email *"
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              error={errors.email}
              fullWidth
              required
            />

            <Input
              label="Password *"
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              error={errors.password}
              fullWidth
              required
            />

            <Input
              label="Phone"
              type="tel"
              name="phone"
              value={formData.phone}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Website"
              type="url"
              name="website"
              value={formData.website}
              onChange={handleChange}
              fullWidth
            />
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
              value={formData.address.street}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Suite"
              type="text"
              name="address.suite"
              value={formData.address.suite}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="City"
              type="text"
              name="address.city"
              value={formData.address.city}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Zipcode"
              type="text"
              name="address.zipcode"
              value={formData.address.zipcode}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Latitude"
              type="text"
              name="address.lat"
              value={formData.address.geo.lat}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Longitude"
              type="text"
              name="address.lng"
              value={formData.address.geo.lng}
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
              value={formData.company.name}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Catch Phrase"
              type="text"
              name="company.catchPhrase"
              value={formData.company.catchPhrase}
              onChange={handleChange}
              fullWidth
            />

            <Input
              label="Business"
              type="text"
              name="company.bs"
              value={formData.company.bs}
              onChange={handleChange}
              fullWidth
            />
          </div>
        </Card>

        <div className="flex justify-end space-x-4">
          <Link to="/users">
            <Button variant="outline">Cancel</Button>
          </Link>
          <Button type="submit" loading={loading}>
            Create User
          </Button>
        </div>
      </form>
    </div>
  );
};
