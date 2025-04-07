package service

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rBurgett/scmsh/internal/constants"
)

type Game struct {
	ID            uuid.UUID
	CurrentPlayer uuid.UUID
	Moves         []Move
	Owner         uuid.UUID
	Players       []Player
	Status        constants.GameStatus
	CreatedAt     time.Time
}

func (g *Game) ExecuteMove(playerID uuid.UUID, playerCardPosition int, targetID uuid.UUID, targetCardPosition int) (Move, error) {
	if g.Status != constants.GameStatusStarted {
		return Move{}, constants.ErrorGameNotStarted
	}
	if playerID != g.CurrentPlayer {
		return Move{}, constants.ErrorPlayerWrongTurn
	}

	p, err := g.GetPlayer(playerID)
	if err != nil {
		return Move{}, err
	}
	pc, err := p.GetCard(playerCardPosition)
	if err != nil {
		return Move{}, err
	}

	t, err := g.GetPlayer(targetID)
	if err != nil {
		return Move{}, err
	}
	tc, err := t.GetCard(targetCardPosition)
	if err != nil {
		return Move{}, err
	}

	m := Move{
		Player:                   playerID,
		PlayerCardPosition:       playerCardPosition,
		PlayerCardType:           pc,
		TargetPlayer:             targetID,
		TargetPlayerCardPosition: targetCardPosition,
		TargetPlayerCardType:     tc,
	}

	winner, err := DetermineMoveWinner(m)
	if err != nil {
		return Move{}, err
	}

	m.Winner = winner

	// remove card(s) from losers
	// handle a win

	g.Moves = append(g.Moves, m)

	return m, nil
}

func (g *Game) GetPlayer(id uuid.UUID) (Player, error) {
	for _, p := range g.Players {
		if p.ID == id {
			return p, nil
		}
	}

	return Player{}, constants.ErrorPlayerNotFound
}

func (g *Game) IsOwner(id uuid.UUID) bool {
	return g.Owner == id
}

func CreateGame(owner Player) (Game, error) {
	owner.Name = strings.TrimSpace(owner.Name)
	if err := owner.Validate(); err != nil {
		return Game{}, err
	}
	owner.Status = constants.PlayerStatusAccepted

	return Game{
		ID:            uuid.New(),
		CurrentPlayer: uuid.Nil,
		Owner:         owner.ID,
		Players:       []Player{owner},
		Status:        constants.GameStatusOpen,
		CreatedAt:     time.Now().UTC(),
	}, nil
}

type Move struct {
	Player                   uuid.UUID
	PlayerCardPosition       int
	PlayerCardType           constants.CardType
	TargetPlayer             uuid.UUID
	TargetPlayerCardPosition int
	TargetPlayerCardType     constants.CardType
	Winner                   uuid.UUID
}

func DetermineMoveWinner(m Move) (uuid.UUID, error) {
	switch m.PlayerCardType {
	case constants.CardTypeDagger:
		switch m.TargetPlayerCardType {
		case constants.CardTypeDagger:
			return uuid.Nil, nil
		case constants.CardTypeShortSword:
			return m.TargetPlayer, nil
		case constants.CardTypeMace:
			return m.TargetPlayer, nil
		case constants.CardTypeBattleAxe:
			return m.TargetPlayer, nil
		case constants.CardTypeSpear:
			return m.TargetPlayer, nil
		case constants.CardTypeLongSword:
			return m.TargetPlayer, nil
		case constants.CardTypeArcher:
			return m.Player, nil
		case constants.CardTypeShield:
			return uuid.Nil, nil
		case constants.CardTypeCrown:
			return m.Player, nil
		}
	case constants.CardTypeShortSword:
		switch m.TargetPlayerCardType {
		case constants.CardTypeDagger:
			return m.Player, nil
		case constants.CardTypeShortSword:
			return uuid.Nil, nil
		case constants.CardTypeMace:
			return m.TargetPlayer, nil
		case constants.CardTypeBattleAxe:
			return m.TargetPlayer, nil
		case constants.CardTypeSpear:
			return m.TargetPlayer, nil
		case constants.CardTypeLongSword:
			return m.TargetPlayer, nil
		case constants.CardTypeArcher:
			return m.Player, nil
		case constants.CardTypeShield:
			return uuid.Nil, nil
		case constants.CardTypeCrown:
			return m.Player, nil
		}
	case constants.CardTypeMace:
		switch m.TargetPlayerCardType {
		case constants.CardTypeDagger:
			return m.Player, nil
		case constants.CardTypeShortSword:
			return m.Player, nil
		case constants.CardTypeMace:
			return uuid.Nil, nil
		case constants.CardTypeBattleAxe:
			return m.TargetPlayer, nil
		case constants.CardTypeSpear:
			return m.TargetPlayer, nil
		case constants.CardTypeLongSword:
			return m.TargetPlayer, nil
		case constants.CardTypeArcher:
			return m.Player, nil
		case constants.CardTypeShield:
			return uuid.Nil, nil
		case constants.CardTypeCrown:
			return m.Player, nil
		}
	case constants.CardTypeBattleAxe:
		switch m.TargetPlayerCardType {
		case constants.CardTypeDagger:
			return m.Player, nil
		case constants.CardTypeShortSword:
			return m.Player, nil
		case constants.CardTypeMace:
			return m.Player, nil
		case constants.CardTypeBattleAxe:
			return uuid.Nil, nil
		case constants.CardTypeSpear:
			return m.TargetPlayer, nil
		case constants.CardTypeLongSword:
			return m.TargetPlayer, nil
		case constants.CardTypeArcher:
			return m.Player, nil
		case constants.CardTypeShield:
			return uuid.Nil, nil
		case constants.CardTypeCrown:
			return m.Player, nil
		}
	case constants.CardTypeSpear:
		switch m.TargetPlayerCardType {
		case constants.CardTypeDagger:
			return m.Player, nil
		case constants.CardTypeShortSword:
			return m.Player, nil
		case constants.CardTypeMace:
			return m.Player, nil
		case constants.CardTypeBattleAxe:
			return m.Player, nil
		case constants.CardTypeSpear:
			return uuid.Nil, nil
		case constants.CardTypeLongSword:
			return m.TargetPlayer, nil
		case constants.CardTypeArcher:
			return m.Player, nil
		case constants.CardTypeShield:
			return uuid.Nil, nil
		case constants.CardTypeCrown:
			return m.Player, nil
		}
	case constants.CardTypeLongSword:
		switch m.TargetPlayerCardType {
		case constants.CardTypeDagger:
			return m.Player, nil
		case constants.CardTypeShortSword:
			return m.Player, nil
		case constants.CardTypeMace:
			return m.Player, nil
		case constants.CardTypeBattleAxe:
			return m.Player, nil
		case constants.CardTypeSpear:
			return m.Player, nil
		case constants.CardTypeLongSword:
			return uuid.Nil, nil
		case constants.CardTypeArcher:
			return m.Player, nil
		case constants.CardTypeShield:
			return uuid.Nil, nil
		case constants.CardTypeCrown:
			return m.Player, nil
		}
	case constants.CardTypeArcher:
		switch m.TargetPlayerCardType {
		case constants.CardTypeDagger:
			return m.Player, nil
		case constants.CardTypeShortSword:
			return m.Player, nil
		case constants.CardTypeMace:
			return m.Player, nil
		case constants.CardTypeBattleAxe:
			return m.Player, nil
		case constants.CardTypeSpear:
			return m.Player, nil
		case constants.CardTypeLongSword:
			return m.Player, nil
		case constants.CardTypeArcher:
			return m.Player, nil
		case constants.CardTypeShield:
			return m.Player, nil
		case constants.CardTypeCrown:
			return m.Player, nil
		}
	case constants.CardTypeShield:
		return m.TargetPlayer, nil
	case constants.CardTypeCrown:
		switch m.TargetPlayerCardType {
		case constants.CardTypeDagger:
			return m.TargetPlayer, nil
		case constants.CardTypeShortSword:
			return m.TargetPlayer, nil
		case constants.CardTypeMace:
			return m.TargetPlayer, nil
		case constants.CardTypeBattleAxe:
			return m.TargetPlayer, nil
		case constants.CardTypeSpear:
			return m.TargetPlayer, nil
		case constants.CardTypeLongSword:
			return m.TargetPlayer, nil
		case constants.CardTypeArcher:
			return m.Player, nil
		case constants.CardTypeShield:
			return uuid.Nil, nil
		case constants.CardTypeCrown:
			return m.Player, nil
		}
	}

	return uuid.Nil, constants.ErrorIllegalMove
}
