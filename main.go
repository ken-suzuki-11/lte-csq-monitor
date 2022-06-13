package main

import (
	"flag"
	"fmt"
	"lte-csq-monitor/lib"
	"os"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"

type CSQInfo struct {
	CSQ       string
	Latitude  float64
	Longitude float64
}

func main() {
	var (
		gpsDevice string
		lteDevice string
	)
	flag.StringVar(&gpsDevice, "g", "nil", "gps device path")
	flag.StringVar(&lteDevice, "l", "nil", "lte device path")
	flag.Parse()

	// gps device open
	gpsDevicePath := "/dev/" + gpsDevice
	_, err := os.Stat(gpsDevicePath)
	if os.IsNotExist(err) {
		fmt.Println("Device Not Found : " + gpsDevicePath)
		os.Exit(-1)
	}
	gps := lib.NewGPSDevice(gpsDevicePath, 9600)
	err = gps.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// lte device open
	lteDevicePath := "/dev/" + lteDevice
	_, err = os.Stat(lteDevicePath)
	if os.IsNotExist(err) {
		fmt.Println("Device Not Found : " + lteDevicePath)
		os.Exit(-1)
	}

	lte := lib.NewLTEDevice(lteDevicePath, 115200)
	err = lte.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	for {
		// GPSデータ読み込み
		googleMapInfo, err := gps.GetGoogleMapValue()
		if err != nil {
			continue
		}
		// LTE の電波強度取得
		csq := lte.GetCSQ()
		// 出力
		fmt.Printf("%s,%f,%f,%s\n",
			time.Now().Format(TimeLayout),
			googleMapInfo.Latitude,
			googleMapInfo.Longitude,
			csq,
		)

		time.Sleep(time.Duration(5) * time.Second)
	}
}
