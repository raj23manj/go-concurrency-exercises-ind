# GOPATH
  * go env GOPATH
  * To edit, vi ~/.bash_profile (normal) / vi ~/.zshrc
  * use source ~/.bash_profile (try for zhrc also) to load in current terminal window
  * Config from tech school
    ```
      export PATH=$PATH:$(go env GOPATH)/bin
    ```
  * Config exisiting:
    ```
      export GOROOT=/usr/local/go
      export GOPATH=/Users/rajeshmanjunath/Projects/go-lang/goworkspace
      export PATH=$GOPATH/bin:$PATH
      export GO111MODULE=auto
      export GOBIN=/Users/rajeshmanjunath/Projects/go-lang/goworkspace/bin
      export PATH=$GOBIN/bin:$PATH
    ```
  * Used docker desktop & psql image for local development, techschool

# Go Tools
  ## online schema generation
    * https://dbdiagram.io/home => for designing db with schema, use export to generate sql code postgres, msql etc
    * https://dbml.dbdiagram.io/docs/

  ## golang-migrate(migration tools) like flyway
    -  brew install golang-migrate


# Go Projects
  ## Tech school
    - https://github.com/techschool/simplebank
    - mydevelopment => https://github.com/raj23manj/techshool-go

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

  ## Context switching due to Asynchronous system call(like network system call, http api call)
    - what happens in general when asynchrounous system call are made ?
      - Asynchronous call happens when the File descriptor that is used to do I/O operation is set to non-blocking mode
      - If the file descriptor that is not ready, (if scoket buffer is empty, and we are trying to read from it  or the socket buffer is full and we are trying to write to it, then the read/write operations does not block but returns error) for I/O operation system call does not block, but returns error. The application has to retry the operation at a later poin in time.
      - Asynchronous IO increases the application complexity.
      - Application havs to create a event loop. (Setup event loops using callbacks functions)
    - How does go handles this situation ?
      - netpoller
        - Go uses netpoller, there is abstraction in Syscall package.
        - Syscall uses netpoller to convert async system call to blocking system call
        - when a goroutine makes a async system call, and the file descriptor is not ready, then the Go scheduler uses netpoller OS thread to part that goroutine
        - Netpoller uses the interface provided by OS (Kqueue(mac), epoll(linux), iocp(win)) to poll on the file descriptor
        - Once the netpoller gets the notification from the OS, when the file descriptor is ready for I/O operation
        - Netpoller notifies the goroutine to retry the I/O operation
        - This way, Complexity os managing async system call is moved from application to GO runtime, which manages it efficiently.
        - For async calls, noextra threads are used like sync calls, but uses exisiting netpoller thread instead
        ```

         Scenario 1:
         core
          |
          M1                                              Netpoller
          |-> G1(makes I/O call)
          P ->(queue) G2, G3, G4

         Scenario 2:
         core
          |
          M1                                              Netpoller (uses interfacr provided from OS to poll on the file descriptor,
                                                          if it receives a notification that a file descriptor is available,
                                                          it looks on it its datastructure if any goroutine needs it to execute)
          |-> G2                                              |
          P ->(queue) G3, G4                                  G1


        Scenario 3:
         core
          |
          M1                                              Netpoller
          |-> G2
          P ->(queue) G3, G4,G1

        ```

  ## Work stealing
    - work stealing helps to balance the goroutines across all the logical processors
    - work gets distributed and gets done more efficiently
    - work stealing rule:
      - if there is no goroutine in local run queue
        * Try to steal from other logical processors, steal half of the goroutines
        * if not found, check the global runnable queue for a G
        * if not found check the netpoller

# Channels
  - channels used to communicate data between the goroutines
  - channels can also help in synchronising the execution of the goroutines(one goroutine can let know the other goroutine in what stage they are in and synchronise their execution)
  - channels are typed, they are ued to send and receive a particular type
  - they are thread safe, channels can be used to send and receive data concurrently by multiple goroutines
  - example:

        var ch chan T
        ch = make(chan T)
        or
        ch := make(chan T) => builtin function
  - pointer operator is used for sending and receiving value from channel
  - the arrow indicates the direction of data flow
    - Send:
        * ch <- v
    - Receive:
        * v = <- ch
  - channels are blocking,
      - sender Goroutines will wait for a recevier to be ready
        * ch <- v
      - receiver Goroutines will wait for a value to be received
        * <- v
  - It is the responsiblity of the channel to make the goroutine runnable again once it has data to send/receve
  - to close a channel `close(ch)`, once sender is close, the receiver is unblocked and proceeds with its computation
    * value, ok = <- ch
      - ok = true, value generated by a wriete. If element is present
      - ok = false, value generated by a close

  ## Range, Buffered channels, channle direction
    - Range
      - for value := range ch {
        ...
      }
      - Iterate over values received from a channel
      - Loop automatically breaks, when a channel is closed
      - range does not return the second boolena value

    - Unbuffered channels
      - There is no buffer between the sender Goroutinr and the Receiver Goroutine
      - ch := make(chan Type)
      - Since there is no buffer the sender Goroutine will block until there is a receiver to receive a value and the receiver Goroutine will block until there is a sender to send a value
      - G1 ------- G2

    - Buffered channels
      - in this there is a buffer between the receiver and the sender
      - channels are given a capacity
      - bufferen channels are in-memory FIFO queue
      - It is asynchronous
      - The sender can keep sending the values without blocking until the receiver being ready, until the buffer gets full. If the buffer gets full then the sender will block
      - THe receiver can keep receiving the values without blocking until the buffer gets gets empty, once empty it will block.
      - ch := make(chan Type, capacity)
      - G1------| | | |-------G2

    - channel direction
      - when using channels as function parametes, you can specify if a channel is meant to only send or receive values.
      - This specificity increases the type-safety of the program
        - func pong(in <-chan string, out chan<- string) {}
        - example:
        ```
          package main

          import "fmt"

          func ping(out chan<- string) {
              // send message on ch1
              out <- "ping"
            }

            func pong(in <-chan string, out chan<- string) {
              // recv message on ch1
              msg := <-in
              msg = msg + " pong"
              // send it on ch2
              out <- msg
            }

            func main() {
              // create ch1 and ch2
              ch1 := make(chan string)
              ch2 := make(chan string)

              // spine goroutine ping and pong
              go ping(ch1)
              go pong(ch1, ch2)

              // recv message on ch2
              msg := <-ch2

              fmt.Println(msg)
            }
          ```
    - channel ownership
      - Default value - CHannels
        - Default value for channels is `nil`
          - var ch chan interface{}
        - reading/writing to a nil channel will block forever
          ```
            var ch channel interface{}
            <-ch
            ch <- struct{}{}
          ```
        - closing on nil channel will panic
          ```
            var ch channel interface{}
            close(ch)
          ```
        - ensure channels are initialized first
          ```
            ch := make(chan int)
          ```
        - How we use the channels is important, to avoid deadlocks and panics
        - we can follow go best practices
          - owner of the channel: The goroutine that creates the channel will be the Go routine that writes to a channel and is also responsible for closing the channel
          - The Go routine that utilizes the channel, will only read from from the channel
        - Ownership of channels avoids
          - Deadlocking by writing to a nil channel
          - closing a nil channel, leads to panic
          - writing to a closed channel, leads to panic
          - closing a channel more than once, leads to panic
          - example
          ```
            package main

            import "fmt"

            func main() {
              // create channel owner goroutine which return channel and
              // writes data into channel and
              // closes the channel when done.

              owner := func() chan int {
                ch := make(chan int)
                go func() {
                  defer close(ch)
                  for i := 0; i < 6; i++ {
                    ch <- i
                  }
                }()
                return ch
              }

              consumer := func(ch <-chan int) {
                // read values from channel
                for v := range ch {
                  fmt.Printf("Received: %d\n", v)
                }
                fmt.Println("Done receiving!")
              }

              ch := owner()
              consumer(ch)
            }

          ```

    ### Channels Deep Dive
      - Buffered channels
        - A hchan struct has below, hchan struct represents a channel
        ```
            |buf|lock|sendx|recvx|...|
              | -> queue
            | | | |
        ```
        - A channel struct has a lock(mutex), buf(buffer). A goroutine needs to acuire the lock on the channel first
        - The sender(G1) first acquires the lock on the hchan struct, then adds elemenst and increments sendx and releases the lock and proceeds with its other computation.
        - THe receiver(G2) if wants to receive for the same channel, it acquires the lock first, dequeues(copies the value from queue into it's local variable) the value and increments recvx and releases the lock
        - Channel Idoms:
          - There is no memory share between goroutines
          - Goroutines copy elements into from hchan
          - hchan is protected by mutex lock
          - "Do not communicate by sharing memory, instead share memory by communicating"
        - Buffer full scenario:
          - example:
          ```
            `G1`

            make(chan int, 3)
            for _, v := range []int{1,2,3,4} {
              ch <- v
            }

            hchan:
            * scenario 1
            ch <- 4 => waits

            |buf|lock|sendx|recvx|senq|recvq|...|
              | -> queue
            |3|2|1|  => channel buffer full, G1 has to wait for the receiver
           ```
          - How does G1 wait happen for the above scenatio?
            - G1 creates a `sudog` struct say `G`
            - This struct will hold the reference of the executing Goroutine `G1`, the next value to sent will be saved in the `elem` field
            - This structure is enqueded in the sendq list
            - example:
            ```
                                  G1   4
                                  |    |
                                | G | elem | | => sudog
                                    |
            |buf|lock|sendx|recvx|senq|recvq|...|
              | -> queue
            |3|2|1|  => full
            ```
            - Now G1 calls `gopark()`, the scheduler will move `G1` out of execution on the OS thread, and other Go routine from the local run queue gets scheduled to run on the OS thread say `G2`
            - Now G2 comes, to dequeue the element from the channel, it dequeues the element and copies it into its local variable, and pops the
              waiting `G1` on the `sendq` and enqueues the value(4) saved in the `elem` field (into the buffer), once this is done, G2 sets the sate of G1 runnable, this is done G@ calling goready function for G1, the G1 gets moved to the runnable state and gets added to the local run queue for execution
            - example:
            ```
            `G2`
                                  G1   4
                                  |    |
                                | G | elem | | => sudog

            |buf|lock|sendx|recvx|senq|recvq|...|
              | -> queue
            |4|3|2|  => deques 1 and pops G1 from sendq and enques the value 4
            ```
        - Buffer Empty Scenario:
          - Scenario: G2 tries to recv on empty channel
            ```
              ch := make(chan int, 3)

              // G1 - goroutine
              func G1(ch chan <- int) {
                for _, v := range []int{1, 2, 3, 4} {
                  ch <- v
                }
              }

              // G2 - goroutine
              func G2(ch <-chan int) {
                for v := range ch {
                  fmt.Println(v)
                }
              }
            ```
          - Explantion:
            - G2 creates itself a sudog struct `G` pointing to G2, with elem poniting to `v` and enqueues it in the `recvq` of hchan struct of the channel
            - G2 calls upon the scheduler using `gopark()` and the scheduler will move G2 out of the OS thread and does a context switching, and scheduler the next go routinr from the local run queue to the OS thread
            - Now `G1` comes and tries to send a value on the channel, first it checks if there are any goroutines waiting in the `recvq` of the channel, and it finds G2
            - Now G1 copies the value directly into the variable of G2 stack. `G1 is directly accessing the stack of v2, and writing into the variable of g2 stack`. This is the only scenario where one goroutine access the stack of another goroutine.
            - Then G1 pops G2 from the recvq and puts G2 into local run queue by making it runnable state, by calling the goready(G2)
            ```
             `G2`
                                        G2   v
                                        |    |
                                      | G | elem | | => sudog
                                        |
            |buf|lock|sendx|recvx|senq|recvq|...|
              | -> queue
            | | | |  => empty


            `G1`
            ch <- 1

                 G1
                 | (copies 1 to v)
            G2|v(stack)|||

            ```

        - Unbuffered channel
          - Send on unbuffered channel
            - when sender goroutine wants to send values
            - if there is corresponding receiver waiting in recvq
            - sender write the value directly into receiver goroutine stack variable
            - sender goroutine puts the receiver goroutine back to runnable state
            - if there is no receiver goroutine in the sendq
            - sender gets parked into sendq
            - data is saved in elem field in sudog struct
            - when receiver comes along, it copies the data
            - puts the sender back to runnable state
          - Receive on unbuffered channel
            - Receiver goroutine wants to receive value
            - if it find a goroutine in wainting in sendq
            - receiver copies the value in elem field to its variable
            - puts the sender goroutine to runnable state
            - if there is no sender goroutine in send q
            - Receiver gets parked into recvq
            - Reference to variable is saved in elem field in sudog struct
            - When sender comes along, it copies the data directly to the receiver stack variable.
            - Puts the receiver to runnable state

# Select
  * Scenario:
    - G1 wants to receive result of computatuion from G2 and G3
      ```
         G1 <---- G2
          ^
          |
          |
          G3
      ```
    - in what order are we going to receive results ?
      * from G2 the G3 or G3 then G2
    - What if G# was much faster than G2 in one instance, and G2 is fater than G3 in another ?
    - Can we do the operation on channel which ever is ready, without considering the order of the results?
    - select statement is like a switch
      ```
        select {
          case <-ch1:
               // block of statements
          case <-ch2:
               // block of statements
          case <-ch3 <- struct{}{}:
               // block of statements
        }
      ```
    - In select, all channel operations are considered simultaneously
    - select waits until some case is ready to proceed, if none of them are ready, then the entire select statement is going to block, until some case is ready
    - when one of the channel is ready, that operation will proceed immediately
    - if multiple channels are ready, main goroutine will pick one of them and execute at random
    - Select is very helpful in implementing
      * Timeouts
        ```
          Timeout waiting on channel

          select {
            case v := <-ch:
              fmt.println(v)
            case <- timeAfter(3 * time.Second):
              fmt.Println("timeout")
          }

          - select waits until there is a event on channel ch or until timeout is reached
          - The time.After function takes in a time.Duration argument and returns a channel that will send the current time after the duration you provide it
        ```

      * Non-blocking communications
        ```
          select {
            case v := <-ch:
              fmt.println(v)
            default:
              fmt.Println("no msg rcvd")
          }

          - send or receive on a channel, but avoid blocking if channel is not ready
          - Default allows you to exit select block without blocking

          - non blocking:
            // if there is no value on channel, do not block.
            for i := 0; i < 2; i++ {
              select {
              case m := <-ch:
                fmt.Println(m)
              default:
                fmt.Println("no message received")
              }
              // Do some processing..
              fmt.Println("processing..")
              time.Sleep(1500 * time.Millisecond)
            }

          - see exmaple code
            for i := 0; i < 2; i++ {
              select {
              case msg1 := <-ch1:
                fmt.Println(msg1)
              case msg2 := <-ch2:
                fmt.Println(msg2)
              }
            }
        ```

      * Empty select
        - select {}
        - this will block for ever
        - select on nil channel will block forever.
        ```
          var ch chan string
          select {
            case v := <-ch:
            case ch <- v:>>
          }
        ```

# Sync Package(apart from wait group)
  - Mutex
    * when to use channels and whe to use mutex
      * channels: used to pass inbetween goroutines
        * passing copy of data
        * distributing unit of work
        * Communicating asynchrounous results
      * Mutex: when data is large and can't be passed using channels
        * caches
        * states
        * access to these data must be cocnurrent safe, so that only one goroutine has access to the data at a time
    * Mutex is used to protect shared resources
    * sync.Mutex - provides exclusive access to shared resources
    * section between lock() and unlock() is called as `critical section`
    * example
      ```
        mu.Lock()
        balance +=amount
        mu.Unlock()

        mu.Lock()
        defer mu.Unlock()
        balance -= amount
      ```
    * if goroutine only wants to read or write then we can use sync.RwMutex. Allows multiple readers, but writers get exclusive lock
      ```
        // write exclusive lock
        mu.Lock()
        balance +=amount
        mu.Unlock()

        // read lock, multiple goroutines can access
        mu.RLock()
        defer mu.RUnlock()
        return balance
      ```

  - Atomic
    * Low level atomic operations on memory
    * Lockless operation
    * Used for atomic operations on counters
    * example
    ```
      atomic.AddUint64(&ops,1)
      value := atomic.LoadUint64(&ops)
    ```

  - Conditional variable
    * Conditional Variable is one of the synchronization mechanisms
    * A condition variable is basically a container of Goroutines that are waiting for a certain condition
    * How to make a goroutine wait till some event(condition) occurs?
      - one way - wait in a loop for condition
        ```
          var sharedRsc = make(map[string]string)
          go func() {
            defer wg.Done()
            mu.Lock()
            for len(sharedRsc) == 0 {
              mu.Unlock()
              time.Sleep(100 * time.Millesecond)
              mu.Lock()
            }
            // Do processing
            fmt.Println(sharedRsc["rsc"])
            mu.Unlock()
          }()

          this is inefficient
        ```
      - We need some way to make goroutine suspend while waiting
      - we need some way to signal the suspended goroutine when particular event has occured
      - sync.Cond
        * conditional variable are type
          - var c *sync.Cond
        * we use constructor method sync.NewCond() to create a conditional varaible, it takes sync.Locker interface as input, which is usually sync.Mutex
          ```
            m := sync.Mutex{}
            c := sync.NewCond(&m)
          ```
        * sync.Cond method has three methods
          - c.Wait()
            * example:
            ```
              c.L.Lock()
              for !condition() {
                c.Wait()
              }
              .... make use of condtion ....
              c.L.Unlock()
            ```
            * wait method suspends the execution of the calling goroutine
            * it automatically unlocks c.L, before suspending the goroutine
            * wait() cannot return unless awoken by BroadCast or Signal
            * Once woken up, it acquires the lock(c.L) again
            * Because c.L is not locked when wait first resumes, the caller typically cannot assume that the condition is true when wait returns(as another go routine might have interfered). Instead, the called should wait in a loop

          - c.Signal()
            * Signal wakes up one goroutine waiting on c (condition), if there is any
            * Signal finds goroutine that has been waiting the longest and notofies that
            * it is allowed but not required for the caller to hold c.L during the call

          - c.Broadcast()
            * Broadcast wakes all the goroutines wainting on c(condition)
            * it is allowed but not required for the caller to hold c.L guring the call

          - example
          ```
          Scenario  1: (using wait and signal)
                    | G2 | => Consumer

                m := sync.Mutex{}
                c := sync.NewCond(&m)
                var sharedRsc = make(map[string]string)
                go func() {
                  defer wg.Done()
                  c.L.Lock()
                  for len(sharedRsc) == 0 {
                    c.Wait() // goroutine is suspended and waits
                  }
                  // Do processing
                  fmt.Println(sharedRsc["rsc"])
                  c.L.Unlock()
                }()


                    |G1| => Producer(signaling)

                go func() {
                  defer wg.Done()
                  c.L.Lock()
                  sharedRsc["rsc"] = "foo"
                  c.Signal()
                  c.L.Unlock()
                }()

          Scenario 2: (using wait & braodcast, with multiple goroutines)
                    | G2 | => Consumer

                m := sync.Mutex{}
                c := sync.NewCond(&m)
                var sharedRsc = make(map[string]string)
                go func() {
                  defer wg.Done()
                  c.L.Lock()
                  for len(sharedRsc) < 1 {
                    c.Wait() // goroutine is suspended and waits
                  }
                  // Do processing
                  fmt.Println(sharedRsc["rsc1"])
                  c.L.Unlock()
                }()

                    | G3 | => Consumer

                m := sync.Mutex{}
                c := sync.NewCond(&m)
                var sharedRsc = make(map[string]string)
                go func() {
                  defer wg.Done()
                  c.L.Lock()
                  for len(sharedRsc) < 2 {
                    c.Wait() // goroutine is suspended and waits
                  }
                  // Do processing
                  fmt.Println(sharedRsc["rsc2"])
                  c.L.Unlock()
                }()

                 |G1| => Producer(broad casting)

                go func() {
                  defer wg.Done()
                  c.L.Lock()
                  sharedRsc["rsc1"] = "foo"
                  sharedRsc["rsc2"] = "foo"
                  c.BroadCast() // signals all goroutines
                  c.L.Unlock()
                }()
          ```

  - Sync Once
    * Run one time initialization functions
    * sync.Once ensures that only one
    * example
    ```
      var wg sync.WaitGroup

    var once sync.Once

    load := func() {
      fmt.Println("Run only once initialization function")
    }

    wg.Add(10)
    for i := 0; i < 10; i++ {
      go func() {
        defer wg.Done()

        //TODO: modify so that load function gets called only once.
        once.Do(load)
      }()
    }
    wg.Wait()
    ```

  - Sync pool
    * sync.Pool
    * create and make avaliable pool of things for use
    ```
      // retrive existing instance or create new instance and return
      b := bufPool.Get().(*bytes.Buffer)
      // release the instance once work done
      bufPool.Put(b)

      example:

      // create pool of bytes.Buffers which can be reused. here creating only one instance
      var bufPool = sync.Pool{
        New: func() interface{} {
          fmt.Println("allocate new bytes.Buffer")
          return new(bytes.Buffer)
        },
      }

      func log(w io.Writer, val string) {
        b := bufPool.Get().(*bytes.Buffer)
        b.Reset()

        b.WriteString(time.Now().Format("15:04:05"))
        b.WriteString(" : ")
        b.WriteString(val)
        b.WriteString("\n")

        w.Write(b.Bytes())
        bufPool.Put(b)
      }

      func main() {
        log(os.Stdout, "debug-string1")
        log(os.Stdout, "debug-string2")
      }
    ```

# Race Detector
  * Go provides race decetor tool for finding race conditions in GO code.
  ```
    $ go test -race mypkg // test the package
    $ go run -race mysrc.go // compile and run the program
    $ go build -race mycmd // build the command
    $ go install -race mypkg // install the package
  ```
  * Binary needs to be race enabled
  * when racy behaviour is detedtec a warning is printed
  * Race enabled binary will be 10 times slower and consume 10 times more memory
  * Integration tests and load tests are good candidates to test with binary wuth race condition
  * Alwyas check ig inbuilt functions are run in a seperate gor-routine or in custom ones also, if so then use channels to communicate data. Read and write to shared data must be done in the main goroutine. see example below
  * example - using timer function, it runs in a seperate goroutine
  ```
    func main() {
    start := time.Now()
    reset := make(chan bool)
    var t *time.Timer
    t = time.AfterFunc(randomDuration(), func() {
        fmt.Println(time.Now().Sub(start))
        reset <- true
      })
      for time.Since(start) < 5*time.Second {
        <-reset
        t.Reset(randomDuration())
      }
    }

    func randomDuration() time.Duration {
      return time.Duration(rand.Int63n(1e9))
    }
  ```

# Web Crawler
  * see example for cuncurrent code

# Concurrency Patterns
  * Pipelines
    * Process streams or batches of data
      - G1 -> G2 -> G3   (ch on each ->)
    * each go routien is a stage, each stage takes data in, perform an operation on it and send the data out
    * Stages:
      * seperate the concerns of each stage
      * process individual stage concurrently
      * A stage could consume and return the same type
        - ex: `func square(in <-chan int) <-chan int{}`
      * this enables composability of pipeline
        - `square(suqre(generator(2,3)))`, assume generator returns `<-chan int`
      * example:
        ```
        func generator(nums ...int) <-chan int {
        out := make(chan int)
        go func() {
            for _, n := range nums {
              out <- n
            }
            close(out)
          }()
          return out
        }

        func square(in <-chan int) <-chan int {
          out := make(chan int)
          go func() {
            for n := range in {
              out <- n * n
            }
            close(out)
          }()
          return out
        }

        func main1() {
          ch := generator(2,3)
          out := square(ch)

          for n := range out {
            fmt.Println(n)
          }
        }

        func main2() {

          for n := square(generator(2,3)) {
            fmt.Println(n)
          }
        }

        ```
    * Fan-out & Fan-in
      * example:
      ```
               |------> G2a ---> Ma ----|
          G1 ---------> G2b ---> Ma --------G3---->G4
               |------> G2c ---> Ma ----|

               Ma -> merge goroutines

      ```
      * Fan-out:
        * Multiple goroutines are started to read data from the single channel
        * Distribute work amongst a group of workers goroutines to parallellize the cpu usage and the I/O usage
        * Helps computational intensive stage to run faster
      * Fan-in:
        * process of combining multiple results into one channel
        * we create merge goroutines, to read data from multiple input channels and send data data to a single output channel
      * example :
      ```
        package main

        import (
          "fmt"
          "sync"
        )

        func generator(nums ...int) <-chan int {
        out := make(chan int)
        go func() {
          for _, n := range nums {
            out <- n
          }
          close(out)
        }()
        return out
      }

      func square(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
          for n := range in {
            out <- n * n
          }
          close(out)
        }()
        return out
      }

      func merge(cs ...<-chan int) <-chan int {
        out := make(chan int)
        var wg sync.WaitGroup

        output := func(c <-chan int) {
          for n := range c {
            out <- n
          }
          wg.Done()
        }

        wg.Add(len(cs))
        for _, c := range cs {
          go output(c)
        }

        go func() {
          wg.Wait()
          close(out)
        }()
        return out
      }

      func main() {
        in := generator(2, 3)

        c1 := square(in) // fan-out
        c2 := square(in) // fan-out

        out := merge(c1, c2) // fan-in
        fmt.Println(<-out)
      }
      ```

  * Cancelling Goroutines
    * In the previous example, we use `close(out)`, `range` in each go routine to indicate work is done and to loop over the values sent from upstream. Similarly we need to wait until all the values are received.
    * In real pipelines:
      * Receiver stages may only need a subset of values to make progress
      * A stage can exit early because an inbound value represents an error in an earlier stage
      * Receiver should not have to wait for the remaining values to arrive
      * we want earlier stages to stop producing values that later stages don't need
    * How to signal goroutine to abandon/terminate what they are doing?
      * Cancellation of goroutines
        * Pass a read-only `done` channel to goroutine
        * close the channel, to send broadcast signal to all goroutine
        * On receiving the signal on done channel, GOroutines needs to abandon their work and terminate.
        * we use `select` to make send/receive operation on channel pre-emptible, by multiplexing on done receive channel
        ```
          select {
            case out <- n:
            case <-done:
              return
          }
        ```
        * example:
        ```
          package main

          import (
            "fmt"
            "sync"
          )

          func generator(done <-chan struct{}, nums ...int) <-chan int {
            out := make(chan int)

            go func() {
              defer close(out)
              for _, n := range nums {
                select {
                case out <- n:
                case <-done:
                  return
                }
              }

            }()
            return out
          }

          func square(done <-chan struct{}, in <-chan int) <-chan int {
            out := make(chan int)
            go func() {
              defer close(out)
              for n := range in {
                select {
                case out <- n * n:
                case <-done:
                  return
                }
              }
            }()
            return out
          }

          func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
            out := make(chan int)
            var wg sync.WaitGroup

            output := func(c <-chan int) {
              defer wg.Done()
              for n := range c {
                select {
                case out <- n:
                case <-done:
                  return
                }
              }
            }

            wg.Add(len(cs))
            for _, c := range cs {
              go output(c)
            }

            go func() {
              wg.Wait()
              close(out)
            }()
            return out
          }

          func main() {
            done := make(chan struct{}) // we use struct{} to send only close signal and not sen any data
            //defer close(done) // cancel gouroutine after recieving one value

            in := generator(done, 2, 3) // go idiom to pass done channel firsr

            c1 := square(done, in)
            c2 := square(done, in)

            out := merge(done, c1, c2)

            // TODO: cancel gouroutine after recieving one value

            fmt.Println(<-out)
            close(done) // cancel gouroutine after recieving one value

            time.Sleep(10 * time.Millisecond) // wait for sometime until all goroutines terminate
            // o/p only 4
            g := runtime.NumGoroutines()
            fmt.Print("no of Goroutines = %d\n", g) // o/p is 1
          }

          // guidelines for pipeline construction

          // stages close their outbound channels when all the send operations are done.
          // stages keep receiving values from inbound channels until those channels are closed or the senders are unblocked.

        ```

# image processing pipeline
  * see video, and example

# context Package
  * In a go service each incoming request is handled in its own goroutine
  * The spawned request goroutines starts other goroutines to access database, rpc calls etc
  * The set of goroutines working together needs access to request specific values, i,e identity of the end user, authorization tokens and request deadline
  * when the request is cancelled by the end user or the request is timedout we want all the group of goroutines working on the request to also cancel or exit quickly
  * Need for context:
    - we need a way to propogate request-scoped data down the call graph
    - we need a way to propogate cancellation signal down the call graph
    - Context package provides API's for cancelling branches of call graph
    - Context package provides a data-bag for transporting request-scoped data through call-graph
  * Context interface
    ```
      type context interface {
        // Done returns a channel that is closed when this context is cancelled or times out
        Done() <- struct
        // Err indicates why this context was cancelled, after the done channel is closed
        Err()
        // Deadline returns the time when this context will be cancelled, if any.
        Deadline() (deadline time.Time, ok bool)
        // value returns the value associated with key or nil if none
        Value(key interface{}) interface{}
      }
    ```
  * A context is safe for simultaneous use by multiple goroutines. We can pass a single context to any number of goroutines and send a cancel signal to all go routines to abandon their work and terminate.
  * Functions provided by context
    - context.Background()
      * returns a empty context
      * it is the root of any context tree
      * it is never cancelled and no value and has no deadline
      * It is typically used by main function
      * Acts as a top level context for incoming request
    - context.Todo()
      * It returns a empty context
      * Todo's intended purpose is to serve as a placeholder
      * it is used when we don't know which context to utilize
  * Context Package ... Cancellation
    * Context is immutable
    * Context package provides functions to add new behaviour
    * To add cancellation behaviour we have function like, these functions generate new instance of the context and add options relative to the behaviour of these functions.
      - context.WithCancel()
        ```
          * takes parent context as the input
          ctx, cancel := context.WIthCancel(context.Background())
          defer cancel()
          * WithCancel returns a copy of parent with a new Done channel and a cancel function. The cancel function is used to close context's done channel
          * Closing the done channel indicates to an operation to abandon its work and return
          * Cancelling the context releases the resources associated with it
          * It is important to call the cancel() function as soon as the operation is complete, if we don't call teh cancel function then their will be a memory leak. The resources won't be freed until the current contex is cancelled or the parent context is cancelled
          * cancel() function does not wait for the operation to stop, it just closes the done channel and returns
          * cancel() function can be called from multiple goroutines simultaneously, after the first call to cancel() the subsequent call does nothing to cancel()

          example:

          ctx, cancel := context.WithCAncel(context.Background())
          ch := generator(ctx)
          ...
          if n == 5 {
            cancel()
          }     // parent goroutine


          child: gorouitne:

          for {
            select {
              case <-ctx.Done():
                return ctx.Err()
              case dst <- n:
                n++
            }
          }

          complete example:

          func main() {

            // TODO: gen generates integers in a separate goroutine and
            // sends them to the returned channel.
            // The callers of gen need to cancel the context once
            // they consume 5 integers
            // so that internal goroutine
            // started by gen is not leaked.
            gen := func(ctx context.Context) <-chan int {
              dst := make(chan int)
              n := 1
              go func() {
                defer close(dst)
                for {
                  select {
                  case <-ctx.Done():
                    return
                  case dst <- n:
                    n++
                  }
                }
              }()
              return dst
            }

            // Create a context that is cancellable.
            ctx, cancel := context.WithCancel(context.Background())

            ch := gen(ctx)
            for n := range ch {
              fmt.Println(n)
              if n == 5 {
                cancel()
              }
            }
          }

        ```
      - context.WithTimeout()
      ```
        duration := 5*time.Millisecond
        ctx, cancel := context.WithTimeout(context.Background(), duration)
        defer cancel()

        * WithTimeout() takes parent context and time duration as input
        * WithTimeout returns a new Context that closes its done channel after the given timeout duration
        * WithTimeout() is useful for setting a deadline on the requests to backend servers

        * example:

        func main() {

            // set a timeout for http get operation.
            req, err := http.NewRequest("GET", "https://andcloud.io", nil)
            if err != nil {
              log.Fatal(err)
            }

            // Create a context with a timeout of 100 milliseconds.
            ctx, cancel := context.WithTimeout(req.Context(), 100*time.Millisecond)
            defer cancel()

            // Bind the new context into the request.
            req = req.WithContext(ctx)

            resp, err := http.DefaultClient.Do(req)
            if err != nil {
              log.Println("ERROR:", err)
              return
            }

            // Close the response body on the return.
            defer resp.Body.Close()

            // Write the response to stdout.
            io.Copy(os.Stdout, resp.Body)
          }


      ```
      - context.WithDeadline()
      ```
        deadline := time.Now().Add(5*time.Millisecond)
        ctx, cancel := context.WIthDeadline(context.Background(), deadline)
        defer cancel()

        * WithDeadline() takes parent context and clock time as input
        * withDeadline returns a new Context that closes its done channel when the machines clock advances past the given deadline

        example:

        func main() {

          // set deadline for goroutine to return computational result.
          deadline := time.Now().Add(5 * time.Millisecond)
          ctx, cancel := context.WithDeadline(context.Background(), deadline)
          defer cancel()

          compute := func() <-chan data {
            ch := make(chan data)
            go func() {
              defer close(ch)

              deadline, ok := ctx.Deadline()
              if ok {
                if deadline.Sub(time.Now().Add(50*time.Millisecond)) < 0 {
                  fmt.Println("not sufficient time given, terminating")
                  return
                }
              }

              // Simulate work.
              time.Sleep(50 * time.Millisecond)

              // Report result.
              select {
              case ch <- data{"123"}:
              case <-ctx.Done():
                fmt.Println("work cancelled")
                return
              }

            }()
            return ch
          }

          // Wait for the work to finish. If it takes too long move on.

          ch := compute()
          d, ok := <-ch
          if ok {
            fmt.Println("work complete", d)
          }

        }

      ```
    * The derived context is passed to child goroutines to facilitate their cancellation

  * Context Package as Data Bag
    * Context package can be used to transport request-scoped data down the call graph
    * context.WithValue() provides a way to associate request scoped values with a Context
    * WithValue() takes a parent context and a key value pair as input and returns a copy of the parent context with the key-pair values
    * ctx.Value() is used to extract data from the context
    ```
      example:

      type userIDKey string
      type database map[string]bool

      var db database = database{
        "jane": true,
      }

      func main() {
        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()
        processRequest(ctx, "jane")
      }

      func processRequest(ctx context.Context, userid string) {
        // send userID information to checkMemberShipStatus for
        // database lookup.
        vctx := context.WithValue(ctx,
          userIDKey("userIDKey"),
          userid)

        ch := checkMemberShipStatus(vctx)
        status := <-ch
        fmt.Printf("membership status of userid : %s : %v\n", userid, status)
      }

      // checkMemberShipStatus - takes context as input.
      // extracts the user id information from context.
      // spins a goroutine to do database lookup
      // sends the result on the returned channel.
      func checkMemberShipStatus(ctx context.Context) <-chan bool {
        ch := make(chan bool)
        go func() {
          defer close(ch)
          // do some database lookup
          userid := ctx.Value(userIDKey("userIDKey")).(string)
          status := db[userid]
          ch <- status
        }()
        return ch
      }
    ```

  * Go Idioms for COntext Package
    * Incoming requests to a server should create a context
      - create context early in processing task or request
      - create a top level context
        func main () {
          ctx := context.Background()
        }
      - http.Request value already contains a Context
      func handleFunc(w http.ResponseWriter, req *http.Request) {
        ctx, cancel = context.WithCancel(req.Context())
      }
    * Outgoing calls to servers should accept a context
      - Higher level calls need to tell lower level calls how long they are willing to wait

      // create a context with a timeout of 100 milliseconds.
      ctx, cance := context.WIthTimeout(req.Context(), 100*time.Millisecond)
      defer cancel()

      // Bind the new context into the request
      req = req.WIthContext(ctx)

      // Do will handle the context level timeout
      resp, err := http.DefaultClient.Do(req)

    * Pass aContext to function performing I/O
      - Any function that is performing I/O should accept a context value as its first parameter and respect any timeout or deadline configured by the caller
      - Any API's that takes a Context, the idiom is to have the first parameter accept the context value
    * Any change to a context value create a new context value that i then propogated forward. The functions that are passed with the new context are affected and not the old functions

    * When a Context is cancelled, all context derived from it are also canceled
      - if a parent context is cancelled, all children derived by that parent context are also canceled
    * Use TODO context if you are unsure about which context to use
      - if a function is not responsible for creating top level context
      - we need a temp top-level context until we figured out where the actual context will come from
    * Use context values only for request-scoped data
      - Do not use the context value to pass data into a function which becomes essential for its successful execution
      - A function should able to execute its logic with an empty Context Value.
# Http Server Timeouts with Context package
  - to timeout request after a specific time

# Interfaces

# Generics
