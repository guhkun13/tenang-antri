-- Revert the category delete fix
-- Drop the cascade constraint and restore the original
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_category_id_fkey;

-- Add back the original foreign key constraint without cascade
ALTER TABLE tickets ADD CONSTRAINT tickets_category_id_fkey 
    FOREIGN KEY (category_id) REFERENCES categories(id);
