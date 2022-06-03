package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

type result struct {
	value string
}

func first(servers ...*httptest.Server) (result, error) {
	c := make(chan result, len(servers))
	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 600)
	defer cancel()
	queryFunc := func(server *httptest.Server, id int) {
		url := server.URL

		r, err := http.NewRequest("get", url, nil)
		r = r.WithContext(ctx)
		// resp, err := http.Get(url)
		log.Printf("goroutine[%d] start request\n", id)
		resp, err := http.DefaultClient.Do(r)

		if err != nil {
			log.Printf("goroutine[%d] http get error: %s\n", id, err)
			return
		}

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("goroutine[%d] get response: %s\n", id, body)
		c <- result{
			value: string(body),
		}
	}
	for i, serv := range servers {
		go queryFunc(serv, i)
	}
	select {
	case result := <-c:
		return result, nil
	case <-time.After(time.Millisecond * 600):
		return result{}, errors.New("timeout")
	}
}

func fakeWeatherServer(name string, interval int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		w.Write([]byte(name + ":ok"))
	}))
}

func show() {
	result, err := first(
		fakeWeatherServer("open-weather-1", 200),
		fakeWeatherServer("open-weather-2", 800),
		fakeWeatherServer("open-weather-3", 200),
	)
	if err != nil {
		log.Println("invoke first error:", err)
		return
	}

	log.Println(result)
	time.Sleep(3000 * time.Millisecond)
}

func doGet(s *httptest.Server) (result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	_ = cancel
	// defer cancel()
	r, err := http.NewRequest("get", s.URL, nil)
	if err != nil {
		log.Println("build request err: ", err)
		return result{}, nil
	}
	r = r.WithContext(ctx)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Println("get response error: ", err)
		return result{}, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return result{value: string(b)}, nil
}

func showCtxTimeoutRequest() {
	r, err := doGet(fakeWeatherServer("open-weather-1", 1500))
	if err != nil {
		log.Println("query err: ", err)
		return
	}
	fmt.Println(r)
}

func fakeServer(name string, delay int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("receive request from: ", r.RemoteAddr)
		// select {}
		d := time.Duration(delay) * time.Millisecond
		time.Sleep(d)
		fmt.Fprintf(w, "Hi, %s delay: %d", r.RemoteAddr, d)
	}))
}

type result2 struct {
	value string
	index int
}

func fetch(timeout int, s ...*httptest.Server) (res []result2, err error) {
	c := make(chan result2)
	res = make([]result2, len(s))
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queryFunc := func(url string, id int) {
		// resp, err := http.Get(url)
		r, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("goroutine[%d] build request err: %v\n", id, err)
			return
		}
		r = r.WithContext(ctx)
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			log.Printf("goroutine[%d] query request err: %v\n", id, err)
			return
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("goroutine[%d] parse err: %v\n", id, err)
			return
		}
		log.Println("queryFunc-: ", string(b))
		c <- result2{
			value: string(b),
			index: id,
		}
	}
	// go queryFunc(s[0].URL)
	for i, v := range s {
		go queryFunc(v.URL, i)
	}

	for {
		select {
		case r := <-c:
			res[r.index].value = r.value
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			err = fmt.Errorf("timeout: %d", timeout)
			return
		}
	}
}

func show2(timeout int) {
	r, err := fetch(
		timeout,
		fakeServer("server1", 200),
		fakeServer("server2", 400),
		fakeServer("server3", 1000),
	)
	if err != nil {
		log.Println("fetch time's up: ", err)
		// return
	}
	fmt.Println("response: ", r)
	// time.Sleep(3000 * time.Millisecond)
}

func main() {
	// show2(300)
	// show()
	// showCtxTimeoutRequest()
}
