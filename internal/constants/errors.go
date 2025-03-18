package constants

import "errors"

var (
	ErrorEmptyStack          = errors.New("empty stack")
	ErrorGameNotStarted      = errors.New("game not started")
	ErrorIllegalMove         = errors.New("illegal move")
	ErrorInvalidDeckCounts   = errors.New("invalid deck counts, must be five groups of five cards")
	ErrorInvalidCardCount    = errors.New("invalid card count")
	ErrorInvalidStack        = errors.New("invalid stack")
	ErrorPlayerInvalidID     = errors.New("invalid user id")
	ErrorPlayerInvalidSecret = errors.New("invalid user secret")
	ErrorPlayerInvalidName   = errors.New("user invalid name")
	ErrorPlayerInvalid       = errors.New("invalid user")
	ErrorPlayerNotFound      = errors.New("user not found")
	ErrorPlayerWrongTurn     = errors.New("wrong user turn")
)
