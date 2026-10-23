package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	handlers *HTTP_Handlers
	server http.Server
	Ctx      context.Context
	Cancel context.CancelFunc
}

func NewServer(h *HTTP_Handlers) *Server {
	return &Server{
		handlers: h,
		server: http.Server{
			Addr: ":4000",
		},
		Ctx: h.Ctx,
		Cancel: h.Cancel,
	}
}

func (s *Server) StartServer() error {
	router := mux.NewRouter()
	s.server.Handler = router

	//exit
	router.Path("/miner/exit").HandlerFunc(s.handlers.HandleExit)
	
	//buy miner and tool
	router.Path("/miner").Methods("POST").HandlerFunc(s.handlers.HandleBuy)

	//get active/inactive miners
	router.Path("/miner").Methods("GET").Queries("status", "active").HandlerFunc(s.handlers.HandleGetActiveMiners)
	router.Path("/miner").Methods("GET").Queries("status", "inactive").HandlerFunc(s.handlers.HandleGetInactiveMiners)
	router.Path("/miner").Methods("GET").HandlerFunc(s.handlers.HandleGetActiveMiners)

	//get info about tools/game/miners
	router.Path("/miner/help").Methods("GET").HandlerFunc(s.handlers.HandleGetInfoAboutGame)
	router.Path("/miner/help_tools").Methods("GET").HandlerFunc(s.handlers.HandleGetInfoTools)
	router.Path("/miner/help_miners").Methods("GET").HandlerFunc(s.handlers.HandleGetInfoMiners)
	
	//get statistics
	router.Path("/miner/info").Methods("GET").HandlerFunc(s.handlers.HandleGetStatistics)

	//PANIC
	//get miners by tier
	router.Path("/miner/{tier}").Methods("GET").HandlerFunc(s.handlers.HandleGetMinersByTier)

	var err error
	go func() {
		for {
			if err != nil {
				return
			}
			select {
			case <- s.Ctx.Done() :
				s.server.Close()
				err = http.ErrServerClosed
				return
			default:
				continue
			}
		}
}()

	print("start!")
	
	if err = s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			print("Сервер остановлен игроком")
		} else {
			return err
		}
	}
	return nil
}