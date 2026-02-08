-- Update counters table status constraint to match new statuses
-- First, drop the existing constraint (it has a system-generated name, so we need to find it)
ALTER TABLE counters DROP CONSTRAINT IF EXISTS counters_status_check;

-- Add the new constraint with updated status values
ALTER TABLE counters ADD CONSTRAINT counters_status_check 
    CHECK (status IN ('offline', 'idle', 'serving', 'paused'));

-- Update existing counters to use new status values
UPDATE counters SET status = 'offline' WHERE status = 'inactive';
UPDATE counters SET status = 'idle' WHERE status = 'active';
