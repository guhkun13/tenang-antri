-- Update existing counters based on is_active and status
-- If counter was inactive (is_active = false), set to OFFLINE
-- If counter was active with status = 'paused', keep as PAUSED
-- If counter was active with other status, set to IDLE
UPDATE counters SET status = 'offline' WHERE is_active = false;
UPDATE counters SET status = 'idle' WHERE is_active = true AND status NOT IN ('paused', 'serving');

-- Drop the is_active column
ALTER TABLE counters DROP COLUMN IF EXISTS is_active;
