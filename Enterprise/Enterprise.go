package enterprise

import (
	"Miner/miner"
	"context"
	"errors"
	"sync"
	"time"
	// Не забудь импортировать свой пакет miner!
	// "твой_проект/miner"
)

type Miner interface {
	Run(ctx context.Context) <-chan miner.Coal
	Info() miner.MinerInfo // Добавили префикс miner., так как структура лежит в том пакете
}

type Enterprise struct {
	Balance   int             // Текущий баланс угля (проще использовать int для математики)
	Miners    []Miner         // Список всех работавших шахтеров
	Workers   map[string]int  // Статистика нанятых шахтеров по классам
	Equipment map[string]bool // Купленное оборудование (Кирки, Вентиляция, Вагонетки)

	startTime  time.Time          // Время старта игры (для финальной статистики)
	CancelFunc context.CancelFunc // Рубильник для остановки всех горутин-шахтеров (пункт 9)
	mtx        sync.Mutex         // Защита от гонки данных
}

// Конструктор теперь возвращает не только Предприятие, но и глобальный контекст игры
func NewEnterprise() (*Enterprise, context.Context) {
	// Создаем глобальный контекст игры с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	ent := &Enterprise{
		Balance: 0,
		Miners:  make([]Miner, 0),
		Workers: map[string]int{"small": 0, "normal": 0, "strong": 0},
		Equipment: map[string]bool{
			"pickaxes":    false,
			"ventilation": false,
			"carts":       false,
		},
		startTime:  time.Now(), // Засекаем время старта
		CancelFunc: cancel,     // Сохраняем "кнопку увольнения" внутри структуры
	}

	return ent, ctx
}

func (ent *Enterprise) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ent.mtx.Lock()
			ent.Balance += 1
			ent.mtx.Unlock()
		}
	}
}

func (ent *Enterprise) HireMiner(ctx context.Context, minerclass string) error {
	ent.mtx.Lock()
	defer ent.mtx.Unlock()
	switch minerclass {
	case "small":
		newMiner := miner.NewSmallMiner()
		if ent.Balance-newMiner.Info().Payment.Amount < 0 {
			return errors.New("balance is negative")
		}
		ent.Balance -= newMiner.Info().Payment.Amount
		ent.Workers["small"]++
		ent.Miners = append(ent.Miners, newMiner)
		ch := newMiner.Run(ctx)
		go func() {
			for coal := range ch {
				ent.mtx.Lock()
				ent.Balance += coal.Amount
				ent.mtx.Unlock()
			}
		}()

	case "normal":
		newMiner := miner.NewNormalMiner()
		if ent.Balance-newMiner.Info().Payment.Amount < 0 {
			return errors.New("balance is negative")
		}
		ent.Balance -= newMiner.Info().Payment.Amount
		ent.Workers["normal"]++
		ent.Miners = append(ent.Miners, newMiner)
		ch := newMiner.Run(ctx)
		go func() {
			for coal := range ch {
				ent.mtx.Lock()
				ent.Balance += coal.Amount
				ent.mtx.Unlock()
			}
		}()

	case "strong":
		newMiner := miner.NewStrongMiner()
		if ent.Balance-newMiner.Info().Payment.Amount < 0 {
			return errors.New("balance is negative")
		}
		ent.Balance -= newMiner.Info().Payment.Amount
		ent.Workers["strong"]++
		ent.Miners = append(ent.Miners, newMiner)
		ch := newMiner.Run(ctx)
		go func() {
			for coal := range ch {
				ent.mtx.Lock()
				ent.Balance += coal.Amount
				ent.mtx.Unlock()
			}
		}()
	default:
		return errors.New("invalid miner class")
	}
	return nil
}

func (ent *Enterprise) PurchaseTools(ctx context.Context, ToolsName string) error {
	ent.mtx.Lock()
	defer ent.mtx.Unlock()
	if purchased, exists := ent.Equipment[ToolsName]; exists && purchased {
		return errors.New("equipment already purchased")
	}
	switch ToolsName {
	case "pickaxes":
		if ent.Balance-3000 < 0 {
			return errors.New("balance is negative")
		}
		ent.Balance -= 3000
		ent.Equipment["pickaxes"] = true
	case "ventilation":
		if ent.Balance-15000 < 0 {
			return errors.New("balance is negative")
		}
		ent.Balance -= 15000
		ent.Equipment["ventilation"] = true
	case "carts":
		if ent.Balance-50000 < 0 {
			return errors.New("balance is negative")
		}
		ent.Balance -= 50000
		ent.Equipment["carts"] = true
	default:
		return errors.New("invalid tool name")
	}
	return nil

}

func (e *Enterprise) GetWorkersCount() map[string]int {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	count := make(map[string]int, len(e.Workers))
	for k, v := range e.Workers {
		count[k] = v
	}
	return count
}

func (ent *Enterprise) GetEquipmentList(flag *bool) []string {
	ent.mtx.Lock()
	ent.mtx.Unlock()
	res := make([]string, 0, len(ent.Equipment))
	for k, v := range ent.Equipment {
		if v == *flag {
			res = append(res, k)
		}
	}
	return res
}
