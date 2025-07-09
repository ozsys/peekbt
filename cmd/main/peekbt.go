// Package main は peekbt 実行ファイルのエントリポイントを提供する。
// 主な役割は CLI コマンドを初期化し、終了コードを適切に返却するだけ。
package main

import (
	"fmt"
	"os"

	"github.com/ozsys/peekbt/cmd/main/commands"
)

// hello は簡易的なウェルカムメッセージを返す。
// 主にサンプル用・ユニットテスト用に利用する。
func hello() string {
	return "Welcome to peekbt!"
}

// goMain はコマンドを実行して終了コードを算出する。
// • エラーが発生した場合 → 標準出力へメッセージを表示し 1 を返す  
// • 正常終了の場合       → 0 を返す
func goMain() int {
	if err := commands.Execute(); err != nil {
		fmt.Println(err) // ヘルプを出力
		return 1         // 異常終了コードを返す
	}
	return 0 // 正常終了
}

// main は実行開始点。
// goMain の戻り値をそのままプロセスの終了ステータスとして返す。
func main() {
	status := goMain()
	os.Exit(status) // プロセス終了コードを設定
}
