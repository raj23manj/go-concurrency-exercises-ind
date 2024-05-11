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
    - As of Go 1.14, Go scheduler implements asynchrounous preemption
    - This prevents long running Goroutines from hogging onto CPU, that could block Goroutines
    - THe asynchrounous preemption is triggered based on a time condition. WHen a goroutine is running for more than 10ms, Go will try to preempt it.

  ## Goroutines states
    - when created it is in `Runnable` state waiting in run queue
    - Then moves to the `Executing` state once the goroutine is scheduled on the OS thread.
    - If the goroutine takes more time(10ms), it is preempted and placed back on the run queue
    - If the goroutine gets blocked, it moves to `Waiting` state
      - goroutine gets blocked due waiting on channel, syscall, waiting for a mutex lock, I/O, event operation
    - Once the I/O opereations(blocking operations) is completed, it is moved back to `Runnable` state

  ## M(OS threads):N(logical local processor) scheduler:
    - Goruntine has mechanism called `MN` scheduler
    - For a cpu core GoRuntime creates a OS thread `M`, GoRoutine also creates a logical processor `P` and associates that with the OS thread `M`
    - The logical processor holds the context for scheduling which can be seen as a local scheduler running on a thread
    - `G` Represents the a goroutine running on the OS thread
    - Each logical processor `P` has a local run queue, where runnable goroutines are queued
    - There is global run queue, once the logical processor exhausts it's queue, it pulls goroutines from global run queue
    - When new goroutines are created there added to the end of the global run queue
    ```
    core
      |                                 Global Run queue -> G10, G11
      |
      M
      |-> G1(running)
      P ->(queue) G2, G3, G4


      Preemption: if G1 takes time, them preemption happens
        - Here it is still the same OS thread, only goroutine is context switched.

      core
      |                                 Global Run queue -> G10, G11
      |
      M
      |-> G2(running)
      P ->(queue) G3, G4, G1

      ```

  ## Context switch due to synchronous system call:
    - what happens in general when synchrounous system call are made ? (like reading or writing to a file with sync flag set)
      - There will be adisk I/O to be performed, so synchrounous system call wait for I/O operation to be completed.
      - Due to this, the OS thread is moved out of CPU to waiting queue for I/O to complete, because of this we will not be able to schedule any goroutine on this thread for execution
      - Implications, Synchronous system call reduces parallelism
    - How does GO scheduler handles this scenario ?
      - Below let us assume, `G1` is blocked due to I/O operation
      - GOscheduler identifies that `G1` has caused `M1` to block, so it brings in a new OS thread either from the thread pool cache or it creates a new OS thread if there are no thread in the thread pool
      - Then Go scheduler will detach the logical processor `P` from `M1` and attach it to the new thread `M2`, `G1` goroutine is still attached to the thread/processor `M1`
      - Once the sync operation of `G1` is complete, Go scheduler moves it to the new thread `M2` on the logical processor `P` into the waiting queue
      - `M1` thrad is put to sleep, and is put back into the thread pool cache, to be reused in the future for same scenario
      ```
      Scenario 1:
      core
      |
      |
      M1                                              M2
      |-> G1(blocked due to I/O)
      P ->(queue) G2, G3, G4

      Scenario 2:
                                                      core
      M1 -> G1                                          |
                                                        M2
                                                        |
                                                        P -> G3, G4

      Scenario 3:

                                                       core
      M1                                                |
                                                        M2
                                                        | -> G3
                                                        P -> G4, G1

      ```
