package gpsd

import "time"

type Mode byte

// Basic struct for checking the type of response we've just received;
// all we care about is the class field
type Report struct {
	Class GpsdResponse
}

// Base contains elements common to most reports
type Base struct {
	Class  GpsdResponse
	Device string
	Time   time.Time
}

// TPV (Time-Position-Velocity) report
type TimePositionVelocityResponse struct {
	Base
	Mode  Mode
	Ept   float64
	Lat   float64
	Lon   float64
	Alt   float64
	Epx   float64
	Epy   float64
	Track float64
	Speed float64
	Climb float64
	Epd   float64
	Eps   float64
	Epc   float64
}

// SKY
type SkyResponse struct {
	Base
	Xdop       float64
	Ydop       float64
	Vdop       float64
	Tdop       float64
	Hdop       float64
	Pdop       float64
	Gdop       float64
	Satellites []Satellite
}

// Satellite (list in SKY)
type Satellite struct {
	PRN            int
	Azimuth        float64 `json:"az"`
	Elevation      float64 `json:"el"`
	SignalStrength float64 `json:"ss"`
	Used           bool
}

// GST (Pseudo-range noise) report

// ATT (Attitude) report; compass/gyroscope sensors

// VERSION
type VersionResponse struct {
	Class         string
	Release       string
	Revision      string
	ProtocolMajor int
	ProtocolMinor int
	Remote        string
}

// WATCH
type WatchResponse struct {
	Class   string
	Enable  bool
	Json    bool
	Nmea    bool
	Raw     int // Type for this?
	Scaled  bool
	Split24 bool
	Pps     bool
	Device  string
	Remote  string
}

// DEVICES

// DEVICE (list in DEVICES)

// ERROR
