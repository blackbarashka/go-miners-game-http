package miner

import (
	"context"
	"time"
)

type SmallMiner struct {
	Energy    int
	Income    Coal
	payment   Coal
	breaktime time.Duration
}

type NormalMiner struct {
	energy    int
	income    Coal
	payment   Coal
	breaktime time.Duration
}

type StrongMiner struct {
	energy    int
	income    Coal
	payment   Coal
	breaktime time.Duration
}

func NewSmallMiner() *SmallMiner {
	return &SmallMiner{
		Energy:    30,
		Income:    Coal{Amount: 1},
		payment:   Coal{Amount: 5},
		breaktime: 3 * time.Second,
	}
}

func NewNormalMiner() *NormalMiner {
	return &NormalMiner{
		energy:    45,
		income:    Coal{Amount: 3},
		payment:   Coal{Amount: 50},
		breaktime: 2 * time.Second,
	}
}

func NewStrongMiner() *StrongMiner {
	return &StrongMiner{
		energy:    60,
		income:    Coal{Amount: 10},
		payment:   Coal{Amount: 450},
		breaktime: 1 * time.Second,
	}
}

func (m *SmallMiner) Info() MinerInfo {
	return MinerInfo{
		Class:     "small",
		Energy:    m.Energy,
		Income:    m.Income,
		Payment:   m.payment,
		BreakTime: m.breaktime,
	}
}
func (m *NormalMiner) Info() MinerInfo {
	return MinerInfo{
		Class:     "normal",
		Energy:    m.energy,
		Income:    m.income,
		Payment:   m.payment,
		BreakTime: m.breaktime,
	}
}
func (m *StrongMiner) Info() MinerInfo {
	return MinerInfo{
		Class:     "strong",
		Energy:    m.energy,
		Income:    m.income,
		Payment:   m.payment,
		BreakTime: m.breaktime,
	}
}

func (m *NormalMiner) Run(ctx context.Context) <-chan Coal {
	ch := make(chan Coal)
	go func() {
		defer close(ch)
		for m.energy > 0 {

			select {
			case ch <- m.income:
				m.energy--
			case <-ctx.Done():
				return
			}
			if m.energy == 0 {
				return
			}
			select {
			case <-time.After(m.breaktime):
			case <-ctx.Done():
				return
			}
		}

	}()
	return ch
}
func (m *SmallMiner) Run(ctx context.Context) <-chan Coal {
	ch := make(chan Coal)
	go func() {
		defer close(ch)
		for m.Energy > 0 {
			select {
			case ch <- m.Income:
				m.Energy--
			case <-ctx.Done():
				return
			}
			if m.Energy == 0 {
				return
			}
			select {
			case <-time.After(m.breaktime):
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}

func (m *StrongMiner) Run(ctx context.Context) <-chan Coal {
	ch := make(chan Coal)
	go func() {
		defer close(ch)
		for m.energy > 0 {
			select {
			case ch <- m.income:
				m.energy--
				m.income.Amount += 3
			case <-ctx.Done():
				return
			}
			if m.energy == 0 {
				return
			}
			select {
			case <-time.After(m.breaktime):
			case <-ctx.Done():
				return

			}
		}
	}()
	return ch
}
