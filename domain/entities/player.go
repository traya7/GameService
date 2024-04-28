package entities

type Player struct {
	ID          string `json:"-"`
	SessionId   string `json:"session"`
	ActiveRound bool
	TotalBet    int
	TotalWin    int
	ActiveLevel int
	ActiveBet   int
	Cashout     int
}
