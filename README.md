# peekbt
A command-line tool for Bluetooth devices discovery, connection handling, and activity tracking.

[![Go Report Card](https://goreportcard.com/badge/github.com/ozsys/peekbt)](https://goreportcard.com/report/github.com/ozsys/peekbt)
[![Coverage Status](https://coveralls.io/repos/github/ozsys/peekbt/badge.svg?branch=main)](https://coveralls.io/github/ozsys/peekbt?branch=main)
## Overview
peekbt is a command-line tool that allows you to scan, observe, and interact with Bluetooth devices around you.  
It is designed for use on Linux-based systems (including Raspberry Pi) and supports both Classic Bluetooth and BLE (Bluetooth Low Energy).  
You can filter scan results, inspect devices in detail, and log their activity from the terminal.

## Usage
```text
peek <COMMAND> [OPTIONS]...
COMMAND
    scan                  Scan nearby Bluetooth devices
    info       <ADDR>     Show device information
OPTIONS
    --rand                Random address only.
    --pub                 Public address only.
    -t, --time <INT>      Scan duration in seconds.
    -j, --json <FILENAME> Write device information in JSON format to <FILENAME> (only available with the "info")command)          

    --help                Print help message and usage.
ADDR
    The target Bluetooth device address. (e.g. 01:23:45:67:89:AB)
```

## Example
```
# 5秒間スキャンし、結果をターミナルに表示
peekbt scan -t 5

# 15秒間スキャンしてパブリックアドレスのみ表示
peekbt scan -t 15 --pub

# MACアドレスを指定して詳細情報を取得（標準出力）
peekbt info 01:23:45:67:89:AB

# 詳細情報を JSON ファイルに書き出し
peekbt info -j device-info.json 01:23:45:67:89:AB
```
