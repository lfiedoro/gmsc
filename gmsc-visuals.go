package gmsc

import (
	gc "code.google.com/p/goncurses"
	"fmt"
	"log"
	"strconv"
)

// find longest string in array
func longestSize(arr []string) int {
	var sz int
	for _, e := range arr {
		if sz < len(e) {
			sz = len(e)
		}
	}

	return sz
}

// print chooseable options in the center of the window
func Present(list *[]string, win *gc.Window) {
	win.Erase()
	win.NoutRefresh()

	row, col := win.MaxYX()
	row = row / 2
	row = row - len(*list)/2
	col = col / 2
	col = col - longestSize(*list)

	for i, v := range *list {
		entry := fmt.Sprintf("%v - %v", i, v)
		win.MovePrint(row+i, col, entry)
	}

	win.NoutRefresh()
	gc.Update()
}

// changle key presses
func Choose(list *[]string, win *gc.Window) (element string, err error) {
	// If there is only one element, choose it by default
	if len(*list) == 1 {
		return (*list)[0], nil
	}

	for {
		key := win.GetChar()
		c := gc.KeyString(key)
		if c == "q" || c == "Q" {
			log.Fatal("Quit!")
		}

		i, err := strconv.Atoi(c)
		if err != nil {
			log.Println(err)
		} else if i >= 0 && i < len(*list) {
			element = (*list)[i]
			return element, nil
		} else {
			log.Println("Index out of bounds %d", key)
		}
	}

	return "", err
}
