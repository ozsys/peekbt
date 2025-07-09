package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-ble/ble"
	"github.com/spf13/cobra"
)

var (
	infoTimeout int
	infoJSON    string
)

func init() {
	infoCmd.Flags().IntVarP(&infoTimeout, "timeout", "t", 10, "Scan timeout in seconds")
	infoCmd.Flags().StringVarP(&infoJSON, "json", "j", "", "Write JSON output to the specified file")
	rootCommand.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info [flags] <ADDR>",
	Short: "Show detailed information for a specific Bluetooth device.",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfoCommand,
}

// runInfoCommand はフラグ取得と各処理の呼び出しだけを行います
func runInfoCommand(cmd *cobra.Command, args []string) error {
	addr := strings.ToLower(args[0])
	if err := validateAddr(addr); err != nil {
		return err
	}

	// BLE デバイス初期化
	if _, err := InitDefaultAdapter(); err != nil {
		return err
	}

	// アドバタイズ取得
	fmt.Printf("Scanning for device %s (timeout %ds)...\n", addr, infoTimeout)
	adv, err := scanAdvertisement(addr, time.Duration(infoTimeout)*time.Second)
	if err != nil {
		return err
	}

	// 構造体組み立て
	info := buildDeviceInfo(adv)

	// JSON or Key:Value
	if infoJSON != "" {
		return writeJSON(info, infoJSON)
	}
	printInfo(info)
	return nil
}

// validateAddr は引数がMACアドレス形式かをチェックします
func validateAddr(addr string) error {
	pat := regexp.MustCompile(`^([0-9a-f]{2}:){5}[0-9a-f]{2}$`)
	if !pat.MatchString(addr) {
		return fmt.Errorf("invalid address format: %s", addr)
	}
	return nil
}

// scanAdvertisement はタイムアウト内に Addr が見つかるまでスキャンします
func scanAdvertisement(addr string, timeout time.Duration) (ble.Advertisement, error) {
	ctx, cancel := NewTimeoutCtx(int(timeout.Seconds()))
	defer cancel()

	ch := make(chan ble.Advertisement, 1)
	go func() {
		ble.Scan(ctx, true, func(a ble.Advertisement) {
			if strings.EqualFold(a.Addr().String(), addr) {
				select {
				case ch <- a:
					cancel()
				default:
				}
			}
		}, nil)
	}()

	select {
	case adv := <-ch:
		return adv, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("device %s not found within %v", addr, timeout)
	}
}

// deviceInfo は出力用の構造体
type deviceInfo struct {
	Address      string   `json:"address"`
	AddressType  string   `json:"addressType"`
	Name         string   `json:"name"`
	RSSI         int      `json:"rssi"`
	ServicesUUID []string `json:"serviceUUIDs"`
	LastSeen     string   `json:"lastSeen"`
	Connectable  bool     `json:"connectable"`
}

// buildDeviceInfo は Advertisement から deviceInfo を組み立てます
func buildDeviceInfo(a ble.Advertisement) deviceInfo {
	uuids := a.Services()
	s := make([]string, len(uuids))
	for i, u := range uuids {
		s[i] = u.String()
	}
	return deviceInfo{
		Address:      a.Addr().String(),
		AddressType:  getAddressType(a.Addr().String()),
		Name:         a.LocalName(),
		RSSI:         a.RSSI(),
		ServicesUUID: s,
		LastSeen:     time.Now().Format(time.RFC3339),
		Connectable:  a.Connectable(),
	}
}

// writeJSON は JSON 形式でファイルに書き出します
func writeJSON(info deviceInfo, file string) error {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(file, data, 0o644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}
	return nil
}

// printInfo は key: value 形式で標準出力します
func printInfo(info deviceInfo) {
	fmt.Printf("Address        : %s\n", info.Address)
	fmt.Printf("Address Type   : %s\n", info.AddressType)
	fmt.Printf("Name           : %s\n", info.Name)
	fmt.Printf("RSSI           : %d dBm\n", info.RSSI)
	fmt.Printf("Services UUIDs : %v\n", info.ServicesUUID)
	fmt.Printf("Last Seen      : %s\n", info.LastSeen)
	fmt.Printf("Connectable    : %t\n", info.Connectable)
}

// getAddressType はアドレスの MSB から種別を返します
func getAddressType(addrStr string) string {
	b, _ := strconv.ParseUint(strings.Split(addrStr, ":")[0], 16, 8)
	switch b & 0xC0 {
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
