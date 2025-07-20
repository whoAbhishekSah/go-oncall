package main

import (
	"log"
	"time"
)

func scheduleChecker() {
	ticker := time.NewTicker(10 * time.Second) // Check every 10 seconds for better granularity
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
		if now.After(schedule.StartTime) && now.Before(schedule.EndTime) {
			currentAssignment := getCurrentAssignmentForSchedule(schedule.ID)
			
			if currentAssignment == nil || shouldRotate(currentAssignment, schedule.RotationPeriod, now) {
				if currentAssignment != nil {
					deactivateAssignment(currentAssignment.ID)
				}
				
				nextUserID := getNextOnCallUser(schedule)
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
		return nil
	}
	
	for _, assignment := range assignments {
		if assignment.ScheduleID == scheduleID {
			return &assignment
		}
	}
	return nil
}

func shouldRotate(assignment *OnCallAssignment, rotationPeriod int, now time.Time) bool {
	return now.After(assignment.EndTime)
}

func getNextOnCallUser(schedule Schedule) int {
	if len(schedule.Participants) == 0 {
		return 0
	}
	
	currentAssignment := getCurrentAssignmentForSchedule(schedule.ID)
	if currentAssignment == nil {
		return schedule.Participants[0]
	}
	
	for i, participant := range schedule.Participants {
		if participant == currentAssignment.UserID {
			return schedule.Participants[(i+1)%len(schedule.Participants)]
		}
	}
	
	return schedule.Participants[0]
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
