package recipes

import (
	"fmt"
	"math"
)

//worker class
type workerHistory struct {
	job  interface{}
	time int
}
type worker struct {
	job     interface{}
	time    int
	history []workerHistory
}

func (w *worker) assign(job interface{}, time int) {
	w.job = job
	w.time = time
}
func (w *worker) finish() interface{} {
	job := w.job
	w.job = nil
	w.time = 0
	return job
}
func (w *worker) work(t int) {
	w.time = int(math.Max(float64(w.time-t), 0))
	if (len(w.history) > 0) && w.history[len(w.history)-1].job == w.job {
		w.history[len(w.history)-1].time += t
	} else {
		w.history = append(w.history, workerHistory{w.job, t})
	}
}

type Group []*worker

//worker group methods
func (G *Group) available() (H Group) {
	for _, g := range *G {
		if (*g).job == nil {
			H = append(H, g)
		}
	}
	return
}
func (G *Group) next() (finished []interface{}) {
	var t int = 0
	for _, g := range *G {
		if (*g).job != nil {
			t = (*g).time
			break
		}
	}
	for _, g := range *G {
		if (*g).job != nil {
			if (*g).time < t {
				t = (*g).time
			}
		}
	}
	for _, g := range *G {
		(*g).work(t)
	}
	for _, g := range *G {
		if (*g).time == 0 && (*g).job != nil {
			finished = append(finished, (*g).finish())
		}
	}
	return
}
func (G *Group) current(j interface{}) bool {
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
func (G *Group) schedule() (h [][]workerHistory) {
	for _, g := range *G {
		h = append(h, []workerHistory{})
		i := len(h) - 1
		for _, m := range (*g).history {
			h[i] = append(h[i], m)
		}
	}
	return
}

type Vertex struct {
	in  []interface{}
	out []interface{}
}

func index(arr []interface{}, s interface{}) (i int) {
	for i < len(arr) {
		if arr[i] == s {
			return i
		}
		i++
	}
	return -1
}
func splice(arr []interface{}, i int) []interface{} {
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
//V is map of vertices ---> time it takes to complete
//E is array of tuples {v1,v2} representing arrow v1 --> v2
//n is number of workers available
func scheduleGraph(V map[interface{}]int, E [][]interface{}, n int) [][]workerHistory {
	// initialize graph
	graph := make(map[interface{}]*Vertex)
	for v, _ := range V {
		graph[v] = &Vertex{}
	}
	for _, e := range E {
		(*graph[e[1]]).in = append((*graph[e[1]]).in, e[0])
		(*graph[e[0]]).out = append((*graph[e[0]]).out, e[1])
	}

	// initialize workers
	workers := make(Group, n)
	for i := 0; i < n; i++ {
		workers[i] = &worker{}
	}

	//run scheduling
	for len(graph) > 0 {

		var S []interface{}
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

//returns array as follows:
//schedule[i] is what person
//i's iterinary, i.e. a list of things
//they need to do, in order.
//schedule[i][j] contains a tuple
//schedule[i][j][0] =  the id of the job
//they should be working on at this point
//schedule[i][j][1] = the time it should
//take them
func Schedule(R *Recipe, n int) (schedule [][][]int) {

	//build and run graph
	V := map[interface{}]int{}
	E := make([][]interface{}, 0)
	for i := 0; i < len(R.Steps); i++ {
		ptr := R.Steps[i]
		if !ptr.IsIngrediant() {
			op := ptr.(*Operation)
			V[op.Id] = int(op.Time.ToBasic().Amount)

			fmt.Println("found op: ", op)

			for j := 0; j < len(op.Requires); j++ {
				E = append(E, []interface{}{op.Requires[j], i})
			}
		} else {
			ing := ptr.(*Ingrediant)

			fmt.Println("found ing = ", ing)

			V[ing.Id] = 10
		}
	}

	fmt.Println("V = ", V)
	fmt.Println("E = ", E)

	H := scheduleGraph(V, E, n)

	//format output
	//initialize output
	schedule = make([][][]int, len(H))
	for i := 0; i < len(H); i++ {
		schedule[i] = make([][]int, len(H[i]))
		for j := 0; j < len(H[i]); j++ {
			var id int
			t := H[i][j].time
			if H[i][j].job == nil {
				id = -1
			} else {
				id = H[i][j].job.(int)
			}
			schedule[i][j] = []int{id, t}
		}
	}
	return
}

//example
// func main() {
// 	v1 := "A"
// 	v2 := "B"
// 	v3 := "C"
// 	v4 := "D"
// 	v5 := "E"
// 	V := map[interface{}]int{v1: 10, v2: 10, v3: 10, v4: 100, v5: 10}
// 	E := [][]interface{}{{v1, v2}, {v2, v3}, {v2, v4}, {v4, v5}, {v3, v5}}
// 	fmt.Println(scheduleGraph(V, E, 2))
// }
