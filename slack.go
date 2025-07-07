package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/slack-go/slack"
)

func sendSlackNotification(user *User, scheduleName string, startTime, endTime time.Time) {
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		log.Println("SLACK_TOKEN not set, skipping Slack notification")
		return
	}
	
	slackChannel := os.Getenv("SLACK_CHANNEL")
	if slackChannel == "" {
		slackChannel = "#oncall"
	}
	
	api := slack.New(slackToken)
	
	message := fmt.Sprintf("🚨 *On-Call Rotation Update*\n\n"+
		"**Schedule:** %s\n"+
		"**New On-Call Person:** %s (%s)\n"+
		"**Start Time:** %s\n"+
		"**End Time:** %s\n\n"+
		"Please ensure you're available during your on-call period!",
		scheduleName,
		user.Email,
		user.SlackHandle,
		startTime.Format("2006-01-02 15:04:05"),
		endTime.Format("2006-01-02 15:04:05"))
	
	// Try to send direct message to user first, fallback to channel
	_, _, _, err := api.SendMessage(user.SlackHandle, slack.MsgOptionText(message, false))
	if err != nil {
		// If direct message fails, send to channel
		_, _, err = api.PostMessage(slackChannel, slack.MsgOptionText(message, false))
		if err != nil {
			log.Printf("Error sending Slack notification: %v", err)
			return
		}
		log.Printf("Slack notification sent to channel %s for %s on schedule %s", slackChannel, user.Email, scheduleName)
	} else {
		log.Printf("Direct Slack notification sent to %s (%s) for schedule %s", user.Email, user.SlackHandle, scheduleName)
	}
}