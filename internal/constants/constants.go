package constants

type CardType int
type PlayerStatus int
type GameStatus int

const (
	CardTypeDagger     CardType = 1
	CardTypeShortSword CardType = 2
	CardTypeMace       CardType = 3
	CardTypeBattleAxe  CardType = 4
	CardTypeSpear      CardType = 5
	CardTypeLongSword  CardType = 6
	CardTypeArcher     CardType = 7
	CardTypeShield     CardType = 8
	CardTypeCrown      CardType = 9

	CardCountDagger     = 4
	CardCountShortSword = 4
	CardCountMace       = 4
	CardCountBattleAxe  = 4
	CardCountSpear      = 2
	CardCountLongSword  = 2
	CardCountArcher     = 2
	CardCountShield     = 2
	CardCountCrown      = 1

	DeckCount     = 5
	DeckCardCount = 5

	PlayerStatusUnaffiliated PlayerStatus = 1
	PlayerStatusRequested    PlayerStatus = 2
	PlayerStatusAccepted     PlayerStatus = 3
	PlayerStatusReady        PlayerStatus = 4
	PlayerStatusLost         PlayerStatus = 5
	PlayerStatusWon          PlayerStatus = 6

	GameStatusOpen    = 1
	GameStatusStarted = 2
	GameStatusDone    = 3
)

var ValidCards = []CardType{
	CardTypeDagger,
	CardTypeShortSword,
	CardTypeMace,
	CardTypeBattleAxe,
	CardTypeSpear,
	CardTypeLongSword,
	CardTypeArcher,
	CardTypeShield,
	CardTypeCrown,
}

func GetCardCount(cardType CardType) int {
	switch cardType {
	case CardTypeDagger:
		return CardCountDagger
	case CardTypeShortSword:
		return CardCountShortSword
	case CardTypeMace:
		return CardCountMace
	case CardTypeBattleAxe:
		return CardCountBattleAxe
	case CardTypeSpear:
		return CardCountSpear
	case CardTypeLongSword:
		return CardCountLongSword
	case CardTypeArcher:
		return CardCountArcher
	case CardTypeShield:
		return CardCountShield
	case CardTypeCrown:
		return CardCountCrown
	}

	return 0
}
