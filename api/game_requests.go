package api

import (
	"errors"

	"github.com/sensepost/sconwar/game"
	"github.com/sensepost/sconwar/storage"
	"gorm.io/gorm"
)

// GetGameDetailRequest is a request with a game uuid
type GetGameDetailRequest struct {
	GameID string `uri:"game_id" binding:"required,uuid" example:"1df69d53-3468-43df-a43b-a9c674240cab"`
}

// Validation validates request values
func (r *GetGameDetailRequest) Validation() error {

	if game.Games[r.GameID] == nil {
		return errors.New("invalid game uuid")
	}

	return nil
}

// JoinPlayerRequest is a request to join a player to a game
type JoinPlayerRequest struct {
	GameID   string `json:"game_id" binding:"required,uuid" example:"1df69d53-3468-43df-a43b-a9c674240cab"`
	PlayerID string `json:"player_id" binding:"required,uuid" example:"6d950e36-b82b-4253-93d7-faa63d3a0e63"`
}

// Validation validates request values
func (r *JoinPlayerRequest) Validation() error {

	if game.Games[r.GameID] == nil {
		return errors.New("invalid game uuid")
	}

	var player storage.Player
	res := storage.Storage.Get().Where("UUID = ?", r.PlayerID).First(&player)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return errors.New("invalid player uuid")
	}

	return nil
}
