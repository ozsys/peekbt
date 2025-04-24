# peekbt
A command-line tool for Bluetooth devices discovery, connection handling, and activity tracking.

## Overview
peekbt is a command-line tool that allows you to scan, observe, and interact with Bluetooth devices around you.  
It is designed for use on Linux-based systems (including Raspberry Pi) and supports both Classic Bluetooth and BLE (Bluetooth Low Energy).  
You can filter scan results, inspect devices in detail, and log their activity from the terminal.

## Usage
```text
peek <COMMAND> [OPTIONS]...
COMMAND
    scan
    info    <ADDR>
    export   [JSON_FILE]
OPTIONS
    -rand      Random address only.
    -pub       Public address only.
    -lim <INT> Limit the number of displayed devices.
ADDR
    The target Bluetooth device address. (e.g. XX:XX:XX:XX:XX:XX)
JSON_FILE
    The output file for exported scan data.
```

## Installation

## About
