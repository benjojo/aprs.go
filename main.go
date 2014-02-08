package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	b, e := ioutil.ReadFile("./aprstestdata.txt")
	if e != nil {
		panic("Cannot read file")
	}

	lines := strings.Split(string(b), "\n")
	for _, v := range lines {
		p, e := ParseAPRSPacket(v)
		if e != nil {
			fmt.Println("Failed to decode packet" + string(e.Error()))
		} else {
			fmt.Println(p)
		}
	}
}
