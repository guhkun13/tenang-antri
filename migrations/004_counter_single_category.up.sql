-- Add category_id column to counters table
ALTER TABLE counters ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL;

-- Drop the counter_categories junction table
DROP TABLE IF EXISTS counter_categories CASCADE;
