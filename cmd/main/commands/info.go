// cmd/main/commands/info.go
package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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
	infoCommand.Flags().IntP("timeout", "t", 10, "Scan timeout in seconds")
	infoCommand.Flags().BoolP("json", "j", false, "Output in JSON format")
	rootCommand.AddCommand(infoCommand)
}

// deviceInfo holds the output structure
type deviceInfo struct {
	Address      string   `json:"address"`
	AddressType  string   `json:"addressType"`
	Name         string   `json:"name"`
	RSSI         int      `json:"rssi"`
	ServicesUUID []string `json:"serviceUUIDs"`
	LastSeen     string   `json:"lastSeen"`
	Connectable  bool     `json:"connectable"`
}

// outputInformation prints key: value format for non-JSON output
func outputInformation(info deviceInfo) {
	fmt.Printf("Address        : %s\n", info.Address)
	fmt.Printf("Address Type   : %s\n", info.AddressType)
	fmt.Printf("Name           : %s\n", info.Name)
	fmt.Printf("RSSI           : %d dBm\n", info.RSSI)
	fmt.Printf("Services UUIDs : %v\n", info.ServicesUUID)
	fmt.Printf("Last Seen      : %s\n", info.LastSeen)
	fmt.Printf("Connectable    : %t\n", info.Connectable)
}

// getAddressType returns the address type based on the MSBs of the first octet
func getAddressType(addrStr string) string {
	firstOctet, _ := strconv.ParseUint(strings.Split(addrStr, ":")[0], 16, 8)
	switch firstOctet & 0xC0 {
	case 0x00:
		return "Public" // 00xxxxxx
	case 0x40:
		return "Resolvable Private" // 01xxxxxx
	case 0x80:
		return "Non-Resolvable Private" // 10xxxxxx
	case 0xC0:
		return "Static Random" // 11xxxxxx
	default:
		return "Unknown"
	}
}

func runInfoCommand(cmd *cobra.Command, args []string) error {
	// Initialize BLE device
	dev, err := linux.NewDevice()
	if err != nil {
		return fmt.Errorf("failed to initialize BLE device: %v", err)
	}
	ble.SetDefaultDevice(dev)

	// Validate address format
	addr := strings.ToLower(args[0])
	macPattern := regexp.MustCompile(`^([0-9a-f]{2}:){5}[0-9a-f]{2}$`)
	if !macPattern.MatchString(addr) {
		return fmt.Errorf("invalid address format: %s", args[0])
	}

	// Get flags
	tout, _ := cmd.Flags().GetInt("timeout")
	useJSON, _ := cmd.Flags().GetBool("json")

	// Perform scan and lookup
	adv, err := findAdvertisementByAddr(addr, time.Duration(tout)*time.Second)
	if err != nil {
		return err
	}

	// Determine address type via helper
	addrStr := adv.Addr().String()
	addrType := getAddressType(addrStr)

	// Convert Services UUIDs to string slice
	bleServices := adv.Services()
	svcStrs := make([]string, len(bleServices))
	for i, u := range bleServices {
		svcStrs[i] = u.String()
	}

	// Prepare info
	info := deviceInfo{
		Address:      addrStr,
		AddressType:  addrType,
		Name:         adv.LocalName(),
		RSSI:         adv.RSSI(),
		ServicesUUID: svcStrs,
		LastSeen:     time.Now().Format("2006-01-02T15:04:05Z07:00"),
		Connectable:  adv.Connectable(),
	}

	// Output based on mode
	if useJSON {
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		fmt.Println(string(data))
	} else {
		outputInformation(info)
	}
	return nil
}

// findAdvertisementByAddr scans up to timeout and returns the first match
func findAdvertisementByAddr(addr string, timeout time.Duration) (ble.Advertisement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := make(chan ble.Advertisement, 1)
	fmt.Printf("Scanning for device %s (timeout %v)...\n", addr, timeout)

	go func() {
		ble.Scan(ctx, true, func(a ble.Advertisement) {
			if strings.EqualFold(a.Addr().String(), addr) {
				select {
				case result <- a:
					cancel()
				default:
				}
			}
		}, nil)
	}()

	select {
	case adv := <-result:
		return adv, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("device %s not found within %v", addr, timeout)
	}
}
