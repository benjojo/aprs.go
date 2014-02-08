package main

import (
	"fmt"
	"strconv"
	"strings"
)

// PacketType can be:
// * Status Report
// * GPGGA
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
	RouteParts := strings.Split(RouteString, ">")
	if len(RouteParts) != 2 {
		e = fmt.Errorf("There was not one > in the route part of the packet, dunno how to decode this")
		return p, e
	}
	p.Callsign = RouteParts[0]
	p.Destination = RouteParts[1]

	LocationOfStatusMarker := strings.Index(input, ":>")
	LocationOfNormalMarker := strings.Index(input, ">")

	if LocationOfStatusMarker > LocationOfNormalMarker {
		p.PacketType = "Status Report"
		RawArray := []byte(input[LocationOfStatusMarker+2 : (LocationOfStatusMarker+2)+(len(input)-LocationOfStatusMarker-2)])
		if len(RawArray) > 6 && strings.ToLower(string(RawArray[6])) == "z" {
			p.GPSTime = input[LocationOfStatusMarker+2 : LocationOfStatusMarker+8]
			p.Status = input[LocationOfStatusMarker+2 : (LocationOfStatusMarker+2)+len(input)-LocationOfStatusMarker-9]
		} else {
			p.Status = input[LocationOfStatusMarker+2 : (LocationOfStatusMarker+2)+len(input)-LocationOfStatusMarker-2]
		}
	}

	// Test if the packet is a GPGGA packet
	if strings.Contains(input, ":$GPGGA,") {
		p.PacketType = "GPGGA"
		GPGGALocation := strings.Index(input, ":$GPGGA,")
		RawData := input[GPGGALocation : GPGGALocation+(len(input)-GPGGALocation)]
		SplitData := strings.Split(RawData, ",")
		if len(SplitData) < 9 {
			e = fmt.Errorf("There was not enough data inside the GPGGA packet to decode it")
			return p, e
		}
		p.GPSTime = SplitData[1]
		DegLatitude := SplitData[2]
		DegLatMin, e := strconv.ParseFloat(DegLatitude[2:2+len(DegLatitude)-2], 64)
		if e != nil {
			e = fmt.Errorf("Could not decode the DegLatMin part of the GPGGA packet")
			return p, e
		}
		DegLatMin = DegLatMin / 60
		// Latitude = degLatitude.Substring(0, 2) + Convert.ToString(degLatMin).Substring(1, Convert.ToString(degLatMin).Length - 1);
		StrDegLatMin := fmt.Sprintf("%f", DegLatMin)
		p.Latitude = fmt.Sprintf("%s%s", DegLatitude[:2], StrDegLatMin[1:len(StrDegLatMin)-1])
	}

	return p, e
}
