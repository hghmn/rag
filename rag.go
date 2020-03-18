package main

import (
    "fmt"
    "reflect"
)

// Built-in events
const (
	Exit = "EXIT"
)

// Runtime<S>(
//   init:   S | () => S,
//   update: (string, S) => S,
//   view:   (S, (string) => void) => void
// )
func Runtime(init, update, view interface{}) {
	// Validate the supplied program
	checkUpdate(init, update)
	checkView(init, view)

	state := init

	fmt.Printf("Running program: %T %v\n", state, state) // reflect.TypeOf(state))

	// Run
	mainloop:
	for {
		events := make([]string, 0)
		dispatch := func(event string) {
			events = append(events, event)
		}

		doView(view, state, dispatch)

		if len(events) > 0 {
			for _, ev := range events {
				switch ev {
				case Exit:
					break mainloop
				}

				state = doUpdate(update, state, ev)
			}
		}
	}

	fmt.Println("Program Exited")
}

func checkUpdate(state, updateFn interface{}) {
	updateFnValue := reflect.ValueOf(updateFn)
	updateFnType := updateFnValue.Type()
	
	// Checking whether the second argument is function or not.
	// And also checking whether its signature is func ({type A}) {type B}.  
	if updateFnType.Kind() != reflect.Func || updateFnType.NumIn() != 2 || updateFnType.NumOut() != 1 {
		panic("Invalid update function type.")
	}

	stateType := reflect.ValueOf(state).Type()
	if !stateType.ConvertibleTo(updateFnType.In(0)) {
		panic("Update function is not compatible with given state type.")
	}

	if updateFnType.In(1).Kind() != reflect.String {
		panic("Update function is not compatible with the event param type")
	}
}

func doUpdate(updateFn interface{}, state interface{}, event string) interface{} {
	fn := reflect.ValueOf(updateFn)
	if fn.Type().NumIn() != 2 { panic("Bad function") }

	res := fn.Call([]reflect.Value{
		reflect.ValueOf(state),
		reflect.ValueOf(event),
	})

	return res[0].Interface()
}

func checkView(state, viewFn interface{}) {
	viewFnValue := reflect.ValueOf(viewFn)
	viewFnType := viewFnValue.Type()

	// Checking whether the second argument is function or not.
	// And also checking whether its signature is func ({type A}) {type B}.  
	if viewFnType.Kind() != reflect.Func || viewFnType.NumIn() != 2 || viewFnType.NumOut() != 0 {
		panic("Invalid view function signature.")
	}

	stateType := reflect.ValueOf(state).Type()
	if !stateType.ConvertibleTo(viewFnType.In(0)) {
		panic("Update function is not compatible with given state type.")
	}

	dispatchType := viewFnType.In(1)
	// dispatchArgType := dispatchArg.Type()
	if dispatchType.Kind() != reflect.Func || dispatchType.NumIn() != 1 {
		panic("Dispatch function signature is invalid")
	}

	if dispatchType.In(0).Kind() != reflect.String {
		panic("Dispatch function signature is invalid with string param")
	}
}

func doView(view interface{}, state interface{}, dispatch func(string)) {
	method := reflect.ValueOf(view)
	method.Call([]reflect.Value{
		reflect.ValueOf(state),
		reflect.ValueOf(dispatch),
	})
}
