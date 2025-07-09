package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
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

type entryDisplay struct {
	entry     deviceEntry
	colorTTL  time.Time // until when to display in green
	highlight string    // "all" for new device
}

var (
	scanTime int
	randOnly bool
	pubOnly  bool
)

var scanCommand = &cobra.Command{
	Use:   "scan",
	Short: "Scan for nearby Bluetooth devices.",
	RunE:  runScanCommand,
}

func init() {
	scanCommand.Flags().IntVarP(&scanTime, "time", "t", 0, "Scan time in seconds (0 = infinite)")
	scanCommand.Flags().BoolVar(&randOnly, "rand", false, "Random address only.")
	scanCommand.Flags().BoolVar(&pubOnly, "pub", false, "Public address only.")
	rootCommand.AddCommand(scanCommand)
}

func runScanCommand(cmd *cobra.Command, args []string) error {
	// 排他フラグチェック
	if randOnly && pubOnly {
		return fmt.Errorf("flags --rand and --pub are mutually exclusive")
	}

	// BLEデバイス初期化
	if err := initBLEDevice(); err != nil {
		return err
	}

	results := make(map[string]deviceEntry)
	displayed := make(map[string]entryDisplay)
	order := make([]string, 0, 16)
	var mu sync.Mutex

	// コンテキスト作成
	ctx, cancel := makeContext(scanTime)
	defer cancel()

	// 終了キー監視
	handleUserCancel(scanTime, cancel)

	// --- 最初に一度だけクリア＆ヘッダを描画 ---
	fmt.Print("\033[2J\033[H")
	drawHeader()

	// 描画ループ開始
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mu.Lock()
				pruneStaleDevices(results, displayed, &order)
				drawBody(displayed, order)
				mu.Unlock()
			}
		}
	}()

	// 実際のスキャン
	err := ble.Scan(ctx, true, func(a ble.Advertisement) {
		addr := a.Addr().String()
		// フィルタ
		firstOctet, _ := strconv.ParseUint(strings.Split(addr, ":")[0], 16, 8)
		isPub := firstOctet&0xC0 == 0x00
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
		// 新規デバイスなら順序追加＆ハイライト「all」
		if _, seen := results[addr]; !seen {
			order = append(order, addr)
			displayed[addr] = entryDisplay{
				entry:     deviceEntry{addr, name, r, time.Now()},
				colorTTL:  time.Now().Add(1 * time.Second),
				highlight: "all",
			}
		} else {
			// 更新のみ
			results[addr] = deviceEntry{addr, name, r, time.Now()}
			// colorTTL は新規時のみ設定
			displayed[addr] = entryDisplay{
				entry:     results[addr],
				colorTTL:  displayed[addr].colorTTL,
				highlight: "",
			}
		}
		results[addr] = deviceEntry{addr, name, r, time.Now()}
		mu.Unlock()
	}, nil)

	// 正常終了判定
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		fmt.Println() // 最後に改行だけ入れる
		return nil
	}
	return err
}

func initBLEDevice() error {
	dev, err := linux.NewDevice()
	if err != nil {
		return fmt.Errorf("failed to initialize BLE device: %v", err)
	}
	ble.SetDefaultDevice(dev)
	return nil
}

func makeContext(seconds int) (context.Context, context.CancelFunc) {
	if seconds == 0 {
		return context.WithCancel(context.Background())
	}
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}

func handleUserCancel(seconds int, cancel context.CancelFunc) {
	if seconds == 0 {
		go func() {
			fmt.Println("Scanning... (press 'e' + Enter to exit)")
			r := bufio.NewReader(os.Stdin)
			for {
				line, _ := r.ReadString('\n')
				if strings.TrimSpace(line) == "e" {
					cancel()
					return
				}
			}
		}()
	} else {
		fmt.Printf("Scanning for %d seconds...\n", seconds)
	}
}

// drawHeader はヘッダ部のみ描画
func drawHeader() {
	fmt.Println("ADDR                 RSSI   NAME")
	fmt.Println(strings.Repeat("-", 50))
}

// drawBody はヘッダ下から各行を上書き
func drawBody(displayed map[string]entryDisplay, order []string) {
	for i, addr := range order {
		disp := displayed[addr]
		entry := disp.entry
		colS, colE := "", ""
		// 新規デバイスのみ緑ハイライト
		if disp.highlight == "all" && time.Now().Before(disp.colorTTL) {
			colS, colE = "\033[32m", "\033[0m"
		}
		// ヘッダ２行分をスキップして i+3 行目へ移動
		fmt.Printf("\033[%d;0H", i+3)
		fmt.Printf("%s%-20s %-6d %-20s%s", colS, entry.addr, entry.rssi, entry.name, colE)
		// 行末クリア
		fmt.Print("\033[K")
	}
}

// pruneStaleDevices は最後受信から10秒経過したデバイスを削除
func pruneStaleDevices(results map[string]deviceEntry, displayed map[string]entryDisplay, order *[]string) {
	cutoff := time.Now().Add(-10 * time.Second)
	newOrder := (*order)[:0]
	for _, addr := range *order {
		if ent, ok := results[addr]; ok && ent.seen.After(cutoff) {
			newOrder = append(newOrder, addr)
		} else {
			delete(results, addr)
			delete(displayed, addr)
		}
	}
	*order = newOrder
}
