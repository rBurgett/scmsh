package service

import (
	"github.com/google/uuid"
	"github.com/rBurgett/scmsh/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestGame_ExecuteMove(t *testing.T) {
	playerID1 := uuid.New()
	playerID2 := uuid.New()

	tests := []struct {
		name               string
		game               Game
		playerID           uuid.UUID
		playerCardPosition int
		targetID           uuid.UUID
		targetCardPosition int
		expected           Move
		expectedError      error
	}{
		{
			name: "valid move",
			game: Game{
				CurrentPlayer: playerID1,
				Players: []Player{
					{
						ID: playerID1,
						Deck: [][]constants.CardType{
							{constants.CardTypeSpear},
							{constants.CardTypeDagger, constants.CardTypeCrown},
						},
					},
					{
						ID: playerID2,
						Deck: [][]constants.CardType{
							{constants.CardTypeMace, constants.CardTypeSpear},
							{constants.CardTypeDagger, constants.CardTypeCrown},
						},
					},
				},
				Status: constants.GameStatusStarted,
			},
			playerID:           playerID1,
			playerCardPosition: 0,
			targetID:           playerID2,
			targetCardPosition: 1,
			expected: Move{
				Player:                   playerID1,
				PlayerCardPosition:       0,
				PlayerCardType:           constants.CardTypeSpear,
				TargetPlayer:             playerID2,
				TargetPlayerCardPosition: 1,
				TargetPlayerCardType:     constants.CardTypeDagger,
				Winner:                   playerID1,
			},
		},
		{
			name: "invalid - game not started",
			game: Game{
				CurrentPlayer: playerID1,
				Players: []Player{
					{
						ID: playerID1,
						Deck: [][]constants.CardType{
							{constants.CardTypeSpear},
							{constants.CardTypeDagger, constants.CardTypeCrown},
						},
					},
					{
						ID: playerID2,
						Deck: [][]constants.CardType{
							{constants.CardTypeMace, constants.CardTypeSpear},
							{constants.CardTypeDagger, constants.CardTypeCrown},
						},
					},
				},
				Status: constants.GameStatusOpen,
			},
			playerID:           playerID1,
			playerCardPosition: 0,
			targetID:           playerID2,
			targetCardPosition: 1,
			expectedError:      constants.ErrorGameNotStarted,
		},
		{
			name: "invalid - player wrong turn",
			game: Game{
				CurrentPlayer: playerID1,
				Players: []Player{
					{
						ID: playerID1,
						Deck: [][]constants.CardType{
							{constants.CardTypeSpear},
							{constants.CardTypeDagger, constants.CardTypeCrown},
						},
					},
					{
						ID: playerID2,
						Deck: [][]constants.CardType{
							{constants.CardTypeMace, constants.CardTypeSpear},
							{constants.CardTypeDagger, constants.CardTypeCrown},
						},
					},
				},
				Status: constants.GameStatusStarted,
			},
			playerID:           playerID2,
			playerCardPosition: 0,
			targetID:           playerID1,
			targetCardPosition: 1,
			expectedError:      constants.ErrorPlayerWrongTurn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.game.ExecuteMove(tt.playerID, tt.playerCardPosition, tt.targetID, tt.targetCardPosition)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestGame_GetPlayer(t *testing.T) {
	playerID1 := uuid.New()

	tests := []struct {
		name          string
		game          Game
		playerID      uuid.UUID
		expected      Player
		expectedError error
	}{
		{
			name: "valid player",
			game: Game{
				Players: []Player{
					{
						ID: playerID1,
					},
				},
			},
			playerID: playerID1,
			expected: Player{
				ID: playerID1,
			},
		},
		{
			name:          "invalid player",
			game:          Game{},
			playerID:      uuid.New(),
			expectedError: constants.ErrorPlayerNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.game.GetPlayer(tt.playerID)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestCreateGame(t *testing.T) {

	validPlayer := Player{
		ID:     uuid.New(),
		Name:   "Ryan",
		Secret: uuid.New(),
	}

	tests := []struct {
		name          string
		player        Player
		expected      Game
		expectedError error
	}{
		{
			name:   "good player",
			player: validPlayer,
			expected: Game{
				CurrentPlayer: uuid.Nil,
				Owner:         validPlayer.ID,
				Players: []Player{
					Player{
						ID:     validPlayer.ID,
						Name:   validPlayer.Name,
						Secret: validPlayer.Secret,
						Status: constants.PlayerStatusAccepted,
					},
				},
				Status: constants.GameStatusOpen,
			},
		},
		{
			name: "invalid player ID",
			player: Player{
				ID:     uuid.Nil,
				Name:   "Ryan",
				Secret: uuid.New(),
			},
			expectedError: constants.ErrorPlayerInvalidID,
		},
		{
			name: "invalid player name",
			player: Player{
				ID:     uuid.New(),
				Name:   "",
				Secret: uuid.New(),
			},
			expectedError: constants.ErrorPlayerInvalidName,
		},
		{
			name: "invalid player secret",
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
			output, err := CreateGame(tt.player)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, output.ID)
			output.ID = uuid.Nil
			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestDetermineMoveWinner(t *testing.T) {
	availableCards := []constants.CardType{
		constants.CardTypeDagger,
		constants.CardTypeShortSword,
		constants.CardTypeMace,
		constants.CardTypeBattleAxe,
		constants.CardTypeSpear,
		constants.CardTypeLongSword,
		constants.CardTypeArcher,
		constants.CardTypeShield,
		constants.CardTypeCrown,
	}
	for _ = range 1000000 {
		player1, _ := uuid.Parse("7b5e96e8-7b7e-44b3-adeb-e12680b04cc4")
		player2, _ := uuid.Parse("ade3aeb4-e644-47e0-b136-b7854cb4980d")
		card1 := availableCards[getRandomExclusive(0, len(availableCards))]
		card2 := availableCards[getRandomExclusive(0, len(availableCards))]
		m := Move{
			Player:               player1,
			PlayerCardType:       card1,
			TargetPlayer:         player2,
			TargetPlayerCardType: card2,
		}
		output, err := DetermineMoveWinner(m)
		require.NoError(t, err)
		if card1 < constants.CardTypeArcher && card2 < constants.CardTypeArcher {
			if card1 == card2 {
				assert.Empty(t, output)
			} else if card1 > card2 {
				assert.Equal(t, player1, output, "expected output for %v and %v to be %s", card1, card2, player1)
			} else {
				assert.Equal(t, player2, output, "expected output for %v and %v to be %s", card1, card2, player2)
			}
		} else if card1 == constants.CardTypeArcher {
			assert.Equal(t, player1, output, "expected output for %v and %v to be %s", card1, card2, player1)
		} else if card1 == constants.CardTypeCrown {
			if card2 == constants.CardTypeCrown {
				assert.Equal(t, player1, output, "expected output for %v and %v to be %s", card1, card2, player1)
			} else if card2 == constants.CardTypeShield {
				assert.Empty(t, output)
			} else if card2 == constants.CardTypeArcher {
				assert.Equal(t, player1, output, "expected output for %v and %v to be %s", card1, card2, player1)
			} else {
				assert.Equal(t, player2, output, "expected output for %v and %v to be %s", card1, card2, player2)
			}
		} else if card1 == constants.CardTypeShield {
			assert.Equal(t, player2, output, "expected output for %v and %v to be %s", card1, card2, player2)
		} else if card2 == constants.CardTypeShield {
			assert.Empty(t, output)
		} else if card2 == constants.CardTypeArcher {
			assert.Equal(t, player1, output, "expected output for %v and %v to be %s", card1, card2, player1)
		} else if card2 == constants.CardTypeCrown {
			assert.Equal(t, player1, output, "expected output for %v and %v to be %s", card1, card2, player1)
		} else {
			t.Errorf("unknown card combination %v and %v", card1, card2)
		}
	}
}

func getRandomExclusive(min int, max int) int {
	return rand.Intn(max-min) + min
}
