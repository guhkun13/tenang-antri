-- Create counter_category table for counter-category many-to-many relationship
CREATE TABLE IF NOT EXISTS counter_category (
    id SERIAL PRIMARY KEY,
    counter_id INTEGER NOT NULL REFERENCES counters(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(counter_id, category_id)
);

-- Create indexes for faster lookups
CREATE INDEX idx_counter_category_counter_id ON counter_category(counter_id);
CREATE INDEX idx_counter_category_category_id ON counter_category(category_id);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_counter_category_updated_at BEFORE UPDATE ON counter_category
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Migrate existing category data from counter_categories to counter_category
INSERT INTO counter_category (counter_id, category_id)
SELECT counter_id, category_id FROM counter_categories
ON CONFLICT DO NOTHING;

-- Drop foreign key constraint on counters.current_staff_id
ALTER TABLE counters DROP CONSTRAINT IF EXISTS counters_current_staff_id_fkey;

-- Drop category_id from counters if it exists (it was added in a later migration)
ALTER TABLE counters DROP COLUMN IF EXISTS category_id;

-- Drop current_staff_id from counters
ALTER TABLE counters DROP COLUMN IF EXISTS current_staff_id;
