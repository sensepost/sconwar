package api

import (
	"errors"

	"github.com/sensepost/sconwar/game"
	"github.com/sensepost/sconwar/storage"
	"gorm.io/gorm"
)

// ActionGamePlayerRequest is a request with the game and player is
type ActionGamePlayerRequest struct {
	GameID   string `json:"game_id" binding:"required,uuid" example:"1df69d53-3468-43df-a43b-a9c674240cab"`
	PlayerID string `json:"player_id" binding:"required,uuid" example:"6d950e36-b82b-4253-93d7-faa63d3a0e63"`
}

// Validation validates request values
func (r *ActionGamePlayerRequest) Validation() error {

	if game.Games[r.GameID] == nil {
		return errors.New("invalid game uuid")
	}

	var player storage.Player
	res := storage.Storage.Get().Where("UUID = ?", r.PlayerID).First(&player)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return errors.New("invalid player uuid")
	}

	// todo: validate that the game has started

	return nil
}

// ActionMoveRequest is a request to move a player
type ActionMoveRequest struct {
	GamePlayerIDs ActionGamePlayerRequest `json:"game_player_id" binding:"required"`
	X             int                     `json:"x" binding:"required" example:"10"`
	Y             int                     `json:"y" binding:"required" example:"9"`
}

// Validation validates request values
func (r *ActionMoveRequest) Validation() error {

	if err := r.GamePlayerIDs.Validation(); err != nil {
		return err
	}

	// todo: validate distance

	return nil
}

// ActionAttackRequest is a request to move a player
type ActionAttackRequest struct {
	GamePlayerIDs ActionGamePlayerRequest `json:"game_player_id" binding:"required"`
	X             int                     `json:"x" binding:"required" example:"10"`
	Y             int                     `json:"y" binding:"required" example:"9"`
}

// Validation validates request values
func (r *ActionAttackRequest) Validation() error {

	if err := r.GamePlayerIDs.Validation(); err != nil {
		return err
	}

	// todo: validate distance

	return nil
}
