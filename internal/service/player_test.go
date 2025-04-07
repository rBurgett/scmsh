package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/rBurgett/scmsh/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlayer_Validate(t *testing.T) {
	tests := []struct {
		name          string
		player        Player
		expectedError error
	}{
		{
			name: "valid player",
			player: Player{
				ID:     uuid.New(),
				Name:   "Ryan",
				Secret: uuid.New(),
			},
		},
		{
			name: "invalid id",
			player: Player{
				ID:     uuid.Nil,
				Name:   "Ryan",
				Secret: uuid.New(),
			},
			expectedError: constants.ErrorPlayerInvalidID,
		},
		{
			name: "invalid name empty",
			player: Player{
				ID:     uuid.New(),
				Name:   "",
				Secret: uuid.New(),
			},
			expectedError: constants.ErrorPlayerInvalidName,
		},
		{
			name: "invalid name too long",
			player: Player{
				ID:     uuid.New(),
				Name:   "1234567890123456789012345678901234567890",
				Secret: uuid.New(),
			},
			expectedError: constants.ErrorPlayerInvalidName,
		},
		{
			name: "invalid secret",
			player: Player{
				ID:     uuid.New(),
				Name:   "Ryan",
				Secret: uuid.Nil,
			},
			expectedError: constants.ErrorPlayerInvalidSecret,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.player.Validate()
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestPlayer_ValidateDeck(t *testing.T) {

}

func TestPlayer_GetCard(t *testing.T) {

}

func TestCreatePlayer(t *testing.T) {
	tests := []struct {
		name          string
		playerName    string
		expected      Player
		expectedError error
	}{
		{
			name:       "valid player name",
			playerName: "Ryan",
			expected: Player{
				Name:   "Ryan",
				Status: constants.PlayerStatusUnaffiliated,
			},
		},
		{
			name:       "valid player name with whitespace",
			playerName: "      Ryan      ",
			expected: Player{
				Name:   "Ryan",
				Status: constants.PlayerStatusUnaffiliated,
			},
		},
		{
			name:          "invalid player name all whitespace",
			playerName:    "            ",
			expectedError: constants.ErrorPlayerInvalidName,
		},
		{
			name:          "invalid player name empty",
			playerName:    "",
			expectedError: constants.ErrorPlayerInvalidName,
		},
		{
			name:          "invalid player name too long",
			playerName:    "1234567890123456789012345678901234567890",
			expectedError: constants.ErrorPlayerInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := CreatePlayer(tt.playerName)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, output.ID)
			output.ID = uuid.Nil
			assert.NotEmpty(t, output.Secret)
			output.Secret = uuid.Nil
			assert.Equal(t, tt.expected, output)
		})
	}
}
