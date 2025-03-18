package service

import (
	"github.com/google/uuid"
	"github.com/rBurgett/scmsh/internal/constants"
	"strings"
)

type Player struct {
	ID     uuid.UUID
	Deck   [][]constants.CardType
	Name   string
	Secret uuid.UUID
	Status constants.PlayerStatus
}

func (p *Player) Validate() error {
	if p.ID == uuid.Nil {
		return constants.ErrorPlayerInvalidID
	}
	if p.Name == "" {
		return constants.ErrorPlayerInvalidName
	}
	if p.Secret == uuid.Nil {
		return constants.ErrorPlayerInvalidSecret
	}

	return nil
}

func (p *Player) ValidateDeck() error {
	cardCounts := map[constants.CardType]int{}
	if len(p.Deck) != constants.DeckCount {
		return constants.ErrorInvalidDeckCounts
	}

	for _, cards := range p.Deck {
		if len(cards) != constants.DeckCardCount {
			return constants.ErrorInvalidDeckCounts
		}
		for _, card := range cards {
			cardCounts[card]++
		}
	}

	for k, v := range cardCounts {
		if v != constants.GetCardCount(k) {
			return constants.ErrorInvalidCardCount
		}
	}

	return nil
}

func (p *Player) GetCard(position int) (constants.CardType, error) {
	if position < 0 || position > len(p.Deck)-1 {
		return 0, constants.ErrorInvalidStack
	}

	stack := p.Deck[position]
	if len(stack) == 0 {
		return 0, constants.ErrorEmptyStack
	}

	return stack[0], nil
}

func CreatePlayer(name string) (Player, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Player{}, constants.ErrorPlayerInvalidName
	}

	return Player{
		ID:     uuid.New(),
		Name:   name,
		Secret: uuid.New(),
		Status: constants.PlayerStatusUnaffiliated,
	}, nil
}
