-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS update_counters_updated_at ON counters;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (order matters for foreign keys)
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS counter_categories;
DROP TABLE IF EXISTS daily_stats;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS counters;
