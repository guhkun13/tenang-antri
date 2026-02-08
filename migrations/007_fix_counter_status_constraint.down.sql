-- Rollback the counters table status constraint
-- Drop the new constraint
ALTER TABLE counters DROP CONSTRAINT IF EXISTS counters_status_check;

-- Add back the original constraint
ALTER TABLE counters ADD CONSTRAINT counters_status_check 
    CHECK (status IN ('active', 'paused', 'inactive'));

-- Revert counters back to original status values
UPDATE counters SET status = 'inactive' WHERE status = 'offline';
UPDATE counters SET status = 'active' WHERE status = 'idle';
