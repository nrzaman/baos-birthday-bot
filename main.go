package main

import (
	"bytes"
	"encoding/json"
	_ "encoding/json"
	"flag"
	_ "flag"
	"fmt"
	_ "fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	_ "io/ioutil"
	"log"
	_ "net/http"
	"os"
	_ "os"
	"os/signal"
	_ "os/signal"
	"strconv"
	_ "strings"
	"syscall"
	_ "syscall"
	"time"
	_ "time"

	_ "github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

var birthdays People

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func check(e error) {
	if e != nil {
		log.Fatal("Error: ", e)
	}
}

type Birthday struct {
	Month int `json:"Month"`
	Day   int `json:"Day"`
}

type Person struct {
	Name     string   `json:"Name"`
	Birthday Birthday `json:"Birthday"`
}

type People struct {
	People []Person `json:"Birthdays"`
}

func ExtractBirthdays() {
	// Open the JSON config file
	content, err := os.Open("./config/birthdays.json")
	check(err)

	defer content.Close()

	// Read all contents
	byteResult, _ := io.ReadAll(content)

	var people People

	//Store contents
	json.Unmarshal(byteResult, &people)
	birthdays = people
}

func main() {
	// Create a new Discord session using the provided bot token.
	//dg, err := discordgo.New("Bot " + Token)
	//check(err)
	// Extract and store birthdays from the JSON config file
	ExtractBirthdays()

	bot, err := Connect(Token)
	check(err)

	// Register the messageCreate func as a callback for MessageCreate events.
	bot.AddHandler(messageCreate)

	// Retrieve current time to compare birthdays to
	now := time.Now()
	currentMonth := now.Month()
	currentDay := now.Day()

	for _, person := range birthdays.People {
		name := person.Name
		month := time.Month(person.Birthday.Month)
		day := person.Birthday.Day

		if int(month) == int(currentMonth) && day == currentDay && name != "Casey" {
			bot.ChannelMessageSend("962434955680579624", fmt.Sprintf("Today is %s's birthday! Please wish them a happy birthday!", name))
		} else if int(currentMonth) == 1 && currentDay == 6 {
			bot.ChannelMessageSend("962434955680579624", fmt.Sprintf("Today is the anniversary of the Capitol Riots. Nothing else special happened today."))
		}
	}

	// Cleanly close down the Discord session.
	defer bot.Close()
	// Wait here until CTRL-C or other term signal is received.
	waitUntilTermination()
}

// Connect starts a Discord session
func Connect(discordToken string) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + discordToken)
	check(err)

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsDirectMessages | discordgo.IntentsDirectMessageReactions
	// Open a websocket connection to Discord and begin listening.
	err = session.Open()
	return session, err
}

func waitUntilTermination() {
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(bot *discordgo.Session, message *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example, but it's a good practice.
	if message.Author.ID == bot.State.User.ID {
		return
	}

	if message.Content == "!upcoming" {
		// Retrieve current time to compare birthdays to
		now := time.Now()
		currentMonth := now.Month()

		// Build the string that contains the list of birthdays for the next 2 months
		var buffer bytes.Buffer
		for i := 0; i < len(birthdays.People); i++ {
			name := birthdays.People[i].Name
			month := time.Month(birthdays.People[i].Birthday.Month)
			day := strconv.Itoa(birthdays.People[i].Birthday.Day)
			if (int(month) >= int(currentMonth) && int(month) < int(currentMonth)+2) || (int(currentMonth) >= 11 && int(month) < 2) {
				buffer.WriteString(name + ", " + month.String() + " " + day + "\n")
			}
		}

		// Send the message to the channel
		bot.ChannelMessageSend(message.ChannelID, buffer.String())
	}

	// Build the string that contains the list of all configured birthdays
	if message.Content == "!all" {
		var buffer bytes.Buffer
		for i := 0; i < len(birthdays.People); i++ {
			name := birthdays.People[i].Name
			month := time.Month(birthdays.People[i].Birthday.Month)
			day := strconv.Itoa(birthdays.People[i].Birthday.Day)
			buffer.WriteString(name + ", " + month.String() + " " + day + "\n")
		}

		// Send the message to the channel
		bot.ChannelMessageSend(message.ChannelID, buffer.String())
	}
}
