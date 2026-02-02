-- Fix category delete by adding ON DELETE CASCADE to tickets table
-- First, drop the existing foreign key constraint
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_category_id_fkey;

-- Add the foreign key constraint with ON DELETE CASCADE
ALTER TABLE tickets ADD CONSTRAINT tickets_category_id_fkey 
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;
