package main

import (
	"context"
	"log"
	"time"

	"github.com/sausheong/ble"
	"github.com/sausheong/ble/linux"
)

type Scanner struct {
	Timeout *time.Duration
	logger  *log.Logger
	device  *linux.Device
}

func NewScanner(timeout time.Duration, logger *log.Logger) Scanner {
	d, err := linux.NewDevice()
	if err != nil {
		logger.Fatal("Can't create new device: ", err)
	}
	ble.SetDefaultDevice(d)
	return Scanner{Timeout: &timeout, logger: logger, device: d}
}

var finish bool = false

type Device struct {
	MACAddress  string
	Name        string
	Timestamp   time.Time
	RSSI        int
	ResponseRaw []byte
}

func (s *Scanner) advHandler(a ble.Advertisement) {
	d := Device{
		MACAddress:  a.Addr().String(),
		Name:        a.LocalName(),
		Timestamp:   time.Now(),
		RSSI:        a.RSSI(),
		ResponseRaw: a.ScanResponseRaw(),
	}
	s.logger.Println(d)
}

func (s *Scanner) advFilter(a ble.Advertisement) bool {
	return true
}

func (s *Scanner) StartScan() {
	finish = false
	for finish {
		ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *s.Timeout))
		ble.Scan(ctx, false, s.advHandler, s.advFilter)
	}
}

func (s *Scanner) StopScan() {
	finish = true
}
