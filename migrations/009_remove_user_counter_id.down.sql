-- Add counter_id column back to users table
ALTER TABLE users ADD COLUMN counter_id INTEGER;
ALTER TABLE users ADD CONSTRAINT fk_users_counter 
    FOREIGN KEY (counter_id) REFERENCES counters(id) ON DELETE SET NULL;

-- Restore data from user_counters (this would need to be done carefully in a real scenario)
-- Note: This is a simplified rollback - in production, you'd need to handle data conflicts
