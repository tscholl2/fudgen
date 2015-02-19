package recipes

import (
	//"fmt"
	"math"
)

//worker class
type WorkerHistory struct {
	job  int
	time int
}
type Worker struct {
	job     int
	time    int
	history []WorkerHistory
}

func (w *Worker) assign(job int, time int) {
	w.job = job
	w.time = time
}
func (w *Worker) finish() int {
	job := w.job
	w.job = 0
	w.time = 0
	return job
}
func (w *Worker) work(t int) {
	w.time = int(math.Max(float64(w.time-t), 0))
	if (len(w.history) > 0) && w.history[len(w.history)-1].job == w.job {
		w.history[len(w.history)-1].time += t
	} else {
		w.history = append(w.history, WorkerHistory{w.job, t})
	}
}

type Group []*Worker

//worker group methods
func (G *Group) available() (H Group) {
	for _, g := range *G {
		if (*g).job == 0 {
			H = append(H, g)
		}
	}
	return
}
func (G *Group) next() (finished []int) {
	var t int = (*(*G)[0]).time
	for _, g := range *G {
		if (*g).job != 0 {
			if (*g).time < t {
				t = (*g).time
			}
		}
	}
	for _, g := range *G {
		(*g).work(t)
	}
	for _, g := range *G {
		if (*g).time == 0 && (*g).job != 0 {
			finished = append(finished, (*g).finish())
		}
	}
	return
}
func (G *Group) current(j int) bool {
	for _, g := range *G {
		if (*g).job == j {
			return true
		}
	}
	return false
}
func (G *Group) busy() bool {
	for _, g := range *G {
		if (*g).job != 0 {
			return true
		}
	}
	return false
}
func (G *Group) schedule() (h [][]WorkerHistory) {
	for _, g := range *G {
		h = append(h, []WorkerHistory{})
		i := len(h) - 1
		for _, m := range (*g).history {
			h[i] = append(h[i], m)
		}
	}
	return
}

type Vertex struct {
	in  []int
	out []int
}

func index(arr []int, s int) (i int) {
	for i < len(arr) {
		if arr[i] == s {
			return i
		}
		i++
	}
	return -1
}
func splice(arr []int, i int) []int {
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

// task algorithm
func Schedule(V map[int]int, E [][]int, n int) [][]WorkerHistory {
	// initialize graph
	graph := make(map[int]*Vertex)
	for v, _ := range V {
		graph[v] = &Vertex{}
	}
	for _, e := range E {
		(*graph[e[1]]).in = append((*graph[e[1]]).in, e[0])
		(*graph[e[0]]).out = append((*graph[e[0]]).out, e[1])
	}

	// initialize workere
	workers := make(Group, n)
	for i := 0; i < n; i++ {
		workers[i] = &Worker{}
	}

	//run scheduling
	for len(graph) > 0 {

		var S []int
		for v, _ := range V {
			if e, ok := graph[v]; ok {
				if len(e.in) == 0 && !workers.current(v) {
					S = append(S, v)
				}
			}
		}

		W := workers.available()

		if (len(S) == 0) && !workers.busy() {
			panic("Cycle?!")
		}
		for i := 0; i < int(math.Min(float64(len(S)), float64(len(W)))); i++ {
			(*(W[i])).assign(S[i], V[S[i]])
		}

		F := workers.next()
		for _, v := range F {
			for _, w := range graph[v].out {
				(*graph[w]).in = splice((*graph[w]).in, index((*graph[w]).in, v))
			}
			delete(graph, v)
		}
	}

	return workers.schedule()
}

/*example
func main() {
	V := map[int]int{1: 10, 2: 10, 3: 10, 4: 10, 5: 10}
	E := [][]int{{1, 2}, {2, 4}, {2, 3}, {4, 5}, {3, 5}}
	fmt.Println(Schedule(V, E, 2))
}
*/
