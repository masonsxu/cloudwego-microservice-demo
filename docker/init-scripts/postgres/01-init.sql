-- -----------------------------------------------------------------------------
-- 01-init.sql
--
-- This script creates the necessary databases for the project.
-- It is executed automatically by the PostgreSQL container on startup.
--
-- The database names are derived from the config.yaml files of the respective
-- microservices.
--
-- Databases to be created:
-- - identity_srv (used by identity_srv)
--
-- The user 'Admin' is created by the entrypoint script via POSTGRES_USER env var.
-- We will set this user as the owner of the new databases.
-- -----------------------------------------------------------------------------

-- Create database for identity and permission services
SELECT 'CREATE DATABASE identity_srv OWNER "Admin"'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'identity_srv')\gexec
SELECT 'CREATE DATABASE permission_srv OWNER "Admin"'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'permission_srv')\gexec

-- Connect to each database and enable required extensions
\c identity_srv
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

\c permission_srv
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";