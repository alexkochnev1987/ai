import React from "react";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import { AuthProvider } from "./contexts/AuthContext";
import { Layout } from "./components/layout/Layout";
import { HomePage } from "./pages/HomePage";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import { UsersPage } from "./pages/UsersPage";
import { UserDetailPage } from "./pages/UserDetailPage";
import { CreateUserPage } from "./pages/CreateUserPage";
import { EditUserPage } from "./pages/EditUserPage";
import { ProfilePage } from "./pages/ProfilePage";
import { HealthPage } from "./pages/HealthPage";

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          {/* Auth routes without layout */}
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          {/* Main routes with layout */}
          <Route
            path="/"
            element={
              <Layout>
                <HomePage />
              </Layout>
            }
          />
          <Route
            path="/users"
            element={
              <Layout>
                <UsersPage />
              </Layout>
            }
          />
          <Route
            path="/users/create"
            element={
              <Layout>
                <CreateUserPage />
              </Layout>
            }
          />
          <Route
            path="/users/:id"
            element={
              <Layout>
                <UserDetailPage />
              </Layout>
            }
          />
          <Route
            path="/users/:id/edit"
            element={
              <Layout>
                <EditUserPage />
              </Layout>
            }
          />
          <Route
            path="/profile"
            element={
              <Layout>
                <ProfilePage />
              </Layout>
            }
          />
          <Route
            path="/health"
            element={
              <Layout>
                <HealthPage />
              </Layout>
            }
          />

          {/* Redirect unknown routes to home */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
