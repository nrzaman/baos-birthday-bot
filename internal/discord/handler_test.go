package bot_test

import (
	"testing"
	"time"

	bot "github.com/nrzaman/baos-birthday-bot/internal/discord"
	"github.com/nrzaman/baos-birthday-bot/util"
)

// MockDiscordClient is a mock implementation of DiscordClient for testing
type MockDiscordClient struct {
	SentMessages []SentMessage
	SendError    error
}

type SentMessage struct {
	ChannelID string
	Message   string
}

func (m *MockDiscordClient) SendMessage(channelID string, message string) error {
	if m.SendError != nil {
		return m.SendError
	}
	m.SentMessages = append(m.SentMessages, SentMessage{
		ChannelID: channelID,
		Message:   message,
	})
	return nil
}

func (m *MockDiscordClient) AddHandler(handler interface{}) {
	// No-op for testing
}

func (m *MockDiscordClient) Close() error {
	return nil
}

// MockTimeProvider for testing
type MockTimeProvider struct {
	CurrentTime time.Time
}

func (m *MockTimeProvider) Now() time.Time {
	return m.CurrentTime
}

func (m *MockTimeProvider) Month() time.Month {
	return m.CurrentTime.Month()
}

func (m *MockTimeProvider) Day() int {
	return m.CurrentTime.Day()
}

// MockBirthdayService for testing
type MockBirthdayService struct {
	BirthdayMessage       string
	CurrentMonthBirthdays string
	AllBirthdays          string
	Birthdays             []struct {
		Name  string
		Month int
		Day   int
	}
}

func (m *MockBirthdayService) IsBirthdayToday(month int, day int) bool {
	return false
}

func (m *MockBirthdayService) GetBirthdayMessage() string {
	return m.BirthdayMessage
}

func (m *MockBirthdayService) ListCurrentMonthBirthdays() string {
	return m.CurrentMonthBirthdays
}

func (m *MockBirthdayService) ListAllBirthdays() string {
	return m.AllBirthdays
}

func (m *MockBirthdayService) GetBirthdays() util.People {
	people := make([]util.Person, len(m.Birthdays))
	for i, b := range m.Birthdays {
		people[i] = util.Person{
			Name: b.Name,
			Birthday: util.Birthday{
				Month: b.Month,
				Day:   b.Day,
			},
		}
	}
	return util.People{People: people}
}

func TestSendBirthdayMessage(t *testing.T) {
	// Arrange
	mockClient := &MockDiscordClient{}
	birthdayService := &MockBirthdayService{}
	handler := bot.NewHandler(mockClient, birthdayService)

	tests := []struct {
		name      string
		channelID string
		message   string
		wantSent  bool
	}{
		{"Send valid message", "123456", "Happy Birthday!", true},
		{"Empty message should not send", "123456", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset sent messages
			mockClient.SentMessages = []SentMessage{}

			// Act
			err := handler.SendBirthdayMessage(tt.channelID, tt.message)

			// Assert
			if err != nil {
				t.Errorf("SendBirthdayMessage() returned error: %v", err)
			}

			if tt.wantSent {
				if len(mockClient.SentMessages) != 1 {
					t.Errorf("Expected 1 message sent, got %d", len(mockClient.SentMessages))
				} else {
					sent := mockClient.SentMessages[0]
					if sent.ChannelID != tt.channelID {
						t.Errorf("Message sent to channel %q; want %q", sent.ChannelID, tt.channelID)
					}
					if sent.Message != tt.message {
						t.Errorf("Message content = %q; want %q", sent.Message, tt.message)
					}
				}
			} else {
				if len(mockClient.SentMessages) != 0 {
					t.Errorf("Expected no messages sent, got %d", len(mockClient.SentMessages))
				}
			}
		})
	}
}
