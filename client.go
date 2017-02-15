package main

import (
	"fmt"
	"time"

	"github.com/chriswalker/gpsd"
)

func watchHandler(r interface{}) {
	resp := r.(*gpsd.WatchResponse)
	fmt.Printf("%s [%t]\n", resp.Class, resp.Json)
}

func tpvHandler(r interface{}) {
	resp := r.(*gpsd.TimePositionVelocityResponse)
	fmt.Printf("%s [%s] %f / %f\n", resp.Class, resp.Time, resp.Lat, resp.Lon)
}

func main() {
	gpsd.HandleFunc(gpsd.Watch, watchHandler)
	gpsd.HandleFunc(gpsd.Tpv, tpvHandler)
	err := gpsd.ConnectAndListen("192.168.1.100:2947")
	if err != nil {
		fmt.Println("Error connecting to gpsd:", err)
	}

	// Just output for 15 secs
	time.Sleep(15000 * time.Millisecond)
}
