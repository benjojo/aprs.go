package main

import (
	"fmt"
)

type APRSPacket struct {
	Callsign      string
	PacketType    string
	Latitude      string
	Longitude     string
	Altitude      string
	GPSTime       string
	RawData       string
	Symbol        string
	Heading       string
	PHG           string
	Speed         string
	Destination   string
	Status        string
	WindDirection string
	WindSpeed     string
	WindGust      string
	WeatherTemp   string
	RainHour      string
	RainDay       string
	RainMidnight  string
	Humidity      string
	Pressure      string
	Luminosity    string
	Snowfall      string
	Raincounter   string
	Error         string
}

func ParseAPRSPacket() (p APRSPacket, e error) {

	return p, e
}
