# draGOnfly-cli
A Golang CLI tool to change polling rate and debounce setting of VGN/VXE Dragonfly mice using the 1K or 4K dongles

It is a clone of my Python script, [dragonpy-cli](https://github.com/Cruleo/dragonpy-cli) but written in Golang as a learning experience.

Please note that this has only been tested on Linux and on a VXE Dragonfly R1 Pro along with a 4K and 1K dongle. Your mileage might vary.

Use this at your own risk. I do not take responsibility for any damage caused.

## Requirements
This script is based on [GoUSB](https://github.com/google/gousb) and therefore requires the dependecies of GoUSB. It was tested on Go 1.22.1.

## Features

- Setting polling rate (up to 1000 Hz with stock dongle, up to 4000 Hz with 4K dongle)
- Setting debounce delay
- Setting Motion Sync on or off

## Supported devices

Theoretically it should work with any VGN mouse connected to the 4K and 1K dongle but it hasn't been tested.

| Device name   | VendorID | ProductID |
|---------------|----------|-----------|
| VGN 4K Dongle | 3554     | f505      |
| VXE 1K Dongle | 3554     | f58a      |
| VXE R1 Pro (Wired) | 3554 | f58c      |

## To-Do

- [ ] Add lift off distance setting (LOD)
- [ ] Add idle timer setting

## Usage

```go run dragonfly-cli.go -pr {125, 250, 500, 1000, 2000, 4000} -db {1, 2, 4, 8, 15, 20} -ms {on, off} -pid PRODUCT_ID```

To get the product id on Linux, just run `lsusb`. For example, this is the outcome for the 4K VGN dongle:

```Bus 003 Device 019: ID 3554:f505 Compx VGN Dragonfly 4K Receiver```

The Vendor ID here is 3554 while the product id is f505.

## Restrictions
From my testing so far I wasn't able to send 2 control requests (2 setting changes) without the first setting getting overridden. Therefore for now, only using 1 flag is supported.

## Udev rule
For Linux, you will require a udev rule to fully access the device without superuser rights. This can easily be done by creating a udev rule under `/etc/udev/rules.d/`. Feel free to look at the [example file](./51-vxe-mouse.rules).
