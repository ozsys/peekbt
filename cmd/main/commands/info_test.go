package commands

import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/go-ble/ble"
)

/* -------------------------------------------------------------
   1. validateAddr のテスト
----------------------------------------------------------------*/
func TestValidateAddr(t *testing.T) {
    valid := []string{
        "aa:bb:cc:dd:ee:ff",
        "01:23:45:67:89:ab",
    }
    invalid := []string{
        "aa:bb:cc:dd:ee",       // 5 区切り
        "gg:hh:ii:jj:kk:ll",   // 16 進文字でない
        "01:23:45:67:89:ab:cd", // 7 区切り
        "0123:4567:89ab",      // フォーマット違い
    }
    for _, addr := range valid {
        if err := validateAddr(addr); err != nil {
            t.Errorf("validateAddr(%s) returned error: %v", addr, err)
        }
    }
    for _, addr := range invalid {
        if err := validateAddr(addr); err == nil {
            t.Errorf("validateAddr(%s) should have failed", addr)
        }
    }
}

/* -------------------------------------------------------------
   2. getAddressType のテスト
----------------------------------------------------------------*/
func TestGetAddressType(t *testing.T) {
    cases := []struct {
        addr string
        want string
    }{
        {"00:00:00:00:00:00", "Public"},
        {"40:00:00:00:00:00", "Resolvable Private"},
        {"80:00:00:00:00:00", "Non-Resolvable Private"},
        {"C0:00:00:00:00:00", "Static Random"},
    }
    for _, c := range cases {
        if got := getAddressType(c.addr); got != c.want {
            t.Errorf("getAddressType(%s) = %s, want %s", c.addr, got, c.want)
        }
    }
}

/* -------------------------------------------------------------
   3. NewTimeoutCtx のテスト
----------------------------------------------------------------*/
func TestNewTimeoutCtx(t *testing.T) {
    // 0 秒 → deadline 無し
    ctx, cancel := NewTimeoutCtx(0)
    t.Cleanup(cancel)
    if _, ok := ctx.Deadline(); ok {
        t.Errorf("context with 0 sec should not have deadline")
    }

    // 1 秒 → deadline あり（±100ms 許容）
    ctx2, cancel2 := NewTimeoutCtx(1)
    t.Cleanup(cancel2)
    dl, ok := ctx2.Deadline()
    if !ok {
        t.Fatalf("expected deadline")
    }
    if d := time.Until(dl); d < 900*time.Millisecond || d > 1100*time.Millisecond {
        t.Errorf("deadline out of range: %v", d)
    }
}

/* -------------------------------------------------------------
   4. buildDeviceInfo のテスト
      ble.Advertisement を最小限実装したスタブでモック
----------------------------------------------------------------*/

type stubAdv struct {
    addr        ble.Addr
    name        string
    rssi        int
    uuids       []ble.UUID
    connectable bool
}

func (s stubAdv) Addr() ble.Addr                       { return s.addr }
func (s stubAdv) RSSI() int                            { return s.rssi }
func (s stubAdv) Services() []ble.UUID                 { return s.uuids }
func (s stubAdv) LocalName() string                    { return s.name }
func (s stubAdv) Connectable() bool                    { return s.connectable }
// 以下はインタフェース充足のためのダミー実装
func (s stubAdv) ManufacturerData() []byte             { return nil }
func (s stubAdv) ServiceData() []ble.ServiceData       { return nil }
func (s stubAdv) TxPowerLevel() int                    { return 0 }
func (s stubAdv) SolicitedServiceUUIDs() []ble.UUID { return nil }
func (s stubAdv) SolicitedService() []ble.UUID      { return nil }
func (s stubAdv) OverflowService() []ble.UUID          { return nil }
func (s stubAdv) Raw() []byte                          { return nil }

func TestBuildDeviceInfo(t *testing.T) {
    mac := "01:23:45:67:89:ab"
    adv := stubAdv{
        addr:        ble.NewAddr(mac),
        name:        "test-device",
        rssi:        -50,
        uuids:       []ble.UUID{ble.MustParse("180D")},
        connectable: true,
    }
    got := buildDeviceInfo(adv)
    if got.Address != mac {
        t.Errorf("unexpected Address: %s", got.Address)
    }
    if got.Name != "test-device" || got.RSSI != -50 {
        t.Errorf("unexpected fields: %+v", got)
    }
    if got.AddressType == "Unknown" {
        t.Errorf("AddressType should not be Unknown")
    }
}

/* -------------------------------------------------------------
   5. writeJSON のテスト
----------------------------------------------------------------*/
func TestWriteJSON(t *testing.T) {
    tmp := t.TempDir()
    file := filepath.Join(tmp, "dev.json")

    info := deviceInfo{
        Address:      "aa:bb:cc:dd:ee:ff",
        AddressType:  "Public",
        Name:         "foo",
        RSSI:         -42,
        ServicesUUID: []string{"180F"},
        LastSeen:     time.Now().Format(time.RFC3339),
        Connectable:  false,
    }

    if err := writeJSON(info, file); err != nil {
        t.Fatalf("writeJSON error: %v", err)
    }

    bs, err := os.ReadFile(file)
    if err != nil {
        t.Fatalf("read file: %v", err)
    }
    var got deviceInfo
    if err := json.Unmarshal(bs, &got); err != nil {
        t.Fatalf("unmarshal: %v", err)
    }
    if got.Address != info.Address || got.Name != info.Name {
        t.Errorf("mismatch: want %+v, got %+v", info, got)
    }
}
