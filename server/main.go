package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

var (
	writeTimeout      = flag.Duration("write", 0*time.Second, "write timeout for server")
	readTimeout       = flag.Duration("read", 0*time.Second, "read timeout for server")
	readHeaderTimeout = flag.Duration("readHeader", 0*time.Second, "read header timeout for server")
	idleTimeout       = flag.Duration("idle", 0*time.Second, "idle timeout for server")

	before = flag.Duration("before", 0*time.Second, "delay before response is written")
	after  = flag.Duration("after", 0*time.Second, "delay after response is written")
)

func main() {
	flag.Parse()

	// keep track of goroutines
	go func() {
		for range time.Tick(time.Second) {
			fmt.Println("--->", runtime.NumGoroutine())
		}
	}()

	// create a server with a single route
	serveMux := &http.ServeMux{}
	serveMux.Handle("/delay", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[server] /delay")

		// ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		// go func() {
		// 	<-ctx.Done()
		// 	fmt.Println("timeout expired!")
		// }()
		// defer cancel()

		// optionally wait before we write response
		time.Sleep(*before)

		// write a non-nil response so that the client must contractually read it;
		// note that a nil response is "automatically" read and the server can exit its goroutine
		w.Write([]byte("done"))

		// optionally wait after we write response
		time.Sleep(*after)

		fmt.Println("[server] done")
	}))

	srv := http.Server{
		Handler: serveMux,
		Addr:    ":8080",

		// WriteTimeout is the maximum duration before timing out
		// writes of the response.
		WriteTimeout: *writeTimeout,

		// ReadHeaderTimeout is the amount of time allowed to read
		// request headers.
		ReadHeaderTimeout: *readHeaderTimeout,

		// ReadTimeout is the maximum duration for reading the entire
		// request, including the body.
		//
		// Because ReadTimeout does not let Handlers make per-request
		// decisions on each request body's acceptable deadline or
		// upload rate, most users will prefer to use
		// ReadHeaderTimeout. It is valid to use them both.
		ReadTimeout: *readTimeout,

		// IdleTimeout is the maximum amount of time to wait for the
		// next request when keep-alives are enabled. If IdleTimeout
		// is zero, the value of ReadTimeout is used. If both are
		// zero, ReadHeaderTimeout is used.
		IdleTimeout: *idleTimeout,
	}

	go func() {
		fmt.Println(srv.ListenAndServe())
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	srv.Close()
}
