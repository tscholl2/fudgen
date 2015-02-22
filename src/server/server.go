package main

import (
	//"./recipes"
	//"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Web")
	})
	http.ListenAndServe(":8080", nil)
}
