// Package main_test は peekbt エントリポイントの簡易テストを含む。
// 目的: goMain の終了コードと出力を検証する。
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

// TestHello は hello 関数の戻り値を検証する。
func TestHello(t *testing.T) {
	got := hello()
	want := "Welcome to peekbt!"
	if got != want {
		t.Errorf("hello() = %q, want %q", got, want)
	}
}

// Example_hello は go doc のサンプルとして hello の出力例を示す。
func Example_hello() {
	fmt.Println(hello())
	// Output:
	// Welcome to peekbt!
}

// TestGoMain_Failure は引数なし実行時にエラーとなり、終了コード 1／出力メッセージを検証する。
func TestGoMain_Failure(t *testing.T) {
	// 引数を「peek」のみの状態にする
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"peek"}

	// 標準出力をパイプでキャプチャ
	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w

	code := goMain()

	w.Close()
	os.Stdout = origStdout
	out, _ := io.ReadAll(r)

	if code != 1 {
		t.Errorf("goMain() = %d, want 1 on failure", code)
	}
	if !bytes.Contains(out, []byte("no command specified")) {
		t.Errorf("output %q does not contain 'no command specified'", string(out))
	}
}

// TestGoMain_Help は --help 実行時に成功コード 0／Usage を含む出力を検証する。
func TestGoMain_Help(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"peek", "--help"}

	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w

	code := goMain()

	w.Close()
	os.Stdout = origStdout
	out, _ := io.ReadAll(r)

	if code != 0 {
		t.Errorf("goMain() = %d, want 0 on help", code)
	}
	if !bytes.Contains(out, []byte("Usage")) {
		t.Errorf("help output %q does not contain 'Usage'", string(out))
	}
}
