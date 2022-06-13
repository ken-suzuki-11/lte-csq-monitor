package lib

import (
	"errors"
	"fmt"
	"github.com/adrianmo/go-nmea"
	"github.com/tarm/serial"
	"os"
	"strings"
	"time"
)

const IsDebug = false

type GoogleMapInfo struct {
	Latitude  float64
	Longitude float64
}

type GPSDevice struct {
	config *serial.Config
	port   *serial.Port
}

func NewGPSDevice(path string, baudrate int) *GPSDevice {
	gps := GPSDevice{}
	gps.config = &serial.Config{
		Name:        path,
		Baud:        baudrate,
		ReadTimeout: 5 * time.Second,
	}
	return &gps
}

func (g *GPSDevice) Open() error {
	p, err := serial.OpenPort(g.config)
	if err != nil {
		fmt.Println("gps serial can't open")
		os.Exit(-1)
	}
	g.port = p
	time.Sleep(time.Second * 2)
	return nil
}

func (g *GPSDevice) readGpsLine(returnString string) (string, error) {
	buf := make([]byte, 1024)
	line := ""
	for i := 0; i < 1024; i++ {
		n, err := g.port.Read(buf)
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

func (g *GPSDevice) GetGoogleMapValue() (GoogleMapInfo, error) {
	for {
		line, err := g.readGpsLine("\n")
		if err != nil {
			return GoogleMapInfo{0, 0}, err
		}
		gpsData, err := nmea.Parse(line)
		if err != nil {
			return GoogleMapInfo{0, 0}, err
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
			return GoogleMapInfo{
				Latitude:  m.Latitude,
				Longitude: m.Longitude,
			}, nil
		}
	}
}
