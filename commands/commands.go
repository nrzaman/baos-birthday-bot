package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nrzaman/baos-birthday-bot/birthdayUtil"
	"io"
	"os"
	"strconv"
	"time"
)

var Birthdays birthdayUtil.People

// ExtractBirthdays This function extracts birthdays from the JSON config file and stores them
// to be referenced later.
func ExtractBirthdays() {
	// Open the JSON config file
	content, err := os.Open("./config/birthdays.json")
	birthdayUtil.Check(err)

	defer content.Close()

	// Read all contents
	byteResult, _ := io.ReadAll(content)

	// Create result variable
	var people birthdayUtil.People

	// Store contents
	json.Unmarshal(byteResult, &people)
	Birthdays = people
}

// MessageCreate This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func MessageCreate(bot *discordgo.Session, message *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example, but it's a good practice.
	if message.Author.ID == bot.State.User.ID {
		return
	}

	// Lists upcoming birthday for the next 2 months (inclusive of current month)
	if message.Content == "!upcoming" {
		fmt.Println("Listing upcoming birthdays.")
		// Retrieve current month to compare birthdays to
		currentMonth := birthdayUtil.GetCurrentMonth()

		// Build the string that contains the list of birthdays for the next 2 months
		var buffer bytes.Buffer
		for i := 0; i < len(Birthdays.People); i++ {
			// Extract the name and date of birth
			name := Birthdays.People[i].Name
			month := time.Month(Birthdays.People[i].Birthday.Month)
			day := strconv.Itoa(Birthdays.People[i].Birthday.Day)

			// Check whether the birthday is within the current month and next month, and add
			// to the string buffer if so
			if (int(month) >= int(currentMonth) && int(month) < int(currentMonth)+2) || (int(currentMonth) >= 11 && int(month) < 2) {
				buffer.WriteString(name + ", " + month.String() + " " + day + "\n")
			}
		}

		// Send the message to the channel
		bot.ChannelMessageSend("962434955680579624", buffer.String())
	}

	// Build the string that contains the list of all configured birthdays
	if message.Content == "!all" {
		fmt.Println("Listing all birthdays.")
		var buffer bytes.Buffer
		for i := 0; i < len(Birthdays.People); i++ {
			name := Birthdays.People[i].Name
			month := time.Month(Birthdays.People[i].Birthday.Month)
			day := strconv.Itoa(Birthdays.People[i].Birthday.Day)
			buffer.WriteString(name + ", " + month.String() + " " + day + "\n")
		}

		// Send the message to the channel
		bot.ChannelMessageSend("962434955680579624", buffer.String())
	}
}
