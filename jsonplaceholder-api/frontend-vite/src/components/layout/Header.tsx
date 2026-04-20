import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";
import { Button } from "../ui/Button";
import { authApi } from "../../services/api";
import styles from "./Header.module.css";

export const Header: React.FC = () => {
  const { user, isAuthenticated, logout, refreshToken } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      if (refreshToken) {
        await authApi.logout(refreshToken);
      }
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      logout();
      navigate("/login");
    }
  };

  return (
    <header className={styles.header}>
      <div className={styles.container}>
        <div className={styles.content}>
          <div className={styles.leftSection}>
            <Link to="/" className={styles.logo}>
              JSONPlaceholder API
            </Link>
            <nav className={styles.nav}>
              <Link to="/users" className={styles.navLink}>
                Users
              </Link>
              <Link to="/health" className={styles.navLink}>
                Health
              </Link>
            </nav>
          </div>

          <div className={styles.rightSection}>
            {isAuthenticated ? (
              <>
                <span className={styles.userInfo}>Welcome, {user?.name}</span>
                <Link to="/profile" className={styles.profileLink}>
                  Profile
                </Link>
                <Button variant="outline" size="sm" onClick={handleLogout}>
                  Logout
                </Button>
              </>
            ) : (
              <>
                <Link to="/login">
                  <Button variant="outline" size="sm">
                    Login
                  </Button>
                </Link>
                <Link to="/register">
                  <Button size="sm">Register</Button>
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};
