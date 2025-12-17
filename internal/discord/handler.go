package bot

import (
	"fmt"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nrzaman/baos-birthday-bot/internal/birthday"
	"github.com/nrzaman/baos-birthday-bot/internal/interfaces"
)

// Handler handles Discord message events with injected dependencies
type Handler struct {
	client          interfaces.DiscordClient
	birthdayService *birthday.Service
}

// NewHandler creates a new Handler with the given dependencies
func NewHandler(client interfaces.DiscordClient, birthdayService *birthday.Service) *Handler {
	return &Handler{
		client:          client,
		birthdayService: birthdayService,
	}
}

// HandleSlashCommand processes slash command interactions
func (h *Handler) HandleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandName := i.ApplicationCommandData().Name

	var response string
	switch commandName {
	case "month":
		fmt.Println("Slash command: Listing the current month's birthdays.")
		response = h.birthdayService.ListCurrentMonthBirthdays()
		if response == "" {
			response = "No birthdays this month!"
		}

	case "all":
		fmt.Println("Slash command: Listing all birthdays.")
		response = h.birthdayService.ListAllBirthdays()
		if response == "" {
			response = "No birthdays configured!"
		}

	case "next":
		fmt.Println("Slash command: Finding next birthday.")
		response = h.getNextBirthday()
		if response == "" {
			response = "No upcoming birthdays found!"
		}

	default:
		response = "Unknown command"
	}

	// Respond to the interaction
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		fmt.Printf("Error responding to slash command: %v\n", err)
	}
}

// HandleMessage processes incoming Discord messages (backward compatibility for !commands)
func (h *Handler) HandleMessage(bot *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if message.Author.ID == bot.State.User.ID {
		return
	}

	// Lists upcoming birthday for the current month
	if message.Content == "!month" {
		fmt.Println("Legacy command: Listing the current month's birthdays.")
		response := h.birthdayService.ListCurrentMonthBirthdays()
		if response == "" {
			response = "No birthdays this month! (Tip: Use /month for slash commands)"
		}
		if err := h.client.SendMessage(message.ChannelID, response); err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		}
	}

	// Build the string that contains the list of all configured birthdays
	if message.Content == "!all" {
		fmt.Println("Legacy command: Listing all birthdays.")
		response := h.birthdayService.ListAllBirthdays()
		if response == "" {
			response = "No birthdays configured! (Tip: Use /all for slash commands)"
		}
		if err := h.client.SendMessage(message.ChannelID, response); err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		}
	}

	// Next birthday
	if message.Content == "!next" {
		fmt.Println("Legacy command: Finding next birthday.")
		response := h.getNextBirthday()
		if response == "" {
			response = "No upcoming birthdays found! (Tip: Use /next for slash commands)"
		}
		if err := h.client.SendMessage(message.ChannelID, response); err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		}
	}
}

// getNextBirthday finds and returns the next upcoming birthday
func (h *Handler) getNextBirthday() string {
	birthdays := h.birthdayService.GetBirthdays()
	if len(birthdays.People) == 0 {
		return ""
	}

	now := time.Now()
	type bdayWithDate struct {
		name      string
		month     time.Month
		day       int
		daysUntil int
	}

	var upcoming []bdayWithDate

	for _, person := range birthdays.People {
		month := time.Month(person.Birthday.Month)
		day := person.Birthday.Day

		// Calculate this year's birthday
		thisYear := time.Date(now.Year(), month, day, 0, 0, 0, 0, now.Location())

		var daysUntil int
		if thisYear.Before(now) {
			// Birthday already passed this year, calculate for next year
			nextYear := time.Date(now.Year()+1, month, day, 0, 0, 0, 0, now.Location())
			daysUntil = int(nextYear.Sub(now).Hours() / 24)
		} else {
			daysUntil = int(thisYear.Sub(now).Hours() / 24)
		}

		upcoming = append(upcoming, bdayWithDate{
			name:      person.Name,
			month:     month,
			day:       day,
			daysUntil: daysUntil,
		})
	}

	// Sort by days until birthday
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].daysUntil < upcoming[j].daysUntil
	})

	// Get the next birthday (or multiple if on same day)
	next := upcoming[0]
	result := fmt.Sprintf("Next birthday: %s on %s %d", next.name, next.month.String(), next.day)

	if next.daysUntil == 0 {
		result += " (Today! 🎉)"
	} else if next.daysUntil == 1 {
		result += " (Tomorrow!)"
	} else {
		result += fmt.Sprintf(" (in %d days)", next.daysUntil)
	}

	// Check if there are multiple birthdays on the same day
	for i := 1; i < len(upcoming) && upcoming[i].daysUntil == next.daysUntil; i++ {
		result += fmt.Sprintf("\nAlso: %s", upcoming[i].name)
	}

	return result
}

// SendBirthdayMessage sends a birthday message to the specified channel
func (h *Handler) SendBirthdayMessage(channelID string, message string) error {
	if len(message) == 0 {
		return nil
	}
	return h.client.SendMessage(channelID, message)
}
