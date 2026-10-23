package company

import (
	"context"
	"errors"
	"miner_game/miner"
	"miner_game/miner/mini"
	"miner_game/miner/normal"
	"miner_game/miner/strong"
	"sync"
	"time"
)

type Company struct {
	Coal miner.Coal
	Info Info
	Tools map[string]Tool
	miners []string

	Ctx context.Context
	Cancel context.CancelFunc
	mtx sync.RWMutex
}

func NewCompany() *Company {
	Ctx, CtxCancel := context.WithCancel(context.Background())
	now := time.Now()
	since := time.Since(now)
	miners := []string{"Mini Miner", "Normal Miner", "Strong Miner", "mini", "normal", "strong"}
	company := Company{
		Coal: miner.Coal{
			Ballance: 0,
		},
		Ctx: Ctx,
		Cancel: CtxCancel,
		Info: Info{
			AllMiners: make(map[int64]miner.Miner),
			StartTime: now,
			GameTime: since.String(),
			AtomicCounters: newAtomicCounters(),
		},
		miners: miners,
		Tools: GetTools(),
	}
	company.Info.Ballance = company.GetBallance()

	company.poolOneCoal()
	return &company
}

func (c *Company) ValidateForMiner(title *string) error {
	if *title == "" {
		return errors.New("title of miner is empty")
	}
	for _, m := range c.miners {
		if *title != m {
			continue
		}
		switch *title {
		case "mini" :
			*title = "Mini Miner"
		case "normal" :
			*title = "Normal Miner"
		case "strong" :
			*title = "Strong Miner"
		}
		return nil
	}
	return ErrMinerNotFound
}

func (c *Company) ValidateForTools(title string) error {
	if title == "" {
		return errors.New("title of tool is empty")
	}
	if _, ok := c.Tools[title]; ok {
		if c.Tools[title].IsBought {
			return ErrToolAlreadyBought
		}
		return nil
	}
	return ErrToolNotFound
}

func (c *Company) AddMiner(miner miner.Miner) {
	c.Info.count(miner) 
	c.Coal.Ballance -= miner.Info().Wage
	ch := miner.Run(c.Ctx)

	go func() {
		for coal := range ch {
		c.mtx.Lock()
		c.Coal.Ballance += coal.Ballance
		c.mtx.Unlock()
	}
	}()
}

func (c *Company) ChoosingMiner(tier string) (miner.Miner) {
	switch tier {
	case "Mini Miner", "mini":
		m := mini.NewMiner()
		return m
	case "Normal Miner", "normal":
		m := normal.NewMiner()
		return m
	case "Strong Miner", "strong":
		m := strong.NewMiner()
		return m
	}
	return nil
}

func (c *Company) poolOneCoal() {
	ch := make(chan miner.Coal)
	wg := &sync.WaitGroup{}
	
	wg.Add(1)
	go c.getOneCoalPerSecond(ch, wg)

	go func() {
		wg.Wait()
		close(ch)
	}()

	go func() {
		for coal := range ch {
		mtx.Lock()
		c.Coal.Ballance += coal.Ballance
		mtx.Unlock()
	}
	}()
}

func (c *Company) getOneCoalPerSecond(ch chan miner.Coal, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <- c.Ctx.Done():
			return
		case <- time.After(1 * time.Second):
			//
		}
		select {
		case <- c.Ctx.Done():
			return
		case ch <- miner.Coal{Ballance: 1}:
			print(c.Coal.Ballance)
			print("\n")
		}
	}
}

func (c *Company) GetInfo() *InfoDTO {
	info := &InfoDTO{}
	info.Ballance = c.GetBallance()
	info.Tools = c.Get_Statistics_Of_Tools()
	info.GameTime = time.Since(c.Info.StartTime).String()
	info.StartTime = c.Info.StartTime
	c.Info.getCounters(info)
	return info
}

func (c *Company) GetBallance() int {
	ballance := c.Coal.Ballance
	return ballance
}