package lib

import (
	"errors"
	"fmt"
	"github.com/tarm/serial"
	"os"
	"strings"
	"time"
)

const ErrorCSQ = "nil,nil"

type LTEDevice struct {
	config *serial.Config
	port   *serial.Port
}

func NewLTEDevice(path string, baudrate int) *LTEDevice {
	lte := LTEDevice{}
	lte.config = &serial.Config{
		Name:        path,
		Baud:        baudrate,
		ReadTimeout: 5 * time.Second,
	}
	return &lte
}

func (l *LTEDevice) Open() error {
	p, err := serial.OpenPort(l.config)
	if err != nil {
		fmt.Println("lte serial can't open")
		os.Exit(-1)
	}
	l.port = p
	time.Sleep(time.Second * 2)
	return nil
}

func (l LTEDevice) parseCSQ(line string) string {
	firstSplitData := strings.Split(line, "+CSQ: ")
	if len(firstSplitData) != 2 {
		return ErrorCSQ
	}
	secondSplitData := strings.Split(firstSplitData[1], "\r\n")
	return secondSplitData[0]
}

func (l *LTEDevice) write() error {
	n, err := l.port.Write([]byte("AT+CSQ\r"))
	if n != len("AT+CSQ\r") {
		return errors.New("csq command write error")
	}
	if err != nil {
		return err
	}
	return nil
}

func (l *LTEDevice) read() string {
	buf := make([]byte, 1024)
	line := ""
	for {
		n, err := l.port.Read(buf)
		if err != nil {
			return ErrorCSQ
		}
		line += string(buf[:n])
		r := strings.HasSuffix(line, "OK\r\n")
		if r {
			csq := l.parseCSQ(line)
			return csq
		}
	}
}

func (l *LTEDevice) GetCSQ() string {
	err := l.write()
	if err != nil {
		return ErrorCSQ
	}
	return l.read()
}
