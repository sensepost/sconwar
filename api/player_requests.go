package api

import (
	"errors"

	"github.com/sensepost/sconwar/game"
	"github.com/sensepost/sconwar/storage"
	"gorm.io/gorm"
)

// PlayerRequest is a request for a player
type PlayerRequest struct {
	PlayerID string `json:"player_id" binding:"required,uuid" example:"6d950e36-b82b-4253-93d7-faa63d3a0e63"`
}

// Validation validates request values
func (r *PlayerRequest) Validation() error {

	var player storage.Player
	res := storage.Storage.Get().Where("UUID = ?", r.PlayerID).First(&player)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return errors.New("invalid player uuid")
	}

	return nil
}

// RegisterPlayerRequest is a request to register a player
type RegisterPlayerRequest struct {
	Name string `json:"name" binding:"required" example:"my name"`
}

// Validation validates request values
func (r *RegisterPlayerRequest) Validation() error {

	var player storage.Player
	res := storage.Storage.Get().Where("name = ?", r.Name).First(&player)
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return errors.New("player name already registered")
	}

	return nil
}

// PlayerGameRequest is a request that needs a player and game id
type PlayerGameRequest struct {
	GameID   string `json:"game_id" binding:"required,uuid" example:"1df69d53-3468-43df-a43b-a9c674240cab"`
	PlayerID string `json:"player_id" binding:"required,uuid" example:"6d950e36-b82b-4253-93d7-faa63d3a0e63"`
}

// Validation validates request values
func (r *PlayerGameRequest) Validation() error {

	var player storage.Player
	res := storage.Storage.Get().Where("UUID = ?", r.PlayerID).First(&player)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return errors.New("invalid player uuid")
	}

	if game.Games[r.GameID] == nil {
		return errors.New("invalid game uuid")
	}

	// todo: maybe validate that this player is actually in this game

	return nil
}
