+++
title = "peekbt"
date  = 2025-07-10
draft = false
+++

# peekbtについて
peekbt は **Linux / Raspberry Pi / Ubuntu** 向けに作られた  
シンプル & 高速な Bluetooth CLI ツールです。
topコマンドのように1画面を更新し続けるようになっているのでログを遡ってデバイスの情報を見る手間が減ります。

* **BLE / Classic** の両方をスキャンに対応しています 
* RSSI・アドレス種別・サービス UUID などを表示できます
* JSON 形式での出力にも対応しています

##
---

## インストール & セットアップ

```bash
# 実行にはbluetooth bluez libbluetooth-devが必要です
sudo apt-get update
sudo apt-get install -y bluetooth bluez libbluetooth-dev

#
# Go 1.21+ が入っていればワンライナーで OK
go install github.com/ozsys/peekbt/cmd/peekbt@latest

# 実行パスに $GOBIN もしくは ~/go/bin が通っているか確認
peekbt --help
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

近くのデバイスを連続スキャンし、端末を行単位で更新表示します。  
表示する情報はデバイスアドレス、RSSI、デバイス名です(デバイス名はついていなければ**no name**になります)。  
新規デバイスは緑色でハイライトされます。    

#### オプション
オプションでアドレスタイプの指定ができます。

* `--time, -t <INT>` — スキャン時間(デフォルト10s)
* `--rand` — ランダムアドレスのみ
* `--pub`  — パブリックアドレスのみ

>[!NOTE]
>パブリックアドレスはアドレス上位1桁が0〜3で始まるアドレスです。  
>ランダムアドレスはアドレス上位1桁が4〜Fで始まるアドレスです。
### `info`

指定アドレスをアクティブスキャンし、アドレスがヒットした場合はアドバタイズパケットを元に情報を提示します。
```bash
peekbt info AA:BB:CC:DD:EE:FF
```
#### オプション
オプションでアドレスタイプの指定ができます。

* `--time, -t <INT>` — スキャン時間(デフォルト10s)
* `--json, -j <FILENAME>` — FILENAMEを指定してjson形式で出力

#### `peekbt scan` の実行結果の例
```bash
Scanning for device fb:2f:14:22:80:2f (timeout 10s)...
Address        : fb:2f:14:22:80:2f
Address Type   : Static Random
Name           : 
RSSI           : -45 dBm
Services UUIDs : []
Last Seen      : 2025-07-11T01:27:19+09:00
Connectable    : false
```
#### `peekbt scan` のスキャン結果について

`peekbt scan` を実行すると、検出されたBluetoothデバイスの情報が以下のように表示されます。

| 表示項目 | 説明 |
| :--- | :--- |
| **`Address`** | 検出したデバイスのMACアドレスです。(例: `AA:BB:CC:DD:EE:FF`) |
| **`Address Type`** | デバイスのアドレス種別です。以下の4種類があります。`public`, `Resolvable Private`, `non-Resolvable Private`, `Static Random`|
| **`Name`** | デバイスに設定されている名称です。（空の場合もあります） |
| **`RSSI`** | 受信電波の強度です (単位: dBm)。<br>数値が `0` に近いほど、電波が強いことを示します。 |
| **`Services UUIDs`** | デバイスが提供しているサービスのUUIDリストです。（なければ `[]` と表示されます） |
| **`Last Seen`** | 最後にデバイスを検出した時刻です。 |
| **`Connectable`** | そのデバイスに接続できるかどうかを `true` / `false` で示します。|
----

#### json形式での出力
```bash
$ peekbt info -j example.json fb:2f:14:22:80:2f

$ cat example.json 
{
  "address": "fb:2f:14:22:80:2f",
  "addressType": "Static Random",
  "name": "",
  "rssi": -45,
  "serviceUUIDs": [],
  "lastSeen": "2025-07-11T01:27:50+09:00",
  "connectable": false
}
```
## peekbtの作成理由

peekbtを作成しようと思ったきっかけは`bluetoothctl`というコマンドが不便に感じたからです。  
`bluetoothctl`は`connect`などの接続関連のUIは素晴らしいですが、デバイスのスキャンやデバイス情報に関しては不満点が残ります。`bluetoothctl`で`scan`を実行すると結果がログのように流れ続けるような仕様になっています。  
そのため目的のデバイス情報が画面上部へどんどん流れていってしまい、見失ってしまいます。  
以下は `bluetoothctl scan` を実行した際のログ例です。
```bash
…
[CHG] Device 00:1C:FC:B6:C4:EA RSSI: -61
[CHG] Device 47:E0:1A:C5:A8:30 RSSI: -61
[DEL] Device 71:27:2A:C7:EF:CA 71-27-2A-C7-EF-CA
[CHG] Device 50:40:36:39:68:5C RSSI: -94
[CHG] Device 70:E5:4E:CF:CF:7C RSSI: -76
[NEW] Device 6B:1C:36:09:A6:74 6B-1C-36-09-A6-74
[DEL] Device 67:99:86:A7:02:25 67-99-86-A7-02-25
[CHG] Device 50:40:36:39:68:5C RSSI: -47
[CHG] Device 70:E5:4E:CF:CF:7C RSSI: -84
…
```
昨今、年間の出荷台数は50億台を突破し、bluetoothデバイスは身の回りに溢れています。そのためアドバタイジングパケットが無数に送信されており、目的のデバイスの情報がすぐに流れてしまい観察に不向きです。そのため`peekbt`では`scan`の情報はtopコマンドのように、terminalの画面内で更新するようにしました。
```bash
ADDR                 RSSI   NAME
--------------------------------------------------
4f:35:62:ba:a2:09    -53    (no name)           
71:7e:87:f2:d2:e6    -63    (no name)           
12:57:42:3b:b9:f0    -73    (no name)           
6d:f6:9d:71:e9:2f    -63    (no name)           
64:38:bd:6c:5c:c5    -65    (no name)           
60:37:05:26:37:5c    -47    (no name)           
7f:0d:ae:e9:f4:37    -55    (no name)           
7a:4e:a3:2e:3c:e9    -63    (no name)  
```
----
## peekbtの今後の発展
* 接続に関するコマンドの実装
* `info`で表示される情報量の増加(企業の固有番号の表示など)
* `json`フォーマット以外での出力
* 各コマンドのオプションの拡充
     
----
## ライセンス
peekbtはMITライセンスの下で公開されています。
