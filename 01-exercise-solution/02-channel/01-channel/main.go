package main

import "fmt"

func main() {
	ch := make(chan int)

	go func(a, b int) {
		c := a + b
		ch <- c
	}(1, 2)

	// demo := func(a, b int) {
	// 	c := a + b
	// 	ch <- c
	// }

	// go demo(1,2)

	// get the value computed from goroutine
	r := <-ch
	fmt.Printf("computed value %v\n", r)
}
