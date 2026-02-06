-- Add is_active column back (for rollback purposes)
ALTER TABLE counters ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true;

-- Update counters back to original state
-- Note: This migration assumes original behavior where active counters were 'active' or 'paused'
-- Rollback will restore approximate original state
UPDATE counters SET is_active = false WHERE status = 'offline';
UPDATE counters SET is_active = true WHERE status IN ('idle', 'serving', 'paused');
