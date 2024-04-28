package handler

import (
	"GameService/domain/entities"
	svc "GameService/service"
	"context"
	"encoding/json"
	"net/http"
)

type ApiHandler struct {
	psvc *svc.PlayerService
	gsvc *svc.GameService
	csvc *svc.CallbackService
}

func NewApi(p Params) *ApiHandler {
	return &ApiHandler{
		psvc: p.Psvc,
		gsvc: p.Gsvc,
		csvc: p.Csvc,
	}
}

// Helpers
func ErrorResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(400)
	w.Write([]byte(message))
}

func JsonResponse(w http.ResponseWriter, data any) {
	json.NewEncoder(w).Encode(data)
}

func (h *ApiHandler) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session_id := r.URL.Query().Get("id")
		if session_id == "" {
			ErrorResponse(w, "invalid or expired session")
			return
		}
		u, err := h.psvc.GetUser(session_id)
		if err != nil {
			ErrorResponse(w, "invalid or expired session")
			return
		}
		ctx := context.WithValue(r.Context(), "user", u)
		next(w, r.WithContext(ctx))
	}
}

func (h *ApiHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("id")
	if user_id == "" {
		ErrorResponse(w, "invalid user id")
		return
	}
	url := "http://localhost:8001/gameservice/game/ludo?id="
	seesion, err := h.psvc.NewSessionId(user_id)
	if err != nil {
		ErrorResponse(w, err.Error())
		return
	}
	JsonResponse(w, map[string]any{"url": url + seesion})
}

func (h *ApiHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*entities.Player)
	game, err := h.gsvc.GetGameInfo("LCX01")
	if err != nil {
		ErrorResponse(w, "Cannot open game")
		return
	}
	response := map[string]any{
		"active_round": user.ActiveRound,
		"active_level": user.ActiveLevel,
		"active_bet":   user.ActiveBet,
		"cashout":      user.Cashout,

		"bets":   game["bets"],
		"levels": game["levels"],
	}
	json.NewEncoder(w).Encode(response)
}

func (h *ApiHandler) BuyRound(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*entities.Player)
	var payload struct {
		Bet int `json:"bet"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		ErrorResponse(w, "cannot handle request")
		return
	}
	err := h.csvc.RequestTakeMoney(user.ID, payload.Bet)
	if err != nil {
		//ErrorResponse(w, err.Error())
		//return
	}
	// SAVE PLAYER BET
	if h.psvc.BuyRound(user, payload.Bet) != nil {
		// TODO: SEND ROLLBACK
		_ = true
		// ELSE RETURN ERROR
		ErrorResponse(w, "Cannot make payment")
		return
	}
	// TODO: h.psvc.SetUserBet()
	w.Write([]byte("ok"))
	return

}
func (h *ApiHandler) Cashout(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*entities.Player)
	if h.psvc.Cashout(user) != nil {
		ErrorResponse(w, "Error in cashout")
		return
	}
	if user.Cashout > 0 {
		err := h.csvc.RequestGiveMoney(user.ID, user.Cashout)
		if err != nil {
			ErrorResponse(w, err.Error())
			return
		}
	}
}

func (h *ApiHandler) FindGame(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Bet int `json:"bet"`
	}
	if json.NewDecoder(r.Body).Decode(&payload) != nil {
		ErrorResponse(w, "cannot handle request")
		return
	}
	user := r.Context().Value("user").(entities.Player)
	_ = user

	// CALL CALLBACK FOR TAKE MONEY
	err := h.csvc.RequestTakeMoney(user.ID, payload.Bet)
	if err != nil {
		ErrorResponse(w, err.Error())
		return
	}

}
