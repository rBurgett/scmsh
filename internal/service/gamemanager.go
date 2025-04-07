package service

import "github.com/rBurgett/scmsh/internal/storage"

type GameManager struct {
	storageClient storage.Client[Game]
}

func NewGameManager(client storage.Client[Game]) *GameManager {
	return &GameManager{
		storageClient: client,
	}
}
