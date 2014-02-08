package main

// This is a direct port of https://gist.github.com/benjojo/0124c7875113831a4274

import (
	"fmt"
	"strconv"
	"strings"
)

// PacketType can be:
// * Status Report
// * GPGGA
// * New Mic-E
// * Old Mic-E
// * Location
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

		// Lat
		DegLatitude := SplitData[2]
		DegLatMin, e := strconv.ParseFloat(DegLatitude[2:2+len(DegLatitude)-2], 64)
		if e != nil {
			e = fmt.Errorf("Could not decode the DegLatMin part of the GPGGA packet")
			return p, e
		}
		DegLatMin = DegLatMin / 60
		StrDegLatMin := fmt.Sprintf("%f", DegLatMin)
		p.Latitude = fmt.Sprintf("%s%s", DegLatitude[:2], StrDegLatMin[1:len(StrDegLatMin)-1])

		if SplitData[3] == "S" {
			p.Latitude = fmt.Sprintf("-%s", p.Latitude)
		}

		// Long
		DegLongitude := SplitData[4]
		DegLonMin, e := strconv.ParseFloat(DegLongitude[3:3+len(DegLongitude)-3], 64)
		if e != nil {
			e = fmt.Errorf("Could not decode the DegLonMin part of the GPGGA packet")
			return p, e
		}
		DegLonMin = DegLonMin / 60
		StrDegLonMin := fmt.Sprintf("%f", DegLonMin)
		p.Longitude = fmt.Sprintf("%s%s", DegLongitude[:3], StrDegLonMin[1:len(StrDegLonMin)-1])

		if SplitData[3] == "W" {
			p.Longitude = fmt.Sprintf("-%s", p.Longitude)
		}

		f, e := strconv.ParseFloat(SplitData[9], 64)
		if e != nil {
			e = fmt.Errorf("Could not decode the Altitude part of the GPGGA packet")
			return p, e
		}
		p.Altitude = fmt.Sprintf("%f", f)
	}

	// Test if the packet is a Mic-E packet
	if strings.Index(input, ":`") != -1 || strings.Index(input, ":'") != -1 {
		// MicEPtr := 0
		if strings.Index(input, ":`") != -1 {
			p.PacketType = "New Mic-E"
			// MicEPtr = strings.Index(input, ":`")
		} else {
			p.PacketType = "Old Mic-E"
			// MicEPtr = strings.Index(input, ":'")
		}
		e = fmt.Errorf("Mic-E is currently not supported")
		return p, e
	}

	// Test to see if its a location packet

	// FirstChr = line.IndexOf(":/")  // With Timestamp
	// SecondChr = line.IndexOf(":!") // Without Timestamp
	// ThirdChr = line.IndexOf(":@")  // With Timestamp and APRS Messaging
	// FourthChr = line.IndexOf(":=") // Without Timestamp and Messaging
	LocationPtr := 0
	TimestampPtr := strings.Index(input, ":!")
	if strings.Index(input, ":@") != -1 { // With Timestamp and APRS Messaging
		LocationPtr = strings.Index(input, ":@")
	}

	if strings.Index(input, ":=") != -1 { // Without Timestamp and Messaging
		LocationPtr = strings.Index(input, ":=")
	}
	LocationPtrStr := LocationSlice(input, int64(LocationPtr), 8)
	TimestampPtrStr := LocationSlice(input, int64(TimestampPtr), 9)

	if (LocationPtr != -1 && (LocationPtrStr == "H" || LocationPtrStr == "Z" || LocationPtrStr == "/")) ||
		(TimestampPtr != -1 && (TimestampPtrStr == "S" || TimestampPtrStr == "N")) {

		p.PacketType = "Location"
		if LocationPtr != -1 && (LocationPtrStr == "H" || LocationPtrStr == "Z" || LocationPtrStr == "/") {
			p.GPSTime = input[TimestampPtr+2 : TimestampPtr+8]
		} else {
			TimestampPtr = TimestampPtr - 7
		}
		p.Symbol = fmt.Sprintf("%s%s", input[LocationPtr+17:LocationPtr+18], input[LocationPtr+27:LocationPtr+28])

		// Lat
		DegLatMin, e := strconv.ParseFloat(input[LocationPtr+11:LocationPtr+16], 64)
		if e != nil {
			// Go bang
			e = fmt.Errorf("Could not make DegLatMin a float")
			return p, e
		}
		DegLatMin = (DegLatMin / 60)

		DegLatMinStr := fmt.Sprintf("%f", DegLatMin)
		p.Latitude = fmt.Sprintf("%s%s", input[LocationPtr+9:LocationPtr+11], DegLatMinStr[1:len(DegLatMinStr)-1])

		if input[LocationPtr+16:LocationPtr+17] == "S" {
			p.Latitude = fmt.Sprintf("-%s", p.Latitude)
		}

		// Long
		DegLonMin, e := strconv.ParseFloat(input[LocationPtr+21:LocationPtr+26], 64)
		if e != nil {
			// Go bang
			e = fmt.Errorf("Could not make DegLonMin a float")
			return p, e
		}
		DegLonMin = (DegLonMin / 60)

		DegLonMinStr := fmt.Sprintf("%f", DegLonMin)
		p.Longitude = fmt.Sprintf("%s%s", input[LocationPtr+19:LocationPtr+21], DegLonMinStr[1:len(DegLonMinStr)-1])

		if input[LocationPtr+16:LocationPtr+17] == "S" {
			p.Longitude = fmt.Sprintf("-%s", p.Longitude)
		}

	}

	return p, e
}

func LocationSlice(line string, LocationPtr int64, ptr int64) string {
	bit := strings.ToUpper(line[LocationPtr+ptr : (LocationPtr+ptr)+1])
	return bit
}
