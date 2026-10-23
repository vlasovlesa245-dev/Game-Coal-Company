package server

import (
	"encoding/json"
	"errors"
	"miner_game/company"
	"net/http"
	"time"
)

type buyDTO struct {
	Miner *string
	Tool *string
}

type ErrorDTO struct {
	Description string
	Time        time.Time
	Help []string
}

func NewErrorDTO(desc string, help []string) ErrorDTO {
	return ErrorDTO{
		Description: desc,
		Time: time.Now(),
		Help: help,
	}
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func WriteError (w http.ResponseWriter, err error, status int) {
	ErrorDTO := NewErrorDTO(err.Error(), company.HelpAboutGame())
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	b, err := json.MarshalIndent(ErrorDTO, "", "    ")
	w.Write(b)
	// err = json.NewEncoder(w).Encode(ErrorDTO.ToString())
	// if err != nil {
	// 	print("Ошибка кодирования JSON:", err)
	// }
}

func (h *HTTP_Handlers) BuyMinerAndTool(w http.ResponseWriter, minerTier string, tool string) {
	if err := h.Company.ValidateForMiner(&minerTier); err != nil {
			WriteError(w, err, http.StatusBadRequest)
			return
		}
	if err := h.Company.ValidateForTools(tool); err != nil {
		if errors.Is(err, company.ErrToolAlreadyBought) {
			WriteError(w, err, http.StatusConflict)
			return
		} else {
			WriteError(w, err, http.StatusBadRequest)
			return
		}
	}
	miner := h.Company.ChoosingMiner(minerTier)
	if h.Company.Coal.Ballance < miner.Info().Wage + h.Company.Tools[tool].Price {
		WriteError(w, company.ErrNotEnoughCoal, http.StatusConflict)
		return
	}

	h.Company.AddMiner(miner)
	h.Company.BuyTool(tool)

	b, err := json.MarshalIndent("Куплен новый шахтёр:", "", "	   ")
	if err != nil {
		panic(err)
	}
	b2, err := json.MarshalIndent(miner.Info(), "", "	")
	if err != nil {
		panic(err)
	}
	b3, err := json.MarshalIndent("Куплено новое оборудование:" + tool, "", "	   ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	b = append(b, b2...)
	b = append(b, b3...)
	if _, err := w.Write(b); err != nil {
		print("failed to write http response:", err)
		return
}
}

func (h *HTTP_Handlers) BuyMinerOnly(w http.ResponseWriter, minerTier string) {
	if err := h.Company.ValidateForMiner(&minerTier); err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}
	miner := h.Company.ChoosingMiner(minerTier)
	if h.Company.Coal.Ballance < miner.Info().Wage {
		WriteError(w, company.ErrNotEnoughCoal, http.StatusConflict)
		return
	}
	h.Company.AddMiner(miner)
	b, err := json.MarshalIndent("Куплен новый шахтёр:", "", "	   ")
	if err != nil {
		panic(err)
	}
	b2, err := json.MarshalIndent(miner.Info(), "", "	")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	b = append(b, b2...)
	if _, err := w.Write(b); err != nil {
		print("failed to write http response:", err)
		return
	}
}

func (h *HTTP_Handlers) BuyToolOnly(w http.ResponseWriter, tool string) {
	if err := h.Company.ValidateForTools(tool); err != nil {
		if errors.Is(err, company.ErrToolAlreadyBought) {
			WriteError(w, err, http.StatusConflict)
		} else {
			WriteError(w, err, http.StatusBadRequest)
		}
		return
	}
	
	if h.Company.Coal.Ballance < h.Company.Tools[tool].Price {
		WriteError(w, company.ErrNotEnoughCoal, http.StatusConflict)
		return
	}
	h.Company.BuyTool(tool)
	b, err := json.MarshalIndent("Куплено новое оборудование:" + tool + "\nEnter", "", "	   ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write(b); err != nil {
		print("failed to write http response:", err)
		return
	}
}