-- Complete schema for OnCall Scheduler with user management
-- Fresh installation schema

-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    slack_handle VARCHAR(255) NOT NULL,
    team_id INTEGER REFERENCES teams(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create team_users junction table for many-to-many relationships
CREATE TABLE IF NOT EXISTS team_users (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id),
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(team_id, user_id)
);

-- Create schedules table with user IDs
CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id),
    name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    rotation_period INTEGER NOT NULL,
    participant_ids TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create oncall_assignments table with user references
CREATE TABLE IF NOT EXISTS oncall_assignments (
    id SERIAL PRIMARY KEY,
    schedule_id INTEGER REFERENCES schedules(id),
    user_id INTEGER REFERENCES users(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    active BOOLEAN DEFAULT TRUE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_slack_handle ON users(slack_handle);
CREATE INDEX IF NOT EXISTS idx_users_team_id ON users(team_id);
CREATE INDEX IF NOT EXISTS idx_teams_name ON teams(name);
CREATE INDEX IF NOT EXISTS idx_team_users_team_id ON team_users(team_id);
CREATE INDEX IF NOT EXISTS idx_team_users_user_id ON team_users(user_id);
CREATE INDEX IF NOT EXISTS idx_schedules_team_id ON schedules(team_id);
CREATE INDEX IF NOT EXISTS idx_schedules_start_time ON schedules(start_time);
CREATE INDEX IF NOT EXISTS idx_schedules_end_time ON schedules(end_time);
CREATE INDEX IF NOT EXISTS idx_oncall_assignments_schedule_id ON oncall_assignments(schedule_id);
CREATE INDEX IF NOT EXISTS idx_oncall_assignments_user_id ON oncall_assignments(user_id);
CREATE INDEX IF NOT EXISTS idx_oncall_assignments_active ON oncall_assignments(active);
CREATE INDEX IF NOT EXISTS idx_oncall_assignments_start_time ON oncall_assignments(start_time);
CREATE INDEX IF NOT EXISTS idx_oncall_assignments_end_time ON oncall_assignments(end_time);