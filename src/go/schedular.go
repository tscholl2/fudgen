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
    job := w.job
	w.job = ""
	w.time = 0
	return job
}
func (w *Worker) work(t int) {
    w.time = int(math.Max(float64(w.time - t),0))
	if (len(w.history) > 0) && w.history[len(w.history) - 1].job == w.job {
		w.history[len(w.history) - 1].time += t
	} else {
		w.history = append(w.history, WorkerHistory{w.job,t})
	}
}

type Group struct {
	workers []Worker
}
//worker group methods
func (W *Group) available() (a []Worker) {
	for _,w := range W.workers {
		if w.job == "" {
			a = append(a,w)
		}
	}
	return
}
func (W *Group) next() (finished []string) {
	var t int = W.workers[0].time
	for _,w := range W.workers {
		if w.job != "" {
			if w.time < t {t = w.time}
		}
	}
	for _,w := range W.workers {w.work(t)}
	for _,w := range W.workers {
		if w.time == 0 && w.job != "" {finished = append(finished,w.finish())}
	}
	return
}
func (W *Group) current(j string) bool {
	for _,w := range W.workers {
		if w.job == j {return true}
	}
	return false
}
func (W *Group) busy() bool {
	for _,w := range W.workers {
		if w.job != "" {return true}
	}
	return false
}
func (W *Group) schedule() (h []WorkerHistory) {
	for _,w := range W.workers {
		for _,m := range w.history {
			h = append(h,m)
		}
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
func task(V map[string]int, E [][]string, n int) []WorkerHistory {
	// initialize graph
	G := make(map[string]Vertex)
	for v,_ := range V {
		G[v] = Vertex{}
	}
	for _,e := range E {
		v := G[e[1]]
		v.in = append(v.in,e[0])
		G[e[1]] = v
		w := G[e[0]]
		w.out = append(w.out,e[1])
		G[e[0]] = w
	}
	
	fmt.Println("G")
	fmt.Println(G)
	
	// initialize workere
	var workers Group
	for i := 0; i < n; i++ {
		workers.workers = append(workers.workers,Worker{})
	}
	
	
	//run scheduling
	for len(G) > 0 {
		
		fmt.Println(len(G))
		
		fmt.Println("Starting...")
		
		var S []string
		for v,_ := range V {
			if len(G[v].in) == 0 && !workers.current(v) {S = append(S,v)}
		}
		fmt.Println("josb to do ")
		fmt.Println(S)
		
		W := workers.available()
		fmt.Println("avaiable:")
		fmt.Println(W)
		
		if (len(S) == 0) && !workers.busy() {panic("Cycle?!")}
		for i := 0; i < int(math.Min(float64(len(S)),float64(len(W)))); i++ {
			fmt.Println("assiging...")
			fmt.Println(W[i])
			W[i].assign(S[i],V[S[i]])
			fmt.Println(W[i])
		}
		
		fmt.Println("assigned:")
		fmt.Println(workers)
		
		F := workers.next()
		for _,v := range F {
			for _,w := range G[v].out {
				to := G[w].in
				to = splice(to, index(to,v))
			}
			delete(G,v)
		}
		
		fmt.Println("finished:")
		fmt.Println(workers)
		fmt.Println(G)
	}
	
	return workers.schedule()
}


func main() {
	V :=map[string]int{"a":10,"b":10}
	E := [][]string{{"a","b"}}
	fmt.Println(task(V,E,3))
}