package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var (
		d = flag.Int("d", 0, "Specific process duration")
		t = flag.Int("t", 0, "Specific terminate duration")
	)
	flag.Parse()

	h, _ := os.Hostname()
	handler := func(w http.ResponseWriter, r *http.Request) {
		a := r.RemoteAddr
		q := r.URL.RequestURI()
		log.Printf("Connection start: %s, URI=%s", a, q)
		fmt.Fprintf(w, "Connection start: duration=%d, hostname=%s, requestURI=%s\n", *d, h, q)
		w.(http.Flusher).Flush()
		for i := *d; i > 0; i-- {
			log.Printf("something to do(%ds)...: %s, URI=%s\n", i, a, q)
			time.Sleep(time.Second * 1)
		}
		log.Printf("Connection end: %s, URI=%s\n", a, q)
		fmt.Fprintf(w, "Connection end: duration=%d, hostname=%s, requestURI=%s\n", *d, h, q)
		w.(http.Flusher).Flush()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	server := &http.Server{Addr: ":80", Handler: mux}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
		sig := <-sigs
		log.Printf("Terminate start: signal=%s\n", sig)
		for i := *t; i > 0; i-- {
			log.Printf("Terminate processing(%ds)...", i)
			time.Sleep(time.Second * 1)
		}

		log.Printf("Terminate end: signal=%s\n", sig)

		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-idleConnsClosed

}
