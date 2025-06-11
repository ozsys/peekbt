package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/spf13/cobra"
)

type deviceEntry struct {
	addr string
	name string
	rssi int
	seen time.Time
}

var scanTime int

var scanCommand = &cobra.Command{
	Use:   "scan",
	Short: "Scan for nearby Bluetooth devices.",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanTime, _ = cmd.Flags().GetInt("time")

		dev, err := linux.NewDevice()
		if err != nil {
			return fmt.Errorf("failed to initialize BLE device: %v", err)
		}
		ble.SetDefaultDevice(dev)

		results := make(map[string]deviceEntry)
		deviceLines := make(map[string]int)              // addr -> line number
		recentlyChanged := make(map[string]time.Time)    // addr -> time of last change

		var ctx context.Context
		var cancel context.CancelFunc
		if scanTime == 0 {
			ctx, cancel = context.WithCancel(context.Background())
		} else {
			ctx, cancel = context.WithTimeout(context.Background(), time.Duration(scanTime)*time.Second)
		}
		defer cancel()

		if scanTime == 0 {
			go func() {
				fmt.Println("Scanning... (press 'e' + Enter to exit)")
				reader := bufio.NewReader(os.Stdin)
				for {
					input, _ := reader.ReadString('\n')
					if strings.TrimSpace(input) == "e" {
						cancel()
						return
					}
				}
			}()
		} else {
			fmt.Printf("Scanning for %d seconds...\n", scanTime)
		}

		// 初期ヘッダー
		fmt.Print("\033[2J") // clear screen
		fmt.Print("\033[H")  // move cursor to top-left
		fmt.Println("ADDR                 RSSI   NAME")
		fmt.Println(strings.Repeat("-", 50))

		lineCounter := 4 // データの描画開始行番号

		// 描画更新用ループ
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					now := time.Now()
					for addr, entry := range results {
						line, ok := deviceLines[addr]
						if !ok {
							deviceLines[addr] = lineCounter
							line = lineCounter
							lineCounter++
						}

						colorStart := ""
						colorEnd := ""
						if t, ok := recentlyChanged[addr]; ok && now.Sub(t) < time.Second {
							colorStart = "\033[32m"
							colorEnd = "\033[0m"
						}

						// move cursor to the device's line
						fmt.Printf("\033[%d;0H%s%-20s %-6d %-20s%s\n",
							line, colorStart, entry.addr, entry.rssi, entry.name, colorEnd)
					}
					time.Sleep(200 * time.Millisecond)
				}
			}
		}()

		// BLEスキャン処理
		return ble.Scan(ctx, false, func(a ble.Advertisement) {
			addr := a.Addr().String()
			rssi := a.RSSI()
			name := a.LocalName()
			if name == "" {
				name = "(no name)"
			}

			prev, seen := results[addr]
			results[addr] = deviceEntry{
				addr: addr,
				name: name,
				rssi: rssi,
				seen: time.Now(),
			}

			if !seen || prev.rssi != rssi {
				recentlyChanged[addr] = time.Now()
			}
		}, nil)
	},
}

func init() {
	scanCommand.Flags().IntP("time", "t", 0, "Scan time in seconds (0 = infinite)")
	rootCommand.AddCommand(scanCommand)
}