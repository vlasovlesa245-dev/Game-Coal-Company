package company

import (
	"miner_game/miner"
	"sync"
	"sync/atomic"
	"time"
)

var mtx sync.RWMutex

type Info struct {
	Ballance int
	AllMiners map[int64]miner.Miner
	Counters Counters
	AtomicCounters map[string]*atomic.Int64
	StartTime time.Time
	GameTime string
}

type InfoDTO struct {
	Ballance int
	Tools StatisticsTools
	Counters Counters
	StartTime time.Time
	GameTime string
}

type Counters struct {
	AllMinersCounter int64
	MiniMinersCounter int64
	NormalMinersCounter int64
	StrongMinersCounter int64
}

func newAtomicCounters() map[string]*atomic.Int64 {
	AtomicCounters := make(map[string]*atomic.Int64)
	AtomicCounters["all"] = &atomic.Int64{}
	AtomicCounters["mini"] = &atomic.Int64{}
	AtomicCounters["normal"] = &atomic.Int64{}
	AtomicCounters["strong"] = &atomic.Int64{}
	return AtomicCounters
}

func (i *Info) count(miner miner.Miner) {
	mtx.Lock()
	defer mtx.Unlock()
	print("----------------------")

	switch miner.Info().MinerTier {
	case "Mini Miner", "mini":
		i.AtomicCounters["mini"].Add(1)
	case "Normal Miner", "normal":
		i.AtomicCounters["normal"].Add(1)
	case "Strong Miner", "strong":
		i.AtomicCounters["strong"].Add(1)
	}
	print("----------------------")
	i.AtomicCounters["all"].Add(1)
	i.AllMiners[i.AtomicCounters["all"].Load()] = miner
}

func (i *Info) getCounters(dto *InfoDTO) {
	dto.Counters.AllMinersCounter = i.AtomicCounters["all"].Load()
	dto.Counters.MiniMinersCounter = i.AtomicCounters["mini"].Load()
	dto.Counters.NormalMinersCounter = i.AtomicCounters["normal"].Load()
	dto.Counters.StrongMinersCounter = i.AtomicCounters["strong"].Load()
}

func (i *Info) GetActiveMiners() []byte {
	mtx.Lock()
	defer mtx.Unlock()
	activeMiners := []byte{}
	for _, m := range i.AllMiners {
		if m.Info().IsActive {
			activeMiners = append(activeMiners, m.InfoToByte()...)
		}
	}
	return activeMiners
}

func (i *Info) GetInactiveMiners() []byte {
	mtx.Lock()
	defer mtx.Unlock()
	inactiveMiners := []byte{}
	for _, m := range i.AllMiners {
		if !m.Info().IsActive {
			inactiveMiners = append(inactiveMiners, m.InfoToByte()...)
		}
	}
	return inactiveMiners
}

func (i *Info) GetMinersByTier(tier string) []byte {
	mtx.Lock()
	defer mtx.Unlock()
	miners := []byte{}
	for _, m := range i.AllMiners {
		if m.Info().IsActive {
			if m.Info().MinerTier == tier {
			miners = append(miners, m.InfoToByte()...)
		}
	}
	}
	return miners
}