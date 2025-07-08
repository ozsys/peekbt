+++ 
title      = "info コマンドの使い方"
date       = "2025-06-20T16:00:00+09:00"
draft      = false
tags       = ["使い方", "info", "Bluetooth"]
+++

# `info` コマンドの基本

`peekbt info <ADDR>` は、指定した Bluetooth デバイスの詳細情報を取得するコマンドです。  
`scan` コマンドでデバイスを見つけたあと、その MAC アドレスを渡すことで以下の情報が得られます。

- **Address**        : デバイスの MAC アドレス  
- **Address Type**   : Public／Resolvable Private／Non-Resolvable Private／Static Random  
- **Name**           : デバイスの名前  
- **RSSI**           : 受信信号強度 (dBm)  
- **Services UUIDs** : 提供されているサービスの UUID 一覧  
- **Last Seen**      : 情報取得時刻  
- **Connectable**    : 接続可能かどうか（true/false）

---

## オプション

| フラグ               | 説明                                             |
|----------------------|--------------------------------------------------|
| `-t`, `--timeout`    | スキャン待機時間を秒単位で指定（デフォルト: 10） |
| `-j`, `--json`       | JSON 形式で出力                                  |
| `-h`, `--help`       | ヘルプを表示                                     |

---

## 使い方例

### 無限に待つ場合

```bash
# デフォルトのタイムアウト 10 秒 → 無制限で待ちたい場合は -t 0 を指定
sudo ./peekbt info -t 0 11:22:33:44:55:66
