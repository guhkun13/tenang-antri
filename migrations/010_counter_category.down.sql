-- Drop trigger
DROP TRIGGER IF EXISTS update_counter_category_updated_at ON counter_category;

-- Drop indexes
DROP INDEX IF EXISTS idx_counter_category_counter_id;
DROP INDEX IF EXISTS idx_counter_category_category_id;

-- Drop table
DROP TABLE IF EXISTS counter_category;

-- Add current_staff_id back to counters
ALTER TABLE counters ADD COLUMN IF NOT EXISTS current_staff_id INTEGER REFERENCES users(id) ON DELETE SET NULL;

-- Note: We cannot restore category_id to counters as it was a single category,
-- but data is still in counter_categories table
