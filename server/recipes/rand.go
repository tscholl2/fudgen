package recipes

import (
	"math/rand"
)

func randomTitle(list []string) (title string) {
	if len(list) == 0 {
		title = "Untitled"
		return
	}
	title = list[rand.Intn(len(list)-1)]
	return
}
