package service

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Pawn struct {
	BelongTo int
	IsLocked bool
	XPOS     int
	YPOS     int
}
type GameQueueMessage struct {
	Mid   string
	Pid   int
	Value int
}

type GameManager struct {
	MOVE_ID        string
	MOVE_ACTION    string
	MOVE_PLAYER    int
	MOVE_LAST_ROLL int

	PLAYERS any
	PAWNS   [16]Pawn
	/*
	   Game_Players
	   Pawn_Position
	   Current_Player
	   Current_Action
	*/

	Queue chan GameQueueMessage
}

func NewGameManager() *GameManager {
	return &GameManager{
		MOVE_ID:        "",
		MOVE_ACTION:    "ROLL_DICE",
		MOVE_PLAYER:    0,
		MOVE_LAST_ROLL: 1,
		PAWNS:          initPawns(),
		Queue:          make(chan GameQueueMessage),
	}
}
func (s *GameManager) InitState() {
	// send Current Data
}
func initPawns() [16]Pawn {
	r := [16]Pawn{}
	x := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			r[x] = Pawn{BelongTo: i, IsLocked: true, XPOS: 0, YPOS: 0}
			x++
		}
	}
	return r
}
func (s *GameManager) GameLoop() {
	for {
		s.MOVE_ID = uuid.NewString()
		// SEND A GO CALL
		go s.GamePlaySync(s.MOVE_ID, s.MOVE_PLAYER)
		// NOTIFY
		s.GameNotifyPlayers(s.MOVE_ID, s.MOVE_ACTION, s.MOVE_PLAYER)

		// WAIT FOR REPLAY
		for {
			m := <-s.Queue

			if m.Mid != s.MOVE_ID {
				continue
			}

			if m.Pid != s.MOVE_PLAYER {
				continue
			}

			if s.MOVE_ACTION == "ROLL_DICE" {
				val := (rand.Int() % 6) + 1
				s.MOVE_LAST_ROLL = val
				if val != 6 && s.IsAllPawnLocked() {
					s.GameNotifyResult(s.MOVE_ACTION, val)
					s.MOVE_PLAYER = (s.MOVE_PLAYER + 1) % 4
					break
				}
				// if only one pawn out
				pawn_index := s.IsOnlyOneOut()
				if pawn_index != -1 {
					s.MovePawn(pawn_index)
					s.GameNotifyResult("PICK_PAWN", pawn_index)

					s.MOVE_ACTION = "ROLL_DICE"
					s.MOVE_PLAYER = (s.MOVE_PLAYER + 1) % 4
					s.GameNotifyResult(s.MOVE_ACTION, s.MOVE_LAST_ROLL)
					break
				}
				s.MOVE_ACTION = "PICK_PAWN"
				s.GameNotifyResult(s.MOVE_ACTION, val)

				break
			}
			if s.MOVE_ACTION == "PICK_PAWN" {
				pawn_index := m.Value
				if pawn_index > 15 {
					log.Println("INVALID PAWN INDEX")
					continue
				}
				pawn := s.PAWNS[pawn_index]
				if pawn.BelongTo != s.MOVE_PLAYER {
					log.Println("\nPAWN NOT BELONG TO PLAYER")
					continue
				}
				if pawn.IsLocked && s.MOVE_LAST_ROLL != 6 {
					log.Println("PAWN IS LOCKED")
					continue
				}
				// MOVE PAWN
				s.PAWNS[pawn_index].IsLocked = false
				s.MovePawn(pawn_index)
			}
			break
		}
	}
}

func (s *GameManager) IsAllPawnLocked() bool {
	for i := 0; i < 4; i++ {
		if !s.PAWNS[s.MOVE_PLAYER*4+i].IsLocked {
			return false
		}
	}
	return true
}

func (s *GameManager) IsOnlyOneOut() int {
	foundOne := -1
	for i := 0; i < 4; i++ {
		if !s.PAWNS[s.MOVE_PLAYER*4+i].IsLocked {
			if foundOne != -1 {
				return -1
			}
			foundOne = s.MOVE_PLAYER*4 + i
		}

	}
	return foundOne
}

func (s *GameManager) MovePawn(i int) {
	// moev pawn to one pos

	// CHECK IF PAWN GONNA KILL OTHER PAWNS

	s.MOVE_ACTION = "ROLL_DICE"
	s.MOVE_PLAYER = (s.MOVE_PLAYER + 1) % 4
	s.GameNotifyResult(s.MOVE_ACTION, i)
}

func (s *GameManager) GamePlaySync(mid string, pid int) {
	time.Sleep(time.Second * 2)
	val := (pid * 4) + (rand.Int() % 4)
	s.Queue <- GameQueueMessage{Mid: mid, Pid: pid, Value: val}
}

func (s *GameManager) GameNotifyPlayers(id, ac string, pl int) {
	fmt.Printf("[%s] WITH ID [_] FOR [%v] ", ac, pl)
}

func (s *GameManager) GameNotifyResult(ac string, val any) {
	fmt.Printf("RESULT WITH VALUE [%v]\n", val)
}
