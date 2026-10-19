package miner

import "time"

type Coal struct {
	Amount int
}

type MinerInfo struct {
	// Класс шахтера.
	Class string
	// Сколько платят.
	Payment Coal
	// Энергия.
	Energy int
	// Добыча.
	Income Coal
	// Перерыв.
	BreakTime time.Duration
}
