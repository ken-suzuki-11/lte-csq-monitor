package main

import (
	"fmt"
	"github.com/tarm/serial"
	"os"
	"strings"
)

const IsDebug = true

type GPS struct {
	Latitude  float64
	Longitude float64
}


func getLteCSQ(s *serial.Port) {
	n, err := s.Write([]byte("AT+CSQ\r"))
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := make([]byte, 1024)
	line := ""

	for {
		n, err = s.Read(buf)
		line += string(buf[:n])


		fmt.Println(line)

		fmt.Printf("%#v\n", line)

		r := strings.HasSuffix(line, "OK\r\n")
                if r {
			fmt.Println(line)
			os.Exit(-1)
                }
	}
}

func main() {
	lte := &serial.Config{
		Name: "/dev/ttyUSB2",
		Baud: 115200,
	}
	ls, err := serial.OpenPort(lte)
	if err != nil {
		fmt.Println("lte serial can't open")
		os.Exit(-1)
	}

	getLteCSQ(ls)
	os.Exit(-1)

}
