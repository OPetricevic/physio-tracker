-- Extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Doctors and authentication
CREATE TABLE IF NOT EXISTS doctors (
    uuid VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS doctor_credentials (
    uuid VARCHAR(255) PRIMARY KEY,
    doctor_uuid VARCHAR(255) NOT NULL UNIQUE REFERENCES doctors(uuid) ON DELETE CASCADE,
    password_hash VARCHAR(255) NOT NULL,
    password_updated_at TIMESTAMP NOT NULL
);

-- Patients (owned by doctor)
CREATE TABLE IF NOT EXISTS patients (
    uuid VARCHAR(255) PRIMARY KEY,
    doctor_uuid VARCHAR(255) NOT NULL REFERENCES doctors(uuid) ON DELETE CASCADE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    address VARCHAR(255),
    date_of_birth DATE,
    sex VARCHAR(20),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_patients_doctor ON patients(doctor_uuid);
CREATE INDEX IF NOT EXISTS idx_patients_name ON patients(LOWER(first_name), LOWER(last_name));

-- Anamneses (session notes)
CREATE TABLE IF NOT EXISTS anamneses (
    uuid VARCHAR(255) PRIMARY KEY,
    patient_uuid VARCHAR(255) NOT NULL REFERENCES patients(uuid) ON DELETE CASCADE,
    note TEXT NOT NULL, -- long text allowed
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_anamneses_patient ON anamneses(patient_uuid, created_at DESC);

-- Auth tokens (e.g., refresh tokens or session tokens)
CREATE TABLE IF NOT EXISTS auth_tokens (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR(255) NOT NULL UNIQUE DEFAULT gen_random_uuid()::text,
    doctor_uuid VARCHAR(255) NOT NULL REFERENCES doctors(uuid) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_auth_tokens_doctor ON auth_tokens(doctor_uuid);
