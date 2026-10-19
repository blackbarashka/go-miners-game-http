package http

import "time"

type HireMinerRequest struct {
	Class string `json:"class"`
}
type BalanceRequest struct {
	Balance int `json:"balance"`
}
type ErrorRequest struct {
	Message   string    `json:"message"`
	TimeError time.Time `json:"time_error"`
}

type PaymentClassRequest struct {
	Class string `json:"class"`
	Pay   int    `json:"pay"`
}

type WorkerClassRequest struct {
	Normal int `json:"normal"`
	Small  int `json:"small"`
	Strong int `json:"strong"`
}

type CostEquipment struct {
	Pickaxe     int `json:"pickaxe"`
	Ventilation int `json:"ventilation"`
	Trolley     int `json:"trolley"`
}

type EquipmentRequest struct {
	Equipment string `json:"equipment"`
}
