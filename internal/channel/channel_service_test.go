package channel_test

import (
	"testing"

	"github.com/nrzaman/baos-birthday-bot/internal/channel"
)

// MockFileReader is a mock implementation of FileReader for testing
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

func TestLoadChannels(t *testing.T) {
	// Arrange
	channelsJSON := `{
		"Channels": [
			{"Name": "general", "ID": "123456789"},
			{"Name": "random", "ID": "987654321"}
		]
	}`

	mockFileReader := &MockFileReader{
		Data: map[string][]byte{
			"./config/channels.json": []byte(channelsJSON),
		},
	}

	service := channel.NewService(mockFileReader)

	// Act
	err := service.LoadChannels("./config/channels.json")

	// Assert
	if err != nil {
		t.Errorf("LoadChannels() returned error: %v", err)
	}
}

func TestGetGeneralChannelID(t *testing.T) {
	// Arrange
	channelsJSON := `{
		"Channels": [
			{"Name": "general", "ID": "123456789"},
			{"Name": "random", "ID": "987654321"}
		]
	}`

	mockFileReader := &MockFileReader{
		Data: map[string][]byte{
			"./config/channels.json": []byte(channelsJSON),
		},
	}

	service := channel.NewService(mockFileReader)
	err := service.LoadChannels("./config/channels.json")
	if err != nil {
		t.Fatalf("Failed to load channels: %v", err)
	}

	// Act
	channelID := service.GetGeneralChannelID()

	// Assert
	expected := "123456789"
	if channelID != expected {
		t.Errorf("GetGeneralChannelID() = %q; want %q", channelID, expected)
	}
}

func TestGetChannelByName(t *testing.T) {
	// Arrange
	channelsJSON := `{
		"Channels": [
			{"Name": "general", "ID": "123456789"},
			{"Name": "random", "ID": "987654321"}
		]
	}`

	mockFileReader := &MockFileReader{
		Data: map[string][]byte{
			"./config/channels.json": []byte(channelsJSON),
		},
	}

	service := channel.NewService(mockFileReader)
	err := service.LoadChannels("./config/channels.json")
	if err != nil {
		t.Fatalf("Failed to load channels: %v", err)
	}

	tests := []struct {
		name        string
		channelName string
		expectedID  string
	}{
		{"Get general channel", "general", "123456789"},
		{"Get random channel", "random", "987654321"},
		{"Get non-existent channel", "nonexistent", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			channelID := service.GetChannelByName(tt.channelName)

			// Assert
			if channelID != tt.expectedID {
				t.Errorf("GetChannelByName(%q) = %q; want %q", tt.channelName, channelID, tt.expectedID)
			}
		})
	}
}

func TestGetGeneralChannelID_NotFound(t *testing.T) {
	// Arrange
	channelsJSON := `{
		"Channels": [
			{"Name": "random", "ID": "987654321"}
		]
	}`

	mockFileReader := &MockFileReader{
		Data: map[string][]byte{
			"./config/channels.json": []byte(channelsJSON),
		},
	}

	service := channel.NewService(mockFileReader)
	err := service.LoadChannels("./config/channels.json")
	if err != nil {
		t.Fatalf("Failed to load channels: %v", err)
	}

	// Act
	channelID := service.GetGeneralChannelID()

	// Assert
	if channelID != "" {
		t.Errorf("GetGeneralChannelID() = %q; want empty string when general channel not found", channelID)
	}
}
