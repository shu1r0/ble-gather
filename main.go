package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	f, err := os.OpenFile("gather.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	mulwriter := io.MultiWriter(os.Stdout, f)
	log.SetFlags(log.Ldate | log.Ldate)
	log.SetOutput(mulwriter)

	s := NewScanner(5 * time.Second)
	go s.StartScan()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit

	s.StopScan()
}
