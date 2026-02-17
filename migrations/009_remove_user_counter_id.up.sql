-- Remove counter_id column from users table
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_counter;
ALTER TABLE users DROP COLUMN IF EXISTS counter_id;
