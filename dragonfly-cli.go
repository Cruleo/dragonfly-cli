package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/gousb"
)

var (
	interfaceArray []*gousb.Interface
	polling_rate   int
	product_id     string
	motion_sync    string
	debounce       int
	dev            *gousb.Device
)

func init() {
	flag.IntVar(&polling_rate, "pr", 0, "Polling rate to set device to (allowed values: 125, 250, 500, 1000, 2000, 4000)")
	flag.IntVar(&debounce, "db", 0, "Debounce delay to set (allowed values: 0, 1, 2, 4, 8, 15, 20)")
	flag.StringVar(&product_id, "pid", "", "Product ID matching the device (default: 4K VGN Dongle)")
	flag.StringVar(&motion_sync, "ms", "", "Turn Motion Sync on or off (allowed values: on, off)")
	flag.Parse()
	if product_id == "" {
		product_id = "f505"
	}
}

func main() {

	// If more than 1 flag were given and the second isn't PID flag, exit
	if flag.NFlag() > 1 {
		if product_id == "f505" {
			fmt.Println("Changing more than 1 setting at one go is not supported")
			os.Exit(0)
		}
	}
	// Open device with VID 3554 and either default PID F505 or one given per command line arg
	openDevice()
	defer dev.Close()

	// Set auto detach to true to automatically detach and attach kernel drivers
	errorVar := dev.SetAutoDetach(true)
	if errorVar != nil {
		log.Fatalf("Couldn't set auto detach on")
	}

	// Get active config that we want to change settings for
	cfg := getActiveConfig()
	defer cfg.Close()

	// Claim all interfaces by defining them and adding them to the interface array
	claimInterfacesForConfig(cfg)
	for x := range interfaceArray {
		defer interfaceArray[x].Close()
	}

	// If motion sync arg is valid, set motion sync to given value
	if is_motion_sync_valid() {
		setMotionSync()
	} else if motion_sync != "" {
		fmt.Println("Invalid motion sync setting received, ignoring...")
	}

	// Verify that polling rate is valid and if it is, set polling rate
	if is_polling_rate_valid() {
		setPollingRate()
	} else if polling_rate != 0 { // if polling rate is not valid and is not 0, then a wrong value was given
		fmt.Println("Invalid polling rate setting received, ignoring...")
	}

	// If debounce arg valid, set to given value
	if is_debounce_valid() {
		setDebounce()
	} else if debounce != 0 {
		fmt.Println("Invalid debounce setting received, ignoring...")
	}

}

func is_polling_rate_valid() bool {
	switch polling_rate {
	case 125, 250, 500, 1000, 2000, 4000:
		return true
	default:
		return false
	}
}

func is_debounce_valid() bool {
	switch debounce {
	case 1, 2, 4, 8, 15, 20:
		return true
	default:
		return false
	}
}

func is_motion_sync_valid() bool {
	switch motion_sync {
	case "on", "off":
		return true
	default:
		return false
	}
}

func openDevice() {
	ctx := gousb.NewContext()
	defer ctx.Close()
	product_id_int, err := strconv.ParseUint(product_id, 16, 16)
	if err != nil {
		log.Fatalf("Couldn't convert Product ID to int")
	}
	device, err := ctx.OpenDeviceWithVIDPID(0x3554, gousb.ID(product_id_int))
	if err != nil {
		log.Fatalf("Couldn't open device with VID %x and PID %x", 0x3554, product_id)
	}
	dev = device
}

func getActiveConfig() *gousb.Config {
	cfgNum, err := dev.ActiveConfigNum()
	if err != nil {
		log.Fatalf("Couldn't get active config num")
	}
	cfg, err := dev.Config(cfgNum)
	if err != nil {
		log.Fatalf("Couldn't get config")
	}
	return cfg
}

func sendHidReport(data []byte) int {
	response, err := dev.Control(0x21, 0x09, 0x208, 1, data)
	if err != nil {
		log.Fatalf("Error sending control")
	}
	return response
}

func claimInterfacesForConfig(cfg *gousb.Config) {
	for _, iface := range cfg.Desc.Interfaces {
		intf, err := cfg.Interface(iface.Number, 0)
		if err != nil {
			log.Fatalf("Failed to obtain interface")
		}
		interfaceArray = append(interfaceArray, intf)
	}
}

func setPollingRate() {
	var data []byte
	switch polling_rate {
	case 125:
		data = []byte{0x08, 0x07, 0x00, 0x00, 0x00, 0x06, 0x08, 0x4d, 0x01, 0x54, 0x00, 0x55, 0x00, 0x00, 0x00, 0x00, 0x41}
	case 250:
		data = []byte{0x08, 0x07, 0x00, 0x00, 0x00, 0x06, 0x04, 0x51, 0x01, 0x54, 0x00, 0x55, 0x00, 0x00, 0x00, 0x00, 0x41}
	case 500:
		data = []byte{0x08, 0x07, 0x00, 0x00, 0x00, 0x06, 0x02, 0x53, 0x01, 0x54, 0x00, 0x55, 0x00, 0x00, 0x00, 0x00, 0x41}
	case 1000:
		data = []byte{0x08, 0x07, 0x00, 0x00, 0x00, 0x06, 0x01, 0x54, 0x01, 0x54, 0x00, 0x55, 0x00, 0x00, 0x00, 0x00, 0x41}
	case 2000:
		data = []byte{0x08, 0x07, 0x00, 0x00, 0x00, 0x06, 0x10, 0x45, 0x01, 0x54, 0x00, 0x55, 0x00, 0x00, 0x00, 0x00, 0x41}
	case 4000:
		data = []byte{0x08, 0x07, 0x00, 0x00, 0x00, 0x06, 0x20, 0x35, 0x01, 0x54, 0x00, 0x55, 0x00, 0x00, 0x00, 0x00, 0x41}
	default:
		log.Fatalln("Unexpected polling rate:", polling_rate)
	}
	response := sendHidReport(data)
	if response == 17 {
		fmt.Println("Polling rate set to", polling_rate)
	}
}

func setDebounce() {
	var data []byte
	switch debounce {
	case 1:
		data = []byte{0x08, 0x07, 0x00, 0x00, 0xa9, 0x0a, 0x01, 0x54, 0x01, 0x54, 0x06, 0x4f, 0x00, 0x55, 0x00, 0x55, 0xea}
	case 2:
		data = []byte{0x08, 0x07, 0x00, 0x00, 0xa9, 0x0a, 0x02, 0x53, 0x01, 0x54, 0x06, 0x4f, 0x00, 0x55, 0x00, 0x55, 0xea}
	case 4:
		data = []byte{0x08, 0x07, 0x00, 0x00, 0xa9, 0x0a, 0x04, 0x51, 0x01, 0x54, 0x06, 0x4f, 0x00, 0x55, 0x00, 0x55, 0xea}
	case 8:
		data = []byte{0x08, 0x07, 0x00, 0x00, 0xa9, 0x0a, 0x08, 0x4d, 0x01, 0x54, 0x06, 0x4f, 0x00, 0x55, 0x00, 0x55, 0xea}
	case 15:
		data = []byte{0x08, 0x07, 0x00, 0x00, 0xa9, 0x0a, 0x15, 0x40, 0x01, 0x54, 0x06, 0x4f, 0x00, 0x55, 0x00, 0x55, 0xea}
	case 20:
		data = []byte{0x08, 0x07, 0x00, 0x00, 0xa9, 0x0a, 0x14, 0x41, 0x01, 0x54, 0x06, 0x4f, 0x00, 0x55, 0x00, 0x55, 0xea}
	default:
		log.Fatalln("Unexpected debounce setting:", debounce)
	}
	response := sendHidReport(data)
	if response == 17 {
		fmt.Println("Debounce set to", debounce)
	}
}

func setMotionSync() {
	var data []byte
	switch motion_sync {
	case "on":
		data = []byte{0x08, 0x07, 0x00, 0x00, 0xa9, 0x0a, 0x00, 0x55, 0x01, 0x54, 0x06, 0x4f, 0x00, 0x55, 0x00, 0x55, 0xea}
	case "off":
		data = []byte{0x08, 0x07, 0x00, 0x00, 0xa9, 0x0a, 0x00, 0x55, 0x00, 0x55, 0x06, 0x4f, 0x00, 0x55, 0x00, 0x55, 0xea}
	default:
		log.Fatalln("Unexpected MotionSync keyword ", motion_sync)
	}
	response := sendHidReport(data)
	if response == 17 {
		fmt.Println("Motion sync has been turned ", motion_sync)
	}
}
