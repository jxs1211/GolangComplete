package main

import (
	"fmt"
	"time"
)

const (
	idCheckTmCost   = 60
	bodyCheckTmCost = 120
	xRayCheckTmCost = 180
)

func idCheck() int {
	time.Sleep(time.Millisecond * time.Duration(idCheckTmCost))
	println("\tidCheck ok")
	return idCheckTmCost
}

func bodyCheck() int {
	time.Sleep(time.Millisecond * time.Duration(bodyCheckTmCost))
	println("\tbodyCheck ok")
	return bodyCheckTmCost
}

func xRayCheck() int {
	time.Sleep(time.Millisecond * time.Duration(xRayCheckTmCost))
	println("\txRayCheck ok")
	return xRayCheckTmCost
}

func airportSecurityCheck() int {
	println("airportSecurityCheck ...")
	total := 0

	total += idCheck()
	total += bodyCheck()
	total += xRayCheck()

	println("airportSecurityCheck ok")
	return total
}

func showConcurrency1() {
	total := 0
	passengers := 30
	for i := 0; i < passengers; i++ {
		total += airportSecurityCheck()
	}
	println("total time cost:", total)
}

const (
	idCheckTimeout = 60
	bodyTimeout    = 120
	xRayTimeout    = 180
)

func idCheck2() int {
	time.Sleep(time.Millisecond * time.Duration(idCheckTimeout))
	fmt.Println("\tidCheck2 ok")
	return idCheckTimeout
}

func bodyCheck2() int {
	time.Sleep(time.Millisecond * time.Duration(bodyTimeout))
	fmt.Println("\tbodyCheck2 ok")
	return bodyTimeout
}

func xRayCheck2() int {
	time.Sleep(time.Millisecond * time.Duration(xRayTimeout))
	fmt.Println("\txRayCheck2 ok")
	return xRayTimeout
}

func AirportCheck() int {
	total := 0
	total += idCheck2()
	total += bodyCheck2()
	total += xRayCheck2()
	return total
}

func showSerialization() {
	println("start")
	total := 0
	passager := 30
	for i := 0; i < passager; i++ {
		total += AirportCheck()
	}
	println("total:", total)
	println("end")
}

func main() {
	showSerialization()
}
