-- Recreate the counter_categories junction table
CREATE TABLE IF NOT EXISTS counter_categories (
    counter_id INTEGER REFERENCES counters(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (counter_id, category_id)
);

-- Remove category_id column from counters table
ALTER TABLE counters DROP COLUMN IF EXISTS category_id;
