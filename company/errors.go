package company

import "errors"

var ErrNotEnoughCoal error = errors.New("Not enough coal on ballance for purchase")
var ErrMinerNotFound error = errors.New("Miner with such title not found")
var ErrToolNotFound error = errors.New("Tool with such title not found")
var ErrToolAlreadyBought error = errors.New("This tool is already bought!")