package strong

import (
	"context"
	"encoding/json"
	"miner_game/miner"
	"sync"
	"sync/atomic"
	"time"
)

type StrongMiner struct {
	MinerTier string
	Wage int
	Energy int
	Income miner.Coal
	Break int
	EnergyLeft atomic.Int64
	IsActive bool
	Progression int
}

func NewMiner() *StrongMiner {
	miner := &StrongMiner{
		MinerTier: "Strong Miner",
		Wage:   450,
		Income: miner.Coal{Ballance: 10},
		Break:  1,
		IsActive: true,
		Progression: 3,
	}
	miner.EnergyLeft.Store(60)
	return miner
}

func (m *StrongMiner) Info() miner.MinerInfo {
	info := miner.MinerInfo{
		MinerTier: "Strong Miner",
		Wage:   450,
		Energy: 60,
		Income: 10,
		Break:  1,
		EnergyLeft: int(m.EnergyLeft.Load()),
		IsActive: m.IsActive,
		Progression: 3,
	}
	return info
}

func (m *StrongMiner) InfoToByte() []byte {
	b, err := json.MarshalIndent(m.Info(), "", "    ")
	if err != nil {
		panic(err)
	}
	return b
}

func Start(miner *StrongMiner, ctx context.Context, coalTransfer chan miner.Coal, wg *sync.WaitGroup) {
	defer func() {
		print("end")
		miner.IsActive = false
		wg.Done()
	}()
	for i := 1; i <= miner.Info().Energy; i++ {
	select {
    	case <- ctx.Done():
    		return
		case <- time.After(time.Duration(miner.Break) * time.Second):
			//
    	}
    select {
    	case <- ctx.Done():
			return
    	case coalTransfer <- miner.Income:
			miner.EnergyLeft.Add(-1)
			miner.Income.Ballance += 3
    	}
    }
}

func (m *StrongMiner) Run(ctx context.Context) <-chan miner.Coal {
	coalTransfer := make(chan miner.Coal)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go Start(m, ctx, coalTransfer, wg)
	go func() {
		wg.Wait()
		close(coalTransfer)
	}()
	return coalTransfer
}