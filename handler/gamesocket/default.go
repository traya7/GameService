package gamesocket

import (
	"GameService/service"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

type GameSocket struct {
	psvc         *service.PlayerService
	LobbyManager *LobbyManager
	Games        map[string]*GameManager
}

func New(ps *service.PlayerService) *GameSocket {
	gs := &GameSocket{
		psvc: ps,
		LobbyManager: &LobbyManager{
			Queue: make(chan LobbyQueueMessage),
			Rooms: map[string]LobbyRoom{},
		},
		Games: map[string]*GameManager{},
	}
	gs.LobbyManager.NewGame = gs.AddNewGame
	go gs.LobbyManager.QueueHandler()
	return gs
}

func (h *GameSocket) LobbyHandle(ws *websocket.Conn) {
	user, err := h.UserAuth(ws)
	if err != nil || user.ActiveRound == false {
		ws.Close()
		return
	}

	wsUser := LobbyUser{UserID: user.ID, Nickname: user.ID[:3] + "**", GameBet: user.ActiveBet, WsConn: ws}
	h.LobbyManager.JoinRoom(wsUser)

	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				h.LobbyManager.DropRoom(wsUser)
				break
			}
			continue
		}
		_ = string(buf[:n])
	}
}

func (h *GameSocket) GameHandle(ws *websocket.Conn) {
	user, err := h.UserAuth(ws)
	if err != nil || user.ActiveRound == false {
		ws.Close()
		return
	}

	game, err := h.GameAuth(ws)
	if err != nil {
		ws.Close()
		return
	}

	if game.PlayerJoin(user, ws) != nil {
		ws.Close()
		return
	}

	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				// h.LobbyManager.DropRoom(wsUser)
				break
			}
			continue
		}
    game.Queue <- GameQueueMessage{UserID: user.ID, Msg: string(buf[:n])}
	}
}

func (h *GameSocket) AddNewGame(game_id string, users []LobbyUser) {
	h.Games[game_id] = NewGameManager(users)
	log.Println("NEW GAME ADDED")
}
