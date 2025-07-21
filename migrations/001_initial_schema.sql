-- Initial schema for oncall service

-- Teams table
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Users table with team reference
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    slack_handle VARCHAR(255) NOT NULL,
    team_id INTEGER REFERENCES teams(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Schedules table
CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id),
    name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    rotation_period INTEGER NOT NULL, -- in seconds
    participant_ids TEXT, -- comma-separated list of user IDs
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- OnCall assignments table
CREATE TABLE oncall_assignments (
    id SERIAL PRIMARY KEY,
    schedule_id INTEGER REFERENCES schedules(id),
    user_id INTEGER REFERENCES users(id),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better query performance
CREATE INDEX idx_users_team ON users(team_id);
CREATE INDEX idx_schedules_team ON schedules(team_id);
CREATE INDEX idx_schedules_times ON schedules(start_time, end_time);
CREATE INDEX idx_oncall_assignments_schedule ON oncall_assignments(schedule_id);
CREATE INDEX idx_oncall_assignments_user ON oncall_assignments(user_id);
CREATE INDEX idx_oncall_assignments_times ON oncall_assignments(start_time, end_time);
CREATE INDEX idx_oncall_assignments_active ON oncall_assignments(active);

-- Comments
COMMENT ON COLUMN schedules.rotation_period IS 'Rotation period in seconds';
COMMENT ON COLUMN schedules.participant_ids IS 'Comma-separated list of user IDs participating in the rotation';