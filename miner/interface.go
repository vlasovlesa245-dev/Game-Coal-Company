package miner

import (
	"context"
)

type Coal struct {
	Ballance int
}

type MinerInfo struct {
	MinerTier string
	Wage int
	Energy int
	Income int
	Break int
	EnergyLeft int
	Progression int
	IsActive bool
}

type Miner interface {
	Run(ctx context.Context) <- chan Coal
	Info() MinerInfo
	InfoToByte() []byte
}