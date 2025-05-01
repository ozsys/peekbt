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
    info    <ADDR>        Show device information
    export   [JSON_FILE]  Export results
OPTIONS
    -rand                 Random address only.
    -pub                  Public address only.
    -lim <INT>            Limit the number of displayed devices.
ADDR
    The target Bluetooth device address. (e.g. XX:XX:XX:XX:XX:XX)
JSON_FILE
    The output file for exported scan data.
```

## Installation

## About
