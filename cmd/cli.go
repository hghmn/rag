package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
    
    "github.com/hghmn/rag"
)

func main() {
	type State struct {
		Counter int
		Init bool
	}

	init := &State{0, true}

	update := func(state *State, event string) *State {
		switch event {
		case "+":
			state = &State{state.Counter + 1, false}
		case "-":
			state = &State{state.Counter - 1, false}
		}

		return state
	}

	scanner := bufio.NewScanner(os.Stdin)
	view := func(state *State, dispatch func(string)) {
		if state.Init {
			fmt.Println("Welcome!")
			fmt.Println("Change the counter by entering '+' or '-', then hitting 'enter'")
		}

		// Print the current state
		fmt.Printf("Count: %d\n", state.Counter)

		// Scans a line from Stdin(Console)
		scanner.Scan()
		text := scanner.Text()
		if len(text) != 0 {
			dispatch(strings.ToUpper(text))
		} else {
			dispatch(/*rag.*/Exit)
		}
	}

	rag.Runtime(init, update, view)
}
