package main

import "time"

const (
	idCheckTmCost   = 60
	bodyCheckTmCost = 120
	xRayCheckTmCost = 180
)

func idCheck(id string) int {
	time.Sleep(time.Millisecond * time.Duration(idCheckTmCost))
	print("\tgoroutine-", id, "-idCheck: idCheck ok\n")
	return idCheckTmCost
}

func bodyCheck(id string) int {
	time.Sleep(time.Millisecond * time.Duration(bodyCheckTmCost))
	print("\tgoroutine-", id, "-bodyCheck: bodyCheck ok\n")
	return bodyCheckTmCost
}

func xRayCheck(id string) int {
	time.Sleep(time.Millisecond * time.Duration(xRayCheckTmCost))
	print("\tgoroutine-", id, "-xRayCheck: xRayCheck ok\n")
	return xRayCheckTmCost
}

func start(id string, f func(string) int, next chan<- struct{}) (chan<- struct{}, chan<- struct{}, <-chan int) {
	queue := make(chan struct{}, 10)
	quit := make(chan struct{})
	result := make(chan int)

	go func() {
		total := 0
		for {
			select {
			case <-quit:
				result <- total
				return
			case v := <-queue:
				total += f(id)
				if next != nil {
					next <- v
				}
			}
		}
	}()
	return queue, quit, result
}

func newAirportSecurityCheckChannel(id string, queue <-chan struct{}) {
	go func(id string) {
		print("goroutine-", id, ": airportSecurityCheckChannel is ready...\n")
		// start xRayCheck routine
		queue3, quit3, result3 := start(id, xRayCheck, nil)

		// start bodyCheck routine
		queue2, quit2, result2 := start(id, bodyCheck, queue3)

		// start idCheck routine
		queue1, quit1, result1 := start(id, idCheck, queue2)

		for {
			select {
			case v, ok := <-queue:
				if !ok {
					close(quit1)
					close(quit2)
					close(quit3)
					total := max(<-result1, <-result2, <-result3)
					print("goroutine-", id, ": airportSecurityCheckChannel time cost:", total, "\n")
					print("goroutine-", id, ": airportSecurityCheckChannel closed\n")
					return
				}
				queue1 <- v
			}
		}
	}(id)
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

func idCheck3(id string) int {
	time.Sleep(time.Millisecond * time.Duration(idCheckTmCost))
	print("\tgoroutine-", id, "-idCheck: idCheck ok\n")
	return idCheckTmCost
}

func bodyCheck3(id string) int {
	time.Sleep(time.Millisecond * time.Duration(bodyCheckTmCost))
	print("\tgoroutine-", id, "-bodyCheck: bodyCheck ok\n")
	return bodyCheckTmCost
}

func xRayCheck3(id string) int {
	time.Sleep(time.Millisecond * time.Duration(xRayCheckTmCost))
	print("\tgoroutine-", id, "-xRayCheck: xRayCheck ok\n")
	return xRayCheckTmCost
}

func showConcurrency3() {
	passengers := 30
	queue := make(chan struct{}, 30)
	newAirportSecurityCheckChannel("channel1", queue)
	newAirportSecurityCheckChannel("channel2", queue)
	newAirportSecurityCheckChannel("channel3", queue)

	time.Sleep(5 * time.Second) // 保证上述三个goroutine都已经处于ready状态
	for i := 0; i < passengers; i++ {
		queue <- struct{}{}
	}
	time.Sleep(5 * time.Second)
	close(queue) // 为了打印各通道的处理时长
	time.Sleep(1000 * time.Second)
}

func start3(id string, f func(string) int, next chan struct{}) (sig chan struct{}, quit chan struct{}, data chan int) {
	sig = make(chan struct{}, 10)
	quit = make(chan struct{})
	data = make(chan int)

	go func() {
		total := 0
		for {
			select {
			case <-quit:
				data <- total
				return
			case v := <-sig:
				total += f(id)
				if next != nil {
					next <- v
				}
			}
		}
	}()
	return sig, quit, data
}

func newAirportSecurityCheckChannel2(s string, sig chan struct{}) {
	println("newAirportSecurityCheckChannel2 ready")
	go func() {
		sig3, quit1, data1 := start3(s, xRayCheck3, nil)
		sig2, quit2, data2 := start3(s, bodyCheck3, sig3)
		sig1, quit3, data3 := start3(s, idCheck3, sig2)
		for {
			v, ok := <-sig
			if !ok {
				close(quit1)
				close(quit2)
				close(quit3)
				total := max(<-data1, <-data2, <-data3)
				println("gouroutine: ", s, ", total: ", total)
				println("gouroutine: ", s, ", closed")
				return
			}
			sig1 <- v
		}
	}()
}

func showConcurrency33() {
	passagers := 30
	sig := make(chan struct{}, 30)
	newAirportSecurityCheckChannel2("channel1", sig)
	newAirportSecurityCheckChannel2("channel2", sig)
	newAirportSecurityCheckChannel2("channel3", sig)

	time.Sleep(5 * time.Second)
	for i := 0; i < passagers; i++ {
		sig <- struct{}{}
	}
	time.Sleep(5 * time.Second)
	println("main done1")
	close(sig)
	time.Sleep(1000 * time.Second)
	println("main done")
}

func main() {
	// showConcurrency3()
	showConcurrency33()
}
