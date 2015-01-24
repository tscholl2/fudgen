package main

import (
	"fmt"
	"math"
)

//worker class
type WorkerHistory struct {
	job string
	time int
}
type Worker struct {
	job string
	time int
	history []WorkerHistory
}
func (w *Worker) assign(job string, time int) {
    w.job = job
	w.time = time
}
func (w *Worker) finish() string {
    job = w.job
	w.job = ''
	w.time = 0
	return job
}
func (w *Worker) work(t int) {
    w.time = math.Max(w.time - t,0)
	if len(history) > 0 && w.history[len(w.history) - 1].time += 1
		w.history[len(w.history) - 1].time += t
	else
		w.history = append(w.history, WorkerHistory{w.job,t})
}
func (w *Worker) toString() string {
    return p.Parents["mother"]
}


//worker group methods
func (W *[]Worker) available() (a []Worker) {
	for _,w := range W {
		if W.job == "" {
			a = append(a,W)
		}
	}
	return
}
func (W *[]Worker) next() (finished []string) {
	var t int = W[0].time
	for _,w := range W {
		if w.job != "" {
			if w.time < t {t = w.time}
		}
	}
	for _,w := range W {w.work(t)}
	for _,w := range W {
		if w.time == 0 && w.job != "" {finished = append(finished,w.finish())}
	}
	return
}
func (W *[]Worker) current(j string) bool {
	for _,w := range W {
		if w.job == j {return true}
	}
	return false
}
func (W *[]Worker) busy() bool {
	for _,w := range W {
		if w.job != "" {return true}
	}
	return false
}
func (W *[]Worker) schedule() (h []WorkerHistory) {
	for _,w := range W {
		h = append(h, w.history)
	}
	return
}
func (W *[]Worker) toString() (out string) {
	for _,w := range W {
		out = append(out, w.toString())
	}
	return
}

type Vertex struct {
	in []string
	out []string
}
func index(arr []string, s string) (i int) {
	for i < len(arr) {
		if arr[i] == s {return i}
		i++
	}
	return -1
}
func splice(arr []string, i int) []string {
	left := arr[:i]
	right := arr[i+1:]
	for _,s := range right {left = append(left,s)}
	return left
}
// task algorithm
func []WorkerHistory task(V map[string]int, E [][]string, n int) {
	// initialize graph
	var G := make(map[string]Vertex)
	for _,e := range E {
		G[e[1]].in = append(G[e[1]].in,e[0])
		G[e[0]].out = append(G[e[0]].out,e[1])
	}
	
	// initialize workere
	var workers [n]Worker
	
	//run scheduling
	for len(G) > 0 {
		var S []string
		for v,_ := range V {
			if len(G[v].in) == 0 && !workers.current(v) {S = append(S,v)}
		}
		W := workers.available()
		if len(S) && !workers.busy() {panic("Cycle?!")}
		for i := 0; i < math.Min(len(S),len(W)); i++ {
			W[i].assign(S[i],V[S[i]])
		}
		F := workers.next()
		for _,v := range F {
			for _,w := range G[v] {
				G[w].in = splice(G[w].in, index(G[w].in,v))
			}
			delete(G,v)
		}
	}
	return workers.schedule()
}


func main() {
	V :=map[string]int{
		"a":10,
		"b":10
	}
	E := [][]string{{"a","b"}}
	fmt.Println(V,E,1)
}