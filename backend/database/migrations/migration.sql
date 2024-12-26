-- Create the database if it doesn't exist
DO $$ BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'user_management') THEN
      CREATE DATABASE user_management;
   END IF;
END $$;

-- Connect to the database
\c user_management;

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    user_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    user_status VARCHAR(1) NOT NULL CHECK (user_status IN ('I', 'A', 'T')),
    department VARCHAR(255) DEFAULT NULL
);