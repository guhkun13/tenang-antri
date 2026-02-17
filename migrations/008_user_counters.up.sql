-- Create user_counters table
CREATE TABLE IF NOT EXISTS user_counters (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    counter_id INTEGER NOT NULL REFERENCES counters(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, counter_id)
);

-- Create index for faster lookups
CREATE INDEX idx_user_counters_user_id ON user_counters(user_id);
CREATE INDEX idx_user_counters_counter_id ON user_counters(counter_id);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_user_counters_updated_at BEFORE UPDATE ON user_counters
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Migrate existing data from users.counter_id to user_counters
INSERT INTO user_counters (user_id, counter_id)
SELECT id, counter_id FROM users WHERE counter_id IS NOT NULL;
