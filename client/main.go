package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

var (
	timeout = flag.Duration("timeout", 0*time.Second, "http client timeout")
	bad     = flag.Bool("bad", false, "makes client bad by not closing the server's response")

	before = flag.Duration("before", 0*time.Second, "before")
	after  = flag.Duration("after", 0*time.Second, "after")
)

func main() {
	flag.Parse()

	// keep track of goroutines
	go func() {
		for range time.Tick(time.Second) {
			fmt.Println("--->", runtime.NumGoroutine())
		}
	}()

	cli := http.Client{
		Timeout: *timeout,
	}
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/delay", nil)

	// make several requests to the server
	for i := 0; i < 10; i++ {
		go func() {
			resp, err := cli.Do(req)
			if err != nil {
				fmt.Println("[client]", err)
				return
			}

			time.Sleep(*before)

			if !*bad {
				fmt.Println("[client] good client - reading and closing")
				ioutil.ReadAll(resp.Body)
				resp.Body.Close()
			}

			time.Sleep(*after)

			fmt.Println("[client] returning")
		}()
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
