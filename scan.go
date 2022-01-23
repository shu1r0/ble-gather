package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/sausheong/ble"
	"github.com/sausheong/ble/linux"
)

type Scanner struct {
	Timeout *time.Duration
	device  *linux.Device
}

func NewScanner(timeout time.Duration) Scanner {
	d, err := linux.NewDevice()
	if err != nil {
		log.Fatal("Can't create new device: ", err)
	}
	ble.SetDefaultDevice(d)
	return Scanner{Timeout: &timeout, device: d}
}

var finish bool = false
var devices map[string]Device

type Device struct {
	MACAddress  string    `json:"address"`
	Name        string    `json:"name"`
	Timestamp   time.Time `json:"timestamp"`
	RSSI        int       `json:"rssi"`
	ResponseRaw []byte    `json:"response"`
}

func (d Device) JSON() (s string, err error) {
	j, err := json.Marshal(d)
	return string(j), err
}

func (s *Scanner) advHandler(a ble.Advertisement) {
	d := Device{
		MACAddress:  a.Addr().String(),
		Name:        a.LocalName(),
		Timestamp:   time.Now(),
		RSSI:        a.RSSI(),
		ResponseRaw: a.ScanResponseRaw(),
	}
	devices[d.MACAddress] = d
	log.Println(d.JSON())
}

func (s *Scanner) advFilter(a ble.Advertisement) bool {
	return true
}

func (s *Scanner) StartScan() {
	finish = false
	for !finish {
		ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *s.Timeout))
		ble.Scan(ctx, false, s.advHandler, s.advFilter)
	}
}

func (s *Scanner) StopScan() {
	finish = true
}
