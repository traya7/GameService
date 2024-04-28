package service

type GameService struct{}

func NewGameService() *GameService {
	return &GameService{}
}

func (s *GameService) GetGameInfo(game_id string) (map[string]any, error) {

	// TODO
	r := map[string]any{
		"bets":   []int{10, 50, 100, 500},
		"levels": 4,
	}
	return r, nil
}
