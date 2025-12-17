package bot_test

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nrzaman/baos-birthday-bot/internal/birthday"
	bot "github.com/nrzaman/baos-birthday-bot/internal/discord"
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

// MockFileReader for testing
type MockFileReader struct {
	Data map[string][]byte
	Err  error
}

func (m *MockFileReader) ReadFile(path string) ([]byte, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Data[path], nil
}

func TestSendBirthdayMessage(t *testing.T) {
	// Arrange
	mockClient := &MockDiscordClient{}
	mockTime := &MockTimeProvider{}
	mockFileReader := &MockFileReader{}

	birthdayService := birthday.NewService(mockTime, mockFileReader)
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

func TestHandleMessage_MonthCommand(t *testing.T) {
	// Arrange
	mockClient := &MockDiscordClient{}
	mockTime := &MockTimeProvider{
		CurrentTime: time.Date(2024, time.March, 15, 10, 0, 0, 0, time.UTC),
	}

	birthdaysJSON := `{
		"Birthdays": [
			{"Name": "Alice", "Birthday": {"Month": 3, "Day": 15}},
			{"Name": "Bob", "Birthday": {"Month": 4, "Day": 20}}
		]
	}`

	mockFileReader := &MockFileReader{
		Data: map[string][]byte{
			"./config/birthdays.json": []byte(birthdaysJSON),
		},
	}

	birthdayService := birthday.NewService(mockTime, mockFileReader)
	birthdayService.LoadBirthdays("./config/birthdays.json")

	handler := bot.NewHandler(mockClient, birthdayService)

	// Create mock Discord message
	session := &discordgo.Session{}
	session.State = discordgo.NewState()
	session.State.User = &discordgo.User{ID: "bot123"}
	message := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Content:   "!month",
			ChannelID: "channel456",
			Author:    &discordgo.User{ID: "user789"},
		},
	}

	// Act
	handler.HandleMessage(session, message)

	// Assert
	if len(mockClient.SentMessages) != 1 {
		t.Fatalf("Expected 1 message sent, got %d", len(mockClient.SentMessages))
	}

	sent := mockClient.SentMessages[0]
	if sent.ChannelID != "channel456" {
		t.Errorf("Message sent to channel %q; want 'channel456'", sent.ChannelID)
	}

	// Should only include March birthdays
	if !contains(sent.Message, "Alice") {
		t.Errorf("Response should include Alice (March birthday)")
	}
	if contains(sent.Message, "Bob") {
		t.Errorf("Response should not include Bob (April birthday)")
	}
}

func TestHandleMessage_AllCommand(t *testing.T) {
	// Arrange
	mockClient := &MockDiscordClient{}
	mockTime := &MockTimeProvider{}

	birthdaysJSON := `{
		"Birthdays": [
			{"Name": "Alice", "Birthday": {"Month": 3, "Day": 15}},
			{"Name": "Bob", "Birthday": {"Month": 4, "Day": 20}}
		]
	}`

	mockFileReader := &MockFileReader{
		Data: map[string][]byte{
			"./config/birthdays.json": []byte(birthdaysJSON),
		},
	}

	birthdayService := birthday.NewService(mockTime, mockFileReader)
	birthdayService.LoadBirthdays("./config/birthdays.json")

	handler := bot.NewHandler(mockClient, birthdayService)

	// Create mock Discord message
	session := &discordgo.Session{}
	session.State = discordgo.NewState()
	session.State.User = &discordgo.User{ID: "bot123"}
	message := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Content:   "!all",
			ChannelID: "channel456",
			Author:    &discordgo.User{ID: "user789"},
		},
	}

	// Act
	handler.HandleMessage(session, message)

	// Assert
	if len(mockClient.SentMessages) != 1 {
		t.Fatalf("Expected 1 message sent, got %d", len(mockClient.SentMessages))
	}

	sent := mockClient.SentMessages[0]

	// Should include all birthdays
	if !contains(sent.Message, "Alice") || !contains(sent.Message, "Bob") {
		t.Errorf("Response should include both Alice and Bob")
	}
}

func TestHandleMessage_IgnoresBotMessages(t *testing.T) {
	// Arrange
	mockClient := &MockDiscordClient{}
	mockTime := &MockTimeProvider{}
	mockFileReader := &MockFileReader{}

	birthdayService := birthday.NewService(mockTime, mockFileReader)
	handler := bot.NewHandler(mockClient, birthdayService)

	// Create mock Discord message from the bot itself
	session := &discordgo.Session{}
	session.State = discordgo.NewState()
	session.State.User = &discordgo.User{ID: "bot123"}
	message := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Content:   "!month",
			ChannelID: "channel456",
			Author:    &discordgo.User{ID: "bot123"}, // Same as bot's ID
		},
	}

	// Act
	handler.HandleMessage(session, message)

	// Assert
	if len(mockClient.SentMessages) != 0 {
		t.Errorf("Bot should ignore its own messages, but sent %d messages", len(mockClient.SentMessages))
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
