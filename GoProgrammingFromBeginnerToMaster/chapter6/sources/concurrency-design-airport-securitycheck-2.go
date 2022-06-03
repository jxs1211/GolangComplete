package main

import "time"

const (
	idCheckTmCost   = 60
	bodyCheckTmCost = 120
	xRayCheckTmCost = 180
)

func idCheck(id int) int {
	time.Sleep(time.Millisecond * time.Duration(idCheckTmCost))
	println("\t goroutine ", id, ": idCheck ok")
	return idCheckTmCost
}

func bodyCheck(id int) int {
	time.Sleep(time.Millisecond * time.Duration(bodyCheckTmCost))
	println("\t goroutine ", id, ": bodyCheck ok")
	return bodyCheckTmCost
}

func xRayCheck(id int) int {
	time.Sleep(time.Millisecond * time.Duration(xRayCheckTmCost))
	println("\t goroutine ", id, ": xRayCheck ok")
	return xRayCheckTmCost
}

func airportSecurityCheck(id int) int {
	print("goroutine-", id, ": airportSecurityCheck ...\n")
	total := 0

	total += idCheck(id)
	total += bodyCheck(id)
	total += xRayCheck(id)

	print("goroutine-", id, ": airportSecurityCheck ok\n")
	return total
}

func start(id int, f func(int) int, queue <-chan struct{}) <-chan int {
	c := make(chan int)
	go func() {
		total := 0
		for {
			_, ok := <-queue
			if !ok {
				c <- total
				return
			}
			total += f(id)
		}
	}()
	return c
}

func max(args ...int) int {
	n := 0
	for _, v := range args {
		if v > n {
			n = v
		}
	}
	return n
}

func showConcurrency1() {
	total := 0
	passengers := 30
	c := make(chan struct{})
	c1 := start(1, airportSecurityCheck, c)
	c2 := start(2, airportSecurityCheck, c)
	c3 := start(3, airportSecurityCheck, c)

	for i := 0; i < passengers; i++ {
		c <- struct{}{}
	}
	close(c)

	total = max(<-c1, <-c2, <-c3)
	println("total time cost:", total)
}

func AirCheck2(id int) int {
	println("start exec goroutine: ", id)
	total := 0

	total += idCheck(id)
	total += bodyCheck(id)
	total += xRayCheck(id)
	println("end exec goroutine: ", id)
	return total
}

func start2(id int, f func(int) int, signal chan struct{}) chan int {
	c := make(chan int)
	go func() {
		total := 0
		for {
			_, ok := <-signal
			if !ok {
				c <- total
				return
			}
			total += f(id)
		}
	}()
	return c
}

func max2(ints ...int) int {
	r := 0
	for i := 0; i < len(ints); i++ {
		if ints[i] > r {
			r = ints[i]
		}
	}
	return r
}

func showConcurrency2() {
	println("start ")
	signal := make(chan struct{})
	ch1 := start2(1, AirCheck2, signal)
	ch2 := start2(2, AirCheck2, signal)
	ch3 := start2(3, AirCheck2, signal)
	passagers := 30
	for i := 0; i < passagers; i++ {
		signal <- struct{}{}
	}
	close(signal)

	total := max2(<-ch1, <-ch2, <-ch3)
	println("total time cost:", total)
	println("end")
}

func main() {
	showConcurrency2()
}
