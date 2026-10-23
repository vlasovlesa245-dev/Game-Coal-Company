package server

import (
	"context"
	"encoding/json"
	"miner_game/company"
	"net/http"
	"github.com/gorilla/mux"
)

type HTTP_Handlers struct {
	Company *company.Company
	Ctx context.Context
	Cancel context.CancelFunc
}

func NewHTTPHandlers(c *company.Company) *HTTP_Handlers {
	return &HTTP_Handlers{
		Company: c,
		Ctx: c.Ctx,
		Cancel: c.Cancel,
	}
}

/*
pattern: /miner
method: POST
info:	miner tier title in json

succeed: 
status code: 201 Created
request body: json with miner

failed:
status code: 400, 409, 500 ...
request body: JSON with error + time
*/

func (h *HTTP_Handlers) HandleBuy(w http.ResponseWriter, r *http.Request) {
	var buyDTO buyDTO
	if err := json.NewDecoder(r.Body).Decode(&buyDTO); err != nil {
		ErrorDTO := NewErrorDTO(err.Error(), company.HelpAboutGame())

		http.Error(w, ErrorDTO.ToString(), http.StatusBadRequest)
		return
	}
	minerTier := buyDTO.Miner
	tool := buyDTO.Tool
	if minerTier != nil && tool != nil {
		h.BuyMinerAndTool(w, *minerTier, *tool)
		return
	}
	if minerTier != nil {
		h.BuyMinerOnly(w, *minerTier)
		return
	}
	if tool != nil {
		h.BuyToolOnly(w, *tool)
		return
	}
	b, err := json.MarshalIndent("Вы не ввели ни название шахтёра, ни оборудования:\n", "", "")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusBadRequest)
	if _, err := w.Write(b); err != nil {
		print("failed to write http response:", err)
		return
	}
}
/*
pattern: /miner
method: GET
info:	---

succeed: 
status code: 200 OK
request body: json with  list of miners

failed:
status code: 400, 409, 500 ...
request body: JSON with error + time
*/

func (h *HTTP_Handlers) HandleGetActiveMiners(w http.ResponseWriter, r *http.Request) {
	list := h.Company.Info.GetActiveMiners()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(list); err != nil {
		print("failed to write http response:", err)
		return
	}
}

/*
pattern: /miner
method: GET
info:	query(inactive)

succeed: 
status code: 200 OK
request body: json with  list of miners

failed:
status code: 400, 409, 500 ...
request body: JSON with error + time
*/

func (h *HTTP_Handlers) HandleGetInactiveMiners(w http.ResponseWriter, r *http.Request) {
	list := h.Company.Info.GetInactiveMiners()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(list); err != nil {
		print("failed to write http response:", err)
		return
	}
}

/*
pattern: /miner
method: GET
info:	tier title in path

succeed: 
status code: 200 OK
request body: json with  list of miners

failed:
status code: 400, 409, 500 ...
request body: JSON with error + time
*/

func (h *HTTP_Handlers) HandleGetMinersByTier(w http.ResponseWriter, r *http.Request) {
	tier := mux.Vars(r)["tier"]
	if err := h.Company.ValidateForMiner(&tier); err != nil {
		ErrorDTO := NewErrorDTO(err.Error(), company.HelpAboutGame())
		http.Error(w, ErrorDTO.ToString(), http.StatusBadRequest)
		return
	}

	list := h.Company.Info.GetMinersByTier(tier)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(list); err != nil {
		print("failed to write http response:", err)
		return
	}
}

/*
pattern: /miner/help_miners
method: GET
info:	---

succeed: 
status code: 200 OK
request body: json with  list of miners

failed:
status code: 400, 409, 500 ...
request body: JSON with error + time
*/

func (h *HTTP_Handlers) HandleGetInfoMiners(w http.ResponseWriter, r *http.Request) {
	msg := company.HelpMiners()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		panic(err)
	}
	
	if _, err := w.Write(b); err != nil {
		print("failed to write http response:", err)
		return
	}
}

/*
pattern: /miner/help
method: GET
info:	---

succeed: 
status code: 200 OK
request body: json with  list of miners

failed:
status code: 400, 409, 500 ...
request body: JSON with error + time
*/

func (h *HTTP_Handlers) HandleGetInfoAboutGame(w http.ResponseWriter, r *http.Request) {
	msg := company.HelpAboutGame()

	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(msg, "", "")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		print("failed to write http response:", err)
		return
	}
}

/*
pattern: /miner/help_tools
method: GET
info:	---

succeed: 
status code: 200 OK
request body: json with  list of miners

failed:
status code: 400, 409, 500 ...
request body: JSON with error + time
*/

func (h *HTTP_Handlers) HandleGetInfoTools(w http.ResponseWriter, r *http.Request) {
	msg := company.HelpTools()
	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		print("failed to write http response:", err)
		return
	}
}

/*
pattern: /miner/info
method: GET
info:	---

succeed: 
status code: 200 OK
request body: json with  list of miners

failed:
status code: 400, 409, 500 ...
request body: JSON with error + time
*/

func (h *HTTP_Handlers) HandleGetStatistics(w http.ResponseWriter, r *http.Request) {
	info := h.Company.GetInfo()
	b, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(b); err != nil {
		print("failed to write http response:", err)
	}
}

/*
pattern: /miner/exit
method: ---
info:	---

succeed: 
status code: 200 OK
request body: json with  list of miners

failed:
status code: 400, 409, 500 ...
request body: JSON with error + time
*/

func (h *HTTP_Handlers) HandleExit(w http.ResponseWriter, r *http.Request) {
	b, err := json.MarshalIndent(h.Company.GetInfo(), "", "    ")
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(b); err != nil {
		print("failed to write http response:", err)
	}
	h.Cancel()
}