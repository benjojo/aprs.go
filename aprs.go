package main

import (
	"fmt"
	"strings"
)

type APRSPacket struct {
	Callsign      string // Done!
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
	Destination   string // Done!
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
}

func ParseAPRSPacket(input string) (p APRSPacket, e error) {
	if input == "" {
		e = fmt.Errorf("Could not parse the packet because the packet line is blank")
		return p, e
	}

	if !strings.Contains(input, ">") {
		e = fmt.Errorf("This libary does not support this kind of packet.")
		return p, e
	}
	p = APRSPacket{}
	CommaParts := strings.Split(input, ",")
	RouteString := CommaParts[0]
	RouteParts := strings.Split(RouteString, ">")[0]
	if len(RouteParts) != 2 {
		e = fmt.Errorf("There was more than one > in the route part of the packet")
		return p, e
	}
	p.Callsign = RouteParts[0]
	p.Destination = RouteParts[1]

	return p, e
}
