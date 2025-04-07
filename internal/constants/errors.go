package constants

import "errors"

var (
	ErrorEmptyStack          = errors.New("empty stack")
	ErrorGameNotFound        = errors.New("game not found")
	ErrorGameNotStarted      = errors.New("game not started")
	ErrorIllegalMove         = errors.New("illegal move")
	ErrorInvalidDeckCounts   = errors.New("invalid deck counts, must be five groups of five cards")
	ErrorInvalidCardCount    = errors.New("invalid card count")
	ErrorInvalidStack        = errors.New("invalid stack")
	ErrorNotFound            = errors.New("not found")
	ErrorPlayerInvalidID     = errors.New("invalid player id")
	ErrorPlayerInvalidSecret = errors.New("invalid player secret")
	ErrorPlayerInvalidName   = errors.New("player invalid name")
	ErrorPlayerInvalid       = errors.New("invalid user")
	ErrorPlayerNotFound      = errors.New("player not found")
	ErrorPlayerNotOwner      = errors.New("player not owner")
	ErrorPlayerWrongTurn     = errors.New("wrong player turn")
)
