package discord

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nrzaman/baos-birthday-bot/util"
	"time"
)

// Worker Checks and posts to the Discord server when it's someone's birthday.
func Worker(bot *discordgo.Session) {
	// Get duration until next day at 9am. Should only need to happen once.
	currentTime := time.Now()
	fmt.Println("Current time is: " + currentTime.String())
	newTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 9, 0, 0, 0, currentTime.Location())
	fmt.Println("New time is: " + newTime.String())
	duration := newTime.Sub(currentTime)
	fmt.Println("Duration is: " + duration.String())

	// If the new time happened before the current time, the worker needs to wait until tomorrow at the same time.
	if duration < 0 {
		newTime = newTime.Add(24 * time.Hour)
		duration = newTime.Sub(currentTime)
		fmt.Println("New duration is: " + duration.String())
	}

	for {
		time.Sleep(duration)
		// Reset until the same time tomorrow.
		duration = 24 * time.Hour
		// List the monthly birthdays if it is the first of the month
		if time.Now().Day() == 1 {
			var buffer bytes.Buffer
			buffer.WriteString("Happy " + time.Now().Month().String() + "! Below are all the birthdays this month:\n" + ListCurrentMonthBirthdays())
			_, err := bot.ChannelMessageSend("962434955680579624", buffer.String())
			util.Check(err)
		}

		// Posts a birthday message if today is a birthday.
		var birthdayMessage = util.BirthdayFinderMessage()
		SendBirthdayMessage(bot, birthdayMessage)
	}
}
