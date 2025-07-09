package commands

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/go-ble/ble"
	"github.com/spf13/cobra"
)

/* ---------- モック Scanner ---------- */
type mockScanner struct {
	fn func(context.Context, bool, ble.AdvHandler, ble.AdvFilter) error
}

func (m mockScanner) Scan(ctx context.Context, b bool,
	h ble.AdvHandler, f ble.AdvFilter) error {
	return m.fn(ctx, b, h, f)
}

/* ---------- 1. pruneStaleDevices ---------- */
func TestPruneStaleDevices(t *testing.T) {
	now := time.Now()
	results := map[string]deviceEntry{
		"AA": {seen: now},
		"BB": {seen: now.Add(-11 * time.Second)},
	}
	displayed := map[string]entryDisplay{"AA": {}, "BB": {}}
	order := []string{"AA", "BB"}

	pruneStaleDevices(results, displayed, &order)

	if len(order) != 1 || order[0] != "AA" {
		t.Fatalf("prune failed, got %v", order)
	}
}

/* ---------- 2. 排他フラグエラー ---------- */
func TestRunScanCommand_MutualEx(t *testing.T) {
	oldRand, oldPub := randOnly, pubOnly
	defer func() { randOnly, pubOnly = oldRand, oldPub }()
	randOnly, pubOnly = true, true

	if err := runScanCommand(&cobra.Command{}, nil); err == nil {
		t.Fatalf("expected error when both flags are set")
	}
}

/* ---------- 3. ヘッダ描画 ---------- */
func TestDrawHeader(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStd := os.Stdout
	os.Stdout = w

	drawHeader()

	w.Close()
	os.Stdout = oldStd
	out, _ := io.ReadAll(r)

	if !bytes.Contains(out, []byte("ADDR")) {
		t.Fatalf("header output missing ADDR column")
	}
}

func TestMakeContext(t *testing.T) {
	ctx, cancel := makeContext(0)
	defer cancel()
	if _, ok := ctx.Deadline(); ok {
		t.Fatalf("0 秒指定は deadline 無しのはず")
	}

	ctx2, cancel2 := makeContext(1)
	defer cancel2()
	if dl, ok := ctx2.Deadline(); !ok || time.Until(dl) > time.Second {
		t.Fatalf("deadline 1s がセットされていない")
	}
}
