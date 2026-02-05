-- Add daily_sequence and queue_date columns to tickets table
ALTER TABLE tickets ADD COLUMN daily_sequence INTEGER;
ALTER TABLE tickets ADD COLUMN queue_date DATE DEFAULT CURRENT_DATE;

-- Create indices to optimize ticket queries
CREATE INDEX idx_date_category ON tickets(queue_date, category_id);
CREATE INDEX idx_date_number ON tickets(queue_date, daily_sequence);

-- Add unique constraint to guarantee no duplicates per category per day
ALTER TABLE tickets ADD CONSTRAINT unique_daily_queue UNIQUE (category_id, queue_date, daily_sequence);

-- Backfill existing data
WITH NumberedTickets AS (
    SELECT 
        id,
        DATE(created_at) as q_date,
        ROW_NUMBER() OVER (PARTITION BY category_id, DATE(created_at) ORDER BY created_at ASC) as seq
    FROM tickets
)
UPDATE tickets
SET 
    daily_sequence = NumberedTickets.seq,
    queue_date = NumberedTickets.q_date
FROM NumberedTickets
WHERE tickets.id = NumberedTickets.id;
