package gamesocket

import (
	"GameService/domain/entities"
	"errors"

	"golang.org/x/net/websocket"
)

func (h *GameSocket) UserAuth(ws *websocket.Conn) (*entities.Player, error) {
	session_id := ws.Request().URL.Query().Get("id")
	if session_id == "" {
		return nil, errors.New("unauthorized access")
	}
	return h.psvc.GetUser(session_id)
}

func (h *GameSocket) GameAuth(ws *websocket.Conn) (*GameManager, error) {
	game_id := ws.Request().URL.Query().Get("gid")
	if game_id == "" {
		return nil, errors.New("unauthorized access")
	}
	game, ok := h.Games[game_id]
	if !ok {
		return nil, errors.New("game not found")
	}
	return game, nil
}
