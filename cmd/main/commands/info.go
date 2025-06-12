package commands

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/spf13/cobra"
)

// infoCommand represents the 'info' CLI command
var infoCommand = &cobra.Command{
	Use:   "info <ADDR>",
	Short: "Show detailed information for a specific Bluetooth device.",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfoCommand,
}

func init() {
	// Register the info command and its flags
	infoCommand.Flags().IntP("timeout", "t", 10, "Scan timeout in seconds")
	rootCommand.AddCommand(infoCommand)
}

// runInfoCommand executes the info logic
func runInfoCommand(cmd *cobra.Command, args []string) error {
	// Initialize BLE device
	dev, err := linux.NewDevice()
	if err != nil {
		return fmt.Errorf("failed to initialize BLE device: %v", err)
	}
	ble.SetDefaultDevice(dev)

	// Normalize and validate address
	addr := strings.ToLower(args[0])
	macPattern := regexp.MustCompile(`^([0-9a-f]{2}:){5}[0-9a-f]{2}$`)
	if !macPattern.MatchString(addr) {
		return fmt.Errorf("invalid address format: %s", addr)
	}

	// Get timeout flag
	tout, _ := cmd.Flags().GetInt("timeout")

	// Perform BLE scan and lookup
	adv, err := findAdvertisementByAddr(addr, time.Duration(tout)*time.Second)
	if err != nil {
		return err
	}

	// Output key: value format
	fmt.Printf("Address        : %s\n", adv.Addr().String())
	fmt.Printf("Name           : %s\n", adv.LocalName())
	fmt.Printf("RSSI           : %d dBm\n", adv.RSSI())
	fmt.Printf("Services UUIDs : %v\n", adv.Services())
	fmt.Printf("Last Seen      : %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("Connectable    : %t\n", adv.Connectable())
	return nil
}

// findAdvertisementByAddr scans up to timeout, logging all advertisements, returning the first match
func findAdvertisementByAddr(addr string, timeout time.Duration) (ble.Advertisement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Channel to receive matching advertisement
	result := make(chan ble.Advertisement, 1)

	// Debug: notify start
	fmt.Printf(">>> [DEBUG] Starting scan for %s (timeout %v)\n", addr, timeout)

	// Start scan in goroutine
	go func() {
		ble.Scan(ctx, true, func(a ble.Advertisement) {
			// Debug: log every advertisement (no clear, logs will remain)
			fmt.Printf(">>> [DEBUG] Got ADV: addr=%s, RSSI=%d, name=%q, services=%v\n",
				a.Addr().String(), a.RSSI(), a.LocalName(), a.Services(),
			)

			// Match address (case-insensitive)
			if strings.EqualFold(a.Addr().String(), addr) {
				fmt.Println(">>> [DEBUG] Address matched, stopping scan")
				select {
				case result <- a:
					cancel()
				default:
				}
			}
		}, nil)
	}()

	// Wait for match or timeout
	select {
	case adv := <-result:
		return adv, nil
	case <-ctx.Done():
		fmt.Println(">>> [DEBUG] Timeout expired, no match found.")
		return nil, fmt.Errorf("device %s not found within %v", addr, timeout)
	}
}
