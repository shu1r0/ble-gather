package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/sausheong/ble"
	"github.com/sausheong/ble/linux"
)

type Scanner struct {
	Timeout *time.Duration

	device  *linux.Device
	finish  bool
	mutex   sync.RWMutex
	devices map[string]Device
}

func NewScanner(timeout time.Duration) Scanner {
	d, err := linux.NewDevice()
	if err != nil {
		log.Fatal("Can't create new device: ", err)
	}
	ble.SetDefaultDevice(d)
	return Scanner{Timeout: &timeout, device: d, finish: false, mutex: sync.RWMutex{}, devices: make(map[string]Device)}
}

type Device struct {
	MACAddress    string    `json:"address"`
	Name          string    `json:"name"`
	Timestamp     time.Time `json:"timestamp"`
	RSSI          int       `json:"rssi"`
	Advertisement string    `json:"advertisement"`
	ResponseRaw   []byte    `json:"response"`
}

func (d Device) JSON() (s string) {
	j, _ := json.Marshal(d)
	return string(j)
}

func (s *Scanner) advHandler(a ble.Advertisement) {
	s.mutex.Lock()
	d := Device{
		MACAddress:    a.Addr().String(),
		Name:          a.LocalName(),
		Timestamp:     time.Now(),
		RSSI:          a.RSSI(),
		Advertisement: hex.EncodeToString(a.LEAdvertisingReportRaw()),
		ResponseRaw:   a.ScanResponseRaw(),
	}
	s.devices[d.MACAddress] = d
	s.mutex.Unlock()
	log.Println(d.JSON())
}

func (s *Scanner) advFilter(a ble.Advertisement) bool {
	return true
}

func (s *Scanner) StartScan() {
	log.Println("start scan")
	s.finish = false
	for !s.finish {
		ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *s.Timeout))
		ble.Scan(ctx, false, s.advHandler, s.advFilter)
	}
}

func (s *Scanner) StopScan() {
	log.Println("stop scan")
	s.finish = true
}
