import React from "react";
import { Button } from "../../common/Button";
import styles from "./UserTable.module.css";
import type { User } from "../../../types/user";

interface UserTableProps {
  users: User[];
  onUserClick: (user: User) => void;
  onUserDelete: (userId: number) => void;
  isLoading?: boolean;
}

export const UserTable: React.FC<UserTableProps> = ({
  users,
  onUserClick,
  onUserDelete,
  isLoading = false,
}) => {
  const handleDeleteClick = (e: React.MouseEvent, userId: number) => {
    e.stopPropagation();
    if (window.confirm("Are you sure you want to delete this user?")) {
      onUserDelete(userId);
    }
  };

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <p>Loading users...</p>
      </div>
    );
  }

  return (
    <div className={styles.tableContainer}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th>Name & Email</th>
            <th>Address</th>
            <th>Phone</th>
            <th>Website</th>
            <th>Company</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {users.map((user) => (
            <tr
              key={user.id}
              className={styles.userRow}
              onClick={() => onUserClick(user)}
              role="button"
              tabIndex={0}
              onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                  e.preventDefault();
                  onUserClick(user);
                }
              }}
            >
              <td>
                <div className={styles.userInfo}>
                  <strong>{user.name}</strong>
                  <div className={styles.email}>{user.email}</div>
                </div>
              </td>
              <td>
                <div className={styles.address}>
                  {user.address.street}, {user.address.suite}
                  <br />
                  {user.address.city}, {user.address.zipcode}
                </div>
              </td>
              <td>{user.phone}</td>
              <td>
                <a
                  href={`https://${user.website}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  onClick={(e) => e.stopPropagation()}
                >
                  {user.website}
                </a>
              </td>
              <td>
                <div className={styles.company}>
                  <strong>{user.company.name}</strong>
                  <div className={styles.catchPhrase}>
                    {user.company.catchPhrase}
                  </div>
                </div>
              </td>
              <td>
                <Button
                  variant="danger"
                  size="sm"
                  onClick={(e) => handleDeleteClick(e, user.id)}
                  aria-label={`Delete user ${user.name}`}
                >
                  Delete
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {users.length === 0 && !isLoading && (
        <div className={styles.emptyState}>
          <p>No users found</p>
        </div>
      )}
    </div>
  );
};

export default UserTable;
