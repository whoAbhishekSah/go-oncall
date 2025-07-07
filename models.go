package main

import "time"

type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	SlackHandle string    `json:"slack_handle"`
	TeamID      int       `json:"team_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Team struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Users     []User    `json:"users"`
	CreatedAt time.Time `json:"created_at"`
}

type Schedule struct {
	ID             int       `json:"id"`
	TeamID         int       `json:"team_id"`
	Name           string    `json:"name"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	RotationPeriod int       `json:"rotation_period"` // in hours
	Participants   []int     `json:"participants"`    // user IDs
	CreatedAt      time.Time `json:"created_at"`
}

type OnCallAssignment struct {
	ID         int       `json:"id"`
	ScheduleID int       `json:"schedule_id"`
	UserID     int       `json:"user_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Active     bool      `json:"active"`
}