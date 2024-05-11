# Go concurrency exercises

Exercises and code walks included in this repository are part of Udemy course "concurrency in Go (Golang)".

https://www.udemy.com/course/concurrency-in-go-golang/?referralCode=5AE5A041D5793C048954

# Commands to run in cli:
  * run:
      - go run <file>.go
  * build and run:
      - go build
      - ./<folder> (see section 2: excercise 8, 3:22)



# Notes:

Section 2:

Go's concurrency ToolSet
  * goroutines
  * channels
  * select
  * sync package

#Goroutines:
   * we can think go routines as user space threads managed by go runtime(it part of the executable).
   * Goroutines are extremly light weight. Goroutines starts with 2KB of stack, which grows and shrinks as required
   * Low cpu overhead - the amount cpu instruction to create a goroutine is 3 instructions per function call, this helps us to create 100 - 1000 gouroutines in the same address space
   * channels are used for communication of data between goroutines. Sharing of memory can be avoided.
   * Context switching is much cheaper, as thread context switching is higher, as goroutines has less state to store.
   * Go runtime can be more selective in what is persisted for retriveal, how it is persisted, and when the persisting needs to occur.
   * Go runtimes creates OS threads, and Goroutines runs in the context of OS threads   => very imp
   * OS schedules the OS threads on to the processors, and go runtime schedules the go routines on the threads => very imp
   * Many goroutines execute in the context of single OS threads.

#Race Condition:
  * Race condition occurs when the order of execution is not guaranteed
  * Concurrent Programs does not execute in the order they are coded(section 2, 9 waitGroups) => very imp.
      - Order of execution of the goroutines is not guaranteed, see below
      - main will run first, then goroutine
      - goroutine will run first and then main thread
      - main thread runs then pauses for some reason, goroutine runs then main thread run
      - func main() {
          var data int
          // anonymous function
          go func() {
            data++
          }()

          if data == 0 {
            fmt.Printf("value is %v\n", data)
          }
        }
      - O/P can be:
        * nothing printed
        * value is 0
        * value is 1
  * WaitGroup from sync package
    - go works on a g=fork and join model
    - example
      * var wg sync.WaitGroup
        wg.Add(1)
        go func() {
          defer wg.Done()
          ...

        }()
        wg.Wait() // blocks the main goroutine

# Concurrent add
  - 01-exercise-solution
    - 01-goroutines
    - counting => using atomic.AddInt64 to add cumulative

# GoRoutines & Closures:
  - Goroutines execute within the same address space they are created in
  - They can directly modify variables in the enclosing lexical block
  - func inc() {
    var i int
    go func () {
      i++ // here go moves the variable I from stach to heap so that goroutine can aceessit
      fmt.Println(i)
    }()
    return
  }
  - Goroutine operates on the current value (excercise-closure02)
    func main() {
    var wg sync.WaitGroup

    // what is the output
    // fix the issue.

    for i := 1; i <= 3; i++ {
      wg.Add(1)
      go func(i int) {
        defer wg.Done()
        fmt.Println(i)
      }(i) // note imp, if we want go routine to operate on a specific value.
    }
    wg.Wait()
  }

# Go Scheduler(M:N scheduler):
  - Go scheduler is part of Go runtime, go runtime is part of Go executable.
  - Go scheduler runs on the user space
  - Go scheduler uses OS threads to schedule goroutines for execution
  - Goroutine runs in the context of OS threads
  - Go runtime create number of worker OS threads, equal to GOMAXPROCS
  - GOMAXPROCS - default value is number of processors on machine
  - Go scheduler distributes runnable goroutines over multiple worker OS threads
  - At any time, N goroutines could be scheduled on M OS threads that runs on at most GOMAXPROCS number of processors

  ## Asynchronous Preemption:
    - As of Go 1.14, Go scheduler implements asynchrounous preemtion
