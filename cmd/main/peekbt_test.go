// Package main_test は peekbt エントリポイントの簡易テストを含む。
// 目的: hello 関数の戻り値を確認し、Example ドキュメントも生成する。
package main

import (
	"fmt"
	"testing"
)

// Example_hello は go doc のサンプルとして hello の出力例を示す。
// 実行時にコメントの Output セクションと結果を突き合わせ検証する。
func Example_hello() {
	fmt.Println(hello())
	// Output:
	// Welcome to peekbt!
}

// TestHello は hello 関数の文字列が期待値どおりか検証する。
func TestHello(t *testing.T) {
	got := hello()
	want := "Welcome to peekbt!"
	if got != want {
		t.Errorf("hello() = %q, want %q", got, want)
	}
}
