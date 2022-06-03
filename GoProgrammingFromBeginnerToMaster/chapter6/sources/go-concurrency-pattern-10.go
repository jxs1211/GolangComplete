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
	queryFunc := func(server *httptest.Server) {
		url := server.URL
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("http get error: %s\n", err)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		c <- result{
			value: string(body),
		}
	}
	for _, serv := range servers {
		go queryFunc(serv)
	}
	return <-c, nil
}

func fakeWeatherServer(name string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s receive a http request\n", name)
		time.Sleep(1 * time.Second)
		w.Write([]byte(name + ":ok"))
	}))
}

func showConcurrency10() {
	result, err := first(
		fakeWeatherServer("open-weather-1"),
		fakeWeatherServer("open-weather-2"),
		fakeWeatherServer("open-weather-3"))
	if err != nil {
		log.Println("invoke first error:", err)
		return
	}

	log.Println(result)
}

type result101 struct {
	value string
}

func first101(servers ...*httptest.Server) (result101, error) {
	c := make(chan result101)
	queryFunc := func(s *httptest.Server) {
		fmt.Println("query start")
		defer fmt.Println("query done")
		resp, err := http.Get(s.URL)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		c <- result101{
			value: string(b),
		}
	}
	for _, s := range servers {
		go queryFunc(s)
	}
	return <-c, nil
}

func first101WithCancel(servers ...*httptest.Server) (result101, error) {
	c := make(chan result101, len(servers))
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		fmt.Println("cancel")
	}()
	queryFunc := func(s *httptest.Server, id int) {
		defer fmt.Printf("goroutine-%d query end\n", id)
		req, err := http.NewRequestWithContext(ctx, "GET", s.URL, nil)
		fmt.Printf("goroutine-%d query start\n", id)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		c <- result101{
			value: string(b),
		}
	}
	for i, s := range servers {
		go queryFunc(s, i)
	}
	select {
	case v := <-c:
		return v, nil
	case <-time.After(1500 * time.Millisecond):
		return result101{}, errors.New("timeout")
	}
	// return <-c, nil
}

func first101WithTimeout(servers ...*httptest.Server) (result101, error) {
	c := make(chan result101, len(servers))
	// defer func() {
	// cancel()
	// fmt.Println("cancel")
	// }()
	queryFunc := func(s *httptest.Server, id int) {
		defer fmt.Printf("goroutine-%d query end\n", id)
		fmt.Printf("goroutine-%d query start\n", id)
		// req, err := http.NewRequestWithContext(ctx, "GET", s.URL, nil)

		// resp, err := http.DefaultClient.Do(req)
		resp, err := http.Get(s.URL)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		c <- result101{
			value: string(b),
		}
	}
	for i, s := range servers {
		go queryFunc(s, i)
	}
	select {
	case r := <-c:
		return r, nil
	case <-time.After(2000 * time.Millisecond):
		return result101{}, errors.New("timeout")
	}
}

func fakeServer101(name string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s receive a http request\n", name)
		time.Sleep(time.Second * 3)
		fmt.Fprintf(w, "%s: ok", name)
	}))
}

func fakeServerWithDiffRespTime(name string, duration int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s receive a http request\n", name)
		time.Sleep(time.Second * time.Duration(duration))
		fmt.Fprintf(w, "%s: ok", name)
	}))
}

func showConcurrency101() {
	result101, err := first101(
		fakeServer101("server1"),
		fakeServer101("server2"),
		fakeServer101("server3"),
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result101)
}

func showConcurrency101WithCancel() {
	result101, err := first101WithCancel(
		fakeServerWithDiffRespTime("server1", 1),
		fakeServerWithDiffRespTime("server2", 2),
		fakeServerWithDiffRespTime("server3", 3),
	)
	if err != nil {
		fmt.Println(err)
		time.Sleep(2 * time.Second)
		return
	}
	fmt.Println(result101)
	time.Sleep(2 * time.Second)
}

func showConcurrency101WithTimeout() {
	result101, err := first101WithTimeout(
		fakeServerWithDiffRespTime("server1", 1),
		fakeServerWithDiffRespTime("server2", 2),
		fakeServerWithDiffRespTime("server3", 4),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result101)
}

func main() {
	// showConcurrency10()
	// showConcurrency101()
	// showConcurrency101WithTimeout()
	showConcurrency101WithCancel()
}
