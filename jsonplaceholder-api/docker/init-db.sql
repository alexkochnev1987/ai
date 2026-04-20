-- Initialize database for JSONPlaceholder API
-- This script is executed when PostgreSQL container starts

-- Create database if not exists (though it's created via POSTGRES_DB env var)
-- CREATE DATABASE IF NOT EXISTS jsonplaceholder;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE jsonplaceholder TO postgres;

-- Create extensions if needed
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- You can add any initial database setup here
-- The application will handle table creation via GORM migrations 