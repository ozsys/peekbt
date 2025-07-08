package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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
	Use:   "info [flags] <ADDR>",
	Short: "Show detailed information for a specific Bluetooth device.",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfoCommand,
}

func init() {
	infoCommand.Flags().IntP("timeout", "t", 10, "Scan timeout in seconds")
	infoCommand.Flags().StringP("json", "j", "", "Write JSON output to the specified file")
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
		return "Public"
	case 0x40:
		return "Resolvable Private"
	case 0x80:
		return "Non-Resolvable Private"
	case 0xC0:
		return "Static Random"
	default:
		return "Unknown"
	}
}

func runInfoCommand(cmd *cobra.Command, args []string) error {
	// 引数とフラグの取得
	addr := strings.ToLower(args[0])
	macPattern := regexp.MustCompile(`^([0-9a-f]{2}:){5}[0-9a-f]{2}$`)
	if !macPattern.MatchString(addr) {
		return fmt.Errorf("invalid address format: %s", args[0])
	}
	tout, _ := cmd.Flags().GetInt("timeout")
	jsonFile, _ := cmd.Flags().GetString("json")

	// BLE デバイス初期化
	dev, err := linux.NewDevice()
	if err != nil {
		return fmt.Errorf("failed to initialize BLE device: %v", err)
	}
	ble.SetDefaultDevice(dev)

	// スキャン開始メッセージ
	fmt.Printf("Scanning for device %s (timeout %ds)...\n", addr, tout)

	// スキャンしてアドバタイズ取得
	adv, err := findAdvertisementByAddr(addr, time.Duration(tout)*time.Second)
	if err != nil {
		return err
	}

	// 情報整形
	addrStr := adv.Addr().String()
	addrType := getAddressType(addrStr)
	bleServices := adv.Services()
	svcStrs := make([]string, len(bleServices))
	for i, u := range bleServices {
		svcStrs[i] = u.String()
	}

	info := deviceInfo{
		Address:      addrStr,
		AddressType:  addrType,
		Name:         adv.LocalName(),
		RSSI:         adv.RSSI(),
		ServicesUUID: svcStrs,
		LastSeen:     time.Now().Format(time.RFC3339),
		Connectable:  adv.Connectable(),
	}

	// JSON出力または標準出力
	if jsonFile != "" {
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		// ファイル末尾に改行を付加
		data = append(data, '\n')
		if err := os.WriteFile(jsonFile, data, 0644); err != nil {
			return fmt.Errorf("failed to write JSON file: %v", err)
		}
		return nil
	}

	outputInformation(info)
	return nil
}

// findAdvertisementByAddr scans up to timeout and returns the first match
func findAdvertisementByAddr(addr string, timeout time.Duration) (ble.Advertisement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := make(chan ble.Advertisement, 1)
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
