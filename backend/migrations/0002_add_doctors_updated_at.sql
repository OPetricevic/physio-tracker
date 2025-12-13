-- Add missing updated_at on doctors
ALTER TABLE doctors
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NULL;
