package commands

import (
    "context"
    "fmt"
    "sort"
    "time"

    "github.com/go-ble/ble"
    "github.com/go-ble/ble/linux"
    "github.com/inancgumus/screen"
    "github.com/spf13/cobra"
)

type deviceEntry struct {
    addr string
    name string
    rssi int
    seen time.Time
}

var scanCommand = &cobra.Command{
    Use:   "scan",
    Short: "Scan for nearby Bluetooth devices.",
    RunE: func(cmd *cobra.Command, args []string) error {
        dev, err := linux.NewDevice()
        if err != nil {
            return fmt.Errorf("failed to initialize BLE device: %v", err)
        }
        ble.SetDefaultDevice(dev)

        // 表示初期化
        screen.Clear()
        screen.MoveTopLeft()

        results := make(map[string]deviceEntry)
        ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), 10*time.Second))

        go func() {
            for {
                screen.MoveTopLeft()
                fmt.Println("Scanning for 10 seconds...\n")
                fmt.Printf("%-20s %-6s %-20s\n", "ADDR", "RSSI", "NAME")
                fmt.Println(strings.Repeat("-", 48))

                keys := make([]string, 0, len(results))
                for k := range results {
                    keys = append(keys, k)
                }
                sort.Strings(keys)

                for _, addr := range keys {
                    dev := results[addr]
                    fmt.Printf("%-20s %-6d %-20s\n", dev.addr, dev.rssi, dev.name)
                }

                time.Sleep(500 * time.Millisecond)
            }
        }()

        return ble.Scan(ctx, false, func(a ble.Advertisement) {
            name := a.LocalName()
            if name == "" {
                name = "(no name)"
            }
            results[a.Addr().String()] = deviceEntry{
                addr: a.Addr().String(),
                name: name,
                rssi: a.RSSI(),
                seen: time.Now(),
            }
        }, nil)
    },
}

func init() {
    rootCommand.AddCommand(scanCommand)
}