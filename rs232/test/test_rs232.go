// test_rs232.go, jpad 2013
// Test serial package: interact with an Arduino @ 115200 baud.

package main

import (
	"fmt"
	"io"
	"time"

	"packages/communication/rs232"
)

func read_raw(r io.Reader) {
	buf := make([]byte, 128)
	for {
		n, err := r.Read(buf)
		if err != nil {
			fmt.Printf("serial.Read(): %s\n", err)
		}
		if n > 0 {
			fmt.Printf("read %d bytes: %v\n", n, buf[:n])
		}
	}
}

func read_raw_polltime(r io.Reader) {
	buf := make([]byte, 128)
	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				time.Sleep(10 * time.Millisecond)
			} else {
				fmt.Printf("serial.Read(): %s\n", err)
			}
		}
		if n > 0 {
			fmt.Printf("read %d bytes: %v\n", n, buf[:n])
		}
	}
}

func read_line(r io.Reader) {
	line := make([]byte, 128)
	buf := make([]byte, 1)
	for {
		line = line[:0] // clear
		// read until '\n' is found
		for {
			n, err := r.Read(buf)
			if err != nil {
				fmt.Printf("serial.Read(): %s\n", err)
			}
			if n > 0 {
				if buf[0] == '\n' {
					break
				}
				line = append(line, buf...)
			}
		}
		fmt.Printf("%s\n", line)
	}
}

func read_line_poll(r io.Reader) {
	line := make([]byte, 128)
	buf := make([]byte, 1)
	for {
		line = line[:0] // clear
		// read until '\n' is found
		for {
			n, err := r.Read(buf)
			if err != nil {
				if err == io.EOF {
					time.Sleep(10 * time.Millisecond)
				} else {
					fmt.Printf("serial.Read(): %s\n", err)
				}
			}
			if n > 0 {
				if buf[0] == '\n' {
					break
				}
				line = append(line, buf...)
			}
		}
		fmt.Printf("%s\n", line)
	}
}

func GetGarbage(p *rs232.Port) {
	n, err := p.BytesAvailable()
	if err != nil {
		fmt.Printf("rs232.BytesAvailable(): %s\n", err)
	}
	if n > 0 {
		garbage := make([]byte, n)
		_, err := p.Read(garbage)
		if err != nil {
			fmt.Printf("rs232.Read(): %s\n", err)
			return
		}
		fmt.Printf("---- garbage (%d bytes):\n", n)
		fmt.Printf("%s", garbage)
		fmt.Print("---- end of garbage\n")
	}
}

func main() {
	fmt.Printf("Package rs232 test starting up...\n")

	options := rs232.Options{
		BitRate:  115200,
		DataBits: 8,
		StopBits: 1,
		Parity:   rs232.PARITY_NONE,
		Timeout:  0,
	}

	serial, err := rs232.Open("/dev/tty.usbmodemfa131", options)
	if err != nil {
		fmt.Printf("rs232.Open(): %s\n", err)
		e := err.(*rs232.Error)
		errType := ""
		switch e.Code {
		case rs232.ERR_DEVICE:
			errType = "ERR_DEVICE"
		case rs232.ERR_ACCESS:
			errType = "ERR_ACCESS"
		case rs232.ERR_PARAMS:
			errType = "ERR_PARAMS"
		}
		fmt.Printf("Error code: %d (%s)\n", e.Code, errType)
		return
	}
	defer serial.Close()
	fmt.Printf("Opened serial port %s\n", serial.String())

	// 'garbage' check
	time.Sleep(500 * time.Millisecond)
	GetGarbage(serial)

	//read_raw(serial)
	read_line(serial)

	fmt.Printf("Package rs232 test shutting down...\n")
}
