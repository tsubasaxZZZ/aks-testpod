package main

import (
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		a := r.RemoteAddr
		q := r.URL.RequestURI()
		log.Printf("Connection start: %s, URI=%s", a, q)
		fmt.Fprintf(w, "Connection suceed: duration=%d, hostname=%s, requestURI=%s\n", *d, h, q)
		for i := *d; i > 0; i-- {
			log.Printf("something to do(%ds)...: %s, URI=%s\n", i, a, q)
			time.Sleep(time.Second * 1)
		}
		log.Printf("Connection end: %s, URI=%s\n", a, q)
	})

	go func() {
		http.ListenAndServe(":80", nil)
	}()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-sigs

	log.Printf("Terminate start: signal=%s\n", sig)
	for i := *t; i > 0; i-- {
		log.Printf("Terminate processing(%ds)...", i)
		time.Sleep(time.Second * 1)
	}
	log.Printf("Terminate end: signal=%s\n", sig)

}
