package main

import "fmt"

type Vertex struct {
	in  []string
	out []string
}
func index(arr []string, s string) (i int) {
	for i < len(arr) {
		if arr[i] == s {
			return i
		}
		i++
	}
	return -1
}
func splice(arr []string, i int) []string {
	if i < 0 || i >= len(arr) {
		return arr
	}
	left := arr[:i]
	right := arr[i+1:]
	for _, s := range right {
		left = append(left, s)
	}
	return left
}

func main() {
	V := map[string]int{"a": 10, "b": 10}
	E := [][]string{{"a", "b"}}
	
	G := make(map[string]*Vertex)
	for v, _ := range V {
		G[v] = &Vertex{}
	}
	for _, e := range E {
		(*G[e[1]]).in = append((*G[e[1]]).in, e[0])
		(*G[e[0]]).out = append((*G[e[0]]).in, e[1])
	}
	
	//delete node
	(*G["b"]).in = splice((*G["b"]).in, index((*G["b"]).in, "b"))
	
	
	
	fmt.Println(G)
	fmt.Println(G["a"])
	fmt.Println(G["b"])
}