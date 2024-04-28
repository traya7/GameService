package handler

import (
	"GameService/handler/gamesocket"
	"GameService/service"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

type Params struct {
	Psvc *service.PlayerService
	Gsvc *service.GameService
	Csvc *service.CallbackService
}

func NewRouter(p Params) *mux.Router {
	m := mux.NewRouter()
	m.Use(enableCORS)

	api := NewApi(p)
	m.HandleFunc("/gameservice/api/open", api.GetGame).Methods("GET")

	m.HandleFunc("/gameservice/api/load", api.Middleware(api.GetInfo)).Methods("POST")
	m.HandleFunc("/gameservice/api/buy", api.Middleware(api.BuyRound)).Methods("POST", "OPTIONS")
	m.HandleFunc("/gameservice/api/cashout", api.Middleware(api.Cashout)).Methods("POST")

	gs := gamesocket.New(p.Psvc)
	m.Handle("/gameservice/ws/lobby", websocket.Handler(gs.LobbyHandle))
	m.Handle("/gameservice/ws/game", websocket.Handler(gs.GameHandle))

	// serve game
	m.HandleFunc("/gameservice/ludo", f)
	m.PathPrefix("/gameservice/assets").Handler(LoadAssets())
	return m
}

func f(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("./games/ludo/index.html"))
	tmpl.Execute(w, nil)
}

func LoadAssets() http.Handler {
	return http.StripPrefix(
		"/gameservice/assets/",
		http.FileServer(http.Dir("./games/ludo/assets/")),
	)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow specified HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Allow specified headers
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		// Continue with the next handler
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}
		next.ServeHTTP(w, r)
	})
}
