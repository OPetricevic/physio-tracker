-- Seed default doctor (Sebastijan)
INSERT INTO doctors (uuid, email, username, first_name, last_name, created_at)
VALUES ('11111111-1111-1111-1111-111111111111', 'sartorius.therapy@gmail.com', 'sartorius.therapy', 'Sebastijan', 'Petricevic', NOW())
ON CONFLICT (uuid) DO NOTHING;

-- Store bcrypt hash using pgcrypto (password: admin)
INSERT INTO doctor_credentials (uuid, doctor_uuid, password_hash, password_updated_at)
VALUES (
  '22222222-2222-2222-2222-222222222222',
  '11111111-1111-1111-1111-111111111111',
  crypt('admin', gen_salt('bf')),
  NOW()
)
ON CONFLICT (doctor_uuid) DO NOTHING;
