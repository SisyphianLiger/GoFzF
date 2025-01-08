# Go's Concurrency Building Blocks

## Goroutines
Basic unit of organization in Go
A Goroutine is a function that is running concurrently

```Go
func main() {
    go sayHello()
}

func sayHello() {
    fmt.Println("Hello")
}

// Anynoumous version
go func() {
    fmt.Println("hello")
}()
```

Coroutines --> subroutinces (functions, closures, or methods) nonpreemtivemeaning they cannot be interupted, that have multiple points throuout which allow for suspension or rentry

Go's runtime observes the behavior of goroutines and automatically suspends them when they block and resumes whenunblocked

This relationship == goroutines are a special class of coroutines

Go uses an M:N scheduler 

M is green threads N is OS threads

Goroutines are somewhat a green thread but also not
The scheduler blocks if not enough green threads


Concurrency Model: fork-join model
fork --> program can split off a child branch of execution from parent
join --> at some point in the furutre the branches will join together

An example of the fork-join model
```Go
var wg sync.WaitGroup
sayHello := func() {
    defer wg.Done() 
    fmt.Println("hello")
}
wg.Add(1)
go sayHello() // This is the "Fork"
wg.Wait() // This is the "Join"
```

### What happens to references to data that are needed in a closure?

```Go
var wg sync.WaitGroup

for _, salutation := range []string{"hello", "greetings", "good day"} {
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println(salutation) // Not gonna be a good day
    }(salutation)
}
wg.Wait()
```

So, what does the code output here? (on average so to speak)

Well, given that salutation is being referenced in the closure, however
there is no copy. They problem here is that the go run time must decide 
what do to with it. Well luckily the go runtime does in fact, understand that it needs to hold onto the refernce, but does so with the last element of the string because the loop finishes before any goroutines begin running, meaning that the refernce on the heap is "good day"

To fix this, for each for loop we take a copy of salutation

```Go
var wg sync.WaitGroup

for _, salutation := range []string{"hello", "greetings", "good day"} {
    wg.Add(1)
    go func(salutation string) { // copy here
        defer wg.Done()
        fmt.Println(salutation) // A great day is upon us
    }(salutation)
}
wg.Wait()
```

### Hey This is a Waitgroup example

A weight group is a bounded set of goroutine actions, effectively we make it such that we wait until all routines are done before we exit

```Go
var wg sync.WaitGroup // struct for WG

wg.Add(1) // Adding 1 "task"
go func() {
    defer wg.Done() // Making sure the WG will be notified when the func is done
    fmt.Println("1st Gorouting")
    time.Sleep(1)
}()

wg.Add(1) 
go func() {
    defer wg.Done() // Making sure the WG will be notified when the func is done
    fmt.Println("2nd Gorouting")
    time.Sleep(2)
}()

// The sleep are here so we can see 1st, then 2nd be completed

wg.Wait() // This is the check point where we make sure goroutines wg.Done have been met
fmt.Println("Finished")
```

Another slick example with for loops
```Go
hello := func(wg *sync.WaitGroup, id int) {
    defer wg.Done() // same thing here
    fmt.Printf("Hello from &v!\n", id)
}


const numGreeters = 5

var wg sync.WaitGroup
wg.Add(numGreeters) // here we add by the total numbers of waiters
for i := 0; i < numGreeters; i++ {
    go hello(&wg, i + 1) // here we call the func with ref to the wg
}
wg.Wait() // same thing here
```

KEY TAKEAWAY:

ADD --> DEFER DONE --> WAIT

## Mutex and RWMutex
```Go
var count int
var lock sync.Mutex
increment := func() {
    lock.Lock()
    defer lock.Unlock()
    count++
    fmt.Printf("Incrementing: %d\n", count)
}

decrement := func() {
    lock.Lock()
    defer lock.Unlock()
    count--
    fmt.Printf("Decrementing : %d\n", count)
}

var arithmetic sync.WaitGroup
for i := 0; i <= 5; i++ {
    arithmetic.Add(1)
    go func() {
        defer arithmetic.Done()
        increment()
}()

for i := 0; i <= 5; i++ {
    arithmetic.Add(1)
    go func() {
        defer arithmetic.Done()
        decrement()
}()

arithmetic.Wait()
fmt.Println("Arithmetic Complete.")

```

### Mutex Locks
Mutex locks are great when you do not know who needs access to the critical section. Meaning it will always perform mutual exclusive access regardless of RW. 

### RWMutex Locks
In the case where we know reads are the primary access point we can use RW Mutexes, which will ensure that the observers in a particular concurrent setup will be able to access the critical section many times over UNLESS there is a writer acquiring the lock. If a R has the locks and a W comes up, the que is set up in such a way such that the W will block until the R's currently are activated, and any R after the W are put after the W not before. 

## Cond
"Short for condition --> a rendezvous point for goroutines waiting for or announcing the occurence of an event."

```Go
c := sync.NewCond(&sync.Mutex{})
c.L.Lock
for conditionTrue() == false {
    c.Wait() // Description below
} 
c.L.Unlock()
c.Signal()
```

What is a Wait method?

Wait call suspends the current goroutine allowing for other goroutines to run the OS thread

--> Upon entering: Ulock is called on Cond variable Locker
--> Upon exiting: Lock is called on the Cond variable Locker

What is a Signal Method?

Notifies Goroutines blocked on a Wait call that the condition has been tirggered on

Signal finds the goroutine that is waiting the longest

What is a Broadcast Method?
Similar to Signal, notifies goroutines that have been blocked.

MAJOR DIFFERENCE: it finds the next goroutine 

Cond is more performant then channels

## Worker Pools
Goroutines are cheap so why would we need to use a worker pool. Well while this is true there are situations in concurrency, such as when we need high throughput, where a worker pool is a better use of memory. 

Take the following code as an example:
```Go

func connectToService() interface{} {
    time.Sleep(1*time.Second)
    return struct{}{}
}

func warmServiceConnCache() *sync.Pool { // 1
    p := &sync.Pool {
        New: connectToService,
    }
    for i := 0; i < 10; i++ {
        p.Put(p.New())
    }
    return p
}
func startNetworkDaemon() *sync.WaitGroup { // 2
    var wg sync.WaitGroup
    wg.Add(1)
    go func() { // 3
        connPool := warmServiceConnCache()
        server, err := net.Listen("tcp", "localhost:8080")
        if err != nil {
            log.Fatalf("cannot listen: %v", err)
        }
        defer server.Close()
        wg.Done()
        for {
            conn, err := server.Accept()
            if err != nil {
                log.Printf("cannot accept connection: %v", err)
                continue
            }
            svcConn := connPool.Get()
            fmt.Fprintln(conn, "")
            connPool.Put(svcConn)
            conn.Close()
        }
    }()
    return &wg
}
```

### So lets think here first the form of a sync.Pool func is as followed.

1. We create a struct sync.Pool that takes some service connection logc. In our case for testing, just an interface that returns a struct{}{}. Then, we generate 10 cached connections to a service. When that is said and done, we have our connection pool

2. Starting the network Daemon we add for each call to the network a Waitgroup add and then our go function here will make a goroutine with the connection pool that will then try to serve a tcp protocol at localhost:8000. This process with a connection pool means that the worker pool can be reused. 


### Advantages of this strategy:
- Used when you have concurrent processes that require objects but dispose of them very rapidly after insantiation
- or wen construction of these objects negatively impact memory

### Caveats:
- If code is not roughly homogenous, then memory costs of spending too much time converting retrieved data from the pool could be to costly versus just instantiating it

I.E. Pool random lengths will not help you much because of the variability of the length

### REMEMBER 
- Instantiating sync.Pool give it a New member variable that is thread-safe
- When you receive an instance fom Get, make no assumptions of the state of the object you recieve back
- Make sure to call Put when you are finished with the object you pulled out of the pool (with defer)
- Objects in the pool must be roughly uniform in data 

## Channels
Think of channels as connections of streams, passing data along a channel. Streams do not need to know eachothers data lines, because of the convergence points

### How to create a Channel
```Go
var dataStream chan interface{} // 1. 
dataStream = make(chan interface{}) // 2. 
```

1. Declare the type of channel an interface{} 
2. Here we instantiate the channel with make

This channel can be RW

### Unidirectional Channels
We can set up Channels unidirectionally 
- A Channel that supposts sending XOR receiving

```Go
var dataStream <- chan interface{} // 1. 
dataStream = make(<- chan interface{}) // 2. 

var dataStream chan<- interface{} // 3. 
dataStream = make(chan<- interface{}) // 4. 
```
#### Example of a Read Only Channel
1. The only difference is the `<-` operator on the LEFTHAND side, LEFT ON READ
2. The only difference is the `<-` operator on the LEFTHAND side

#### Example of a Write Only Channel
3. The only difference is the `->` operator on the RIGHTHAND side, LEFT ON READ
4. The only difference is the `->` operator on the RIGHTHAND side

Go implicitily converts bidirectional channels to unidirectional channels when needed

##### EXAMPLE
```Go
var receiveChan <-chan interface{}
var sendChan chan<- interface{}
var dataStream := make(chan interface{}
```
Channels are types, here we can place anything in it bt we can define it with types

##### Writing to a chan
```Go
var stringStream := make(chan string) 
go func() {
    stringStream <- "Hello Channel" // The <- indicates we are writing to the CHannel
}()

fmt.Println(<-stringStream) // Here becase the data is stored into the chan we can print it

```

### Channel Properties
- Channels are blocking, any goroutine that attempts to write to a channel that is full will wait until the channel is empty

- This causes the main gorouting and anonymous goroutine to block deterministically

- If no value is placed onto a channel it can cause a deadlock


COMMON CHANNEL PATTERN
```Go
	intStream := make(chan int) 
	go func() {
		defer close(intStream) // Make sure we close the channel after the routine
		for i := 1; i <= 5; i++ {
			intStream <- i
		}
	}()
	for integer := range intStream {
		fmt.Printf("%v ", integer) // can still read intStream and print :)
	}
```

#### A way to make multiple goroutines unblock at once
```Go
begin := make(chan interface{})
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(i int) {
        defer wg.Done()
        <- begin // Goroutine waits until it can continue
        fmt.Printf("%v has begun\n", i)
    }(i)
    fmt.Println("Unblocking goroutines...")
    close(begin) // channel is closed releasing gorotuines from block
    wg.Wait() 
```

This is very similar to the `sync.Cond` however, channels are composable.

### Buffered Channels
Channels that are given a capacity when instantiated. This means that there is some n given to the channel and that even if no reads are performed on the channel, a goroutine can still perform n writes.

```Go
var dataStream chan interface{}
dataStream = make(chan interface{}, 4) // 1
```
1. This is a buffered channel that has a capacity of 4, so regardless of where it is being read from only 4 goroutines may be on it.

#### More Complex Example of Channel Buffer
```Go
var stdoutBuff bytes.Buffer // 1
defer stdoutBuff.WriteTo(os.Stdout) // 2

intStream := make(chan int, 4) // 3
go func() {
    defer close(intStream)
    defer fmt.Fprintln(&stdoutBuff, "Producer Done.") // 4
    for i := 0; i < 5; i++ {
        fmt.Fprintln(&stdoutBuff, "Sending: %d\n", i)
        intStream <- i // 5
    }
}()

for integer := range intStream {
    fmt.Fprintln(&stdoutBuff, "Recieved: %v. \n", integer) // Channel is now closed so we read the reads
}

```

1. We create an in-memory buffer to mitigate nondeterministic output
2. ensure the buffer is written out to stdout before exiting
3. create a buffered channel with cap of 4
4. Occures when the function has closed the channel
5. Sending integer i to the channel


### Nil Channels
Default values for channels is nil

#### Reading from a nil channel
```Go
var dataStream chan interface{}
<-dataStream
```

Produces a "deadlock" for read channels

#### Writing to a nil send channel
```Go
var dataStream chan interface{}
dataStream <- struct{}{}
```

Produces a "deadlock" for send channel

#### Closing a Channel
```Go
var dataStream chan interface{}
close(dataStream)
```

Will cause a panic

##### Botom Line
Make sure all channels are initialized first

### Table of Interactions for Channels
| Operation | Channel State         | Result |
| --------- | -------------         | ------ |
| Read      | nil                   | Block |
|           | Open and Not Empty    | Value |
|           | Open and Empty        | Block|
|           | Closed                | Defaule value, false|
|           | Write Only            | Compiliation Error |
| Write     | nil                   | Block |
|           | Open and Full         | Block |
|           | Open and Not Full     | Write Value |
|           | Closed                | panic |
|           | Receive Only          | Compliation Error |
| close     | nil                   | panic |
|           | Open and Not Empty    | closes Channel; reads succeed until empty --> default |
|           | Open and Empty        | closes channel; default value|
|           | Closed                | closed |
|           | Recieve Only          | Compiliation Error |


To make sure we make Channels consistant 

```Go
chanOwner := func() <- chan int {
    resultStream := make(chan int, 5) // 1
    go func() { // 2
        defer close(resultStream) // 3
        for i := 0; i <= 5; i++; {
            resultStream <- i
        }
    }()
    return resultStream // 4
}

resultStream := chanOwner()
for result := range resultStream { // 5
    fmt.Printf("Recieved: %d\n", result)
}
fmt.Println("Done receiving1")
```

1. Instantiate a buffered Channel
2. Make an anonymous goroutine that performs writes on resultStream
3. Ensure once updated we close Channel
4. return channel because it's declared as read only it implicitly converts to
5. range over only caring about blocking/closed channels

A Consumer channel should only have access to a read channel, and producer write

## The select Statement
The `select` keyword is pretty important to go, so it will be very important for channels

```Go
var c1, c2 <- chan interface{}
var c3 chan <- interface{}
select {
    case <- c1:
        // Do something cool
    case <- c2:
        // Do something cool
    case c3 <- struct{}{}:
        // Do something cool
    }
```
select blocks are not tested sequentially, and execution wont automatically fall through if no criteria is met.

Actions occur simultanously


```Go
start := time.Now()
c := make(chan interface{})
go func() {
    time.Sleep(5*time.Second)
    close(c) // 1
}()

fmt.Println("Blocking on read...")
select {
    case <-c: // 2
    fmt.Printf("Unblocked %v later. \n")
}
```
1. Close the channel after waiting 5 seconds
2. Attempt to read the channel 

### What happens when multiple channels have something to read 
```Go
    c1 := make(chan interface{}); close(c1)
    c2 := make(chan interface{}); close(c2)

    var c1Count, c2Count int
    for i := 1000; i >= 0; i-- {
        select {
            case <- c1:
                c1Count++
            case <- c2:
                c2Count++
        }
    }
```

So with the select case...apparently the Go runtime cannot know anything about the intent of the channels, but will try to run it in the average case, weighing certain options when needed. 

