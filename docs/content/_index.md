+++
title = "peekbt"
date  = 2025-07-10
draft = false
+++

peekbt は **Linux / Raspberry Pi** 向けに作られた  
シンプル & 高速な Bluetooth CLI ツールです。  

* **BLE / Classic** の両方をスキャン  
* RSSI・アドレス種別・サービス UUID などを即座に表示  
* JSON 形式でのエクスポート・ログ取りもサポート  

インストールしてすぐに “測位・デバッグ・製品検証” に利用できます。

---

## インストール & セットアップ

```bash
# Go 1.21+ が入っていればワンライナーで OK
go install github.com/ozsys/peekbt/cmd/peekbt@latest

# 実行パスに $GOBIN もしくは ~/go/bin が通っているか確認
peekbt --help
```

**Raspberry Pi では追加で以下のパッケージが必要です。**

```bash
sudo apt-get update
sudo apt-get install -y bluetooth bluez libbluetooth-dev
```

---

## クイックスタート

| 操作                 | コマンド例                                 | 補足                         |
|----------------------|--------------------------------------------|------------------------------|
| 10 秒間スキャン      | `peekbt scan -t 10`                        | アドレス・RSSI を一覧表示    |
| ランダムアドレスのみ | `peekbt scan --rand`                       | MSB 0xC0 ビット判定          |
| 端末詳細             | `peekbt info AA:BB:CC:DD:EE:FF`            | サービス UUID・TxPower 等    |
| JSON で保存          | `peekbt scan -t 30 \| tee devices.json`    | 解析ツールへパイプ可能       |

---

## 主なコマンド

### `scan`

近傍デバイスを連続スキャンし、端末を行単位で更新表示。

* `--time, -t <N>` — スキャン時間（0 で無限）
* `--rand` — ランダムアドレスのみ
* `--pub`  — パブリックアドレスのみ

### `info`

指定アドレスをアクティブスキャンし、**1 回だけ** アドバタイズを解析。

```bash
peekbt info -t 15 AA:BB:CC:DD:EE:FF
```

---

## よくある質問（FAQ）

<details>
<summary>Windows で動きますか？</summary>

いいえ。現状は BlueZ に依存しているため Linux 系 OS のみサポートしています。
</details>

<details>
<summary>RSSI が正確でないように見えます</summary>

アンテナ／チップセット差や反射など環境要因が大きいため、**相対値** として利用してください。
</details>

---

## コントリビュート

1. `develop` ブランチへ PR  
2. `go test ./...` & `just fmt vet` が通ること  
3. 新コマンドの場合は `docs` 下に使い方ページを追加

---

## ライセンス

Apache 2.0  
© 2025 ozsys
