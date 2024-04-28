package gamesocket

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type LobbyUser struct {
	UserID    string
	Nickname  string
	GameBet   int
	IsBot     bool
	BotRoomID string
	WsConn    *websocket.Conn
}

type LobbyQueueMessage struct {
	Action string
	User   LobbyUser
}

type LobbyRoom struct {
	RoomID  string
	RoomBet int
	Names   []string
	Users   []LobbyUser
}
type LobbyManager struct {
	Queue chan LobbyQueueMessage
	Rooms map[string]LobbyRoom

	NewGame func(string, []LobbyUser)
}

func (m *LobbyManager) QueueHandler() {
	for {
		msg := <-m.Queue
		if msg.Action == "JOIN" {
			room, err := m.GetRoom(msg.User.GameBet)
			if err != nil {
				log.Println("NEW ROOM")
				room = m.CreateRoom(msg.User.GameBet)
				go m.AddBot(room.RoomID)
			}
			room.Users = append(room.Users, msg.User)
			room.Names = append(room.Names, msg.User.Nickname)
			m.Rooms[room.RoomID] = room
			m.UpdateNotify(room)
		}

		if msg.Action == "BOT" {
			if room, ok := m.Rooms[msg.User.BotRoomID]; ok {
				room.Users = append(room.Users, msg.User)
				room.Names = append(room.Names, msg.User.Nickname)
				m.Rooms[room.RoomID] = room
				m.UpdateNotify(room)
				go m.AddBot(room.RoomID)
			}
		}

		if msg.Action == "DROP" {
			if room, err := m.GetRoom(msg.User.GameBet); err == nil {
				for i, wu := range room.Users {
					if wu.UserID == msg.User.UserID {
						room.Users = append(room.Users[:i], room.Users[i+1:]...)
						room.Names = []string{}
						for _, wu := range room.Users {
							room.Names = append(room.Names, wu.Nickname)
						}
						break
					}
				}
				m.Rooms[room.RoomID] = room
				m.UpdateNotify(room)
			}
		}
	}
}

// QUEUE HELPERS
func (m *LobbyManager) JoinRoom(u LobbyUser) {
	m.Queue <- LobbyQueueMessage{Action: "JOIN", User: u}
}

func (m *LobbyManager) DropRoom(u LobbyUser) {
	m.Queue <- LobbyQueueMessage{Action: "DROP", User: u}
}

func (m *LobbyManager) AddBot(rid string) {
	time.Sleep(time.Second * 3)
	u := LobbyUser{
		UserID:    "00000000000",
		Nickname:  "BOT",
		IsBot:     true,
		BotRoomID: rid,

		GameBet: 0,
		WsConn:  nil,
	}
	m.Queue <- LobbyQueueMessage{Action: "BOT", User: u}
}

func (m *LobbyManager) UpdateNotify(r LobbyRoom) {
	d := map[string]any{}
	if len(r.Users) == 4 {
		d = map[string]any{
			"a": "s",
			"c": r.Names,
			"i": r.RoomID,
		}
		m.StartRoom(r)
	} else {
		d = map[string]any{
			"a": "u",
			"c": r.Names,
		}
	}

	for _, v := range r.Users {
		if v.IsBot {
			continue
		}
		d["u"] = v.Nickname
		data, _ := json.Marshal(d)
		v.WsConn.Write(data)
	}
}

// MANAGER FUNCTIONS
func (m *LobbyManager) StartRoom(room LobbyRoom) {
	botCount := 0
	for _, v := range room.Users {
		if v.IsBot {
			botCount++
		}
	}
	if botCount != 4 {
		m.NewGame(room.RoomID, room.Users)
	}
	delete(m.Rooms, room.RoomID)
	log.Println("ROOM DELETED")
}

func (m *LobbyManager) GetRoom(bet int) (LobbyRoom, error) {
	for _, v := range m.Rooms {
		if v.RoomBet == bet {
			return v, nil
		}
	}
	return LobbyRoom{}, errors.New("NO_ROOM")
}

func (m *LobbyManager) CreateRoom(bet int) LobbyRoom {
	return LobbyRoom{
		RoomID:  uuid.NewString(),
		RoomBet: bet,
		Names:   []string{},
		Users:   []LobbyUser{},
	}
}
