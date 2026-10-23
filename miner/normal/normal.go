package normal

import (
	"context"
	"encoding/json"
	"miner_game/miner"
	"sync"
	"sync/atomic"
	"time"
)

type NormalMiner struct {
	MinerTier string
	Wage int
	Energy int
	Income miner.Coal
	Break int
	EnergyLeft atomic.Int64
	IsActive bool
}

func NewMiner() *NormalMiner {
	miner := &NormalMiner{
		MinerTier: "Normal Miner",
		Wage:   50,
		Income: miner.Coal{Ballance: 3},
		Break:  2,
		IsActive: true,
	}
	miner.EnergyLeft.Store(45)
	return miner
}

func (m *NormalMiner) Info() miner.MinerInfo {
	info := miner.MinerInfo{
		MinerTier: "Normal Miner",
		Wage:   50,
		Energy: 45,
		Income: 3,
		Break:  2,
		EnergyLeft: int(m.EnergyLeft.Load()),
		IsActive: m.IsActive,
	}
	return info
}

func (m *NormalMiner) InfoToByte() []byte {
	b, err := json.MarshalIndent(m.Info(), "", "    ")
	if err != nil {
		panic(err)
	}
	return b
}

func Start(miner *NormalMiner, ctx context.Context, coalTransfer chan miner.Coal, wg *sync.WaitGroup) {
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
    	}
    }
}
func (m *NormalMiner) Run(ctx context.Context) <-chan miner.Coal {
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