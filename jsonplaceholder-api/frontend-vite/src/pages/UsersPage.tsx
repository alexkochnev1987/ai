import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { usersApi } from "../services/api";
import { useAuth } from "../contexts/AuthContext";
import { Button } from "../components/ui/Button";
import { X } from "lucide-react";
import type { UserResponse, PaginatedResponse } from "../types/api";
import styles from "./UsersPage.module.css";

export const UsersPage: React.FC = () => {
  const [users, setUsers] = useState<UserResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);
  const [deleteLoading, setDeleteLoading] = useState<number | null>(null);

  const { isAuthenticated } = useAuth();
  const limit = 10;

  const fetchUsers = async (page: number) => {
    setLoading(true);
    setError("");

    try {
      const response: PaginatedResponse<UserResponse[]> =
        await usersApi.getUsers(page, limit);
      setUsers(response.data);
      setTotalPages(response.pagination.totalPages);
      setTotal(response.pagination.total);
    } catch (err) {
      console.error("Error fetching users:", err);
      setError("Failed to load users");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers(currentPage);
  }, [currentPage]);

  const handleDelete = async (userId: number) => {
    if (!isAuthenticated) {
      alert("You must be logged in to delete users");
      return;
    }

    if (!confirm("Are you sure you want to delete this user?")) {
      return;
    }

    setDeleteLoading(userId);
    try {
      await usersApi.deleteUser(userId);
      // Refresh the current page
      await fetchUsers(currentPage);
    } catch (err) {
      console.error("Error deleting user:", err);
      alert("Failed to delete user");
    } finally {
      setDeleteLoading(null);
    }
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  if (loading && users.length === 0) {
    return (
      <div className={styles.loading}>
        <div className="spinner"></div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.headerInfo}>
          <h1>Users</h1>
          <p>Total: {total} users</p>
        </div>
        {isAuthenticated && (
          <Link to="/users/create">
            <Button>Create User</Button>
          </Link>
        )}
      </div>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.tableContainer}>
        <table className={styles.table}>
          <thead className={styles.tableHeader}>
            <tr>
              <th className={styles.tableHeaderCell}>Name / Email</th>
              <th className={styles.tableHeaderCell}>Address</th>
              <th className={styles.tableHeaderCell}>Phone</th>
              <th className={styles.tableHeaderCell}>Website</th>
              <th className={styles.tableHeaderCell}>Company</th>
              <th className={styles.tableHeaderCell}>Action</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id} className={styles.tableRow}>
                <td className={`${styles.tableCell} ${styles.nameCell}`}>
                  <div className={styles.userName}>{user.name}</div>
                  <div className={styles.userEmail}>{user.email}</div>
                </td>
                <td className={`${styles.tableCell} ${styles.addressCell}`}>
                  {user.address && (
                    <div className={styles.addressText}>
                      {user.address.street} {user.address.suite}
                      <br />
                      {user.address.city}, {user.address.zipcode}
                    </div>
                  )}
                </td>
                <td className={`${styles.tableCell} ${styles.phoneCell}`}>
                  {user.phone && (
                    <div className={styles.phoneText}>{user.phone}</div>
                  )}
                </td>
                <td className={`${styles.tableCell} ${styles.websiteCell}`}>
                  {user.website && (
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
                  )}
                </td>
                <td className={`${styles.tableCell} ${styles.companyCell}`}>
                  {user.company && (
                    <div>
                      <div className={styles.companyName}>
                        {user.company.name}
                      </div>
                      {user.company.catchPhrase && (
                        <div className={styles.companyCatchPhrase}>
                          {user.company.catchPhrase}
                        </div>
                      )}
                    </div>
                  )}
                </td>
                <td className={`${styles.tableCell} ${styles.actionCell}`}>
                  <div className={styles.actionButtons}>
                    <Link to={`/users/${user.id}`}>
                      <Button variant="outline" size="sm">
                        View
                      </Button>
                    </Link>
                    {isAuthenticated && (
                      <button
                        className={styles.deleteButton}
                        onClick={() => handleDelete(user.id)}
                        disabled={deleteLoading === user.id}
                        title="Delete user"
                      >
                        <X size={16} />
                      </button>
                    )}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className={styles.pagination}>
          <Button
            variant="outline"
            size="sm"
            onClick={() => handlePageChange(currentPage - 1)}
            disabled={currentPage === 1}
          >
            Previous
          </Button>

          <div className={styles.paginationNumbers}>
            {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
              const page = i + 1;
              return (
                <Button
                  key={page}
                  variant={currentPage === page ? "primary" : "outline"}
                  size="sm"
                  onClick={() => handlePageChange(page)}
                >
                  {page}
                </Button>
              );
            })}
          </div>

          <Button
            variant="outline"
            size="sm"
            onClick={() => handlePageChange(currentPage + 1)}
            disabled={currentPage === totalPages}
          >
            Next
          </Button>
        </div>
      )}
    </div>
  );
};

export default UsersPage;
