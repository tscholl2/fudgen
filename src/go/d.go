package main

import (
	"fmt"
	"math"
)

//worker class
type WorkerHistory struct {
	job  string
	time int
}
type Worker struct {
	job     string
	time    int
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
		if (*g).job == "" {
			H = append(H, g)
		}
	}
	return
}


func main() {
	arr := make(Group, 1)
	arr[0] = &Worker{}
	fmt.Println("nothing...")
	fmt.Println(arr)
	arr2 := arr.available()
	fmt.Println(arr2)
}
