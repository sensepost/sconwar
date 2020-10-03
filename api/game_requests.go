package api

import (
	"errors"

	"github.com/sensepost/sconwar/game"
)

type GetGameDetailRequest struct {
	GameID string `uri:"game_id" binding:"required,uuid" example:"1df69d53-3468-43df-a43b-a9c674240cab"`
}

func (r *GetGameDetailRequest) Validation() error {

	if game.Games[r.GameID] == nil {
		return errors.New("unable to find a game by that id")
	}

	return nil
}

type JoinPlayerRequest struct {
	GameID   string `json:"game_id" binding:"required,uuid" example:"1df69d53-3468-43df-a43b-a9c674240cab"`
	PlayerID string `json:"player_id" binding:"required,uuid" example:"6d950e36-b82b-4253-93d7-faa63d3a0e63"`
}

func (r *JoinPlayerRequest) Validation() error {

	if game.Games[r.GameID] == nil {
		return errors.New("unable to find a game by that id")
	}

	// todo: validate player_id

	return nil
}
