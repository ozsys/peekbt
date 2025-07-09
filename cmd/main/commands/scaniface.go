package commands

import (
	"context"

	"github.com/go-ble/ble"
)

type Scanner interface {
	// 第 4 引数は ble.AdvFilter 固定
	Scan(ctx context.Context, allowDup bool,
		h ble.AdvHandler, f ble.AdvFilter) error
}

type bleScanner struct{}

func (bleScanner) Scan(ctx context.Context, allowDup bool,
	h ble.AdvHandler, f ble.AdvFilter) error {
	return ble.Scan(ctx, allowDup, h, f)
}

var DefaultScanner Scanner = bleScanner{}
