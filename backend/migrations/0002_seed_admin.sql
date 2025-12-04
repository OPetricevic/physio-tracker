-- Seed a default doctor (local offline use)
INSERT INTO doctors (uuid, email, first_name, last_name, created_at)
VALUES ('11111111-1111-1111-1111-111111111111', 'sebastijan@gmail.com', 'Admin', 'User', NOW())
ON CONFLICT (uuid) DO NOTHING;

-- Store a simple password hash placeholder (local-only; replace with real hash later)
INSERT INTO doctor_credentials (doctor_uuid, password_hash, password_updated_at)
VALUES ('11111111-1111-1111-1111-111111111111', 'admin', NOW())
ON CONFLICT (doctor_uuid) DO NOTHING;
