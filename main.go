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
			//return line, nil
			return strings.TrimRight(line, returnString), nil
		}
	}
	return "", errors.New("reach the read limit")
}

func getGoogleMapValue(line string) (string, error) {
	gpsData, err := nmea.Parse(line)
	if err != nil {
		return "", err
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

		return fmt.Sprintf("%f,%f", m.Latitude, m.Longitude), nil
	}
	return "", errors.New("not rmc data")
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

		googleMapValue, err := getGoogleMapValue(line)
		if err != nil {
			continue
		}
		fmt.Printf("GPS Coordinates : %s\n", googleMapValue)

		fmt.Println()
	}
}
