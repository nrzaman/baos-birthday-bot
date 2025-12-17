package bot

import (
	"bytes"
	"fmt"
	"time"

	"github.com/nrzaman/baos-birthday-bot/internal/birthday"
	"github.com/nrzaman/baos-birthday-bot/internal/interfaces"
)

// Worker handles scheduled birthday checks
type Worker struct {
	client          interfaces.DiscordClient
	birthdayService *birthday.Service
	timeProvider    interfaces.TimeProvider
	channelID       string
	stopChan        chan struct{}
}

// NewWorker creates a new Worker with the given dependencies
func NewWorker(client interfaces.DiscordClient, birthdayService *birthday.Service, timeProvider interfaces.TimeProvider, channelID string) *Worker {
	return &Worker{
		client:          client,
		birthdayService: birthdayService,
		timeProvider:    timeProvider,
		channelID:       channelID,
		stopChan:        make(chan struct{}),
	}
}

// Start begins the worker's scheduled tasks
func (w *Worker) Start() {
	// Get duration until next day at 9am
	currentTime := w.timeProvider.Now()
	fmt.Println("Current time is: " + currentTime.String())

	newTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 9, 0, 0, 0, currentTime.Location())
	fmt.Println("New time is: " + newTime.String())

	duration := newTime.Sub(currentTime)
	fmt.Println("Duration is: " + duration.String())

	// If the new time happened before the current time, wait until tomorrow at the same time
	if duration < 0 {
		newTime = newTime.Add(24 * time.Hour)
		duration = newTime.Sub(currentTime)
		fmt.Println("New duration is: " + duration.String())
	}

	for {
		select {
		case <-w.stopChan:
			return
		case <-time.After(duration):
			// Reset until the same time tomorrow
			duration = 24 * time.Hour
			w.performDailyCheck()
		}
	}
}

// Stop stops the worker
func (w *Worker) Stop() {
	close(w.stopChan)
}

// performDailyCheck performs the daily birthday check and sends messages
func (w *Worker) performDailyCheck() {
	now := w.timeProvider.Now()

	// List the monthly birthdays if it is the first of the month
	if now.Day() == 1 {
		var buffer bytes.Buffer
		buffer.WriteString("Happy " + now.Month().String() + "! Below are all the birthdays this month:\n" + w.birthdayService.ListCurrentMonthBirthdays())
		if err := w.client.SendMessage(w.channelID, buffer.String()); err != nil {
			fmt.Printf("Error sending monthly birthday message: %v\n", err)
		}
	}

	// Posts a birthday message if today is a birthday
	birthdayMessage := w.birthdayService.GetBirthdayMessage()
	if len(birthdayMessage) > 0 {
		if err := w.client.SendMessage(w.channelID, birthdayMessage); err != nil {
			fmt.Printf("Error sending birthday message: %v\n", err)
		}
	}
}
