package service

import (
	"GameService/domain/entities"
	"GameService/domain/repositories"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type PlayerService struct {
	repo *repositories.PlayerRepository
}

func NewPlayerService(r *repositories.PlayerRepository) *PlayerService {
	return &PlayerService{
		repo: r,
	}
}

func (s *PlayerService) NewSessionId(user_id string) (string, error) {
	_, err := s.repo.GetPlayerBy("id", user_id)
	if err != nil {
		if err == repositories.ErrNotFound {
			if err := s.repo.CreateNewPlayer(user_id); err != nil {
				return "", errors.New("internal error 1")
			}
		} else {
			return "", errors.New("internal error")
		}
	}
	buf := []byte(user_id + fmt.Sprint(time.Now().Unix()))
	h := md5.New()
	h.Write(buf)
	session := hex.EncodeToString(h.Sum(nil))
	if s.repo.SetPlayerSession(user_id, session) != nil {
		return "", errors.New("internal error")
	}
	return session, nil
}

func (s *PlayerService) GetUser(session_id string) (*entities.Player, error) {
	user, err := s.repo.GetPlayerBy("session_id", session_id)
	if err != nil {
		return nil, errors.New("internal error")
	}
	return user, nil
}

func (s *PlayerService) BuyRound(u *entities.Player, betValue int) error {
	totalBet := u.TotalBet + betValue
	if err := s.repo.SetPlayerBet(u.ID, totalBet, betValue); err != nil {
		return errors.New("internal error")
	}
  // SAVE Activites

	return nil
}

func (s *PlayerService) Cashout(u *entities.Player) error {
	totalWin := u.TotalWin + u.Cashout
	if err := s.repo.SetPlayerWin(u.ID, totalWin); err != nil {
		return errors.New("internal error")
	}
	return nil
}
