# Go concurrency exercises

Exercises and code walks included in this repository are part of Udemy course "concurrency in Go (Golang)".

https://www.udemy.com/course/concurrency-in-go-golang/?referralCode=5AE5A041D5793C048954

# Commands to run in cli: 


# Notes: 

Section 2:

Go's concurrency ToolSet
  * goroutines
  * channels
  * select
  * sync package

Goroutines: 
   * we can think go routines as user space threads managed by go runtime(it part of the executable).
   * Goroutines are extremly light weight. Goroutines starts with 2KB of stack, which grows and shrinks as required
   * Low cpu overhead - the amount cpu instruction to create a goroutine is 3 instructions per function call, this helps us to create 100 - 1000 gouroutines in the same address space
   * channels are used for communication of data between goroutines. Sharing of memory can be avoided.
   * Context switching is much cheaper, as thread context switching is higher, as goroutines has less state to store.
   * Go runtime can be more selective in what is persisted for retriveal, how it is persisted, and when the persisting needs to occur.
   * Go runtimes creates OS threads, and Goroutines runs in the context of OS threads   => very imp
   * OS schedules the OS threads on to the processors, and go runtime schedules the go routines on the threads => very imp
   * Many goroutines execute in the context of single OS threads.
