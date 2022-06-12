package main

import (
	"errors"
	"fmt"
	"github.com/adrianmo/go-nmea"
	"github.com/tarm/serial"
	"log"
	"os"
	"strings"
)

const IsDebug = true

type GPS struct {
	Latitude  float64
	Longitude float64
}

func readGpsLine(s *serial.Port, returnString string) (string, error) {
	buf := make([]byte, 1024)
	line := ""
	for i := 0; i < 1024; i++ {
		n, err := s.Read(buf)
		if err != nil {
			return "", err
		}
		data := string(buf[:n])
		line += data
		r := strings.HasSuffix(data, returnString)
		if r {
			return strings.TrimRight(line, returnString), nil
		}
	}
	return "", errors.New("reach the read limit")
}

func getGoogleMapValue(line string) (GPS, error) {
	gpsData, err := nmea.Parse(line)
	if err != nil {
		return GPS{0.0, 0.0}, err
	}
	if gpsData.DataType() == nmea.TypeRMC {
		m := gpsData.(nmea.RMC)
		if IsDebug {
			fmt.Printf("Raw sentence: %v\n", m)
			fmt.Printf("Time: %s\n", m.Time)
			fmt.Printf("Validity: %s\n", m.Validity)
			fmt.Printf("Speed: %f\n", m.Speed)
			fmt.Printf("Course: %f\n", m.Course)
			fmt.Printf("Date: %s\n", m.Date)
			fmt.Printf("Variation: %f\n", m.Variation)
		}
		return GPS{
			Latitude:  m.Latitude,
			Longitude: m.Longitude,
		}, nil
	}
	return GPS{0.0, 0.0}, errors.New("not rmc data")
}

func main() {
	c := &serial.Config{
		Name: "/dev/ttyUSB0",
		Baud: 9600,
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	_, err = readGpsLine(s, "\n")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	for {
		// GPSデータ読み込み
		line, err := readGpsLine(s, "\n")
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		gpsInfo, err := getGoogleMapValue(line)
		if err != nil {
			continue
		}
		fmt.Printf("GPS Coordinates : %f,%f\n", gpsInfo.Latitude, gpsInfo.Longitude)

		fmt.Println()
	}
}
