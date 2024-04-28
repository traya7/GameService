package repositories

import (
	"GameService/domain/entities"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("No found entity")

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}

func (r *PlayerRepository) GetPlayerBy(column, value string) (*entities.Player, error) {

	// Make the query to db
	query := fmt.Sprintf("SELECT * FROM players WHERE %s=?", column)
	row := r.db.QueryRow(query, value)

	var id string
	var sessionID string
	var activeRound bool
	var totalBet int
	var totalWin int
	var activeLevel int
	var activeBet int
	var cashout int

	// Scan the values from the selected row into variables
	err := row.Scan(&id, &sessionID, &totalBet, &totalWin, &activeRound, &activeLevel, &activeBet, &cashout)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	// Create a Player object with the retrieved values
	player := &entities.Player{
		ID:          id,
		SessionId:   sessionID,
		TotalBet:    totalBet,
		TotalWin:    totalWin,
		ActiveRound: activeRound,
		ActiveLevel: activeLevel,
		ActiveBet:   activeBet,
		Cashout:     cashout,
	}

	return player, nil
}

func (r *PlayerRepository) CreateNewPlayer(user_id string) error {
	insertQuery := "INSERT INTO players (id) VALUES (?)"
	_, err := r.db.Exec(insertQuery, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PlayerRepository) SetPlayerSession(user_id, session string) error {
	updateQuery := "UPDATE players SET session_id=? WHERE id=?"
	// Execute the update query
	result, err := r.db.Exec(updateQuery, session, user_id)
	if err != nil {
		return err
	}
	// Check the number of rows affected by the update operation
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (r *PlayerRepository) SetPlayerBet(user_id string, totalBet, betValue int) error {
	// Make and Execute the update query
	updateQuery := "UPDATE players SET total_bet=?, active_bet=?, active_round=1, active_level=0, cashout=0 WHERE id=?"
	result, err := r.db.Exec(updateQuery, totalBet, betValue, user_id)
	if err != nil {
		return err
	}
	// Check the number of rows affected by the update operation
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (r *PlayerRepository) SetPlayerWin(user_id string, totalWin int) error {
	// Make and Execute the update query
	updateQuery := "UPDATE players SET total_win=?, active_bet=0, active_round=0, active_level=0, cashout=0 WHERE id=?"
	result, err := r.db.Exec(updateQuery, totalWin, user_id)
	if err != nil {
		return err
	}
	// Check the number of rows affected by the update operation
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
