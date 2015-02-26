package main

import (
	"./recipes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type message struct {
	Error  error          `json:error`
	Recipe recipes.Recipe `json:recipe`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Web")
	})
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		m := getRandomRecipe()
		bytes, err := json.Marshal(m)
		fmt.Println(err)
		if err != nil {
			bytes = []byte("{error:'I dont even know'}")
		}
		w.Write(bytes)
	})
	http.ListenAndServe(":8080", nil)
}

func getRandomRecipe() (m message) {
	b, err := ioutil.ReadFile("../recipes/r3.yml")
	if err != nil {
		panic(err)
	}
	r, err := recipes.ParseYaml(string(b))
	m.Error = err
	m.Recipe = r
	return
}
