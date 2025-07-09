package commands

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-ble/ble"
	"github.com/spf13/cobra"
)

/* -------------------------------------------------------------
   1. validateAddr
----------------------------------------------------------------*/
func TestValidateAddr(t *testing.T) {
	valid := []string{
		"aa:bb:cc:dd:ee:ff",
		"01:23:45:67:89:ab",
	}
	invalid := []string{
		"aa:bb:cc:dd:ee",
		"gg:hh:ii:jj:kk:ll",
		"01:23:45:67:89:ab:cd",
		"0123:4567:89ab",
	}
	for _, a := range valid {
		if err := validateAddr(a); err != nil {
			t.Errorf("valid %s rejected: %v", a, err)
		}
	}
	for _, a := range invalid {
		if err := validateAddr(a); err == nil {
			t.Errorf("invalid %s accepted", a)
		}
	}
}

/* -------------------------------------------------------------
   2. getAddressType
----------------------------------------------------------------*/
func TestGetAddressType(t *testing.T) {
	cases := []struct {
		addr, want string
	}{
		{"00:00:00:00:00:00", "Public"},
		{"40:00:00:00:00:00", "Resolvable Private"},
		{"80:00:00:00:00:00", "Non-Resolvable Private"},
		{"C0:00:00:00:00:00", "Static Random"},
	}
	for _, c := range cases {
		if got := getAddressType(c.addr); got != c.want {
			t.Errorf("%s => %s, want %s", c.addr, got, c.want)
		}
	}
}

/* -------------------------------------------------------------
   3. NewTimeoutCtx
----------------------------------------------------------------*/
func TestNewTimeoutCtx(t *testing.T) {
	ctx, cancel := NewTimeoutCtx(0)
	t.Cleanup(cancel)
	if _, ok := ctx.Deadline(); ok {
		t.Fatalf("deadline should be absent for 0 sec")
	}

	ctx2, cancel2 := NewTimeoutCtx(1)
	t.Cleanup(cancel2)
	if dl, ok := ctx2.Deadline(); !ok || time.Until(dl) > time.Second {
		t.Fatalf("deadline 1s not set")
	}
}

/* -------------------------------------------------------------
   4. buildDeviceInfo
----------------------------------------------------------------*/
type stubAdv struct {
	addr ble.Addr
	name string
	rssi int
}

func (s stubAdv) Addr() ble.Addr                    { return s.addr }
func (s stubAdv) RSSI() int                         { return s.rssi }
func (s stubAdv) Services() []ble.UUID              { return nil }
func (s stubAdv) LocalName() string                 { return s.name }
func (s stubAdv) Connectable() bool                 { return true }
func (s stubAdv) ManufacturerData() []byte          { return nil }
func (s stubAdv) ServiceData() []ble.ServiceData    { return nil }
func (s stubAdv) TxPowerLevel() int                 { return 0 }
func (s stubAdv) SolicitedServiceUUIDs() []ble.UUID { return nil }
func (s stubAdv) SolicitedService() []ble.UUID      { return nil }
func (s stubAdv) OverflowService() []ble.UUID       { return nil }
func (s stubAdv) Raw() []byte                       { return nil }

func TestBuildDeviceInfo(t *testing.T) {
	mac := "01:23:45:67:89:ab"
	adv := stubAdv{addr: ble.NewAddr(mac), name: "dev", rssi: -50}
	if got := buildDeviceInfo(adv); got.Address != mac || got.Name != "dev" {
		t.Errorf("unexpected buildDeviceInfo result: %+v", got)
	}
}

/* -------------------------------------------------------------
   5. writeJSON
----------------------------------------------------------------*/
func TestWriteJSON(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "d.json")
	info := deviceInfo{Address: "aa:bb:cc:dd:ee:ff", Name: "foo"}

	if err := writeJSON(info, p); err != nil {
		t.Fatalf("writeJSON: %v", err)
	}
	b, _ := os.ReadFile(p)
	var got deviceInfo
	_ = json.Unmarshal(b, &got)
	if got.Address != info.Address || got.Name != info.Name {
		t.Errorf("mismatch: %+v vs %+v", got, info)
	}
}

/* -------------------------------------------------------------
   6. printInfo 出力
----------------------------------------------------------------*/
func TestPrintInfo(t *testing.T) {
	info := deviceInfo{Address: "aa:bb:cc:dd:ee:ff"}
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	printInfo(info)
	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	if !bytes.Contains(out, []byte("aa:bb:cc:dd:ee:ff")) {
		t.Fatalf("addr missing in output: %s", out)
	}
}

/* -------------------------------------------------------------
   7. runInfoCommand : invalid MAC
----------------------------------------------------------------*/
func TestRunInfoCommand_InvalidAddr(t *testing.T) {
	if err := runInfoCommand(&cobra.Command{}, []string{"foo"}); err == nil {
		t.Fatalf("expected error on invalid addr")
	}
}
