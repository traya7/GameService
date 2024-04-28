package gamesocket

import (
	"GameService/domain/entities"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"time"

	"golang.org/x/net/websocket"
)

type Pawn struct {
	ID    int     `json:"id"`
	Color int     `json:"c"`
	Xpos  float64 `json:"xp"`
	Ypos  float64 `json:"yp"`
	End   bool    `json:"-"`
}

type GamePlayer struct {
	UserID   string
	Nikename string
	Color    int
	IsBot    bool
	WsConn   *websocket.Conn
}

type GameQueueMessage struct {
	UserID string
	Msg    string
}

type GameManager struct {
	Players       map[int]*GamePlayer
	Pawns         [16]Pawn
	MoveID        string
	CurrentMove   int
	CurrentAction string

	Queue chan GameQueueMessage
}

func NewGameManager(us []LobbyUser) *GameManager {
	plys := map[int]*GamePlayer{}
	plys[0] = playerFromLobbyUser(us[0], 0)
	plys[1] = playerFromLobbyUser(us[1], 1)
	plys[2] = playerFromLobbyUser(us[2], 2)
	plys[3] = playerFromLobbyUser(us[3], 3)
	g := &GameManager{
		Players:     plys,
		Pawns:       InitPawns(),
		CurrentMove: 0,
		Queue:       make(chan GameQueueMessage),
	}
	go g.GameLoop()

	return g
}

func (s *GameManager) GameLoop() {
	for {
		msg := <-s.Queue
		if s.Players[s.CurrentMove].UserID != msg.UserID {
			continue
		}
		s.Notfiy("r", true)
		time.Sleep(time.Second * 2)
		val := (rand.Int() % 6) + 1
		s.Notfiy("d", val)

		time.Sleep(time.Second * 2)
		s.CurrentMove = (s.CurrentMove + 1) % 4
		s.Notfiy("d", val)
		// IF BOT SYSTEM PICK AFTER SECONDS

	}
}

func (s *GameManager) PlayerJoin(user *entities.Player, ws *websocket.Conn) error {
	p, err := s.findPlayer(user.ID)
	if err != nil {
		return errors.New("Cannot join the room")
	}
	p.WsConn = ws
	s.sendInitData(p)
	return nil
}

func (s *GameManager) findPlayer(id string) (*GamePlayer, error) {
	for _, v := range s.Players {
		if v.UserID == id {
			return v, nil
		}
	}

	return nil, errors.New("not found player")
}

func playerFromLobbyUser(u LobbyUser, i int) *GamePlayer {
	return &GamePlayer{
		UserID:   u.UserID,
		Nikename: u.Nickname,
		Color:    (i),
		IsBot:    u.IsBot,
	}
}

func (s *GameManager) Notfiy(t string, v any) {
	data, _ := json.Marshal(map[string]any{
		"a":  t,
		"v":  v,
		"cm": s.CurrentMove,
	})

	for _, v := range s.Players {
		if !v.IsBot {
			v.WsConn.Write(data)
		}
	}
}

func (s *GameManager) sendInitData(p *GamePlayer) {
	data, _ := json.Marshal(map[string]any{
		"a": "i",
		"v": map[string]any{
			"b": s.Players[0].Nikename,
			"r": s.Players[1].Nikename,
			"g": s.Players[2].Nikename,
			"y": s.Players[3].Nikename,
			"u": p.Nikename,
		},
		"p":  s.Pawns,
		"cm": s.GetColor(s.CurrentMove),
	})
	p.WsConn.Write(data)
	log.Println("Message Init Sent")
}

func (s *GameManager) GetColor(i int) string {
	if i == 0 {
		return "blue"
	}
	if i == 1 {
		return "red"
	}
	if i == 2 {
		return "green"
	}
	return "yellow"
}

func InitPawns() [16]Pawn {
	return [16]Pawn{
		{ID: 0, Color: 0, Xpos: 1.5, Ypos: 1.5, End: false},
		{ID: 1, Color: 0, Xpos: 1.5, Ypos: 3.5, End: false},
		{ID: 2, Color: 0, Xpos: 3.5, Ypos: 1.5, End: false},
		{ID: 3, Color: 0, Xpos: 3.5, Ypos: 3.5, End: false},
		{ID: 0, Color: 1, Xpos: 10.5, Ypos: 1.5, End: false},
		{ID: 1, Color: 1, Xpos: 10.5, Ypos: 3.5, End: false},
		{ID: 2, Color: 1, Xpos: 12.5, Ypos: 1.5, End: false},
		{ID: 3, Color: 1, Xpos: 12.5, Ypos: 3.5, End: false},
		{ID: 0, Color: 2, Xpos: 10.5, Ypos: 10.5, End: false},
		{ID: 1, Color: 2, Xpos: 10.5, Ypos: 12.5, End: false},
		{ID: 2, Color: 2, Xpos: 12.5, Ypos: 10.5, End: false},
		{ID: 3, Color: 2, Xpos: 12.5, Ypos: 12.5, End: false},
		{ID: 0, Color: 3, Xpos: 1.5, Ypos: 10.5, End: false},
		{ID: 1, Color: 3, Xpos: 1.5, Ypos: 12.5, End: false},
		{ID: 2, Color: 3, Xpos: 3.5, Ypos: 10.5, End: false},
		{ID: 3, Color: 3, Xpos: 3.5, Ypos: 12.5, End: false},
	}
}
