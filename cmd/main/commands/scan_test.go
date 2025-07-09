package commands

import (
    "testing"
    "time"

    "github.com/spf13/cobra"
)

/* -------------------------------------------------------------
   1. pruneStaleDevices のテスト
----------------------------------------------------------------*/
func TestPruneStaleDevices(t *testing.T) {
    now := time.Now()
    // set up: 3 devices, one is stale (>10s ago)
    results := map[string]deviceEntry{
        "AA:BB:CC:00:00:01": {addr: "AA:BB:CC:00:00:01", rssi: -40, seen: now},
        "AA:BB:CC:00:00:02": {addr: "AA:BB:CC:00:00:02", rssi: -50, seen: now.Add(-11 * time.Second)}, // stale
        "AA:BB:CC:00:00:03": {addr: "AA:BB:CC:00:00:03", rssi: -60, seen: now.Add(-5 * time.Second)},  // recent
    }
    displayed := map[string]entryDisplay{
        "AA:BB:CC:00:00:01": {entry: results["AA:BB:CC:00:00:01"]},
        "AA:BB:CC:00:00:02": {entry: results["AA:BB:CC:00:00:02"]},
        "AA:BB:CC:00:00:03": {entry: results["AA:BB:CC:00:00:03"]},
    }
    order := []string{
        "AA:BB:CC:00:00:01",
        "AA:BB:CC:00:00:02",
        "AA:BB:CC:00:00:03",
    }

    pruneStaleDevices(results, displayed, &order)

    if len(order) != 2 {
        t.Fatalf("expected 2 devices after prune, got %d", len(order))
    }
    for _, addr := range order {
        if addr == "AA:BB:CC:00:00:02" {
            t.Errorf("stale device still present in order slice")
        }
    }
    if _, ok := results["AA:BB:CC:00:00:02"]; ok {
        t.Errorf("stale device still present in results map")
    }
    if _, ok := displayed["AA:BB:CC:00:00:02"]; ok {
        t.Errorf("stale device still present in displayed map")
    }
}

/* -------------------------------------------------------------
   2. runScanCommand の排他フラグエラーテスト
      (InitDefaultAdapter や BLE スキャンを呼ぶ前にエラーが返るルート)
----------------------------------------------------------------*/
func TestRunScanCommand_MutuallyExclusiveFlags(t *testing.T) {
    // 保存 & 復元
    oldRand, oldPub := randOnly, pubOnly
    t.Cleanup(func() {
        randOnly = oldRand
        pubOnly = oldPub
    })

    randOnly = true
    pubOnly = true
    cmd := &cobra.Command{}
    if err := runScanCommand(cmd, nil); err == nil {
        t.Fatalf("expected error when both --rand and --pub are set, got nil")
    }
}
