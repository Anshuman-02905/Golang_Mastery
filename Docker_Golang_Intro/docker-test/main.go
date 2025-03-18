package main

import (
	"fmt"
	"html"     //print HTML response
	"net/http" //Start the server
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { //ROOT route
		fmt.Fprintf(w, "HELLO, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	http.ListenAndServe(":8081", nil)

}
