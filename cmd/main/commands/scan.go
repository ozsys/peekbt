package commands

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/go-ble/ble"
    "github.com/go-ble/ble/linux"
    "github.com/spf13/cobra"
)

var scanCommand = &cobra.Command{
    Use:   "scan",
    Short: "Scan for nearby Bluetooth devices.",
    RunE: func(cmd *cobra.Command, args []string) error {
        dev, err := linux.NewDevice()
        if err != nil {
            return fmt.Errorf("failed to initialize BLE device: %v", err)
        }
        ble.SetDefaultDevice(dev)

        ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), 10*time.Second))
        fmt.Println("Scanning for 10 seconds...")

        err = ble.Scan(ctx, false, func(a ble.Advertisement) {
            fmt.Printf("Found: %s [%s] RSSI:%d\n", a.LocalName(), a.Addr(), a.RSSI())
        }, nil)
        if err != nil {
            return fmt.Errorf("scan failed: %v", err)
        }

        return nil
    },
}

func init() {
    rootCommand.AddCommand(scanCommand)
}