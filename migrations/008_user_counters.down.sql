-- Drop trigger
DROP TRIGGER IF EXISTS update_user_counters_updated_at ON user_counters;

-- Drop indexes
DROP INDEX IF EXISTS idx_user_counters_user_id;
DROP INDEX IF EXISTS idx_user_counters_counter_id;

-- Drop table
DROP TABLE IF EXISTS user_counters;
