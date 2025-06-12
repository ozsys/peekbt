package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
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
	highlight string    // "all" for full line, "rssi" for RSSI only
}

var scanTime int

var scanCommand = &cobra.Command{
	Use:   "scan",
	Short: "Scan for nearby Bluetooth devices.",
	RunE:  runScanCommand,
}

func init() {
	scanCommand.Flags().IntP("time", "t", 0, "Scan time in seconds (0 = infinite)")
	rootCommand.AddCommand(scanCommand)
}

func runScanCommand(cmd *cobra.Command, args []string) error {
	scanTime, _ = cmd.Flags().GetInt("time")

	// Initialize BLE device
	dev, err := linux.NewDevice()
	if err != nil {
		return fmt.Errorf("failed to initialize BLE device: %v", err)
	}
	ble.SetDefaultDevice(dev)

	// Data structures
	results := make(map[string]deviceEntry)
	displayed := make(map[string]entryDisplay)
	order := make([]string, 0)
	var mu sync.Mutex

	// Context setup
	var ctx context.Context
	var cancel context.CancelFunc
	if scanTime == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(scanTime)*time.Second)
	}
	defer cancel()

	// Exit on 'e' + Enter if infinite
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

	// Initial draw
	fmt.Print("\033[2J") // clear screen
	fmt.Print("\033[H")  // move cursor to top-left
	fmt.Println("ADDR                 RSSI   NAME")
	fmt.Println(strings.Repeat("-", 50))

	// Draw loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				mu.Lock()
				for i, addr := range order {
					disp := displayed[addr]
					entry := disp.entry
					colorStart, colorEnd := "", ""
					if time.Now().Before(disp.colorTTL) {
						colorStart = "\033[32m"
						colorEnd = "\033[0m"
					}

					// move cursor to line i+4
					fmt.Printf("\033[%d;0H", i+4)

					// print with highlight scope
					switch disp.highlight {
					case "all":
						fmt.Printf("%s%-20s %-6d %-20s%s", colorStart, entry.addr, entry.rssi, entry.name, colorEnd)
					default:
						fmt.Printf("%-20s %-6d %-20s", entry.addr, entry.rssi, entry.name)
					}
					// clear rest of line
					fmt.Print("\033[K")
				}
				mu.Unlock()

				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	// BLE scan (active mode)
	return ble.Scan(ctx, true, func(a ble.Advertisement) {
		addr := a.Addr().String()
		rssi := a.RSSI()
		name := a.LocalName()
		if name == "" {
			name = "(no name)"
		}

		mu.Lock()
		defer mu.Unlock()

		// On first seen, record order
		if _, seen := results[addr]; !seen {
			order = append(order, addr)
		}

		// Update entry
		newEntry := deviceEntry{addr, name, rssi, time.Now()}
		results[addr] = newEntry

		// Determine highlight scope
		highlightScope := ""
		if _, seen := displayed[addr]; !seen {
			highlightScope = "all"
		} else if displayed[addr].entry.rssi != rssi {
			highlightScope = "rssi"
		}

		displayed[addr] = entryDisplay{
			entry:     newEntry,
			colorTTL:  time.Now().Add(1 * time.Second),
			highlight: highlightScope,
		}

	}, nil)
}
