package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/spf13/cobra"
)

// deviceEntry holds the basic scan result
// entryDisplay holds display state for a device
// scanTime is the scan duration in seconds
var (
	scanTime int
)

// scanCommand represents the 'scan' CLI command
var scanCommand = &cobra.Command{
	Use:   "scan",
	Short: "Scan for nearby Bluetooth devices.",
	RunE:  runScanCommand,
}

func init() {
	scanCommand.Flags().IntP("time", "t", 0, "Scan time in seconds (0 = infinite)")
	scanCommand.Flags().Bool("rand", false, "Random address only")
	scanCommand.Flags().Bool("pub", false, "Public address only")
	rootCommand.AddCommand(scanCommand)
}

type deviceEntry struct {
	addr string
	name string
	rssi int
	seen time.Time
}

type entryDisplay struct {
	entry     deviceEntry
	colorTTL  time.Time // until when to display in green
	highlight string    // "all" for full line, "rssi" for RSSI only
}

func runScanCommand(cmd *cobra.Command, args []string) error {
	// フラグ取得
	scanTime, _ = cmd.Flags().GetInt("time")
	randOnly, _ := cmd.Flags().GetBool("rand")
	pubOnly, _ := cmd.Flags().GetBool("pub")
	if randOnly && pubOnly {
		return fmt.Errorf("flags --rand and --pub are mutually exclusive")
	}

	// BLE デバイス初期化
	dev, err := linux.NewDevice()
	if err != nil {
		return fmt.Errorf("failed to initialize BLE device: %v", err)
	}
	ble.SetDefaultDevice(dev)

	// 内部データ構造
	results := make(map[string]deviceEntry)
	displayed := make(map[string]entryDisplay)
	order := make([]string, 0, 16)
	var mu sync.Mutex

	// コンテキスト準備
	var ctx context.Context
	var cancel context.CancelFunc
	if scanTime == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(scanTime)*time.Second)
	}
	defer cancel()

	// 無限スキャン時は 'e'+Enter で終了
	if scanTime == 0 {
		go func() {
			fmt.Println("Scanning... (press 'e' + Enter to exit)")
			r := bufio.NewReader(cmd.InOrStdin())
			for {
				line, _ := r.ReadString('\n')
				if strings.TrimSpace(line) == "e" {
					cancel()
					return
				}
			}
		}()
	} else {
		fmt.Printf("Scanning for %d seconds...\n", scanTime)
	}

	// 初回描画
	fmt.Print("\033[2J\033[H")
	fmt.Println("ADDR                 RSSI   NAME")
	fmt.Println(strings.Repeat("-", 50))

	// 描画ループ
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				mu.Lock()
				for i, addr := range order {
					d := displayed[addr]
					colS, colE := "", ""
					if time.Now().Before(d.colorTTL) {
						colS, colE = "\033[32m", "\033[0m"
					}
					// カーソル移動＆上書き
					fmt.Printf("\033[%d;0H", i+4)
					if d.highlight == "all" {
						fmt.Printf("%s%-20s %-6d %-20s%s", colS, d.entry.addr, d.entry.rssi, d.entry.name, colE)
					} else {
						fmt.Printf("%-20s %-6d %-20s", d.entry.addr, d.entry.rssi, d.entry.name)
					}
					fmt.Print("\033[K")
				}
				mu.Unlock()
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	// スキャン実行
	err = ble.Scan(ctx, true, func(a ble.Advertisement) {
		addr := a.Addr().String()
		// フィルタリング
		oct := strings.Split(addr, ":")[0]
		first, _ := strconv.ParseUint(oct, 16, 8)
		isPub := (first & 0xC0) == 0x00
		if pubOnly && !isPub {
			return
		}
		if randOnly && isPub {
			return
		}
		r := a.RSSI()
		name := a.LocalName()
		if name == "" {
			name = "(no name)"
		}

		mu.Lock()
		if _, seen := results[addr]; !seen {
			order = append(order, addr)
		}
		results[addr] = deviceEntry{addr, name, r, time.Now()}
		hi := ""
		if old, seen := displayed[addr]; !seen {
			hi = "all"
		} else if old.entry.rssi != r {
			hi = "rssi"
		}
		displayed[addr] = entryDisplay{results[addr], time.Now().Add(1 * time.Second), hi}
		mu.Unlock()
	}, nil)

	// 正常終了判定
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		fmt.Println()
		return nil
	}
	return err
}
