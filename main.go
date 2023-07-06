package main

import (
	_ "encoding/json"
	"flag"
	_ "flag"
	"fmt"
	_ "fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/bwmarrin/discordgo"
	"github.com/nrzaman/baos-birthday-bot/birthdayUtil"
	_ "github.com/nrzaman/baos-birthday-bot/birthdayUtil"
	_ "io/ioutil"
	_ "net/http"
	"os"
	_ "os"
	"os/signal"
	_ "os/signal"
	_ "strings"
	"syscall"
	_ "syscall"
	"time"
	_ "time"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Extract and store birthdays from the JSON config file
	ExtractBirthdays()

	// Create a new Discord session using the provided bot token.
	bot, err := Connect(Token)
	birthdayUtil.Check(err)

	for _, person := range birthdays.People {
		// Get each person's name and date of birth
		name := person.Name
		month := time.Month(person.Birthday.Month)
		day := person.Birthday.Day

		// Check the person's birthday against current day
		if birthdayUtil.IsBirthdayCurrentDay(int(month), day) && name != "Casey" {
			bot.ChannelMessageSend("962434955680579624", fmt.Sprintf("Today is %s's birthday! Please wish them a happy birthday!", name))
		} else if birthdayUtil.IsBirthdayCurrentDay(1, 6) {
			// Special handling for Casey
			bot.ChannelMessageSend("962434955680579624", fmt.Sprintf("Today is the anniversary of the Capitol Riots. Nothing else special happened today."))
		}
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	bot.AddHandler(MessageCreate)

	// Wait here until CTRL-C or other term signal is received.
	waitUntilTermination()
	bot.Close()
	fmt.Println("Bot terminated.")
}

// Connect starts a Discord session
func Connect(discordToken string) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + discordToken)
	birthdayUtil.Check(err)

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsMessageContent
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
