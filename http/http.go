package http

import (
	"Miner/Enterprise"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	enterprise *enterprise.Enterprise
	ctx        context.Context
}

func NewHandler(ent *enterprise.Enterprise, ctx context.Context) *Handler {
	return &Handler{
		enterprise: ent,
		ctx:        ctx,
	}
}

func (h *Handler) HireMiner(w http.ResponseWriter, r *http.Request) {
	var req HireMinerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorHttp(w, err, http.StatusBadRequest)
		return
	}
	err := h.enterprise.HireMiner(h.ctx, req.Class)
	if err != nil {
		ErrorHttp(w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func ErrorHttp(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	newError := ErrorRequest{
		Message:   err.Error(),
		TimeError: time.Now(),
	}

	err = json.NewEncoder(w).Encode(newError)
	if err != nil {
		panic(err)
	}

}

func (h *Handler) GetMinerPayment(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	class := q.Get("class")
	switch class {
	case "small":
		newReq := PaymentClassRequest{
			Class: class,
			Pay:   5,
		}
		err := json.NewEncoder(w).Encode(newReq)
		if err != nil {
			panic(err)
		}
	case "normal":
		newReq := PaymentClassRequest{
			Class: class,
			Pay:   50,
		}
		err := json.NewEncoder(w).Encode(newReq)
		if err != nil {
			panic(err)
		}
	case "strong":
		newReq := PaymentClassRequest{
			Class: class,
			Pay:   450,
		}
		err := json.NewEncoder(w).Encode(newReq)
		if err != nil {
			panic(err)
		}
	default:
		ErrorHttp(w, errors.New("invalid class parameter"), http.StatusBadRequest)
	}
}

func (h *Handler) GetWorkers(w http.ResponseWriter, r *http.Request) {
	counts := h.enterprise.GetWorkersCount()
	q := r.URL.Query().Get("class")
	if q != "" {
		count, ok := counts[q]
		if !ok {
			ErrorHttp(w, errors.New(""), http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]int{q: count})
		return
	}
	req := WorkerClassRequest{
		Normal: counts["normal"],
		Small:  counts["small"],
		Strong: counts["strong"],
	}
	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) CostEquipment(w http.ResponseWriter, r *http.Request) {
	req := CostEquipment{
		Pickaxe:     3000,
		Ventilation: 15000,
		Trolley:     50000,
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) HTTPGetBalance(w http.ResponseWriter, r *http.Request) {
	balance := BalanceRequest{
		Balance: h.enterprise.Balance,
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(balance)
	if err != nil {
		panic(err)
	}
}
func (h *Handler) BuyEquipment(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("equipment")
	if err := h.enterprise.PurchaseTools(h.ctx, q); err != nil {
		ErrorHttp(w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PurchasedEquipment(w http.ResponseWriter, r *http.Request) {
	q, err := parseBoolPtr(r.URL.Query().Get("purchased"))
	if err != nil {
		ErrorHttp(w, err, http.StatusBadRequest)
		return
	}
	flagVal := false
	if q != nil {
		flagVal = *q
	}
	req := make([]EquipmentRequest, 0, 3)
	res := h.enterprise.GetEquipmentList(&flagVal)
	if len(res) == 0 {
		ErrorHttp(w, errors.New("no found"), http.StatusNotFound)
		return
	}
	for _, v := range res {
		req = append(req, EquipmentRequest{Equipment: v})
	}
	w.Header().Set("Content-Type", "application/json")
	err_ := json.NewEncoder(w).Encode(req)
	if err_ != nil {
		panic(err_)
	}

}
func parseBoolPtr(s string) (*bool, error) {
	if s == "" {
		return nil, nil // фильтр не задан
	}
	v, err := strconv.ParseBool(s) // принимает "1","0","t","f","true","false", ...
	if err != nil {
		return nil, err
	}
	return &v, nil
}
