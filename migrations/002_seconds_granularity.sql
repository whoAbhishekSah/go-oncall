-- Update rotation_period to support seconds granularity
-- This migration changes the rotation_period column comment and adds better indexing

-- Add a comment to clarify that rotation_period is now in seconds
COMMENT ON COLUMN schedules.rotation_period IS 'Rotation period in seconds';

-- Create additional indexes for better performance with frequent rotations
CREATE INDEX IF NOT EXISTS idx_schedules_rotation_period ON schedules(rotation_period);
CREATE INDEX IF NOT EXISTS idx_oncall_assignments_times ON oncall_assignments(start_time, end_time, active);