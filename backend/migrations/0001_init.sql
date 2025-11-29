-- Doctors and authentication
CREATE TABLE IF NOT EXISTS doctors (
    uuid UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS doctor_credentials (
    doctor_uuid UUID PRIMARY KEY REFERENCES doctors(uuid) ON DELETE CASCADE,
    password_hash TEXT NOT NULL,
    password_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Patients (owned by doctor)
CREATE TABLE IF NOT EXISTS patients (
    uuid UUID PRIMARY KEY,
    doctor_uuid UUID NOT NULL REFERENCES doctors(uuid) ON DELETE CASCADE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    phone TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_patients_doctor ON patients(doctor_uuid);
CREATE INDEX IF NOT EXISTS idx_patients_name ON patients(LOWER(first_name), LOWER(last_name));

-- Anamneses (session notes)
CREATE TABLE IF NOT EXISTS anamneses (
    uuid UUID PRIMARY KEY,
    patient_uuid UUID NOT NULL REFERENCES patients(uuid) ON DELETE CASCADE,
    note TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_anamneses_patient ON anamneses(patient_uuid, created_at DESC);

-- Auth tokens (e.g., refresh tokens or session tokens)
CREATE TABLE IF NOT EXISTS auth_tokens (
    id BIGSERIAL PRIMARY KEY,
    doctor_uuid UUID NOT NULL REFERENCES doctors(uuid) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_auth_tokens_doctor ON auth_tokens(doctor_uuid);
