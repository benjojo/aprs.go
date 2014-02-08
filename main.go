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
			panic("Failed to recode packet" + string(e))
		}
		fmt.Println(p)
	}
}
