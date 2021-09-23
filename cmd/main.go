package main

import (
	"flag"
	"fmt"
	"net/http"

	iorate "github.com/myml/io-rate"
)

var addr = ":8080"
var dir = "./"
var speed = 1024
var global = true

func main() {
	flag.IntVar(&speed, "s", speed, "speed B/s")
	flag.StringVar(&dir, "d", dir, "dir")
	flag.StringVar(&addr, "l", addr, "http listen addr")
	flag.BoolVar(&global, "g", global, "global Speed Limit")
	flag.Parse()

	fmt.Println("dir:", dir)
	fmt.Println("speed:", speed)
	fmt.Println("global limit:", global)
	fmt.Println("listen:", addr)
	h := http.FileServer(http.Dir(dir))
	if global {
		limiter := iorate.NewLimiter(speed * 1024)
		http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
			rw = iorate.NewResponseWriteByLimit(rw, limiter)
			h.ServeHTTP(rw, r)
		})
	} else {
		http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
			rw = iorate.NewResponseWrite(rw, speed*1024)
			h.ServeHTTP(rw, r)
		})
	}
	http.ListenAndServe(addr, nil)
}
