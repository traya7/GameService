package handler
//
// import (
// 	"GameService/domain/entities"
// 	"GameService/service"
// 	"errors"
// 	"io"
//
// 	"golang.org/x/net/websocket"
// )
//
// type GameSocket struct {
// 	psvc *service.PlayerService
// }
//
// func NewGameSocket(s *service.PlayerService) *GameSocket {
// 	return &GameSocket{
// 		psvc: s,
// 	}
// }
//
// func (h *GameSocket) auth(ws *websocket.Conn) (*entities.Player, error) {
// 	session_id := ws.Request().URL.Query().Get("id")
// 	if session_id == "" {
// 		return nil, errors.New("unauthorized access")
// 	}
//
// 	return h.psvc.GetUser(session_id)
// }
//
// func (h *GameSocket) Handle(ws *websocket.Conn) {
// 	user, err := h.auth(ws)
// 	if err != nil || user.ActiveRound == false {
// 		ws.Close()
// 		return
// 	}
// 	wsUser := WsUser{UserID: user.ID, Nickname: user.ID[:4], GameBet: user.ActiveBet, WsConn: ws}
// 	_ = wsUser
//
// 	buf := make([]byte, 1024)
// 	for {
// 		n, err := ws.Read(buf)
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			continue
// 		}
// 		_ = string(buf[:n])
// 	}
// }
