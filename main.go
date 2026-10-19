package main

import (
	"Miner/Enterprise"
	myhttp "Miner/http"
	"log"
	"net/http"
)

func main() {
	ent, ctx := enterprise.NewEnterprise()
	go ent.Start(ctx)
	handler := myhttp.NewHandler(ent, ctx)
	mux := http.NewServeMux()
	mux.HandleFunc("/miners", handler.HireMiner)             // POST: нанять шахтера
	mux.HandleFunc("/workers", handler.GetWorkers)           // GET: посмотреть рабочих
	mux.HandleFunc("/payment", handler.GetMinerPayment)      // GET: узнать зарплату/стоимость найма
	mux.HandleFunc("/equipment/cost", handler.CostEquipment) // GET: цены оборудования
	mux.HandleFunc("/equipment/buy", handler.BuyEquipment)   // POST: купить оборудование
	mux.HandleFunc("/equipment", handler.PurchasedEquipment) // GET: список оборудования
	mux.HandleFunc("/balance", handler.HTTPGetBalance)

	// 5. Запускаем HTTP-сервер на порту 8080
	log.Println("Server is running on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
