ALTER TABLE refresh_tokens 
DROP COLUMN IF EXISTS revocation_reason,
DROP COLUMN IF EXISTS revoked_at;
