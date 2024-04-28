package gamesocket

import (
	"encoding/json"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type LobbyUser struct {
	UserID   string
	Nickname string
	GameBet  int
	IsBot    bool
	WsConn   *websocket.Conn
}
type LobbyQueueMessage struct {
	Action string
	User   LobbyUser
}
type LobbyRoom struct {
	RoomBet int
	Names   []string
	Users   []LobbyUser
}
type LobbyManager struct {
	Queue   chan LobbyQueueMessage
	Rooms   map[int]LobbyRoom
	NewGame func(string, []LobbyUser)
}

func (m *LobbyManager) QueueHandler() {
	for {
		msg := <-m.Queue
		if msg.Action == "JOIN" {
			room, ok := m.Rooms[msg.User.GameBet]
			if !ok {
				room = LobbyRoom{RoomBet: msg.User.GameBet, Users: []LobbyUser{}}
			}
			room.Users = append(room.Users, msg.User)
			room.Names = append(room.Names, msg.User.Nickname)
			m.Rooms[msg.User.GameBet] = room
			m.UpdateNotify(room)

			if len(room.Users) == 4 {
				// START THE GAME

				gid := uuid.NewString()
				m.NewGame(gid, room.Users)
				delete(m.Rooms, room.RoomBet)
				m.StartNotify(room, gid)
			}
		}
		if msg.Action == "DROP" {
			room, ok := m.Rooms[msg.User.GameBet]
			if !ok {
				room = LobbyRoom{Users: []LobbyUser{}}
			}
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
			m.Rooms[msg.User.GameBet] = room
			m.UpdateNotify(room)
		}
	}
}

func (m *LobbyManager) JoinRoom(u LobbyUser) {
	m.Queue <- LobbyQueueMessage{Action: "JOIN", User: u}
}

func (m *LobbyManager) DropRoom(u LobbyUser) {
	m.Queue <- LobbyQueueMessage{Action: "DROP", User: u}
}

func (m *LobbyManager) UpdateNotify(r LobbyRoom) {
	d := map[string]any{
		"a": "u",
		"c": r.Names,
	}

	for _, v := range r.Users {
		d["u"] = v.Nickname
		data, _ := json.Marshal(d)
		v.WsConn.Write(data)
	}
}

func (m *LobbyManager) StartNotify(r LobbyRoom, gid string) {
	data, _ := json.Marshal(map[string]any{
		"a": "s",
		"c": r.Names,
		"i": gid,
	})

	for _, v := range r.Users {
		v.WsConn.Write(data)
	}
}
