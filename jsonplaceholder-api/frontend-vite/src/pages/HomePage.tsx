import React from "react";
import { Link } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";
import { Card } from "../components/ui/Card";
import { Button } from "../components/ui/Button";
import styles from "./HomePage.module.css";

export const HomePage: React.FC = () => {
  const { isAuthenticated, user } = useAuth();

  return (
    <div className={styles.container}>
      {/* Hero Section */}
      <div className={styles.hero}>
        <h1 className={styles.title}>JSONPlaceholder API Frontend</h1>
        <p className={styles.subtitle}>
          A modern React frontend for the JSONPlaceholder API with full CRUD
          operations
        </p>
        {!isAuthenticated ? (
          <div className={styles.actions}>
            <Link to="/login">
              <Button size="lg">Get Started</Button>
            </Link>
            <Link to="/users">
              <Button variant="outline" size="lg">
                Browse Users
              </Button>
            </Link>
          </div>
        ) : (
          <div className={styles.welcome}>Welcome back, {user?.name}!</div>
        )}
      </div>

      {/* Features Grid */}
      <div className={styles.features}>
        <Card>
          <div className={styles.feature}>
            <div className={styles.featureIconBlue}>
              <svg
                width="24"
                height="24"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"
                />
              </svg>
            </div>
            <h3 className={styles.featureTitle}>User Management</h3>
            <p className={styles.featureDescription}>
              View, create, edit, and delete users with full CRUD operations
            </p>
            <Link to="/users">
              <Button variant="outline" size="sm">
                View Users
              </Button>
            </Link>
          </div>
        </Card>

        <Card>
          <div className={styles.feature}>
            <div className={styles.featureIconGreen}>
              <svg
                width="24"
                height="24"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <h3 className={styles.featureTitle}>Authentication</h3>
            <p className={styles.featureDescription}>
              Secure login and registration with JWT token management
            </p>
            {!isAuthenticated ? (
              <Link to="/login">
                <Button variant="outline" size="sm">
                  Sign In
                </Button>
              </Link>
            ) : (
              <Link to="/profile">
                <Button variant="outline" size="sm">
                  View Profile
                </Button>
              </Link>
            )}
          </div>
        </Card>

        <Card>
          <div className={styles.feature}>
            <div className={styles.featureIconPurple}>
              <svg
                width="24"
                height="24"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                />
              </svg>
            </div>
            <h3 className={styles.featureTitle}>API Health</h3>
            <p className={styles.featureDescription}>
              Monitor API status and view available endpoints
            </p>
            <Link to="/health">
              <Button variant="outline" size="sm">
                Check Status
              </Button>
            </Link>
          </div>
        </Card>
      </div>

      {/* Quick Actions */}
      {isAuthenticated && (
        <Card>
          <h2 className={styles.quickActionsTitle}>Quick Actions</h2>
          <div className={styles.quickActionsGrid}>
            <Link to="/users/create">
              <Button>Create New User</Button>
            </Link>
            <Link to="/users">
              <Button variant="outline">Browse All Users</Button>
            </Link>
            <Link to="/profile">
              <Button variant="outline">View My Profile</Button>
            </Link>
            <Link to="/health">
              <Button variant="outline">API Health</Button>
            </Link>
          </div>
        </Card>
      )}

      {/* API Information */}
      <Card>
        <h2 className={styles.aboutTitle}>About This Application</h2>
        <div className={styles.aboutContent}>
          <p>
            This is a modern React frontend built with TypeScript and CSS
            Modules that interfaces with a Go-based JSONPlaceholder API. The
            application demonstrates:
          </p>
          <ul>
            <li>Complete CRUD operations for user management</li>
            <li>JWT-based authentication and authorization</li>
            <li>Responsive design with modern UI components</li>
            <li>Real-time API health monitoring</li>
            <li>Form validation and error handling</li>
            <li>Pagination and search functionality</li>
          </ul>
          <p>
            The backend API provides endpoints for user management,
            authentication, and system health checks, following RESTful
            principles and modern security practices.
          </p>
        </div>
      </Card>
    </div>
  );
};
