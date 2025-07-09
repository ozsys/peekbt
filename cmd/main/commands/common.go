package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

var defaultDev ble.Device

// InitDefaultAdapter は一度だけ linux.NewDevice() を呼び出し、以後は再利用します。
func InitDefaultAdapter() (ble.Device, error) {
	// すでに作成済みなら再利用
	if defaultDev != nil {
		return defaultDev, nil
	}

	// 初回のみアダプタを生成
	dev, err := linux.NewDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize BLE device: %w", err)
	}
	ble.SetDefaultDevice(dev)
	defaultDev = dev
	return dev, nil
}

// NewTimeoutCtx 秒指定で Context を生成（0 ならキャンセル型）
func NewTimeoutCtx(seconds int) (context.Context, context.CancelFunc) {
	if seconds == 0 {
		return context.WithCancel(context.Background())
	}
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}
