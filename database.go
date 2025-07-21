package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func initDB() {
	log.Println("Database initialization - tables should be created via migration script")
	log.Println("Run ./migrate.sh docker or ./migrate.sh prod to create tables")
}

// User functions
func createUser(email, slackHandle string, teamID int) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO users (email, slack_handle, team_id) VALUES ($1, $2, $3) RETURNING id", 
		email, slackHandle, teamID).Scan(&id)
	return id, err
}

func getUsers() ([]User, error) {
	rows, err := db.Query("SELECT id, email, slack_handle, team_id, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.SlackHandle, &user.TeamID, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func getUsersByTeamID(teamID int) ([]User, error) {
	rows, err := db.Query("SELECT id, email, slack_handle, team_id, created_at FROM users WHERE team_id = $1", teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.SlackHandle, &user.TeamID, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func getUserByID(userID int) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, email, slack_handle, team_id, created_at FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Email, &user.SlackHandle, &user.TeamID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Team functions
func createTeam(name string) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO teams (name) VALUES ($1) RETURNING id", name).Scan(&id)
	return id, err
}

func getTeams() ([]Team, error) {
	rows, err := db.Query("SELECT id, name, created_at FROM teams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var team Team
		err := rows.Scan(&team.ID, &team.Name, &team.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		// Get users for this team
		users, err := getUsersByTeamID(team.ID)
		if err != nil {
			return nil, err
		}
		team.Users = users
		
		teams = append(teams, team)
	}
	return teams, nil
}

func addUserToTeam(userID, teamID int) error {
	_, err := db.Exec("INSERT INTO team_users (team_id, user_id) VALUES ($1, $2) ON CONFLICT (team_id, user_id) DO NOTHING", 
		teamID, userID)
	return err
}

// Schedule functions
func createSchedule(teamID int, name string, startTime, endTime time.Time, rotationPeriod int, participants []int) (int, error) {
	// Convert participant IDs to string for storage
	participantStrings := make([]string, len(participants))
	for i, id := range participants {
		participantStrings[i] = strconv.Itoa(id)
	}
	participantList := strings.Join(participantStrings, ",")
	
	var id int
	err := db.QueryRow("INSERT INTO schedules (team_id, name, start_time, end_time, rotation_period, participant_ids) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", 
		teamID, name, startTime, endTime, rotationPeriod, participantList).Scan(&id)
	return id, err
}

func getSchedules() ([]Schedule, error) {
	rows, err := db.Query("SELECT id, team_id, name, start_time, end_time, rotation_period, COALESCE(participant_ids, '') as participant_ids, created_at FROM schedules")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var schedule Schedule
		var participantList string
		err := rows.Scan(&schedule.ID, &schedule.TeamID, &schedule.Name, &schedule.StartTime, 
			&schedule.EndTime, &schedule.RotationPeriod, &participantList, &schedule.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		// Convert participant string back to int slice
		if participantList != "" {
			participantStrings := strings.Split(participantList, ",")
			schedule.Participants = make([]int, len(participantStrings))
			for i, idStr := range participantStrings {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return nil, fmt.Errorf("error converting participant ID: %v", err)
				}
				schedule.Participants[i] = id
			}
		}
		
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

// OnCall Assignment functions
func getCurrentOnCallAssignments() ([]OnCallAssignment, error) {
	query := `
		SELECT a.id, a.schedule_id, a.user_id, a.start_time, a.end_time, a.active 
		FROM oncall_assignments a
		INNER JOIN (
			SELECT schedule_id, MAX(start_time) as max_start_time
			FROM oncall_assignments 
			GROUP BY schedule_id
		) latest ON a.schedule_id = latest.schedule_id AND a.start_time = latest.max_start_time`
	fmt.Printf("Executing query: %s\n", query)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []OnCallAssignment
	for rows.Next() {
		var assignment OnCallAssignment
		err := rows.Scan(&assignment.ID, &assignment.ScheduleID, &assignment.UserID, 
			&assignment.StartTime, &assignment.EndTime, &assignment.Active)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}
	return assignments, nil
}

func createOnCallAssignment(scheduleID int, userID int, startTime, endTime time.Time) error {
	_, err := db.Exec("INSERT INTO oncall_assignments (schedule_id, user_id, start_time, end_time, active) VALUES ($1, $2, $3, $4, $5)", 
		scheduleID, userID, startTime, endTime, true)
	return err
}

func deactivateAssignment(assignmentID int) error {
	_, err := db.Exec("UPDATE oncall_assignments SET active = false WHERE id = $1", assignmentID)
	return err
}
