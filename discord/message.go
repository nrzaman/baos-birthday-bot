package discord

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nrzaman/baos-birthday-bot/util"
	"strconv"
	"time"
)

// MessageCreate This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func MessageCreate(bot *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example, but it's a good practice.
	if message.Author.ID == bot.State.User.ID {
		return
	}

	// Lists upcoming birthday for the current month
	if message.Content == "!month" {
		fmt.Println("Listing the current month's birthdays.")
		// Send the message to the channel
		_, err := bot.ChannelMessageSend(message.ChannelID, ListCurrentMonthBirthdays())
		util.Check(err)
	}

	// Build the string that contains the list of all configured birthdays
	if message.Content == "!all" {
		fmt.Println("Listing all birthdays.")
		// Send the message to the channel
		_, err := bot.ChannelMessageSend(message.ChannelID, ListAllBirthdays())
		util.Check(err)
	}
}

// SendBirthdayMessage This function will send a birthday message to the discord if today is a birthday.
func SendBirthdayMessage(bot *discordgo.Session, birthdayMessage string) {
	if len(birthdayMessage) != 0 {
		_, err := bot.ChannelMessageSend(util.GeneralChannelID, birthdayMessage)
		util.Check(err)
	}
}

// ListCurrentMonthBirthdays Lists the current month's birthdays.
func ListCurrentMonthBirthdays() string {
	// Retrieve current month to compare birthdays to
	currentMonth := util.GetCurrentMonth()

	// Build the string that contains the list of birthdays for the next 2 months
	var buffer bytes.Buffer
	for _, person := range util.Birthdays.People {
		// Extract the name and date of birth
		name := person.Name
		month := time.Month(person.Birthday.Month)
		day := person.Birthday.Day

		// Check whether the birthday is within the current month and next month, and add
		// to the string buffer if so
		if int(month) == int(currentMonth) {
			buffer.WriteString(name + ", " + month.String() + " " + strconv.Itoa(day) + "\n")
		}
	}

	return buffer.String()
}

// ListAllBirthdays Lists all birthdays.
func ListAllBirthdays() string {
	var buffer bytes.Buffer
	for _, person := range util.Birthdays.People {
		name := person.Name
		month := time.Month(person.Birthday.Month)
		day := strconv.Itoa(person.Birthday.Day)
		buffer.WriteString(name + ", " + month.String() + " " + day + "\n")
	}

	return buffer.String()
}
