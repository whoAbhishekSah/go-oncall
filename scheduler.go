package main

import (
	"fmt"
	"log"
	"time"
)

func scheduleChecker() {
	ticker := time.NewTicker(5 * time.Second) // Check every 10 seconds for better granularity
	defer ticker.Stop()
	
	log.Println("Schedule checker started (checking every 10 seconds)")
	
	for {
		select {
		case <-ticker.C:
			log.Println("Checking and Updating schedules")
			checkAndUpdateSchedules()
		}
	}
}

func checkAndUpdateSchedules() {
	schedules, err := getSchedules()
	if err != nil {
		log.Printf("Error getting schedules: %v", err)
		return
	}
	
	now := time.Now()
	
	for _, schedule := range schedules {
		if err := validateScheduleParticipants(schedule); err != nil {
			log.Printf("Invalid schedule participants: %v", err)
			continue
		}

		fmt.Println("is now after schedule start time", now.After(schedule.StartTime))
		fmt.Println("is now before schedule end time", now.Before(schedule.EndTime))

		if now.After(schedule.StartTime) && now.Before(schedule.EndTime) {
			fmt.Println("now is between the schedule start and end")
			currentAssignment := getCurrentAssignmentForSchedule(schedule.ID)
			
			if currentAssignment == nil || shouldRotate(currentAssignment, schedule.RotationPeriod, now) {
				if currentAssignment != nil {
					deactivateAssignment(currentAssignment.ID)
				}
				
				nextUserID := getNextOnCallUser(schedule)
				fmt.Println("next user id", nextUserID)
				if nextUserID != 0 {
					rotationStart := calculateRotationStart(schedule.StartTime, schedule.RotationPeriod, now)
					rotationEnd := rotationStart.Add(time.Duration(schedule.RotationPeriod) * time.Second)
					
					err := createOnCallAssignment(schedule.ID, nextUserID, rotationStart, rotationEnd)
					if err != nil {
						log.Printf("Error creating assignment: %v", err)
						continue
					}
					
					user, err := getUserByID(nextUserID)
					if err != nil {
						log.Printf("Error getting user: %v", err)
						continue
					}
					
					log.Printf("New on-call assignment: %s (%s) for schedule %s", user.Email, user.SlackHandle, schedule.Name)
					
					go sendSlackNotification(user, schedule.Name, rotationStart, rotationEnd)
				}
			}
		}
	}
}

func getCurrentAssignmentForSchedule(scheduleID int) *OnCallAssignment {
	assignments, err := getCurrentOnCallAssignments()
	if err != nil {
		log.Printf("Error getting current assignments: %v", err)
		return nil
	}

	fmt.Printf("Found %d current assignments\n", len(assignments))
	for _, assignment := range assignments {
		fmt.Printf("Checking assignment - ID: %d, ScheduleID: %d, UserID: %d, StartTime: %v, EndTime: %v, Active: %v\n",
			assignment.ID, assignment.ScheduleID, assignment.UserID, assignment.StartTime, assignment.EndTime, assignment.Active)
		if assignment.ScheduleID == scheduleID {
			return &assignment
		}
	}
	return nil
}

func shouldRotate(assignment *OnCallAssignment, rotationPeriod int, now time.Time) bool {
	shouldRotate := now.After(assignment.EndTime)
	fmt.Printf("Checking rotation - Now: %v, Assignment End: %v, Should Rotate: %v\n", 
		now, assignment.EndTime, shouldRotate)
	return shouldRotate
}

func validateScheduleParticipants(schedule Schedule) error {
	for _, userID := range schedule.Participants {
		user, err := getUserByID(userID)
		if err != nil {
			return fmt.Errorf("participant user %d not found: %v", userID, err)
		}
		fmt.Printf("Validated participant: ID %d, Email %s\n", user.ID, user.Email)
	}
	return nil
}

func getNextOnCallUser(schedule Schedule) int {
	if len(schedule.Participants) == 0 {
		fmt.Println("No participants in schedule")
		return 0
	}

	fmt.Printf("Schedule participants: %v\n", schedule.Participants)

	currentAssignment := getCurrentAssignmentForSchedule(schedule.ID)
	if currentAssignment == nil {
		fmt.Printf("No current assignment, starting rotation with first user: %d\n", schedule.Participants[0])
		return schedule.Participants[0]
	}

	fmt.Printf("Current assignment: ID=%d, UserID=%d, StartTime=%v, EndTime=%v, Active=%v\n",
		currentAssignment.ID, currentAssignment.UserID, currentAssignment.StartTime,
		currentAssignment.EndTime, currentAssignment.Active)

	// Find current user's position
	currentIndex := -1
	for i, participant := range schedule.Participants {
		if participant == currentAssignment.UserID {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		fmt.Printf("Current user %d not found in participants list %v, starting with first user\n",
			currentAssignment.UserID, schedule.Participants)
		return schedule.Participants[0]
	}

	// Get next user
	nextIndex := (currentIndex + 1) % len(schedule.Participants)
	nextUser := schedule.Participants[nextIndex]
	fmt.Printf("Current user %d at index %d, rotating to next user %d at index %d\n",
		currentAssignment.UserID, currentIndex, nextUser, nextIndex)
	return nextUser
}

func calculateRotationStart(scheduleStart time.Time, rotationPeriod int, now time.Time) time.Time {
	if now.Before(scheduleStart) {
		return scheduleStart
	}
	
	elapsed := now.Sub(scheduleStart)
	rotationDuration := time.Duration(rotationPeriod) * time.Second
	
	rotationsSinceStart := elapsed / rotationDuration
	
	return scheduleStart.Add(rotationsSinceStart * rotationDuration)
}
