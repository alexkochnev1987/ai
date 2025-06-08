import React, { useState, useEffect } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { usersApi } from "../services/api";
import { useAuth } from "../contexts/AuthContext";
import { Button } from "../components/ui/Button";
import type { UserResponse } from "../types/api";
import styles from "./UserDetailPage.module.css";

export const UserDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [user, setUser] = useState<UserResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [deleteLoading, setDeleteLoading] = useState(false);

  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchUser = async () => {
      if (!id) return;

      setLoading(true);
      setError("");

      try {
        const response = await usersApi.getUser(parseInt(id));
        setUser(response.data);
      } catch (err) {
        console.error("Error fetching user:", err);
        setError("Failed to load user");
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [id]);

  const handleDelete = async () => {
    if (!user || !isAuthenticated) return;

    if (!confirm("Are you sure you want to delete this user?")) {
      return;
    }

    setDeleteLoading(true);
    try {
      await usersApi.deleteUser(user.id);
      navigate("/users");
    } catch (err) {
      console.error("Error deleting user:", err);
      alert("Failed to delete user");
    } finally {
      setDeleteLoading(false);
    }
  };

  if (loading) {
    return (
      <div className={styles.loading}>
        <div className="spinner"></div>
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className={styles.container}>
        <div className={styles.error}>{error || "User not found"}</div>
        <div className={styles.backToListButton}>
          <Link to="/users">
            <Button>Back to Users</Button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.modal}>
        {/* Modal Header */}
        <div className={styles.modalHeader}>
          <h1 className={styles.modalTitle}>{user.name}</h1>
          <p className={styles.modalSubtitle}>
            @{user.username} • ID: {user.id}
          </p>
        </div>

        {/* Modal Body */}
        <div className={styles.modalBody}>
          {/* Personal Information */}
          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>Personal Information</h2>
            <div className={styles.fieldGroup}>
              <span className={styles.fieldLabel}>Full Name:</span>
              <p className={styles.fieldValue}>{user.name}</p>
            </div>
            <div className={styles.fieldGroup}>
              <span className={styles.fieldLabel}>Username:</span>
              <p className={styles.fieldValue}>@{user.username}</p>
            </div>
            <div className={styles.fieldGroup}>
              <span className={styles.fieldLabel}>Email:</span>
              <p className={styles.fieldValue}>{user.email}</p>
            </div>
            {user.phone && (
              <div className={styles.fieldGroup}>
                <span className={styles.fieldLabel}>Phone:</span>
                <p className={styles.fieldValue}>{user.phone}</p>
              </div>
            )}
            {user.website && (
              <div className={styles.fieldGroup}>
                <span className={styles.fieldLabel}>Website:</span>
                <a
                  href={
                    user.website.startsWith("http")
                      ? user.website
                      : `https://${user.website}`
                  }
                  target="_blank"
                  rel="noopener noreferrer"
                  className={styles.websiteLink}
                >
                  {user.website}
                </a>
              </div>
            )}
          </div>

          {/* Address Information */}
          {user.address && (
            <div className={styles.section}>
              <h2 className={styles.sectionTitle}>Address</h2>
              <div className={styles.fieldGroup}>
                <span className={styles.fieldLabel}>Street:</span>
                <p className={styles.fieldValue}>
                  {user.address.street} {user.address.suite}
                </p>
              </div>
              <div className={styles.fieldGroup}>
                <span className={styles.fieldLabel}>City:</span>
                <p className={styles.fieldValue}>{user.address.city}</p>
              </div>
              <div className={styles.fieldGroup}>
                <span className={styles.fieldLabel}>Zipcode:</span>
                <p className={styles.fieldValue}>{user.address.zipcode}</p>
              </div>
              {user.address.geo &&
                (user.address.geo.lat || user.address.geo.lng) && (
                  <div className={styles.fieldGroup}>
                    <span className={styles.fieldLabel}>Coordinates:</span>
                    <p className={styles.fieldValue}>
                      {user.address.geo.lat}, {user.address.geo.lng}
                    </p>
                  </div>
                )}
            </div>
          )}

          {/* Company Information */}
          {user.company && (
            <div className={styles.section}>
              <h2 className={styles.sectionTitle}>Company</h2>
              <div className={styles.fieldGroup}>
                <span className={styles.fieldLabel}>Company Name:</span>
                <p className={styles.fieldValue}>{user.company.name}</p>
              </div>
              {user.company.catchPhrase && (
                <div className={styles.fieldGroup}>
                  <span className={styles.fieldLabel}>Catch Phrase:</span>
                  <p className={styles.fieldValueItalic}>
                    "{user.company.catchPhrase}"
                  </p>
                </div>
              )}
              {user.company.bs && (
                <div className={styles.fieldGroup}>
                  <span className={styles.fieldLabel}>Business:</span>
                  <p className={styles.fieldValue}>{user.company.bs}</p>
                </div>
              )}
            </div>
          )}
        </div>

        {/* Modal Footer */}
        <div className={styles.modalFooter}>
          <Link to="/users">
            <Button variant="outline">Back to Users</Button>
          </Link>
          {isAuthenticated && (
            <>
              <Link to={`/users/${user.id}/edit`}>
                <Button variant="secondary">Edit User</Button>
              </Link>
              <Button
                variant="danger"
                onClick={handleDelete}
                loading={deleteLoading}
              >
                Delete User
              </Button>
            </>
          )}
        </div>
      </div>
    </div>
  );
};
