package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"html"
	"io"
	"log"
	"net/http"
)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
	//w.Write([]byte("hello!"))
}

func main() {
	http.HandleFunc("/hello", HelloServer)

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		w.Write([]byte("hello"))
	})

	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello m!"
	})
	m.Get("/test", func() string {
		return "Hello test!"
	})
	http.Handle("/", m)
	err := http.ListenAndServe(":8401", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		println("start Listening")
	}
}
