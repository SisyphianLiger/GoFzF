# Concurrency Patterns in Go
## Confinement
Confinement is the idea that we ensure information is oly ever available from one concurrent process.
Either Ad hoc or lexical

### Ad hoc confinement
```Go
    data := make([]int, 4)

    loopData := func(handleData chan <- int) {
        defer close(handleData)
        for i := range data {
            handleData <- data[i]
        }
    }

    handleData := make(chan int)
    go loopData(handleData)

    for num := range handleData {
        fmt.Println(num)
    }
```

Ok Ok, this is a horirble design pattern need, basically, there is no real encapsulation of data here because data is make outside of the function, now of course the for loop is the only thing using the data.

TBH this is not the way and using a closure is so lets go to 

### Lexical Confinement
```Go
    chanOwner := func () <-chan int {
        results := make(chan int, 5) // 1 
        go func() {
            for i := 0; i <= 5; i++ {
                results <- 1
        }()
        return results

    consumer := func(results <- chan int) { // 3 
        for result := range results {
            fmt.Println("Received: %d\n", result)
        }
        fmt.Println("Done Receiving")
    }

    results := chanOwner() // 2
    consumer(results)

```

This is a more FP approach tbh, and the go func() is a closure. 

1. Make it certain that the chanOwner create in the func the results channel
2. results is stored as a read channel from `<- chan int` 
3. There is a read only copy of the int channel, this "confines usage of the channel" to read

One more example
```Go
printData := func(wg *sync.WaitGroup, data []byte) {
    defer wg.Done()
    
    var buff bytes.Buffer
    for _,b := range data {
        fmt.Fprintf(&buff, "%c", b)
    }
    fmt.Println(buff.String())
}

    var wg sync.WaitGroup
    wg.Add(2)
    data := []byte("golang")
    go printData(&wg, data[:3])
    go printData(&wg, data[3:])

    wg.Wait()
```

So what is so special about lexical versus ad hoc? well think about it, let's say we used ad hoc...this would mean that the data we pass into our function would have access to a gloabl variable, and maybe with two functions accesing that data it would be fine, but imagine 20 or some large N, it becomes much more difficult to discern where the problem may lie. So with Lexical Confinement and closures, we ensure that the data is being passed is safe when being used.

## The for-select Loop
```Go
    for { // Either loop infinitely or range over
        select {
            // do some channel work 
        }
    }
```
When do we see this pattern?
### Sending iteration variables out on a channel
```Go
for _, s := range []string{"a", "b", "c"} {
    select {
    case <-done:
        return
    case stringStream <- s:
    }
}
```

### Loopiing infinitely waiting to be stopped
```Go
for {
    select {
    case <-done:
        return
    default:
    }
    // Do non-preemptable work
}
```

## Preventing Goroutine Leaks
    - Go routines DO cose resources
    - Go routines are not GC'ed by the runtime

### Three types of termination
    - Completed work
    - Cannot continue due t unrecoverable error
    - told to stop working

### Goroutine Leak Example

```Go
    doWork := func(strings <-chan string) <-chan interface{} {
        completed := make(chan interface{})
        go func() {
            defer fmt.Println("doWork exited")
            defer close(completed)
            for s := range strings {
                // Do something
                fmt.Println(s)
            }
        }()
        return completed
    }

    dowWork(nil)
    // Perhaps more work is doen here
    fmt.Println("Done.")
```

Lets fix this
```Go
    doWork := func(done <-chan interface{}, strings <-chan string,) <-chan interface{} { // 1
        terminated := make(chan interface{})
        go func() {
            defer fmt.Println("doWork exited")
            defer close(terminated)
            for {
                select {
                    case s := <- strings:
                        // Do something interesting
                        fmt.Println(s)
                    case <-done: // 2
                        return
                }
            }
        }()
        return terminated 
    }

    done := make(chan interface{}) 
    terminated := doWork(done, nil)

    go func() { // 3
        // Cancel the op after 1 sec
        time.Sleep(1 * time.Second)
        fmt.Println("Cancelling doWork goroutine...")
        close(done)
    }()

    <-termiated // 4
    fmt.Println("Done")
```

1. Pass Done Channel 
2. This is our for select pattern that has done in it to terminate
3. Make another goroutine that will cancel the goroutine spanwed in dowork if more than 1 second pass
4. join goroutine spanwed from doWorkd


```Go
	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = func(channels ...<-chan interface{}) <-chan interface{} { // 1
		switch len(channels) {
		case 0: // 2
			return nil
		case 1: // 3
			return channels[0]
		}
		orDone := make(chan interface{})
		go func() { // 4
			defer close(orDone)
			switch len(channels) {
			case 2: // 5
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default: // 6
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...): // 6
				}
			}
		}()
		return orDone
	}
```

1. We have a func that takes in a variadic slice of channels and returns a single chan
2. Base case for recursive function --> This one is if nil then end
3. Second term criteria is that if slice only has one then return ele
4. Main body: we create a goroutine to wait for messages on our channels without blocking
5. because of our recursion every call needs to have at least 2 channels --> so this is a special case where there is a call to or with only two channels
6. we recursively create an or-channel from all the channels in our slice after the third index, this recurrence relation will destructure the rest of the slice into or-channels to form a tree which the first signal will return.

We pass the orDone channel, so when the goroutines up the tree exit goroutines down the tree also exit

## Error Handling
```Go
checkStatus := func( done <-chan interface{}, urls...string,) <-chan *http.Response {
    responses := make(chan *http.Response)
    go func() {
        defer close(responses)
        for _, url := range urls {
            resp, err := http.Get(url)
            if err != nil {
                fmt.Println(err) // 1.
                continue
            }
            select {
                case <-done: 
                    return
                case responses <- resp:
                }
            }
    }

    done := make(chan interface{})
    defer close(done)
    
    urls := []string{"https://www.google.com", "https://badhost"}
    for response := range checkStatus(done, urls...) {
        fmt.Printf("Response: %v\n", response.Status)
    }

```

1. Here we see the goroutine doing its best to signal that there is an erro. What can it do?
It can't pass it back! How many errors is too many? Will it continue making request?

We need a "Result type"

```Go
type Result struct { // 1
    Error error
    Response *http.Response
}

    checkStatus := func( done <-chan interface{}, urls...string,) <-chan Result { // 2
        responses := make(chan Result)
        go func() {
            defer close(responses)

            for _, url := range urls {

                var result Result
                resp, err := http.Get(url)
                result = Result{Error: err, Response: resp} // 3

                select {
                case <-done: 
                    return
                case results <- result: // 4
                }
            }
        }()
        return results    
    }

    done := make(chan interface{})
    defer close(done)
    
    urls := []string{"https://www.google.com", "https://badhost"}

    for result := range checkStatus(done, urls...) {
        if result.Error != nil { // 5
            fmt.Printf("Response: %v\n", response.Status)
            continue
        }
        fmt.Printf("Response:  %v\n", result.Response.Status
    }

```
1. Now we use a Result Type that contains Error and http response
2. the channel now returns this result 
3. Create Result instance that holds the info
4. Send Result data to our channel
5. Now because the channel has been closed we can iterate printing out the statuses

## Pipelines
Generators 

```Go 
    repeat := func(done <-chan interface{}, values ...interface{},) <-chan interface{} {
        valueStream := make(chan interface{})
        go func() {
            defer close(valueStream)
            for {
                for _,v := range values {
                    select {
                        case <-done:
                            return
                        case valueStream <- v:
                    }
                }   
            }
        }()
        return valueStream
    }

    take := func(done <-chan interface{}, valueStream <-chan interface{}, num int
    ) <-chan interface{} {
        takeStream := make(chan interface{})
        go func() {
           defer close(takeStream)
            for i := 0; i < num; i++ {
                select {
                    case <-done:
                        return
                    case takeStream <- <-valueStream:
                }
            }
        }()
    }

    done := make(chan interface{})
    defer close(done)
    
    for num := range take(done, repeat(done, 1), 10) {
        fmt.Printf("%v ", num) 
    }


```
    
This output because there is only 1 number produces 1111111111

So to make it more...efficient?

```Go
    
    take := func(done <-chan interface{}, valueStream <-chan interface{}, num int
    ) <-chan interface{} {
            takeStream := make(chan interface{})
            go func() {
                defer close(takeStream)
                for i := 0; i < num; i++ {
                    select {
                    case <-done:
                        return
                    case takeStream <- <-valueStream:
                    }
                }
            }()
        }


    repeatFn := func(done <-chan interface{}, fn func() interface{},) <-chan interface{} {
        valueStream := make(chan interface{})
        go func() {
            defer close(valueStream)
            for {
                select {
                    case <- done:
                        return
                    case valueStream <- fn():
                    }
                }
            }()
            return valueStream
        }
    
    // Then we apply that to the already done take func
    done := make(chan interface{})
    defer close(done)

    rand := func() interface {} { return rand.Int()}
    
    for num := range take(done, repeatFn(done, rand), 10) {
        fmt.Println(num)
    }
```

Now with this pipe we can use the take function to apply some function to a value some amount of times. 
There is a caveat to this in that using interface{} as opposed to typed generatos is marginally faster. And there is always 

## Fan-Out, Fan-In
Sometimes we want to add more routines to a compute area and less to others

Lets take a loot at a heavy computation func PRIMES!!!!

```Go
done := make(chan interface{})
    defer close(done)

    start := time.Now()
    rand := func() interface{} { return rand.Intn(50000000) }

    randIntStream := toInt(done, repeatFn(done, rand))

    numFinders := runtime.NumCPU()
    fmt.Printf("Spinning up %d prime finders.\n", numFinders)

    finders := make([]<-chan interface{}, numFinders) // 1.
    fmt.Println("Primes:")

    for i := 0; i < numFinders; i++ {
    finders[i] = primeFinder(done, randIntStream) // 2.
    }

    for prime := range take(done, fanIn(done, finders...), 10) { // 3.
    fmt.Printf("\t%d\n", prime)
    }

    fmt.Printf("Search took: %v", time.Since(start))
```

First, to understand the algorithm for primeFinder, it takes N random integer value(s) with the following snippet of could. In the case of the example program it will just take one int, but can take many ints from the randIntStream

1. This is out fan-in channel. It is a slice of channels made with a max length == to the number of cores on our computer

2. when we fan-out as shown int he primeFinder func, we are returning a channel meaning a value from that channel is of type <-chan int. In theory, this pipeline could process an infinite amount of numbers but in practice it will just return 1 prime number back

3. The fan-in func creates a waitgroup such that all channels coming in must be done computing before moving forward, As such, the take here blocks until all tasks are completed.

And lastly the fan-in func looks like the following:

```Go
    fanIn := func(done <-chan interface{}, channels ...<-chan interface{},) <-chan interface{} { // <1>
        var wg sync.WaitGroup // <2>

        multiplexedStream := make(chan interface{})

        multiplex := func(c <-chan interface{}) { // <3>
            defer wg.Done()
            for i := range c {
                select {
                case <-done:
                    return
                case multiplexedStream <- i:
                }
            }
        }

        // Select from all the channels
        wg.Add(len(channels)) // <4>
        for _, c := range channels {
            go multiplex(c)
        }

        // Wait for all the reads to complete
        go func() { // <5>
            wg.Wait()
            close(multiplexedStream)
        }()

        return multiplexedStream
    }
```

This code is a fan-in function, by design it takes in multiple channels and a done channel such that it will create a waitgroup to synchronize all the work currently being performed from the slice of channels. It does this by both creating a waitgroup but also by having a inner function that adds each channel to a wait group by way of multiplexing. There is also another goroutine whose purpose it is to wait for all goroutines to be finished and remember to close the multiplexed Stream channel as well, then it returns a read only multiplexedStream


## The or-done-channel
Sometimes you need to work with channels in parts of your system unconnected, (no pipes) where you do not know exactly when you can wrap things up. So when you rand over a channel, you can get lots of extra code having to make multiple nested if loops to check

A single gorouting can fix this issue
```Go
    orDone := func(done, c <-chan interface{}) <-chan interface{} {
        valStream := make(chan interface{})
        go func() {
            defer close(valStream)
            for {
                select {
                case <-done:
                    return
                case v, ok := <-c:
                    if ok == false {
                        return
                    }
                    select {
                    case valStream <- v:
                    case <-done:
                    }
                }
            }
        }()
        return valStream
 
    }

    for val := range orDone(done, myChan) {
        // Do something with val
    }
```
now we can handle internal logic of channels outside of our for-loop

## The tee-channel
Sometimes you may want to split values coming in from a channel sot hat you can send them in two diff directions. Like sending somet that compute and some that log for later

 Just like the Unix tee command that reads from standard input and writes to both standard output AND a file, a tee channel in Go splits a single input channel into two output channels.

```Go
    tee := func(
        done <-chan interface{}, in <-chan interface{},) 
        (_, _ <-chan interface{}) {
            
            out1 := make(chan interface{})
            out2 := make(chan interface{})

            go func() {
                defer close(out1)
                defer close(out2)
                // orDone is from above
                for val := range orDone(done, in) {

                    var out1, out2 = out1, out2 // 1,

                    for i := 0; i < 2; i++ { // 2. 
                        select {
                        case <-done:
                        case out1<-val:
                            out1 = nil // 3.
                        case out2<-val:
                            out2 = nil  // 3. 
                        }
                    }
                }
            }()
            return out1, out2
        }
```
1. We want to use local versions of out1, and out2 so we shadow these variables
2. we are going to use one select statement so that writes to out1 and out2 do not block eachother ensuring both are written to, and perform 2 iterations to make sure
3. once written we set its shadowed copy to nil so that further writes will block and the other channel may continue

### Why the shadowing?
Because we use these copies to signal to the select channel when to unblock a value. In otherwords it could be the case that out1<-val has recieved a value, but if we didn't have a local variable, then there is a chance we write to it again. But by going through the local variables, instead, we can nil it, such that the select statment will choose the non-niled value next, and we couple this with two iterations.

## The bridge-channel
Sometimes you want to consume values from a sequence of channels like
`<-chan <-chan interface{}`

For example a pipeline whose lifetime is intermittent.

The consumer may not care about the fact that its values come from a sequece of channels, so dealing with a channel of channels can be cumbersome. We can however define a function that can destructure the channel of channels into a simple channel, called `bridging`

```Go
bridge := func(done <-chan interface{}, chanStream <-chan <-chan interface{}, )
        <- chan interface{} {
        valStream := make(chan interface{}) // 1. 
        go func() {
            defer close(valStream)
            for { // 2 
                var stream <-chan interface{}
                select {
                case maybeStream, ok := <-chanStream:
                    if ok == false {
                        return
                    }
                    stream = maybeStream
                case <-done: 
                    return
                }
                for val := range orDone(done, stream) { // 3
                    select {
                        case valStream <-val:
                        case <-done: 
                }   

            }
        }
    }()
    return valStream
}
```

1. Channel that will return all values from bridging
2. loop that pulls channels off the chan stream providing them to a nested loop for use
3. The loop reads values off the channel given and repeats those values to valStream. if the stream is closed then we break out of the loop performing the reads from this channel and continue with the next iteration of the loop, selecting channels to read from. 

## Queing
Accepting work for your pipeline even though the pipeline is not yet ready for more

Once your stage has completed some work, it stores it in a temporary location in memory so that other stages can retreive it later and your stage doesn't need to hold oa ref to it


#### CAVEAT
Adding a que prematurely can hide synchronization issues such as deadlocks, livelocks, and more

### When do we use a queue?
A cue will almost never speed up the total runtime of your program; it will only allow the program to behave differently

```Go
done := make(chan interface{})
defer close(done)

zeros := take(done, 3, repeat(done, 0))
short := sleep(done, 1*time.Second, zeros)
long := sleep(done, 4*time.Second, short)
pipeline := long

```

The pipeline does the following
1. repeats stage that generates endless streams of 0s
2. Stage that cancels the previous stages after seeing three items
3. a short stage that sleeps one second
4. a long stage that sleeps four seconds

We use a que here, because it provides a sense of "linearity" where here, we can process the short timeframe, with a que quicker than letting everything go all at once. But not from a performance perspective but rather an ordering perspective.

### Applicability of Queues
1. If batching requests in a stage saves time
    - I.E. Sending from In memory to Disk
2. If delays in a stage produce a feedback loop into the system
    - negative feedback loop --> downdward spireals that exist because the recurrent relation exists between 
        the pipeline and its upstream systems; the rate at which upstream stages or systems submit new requests is somehow linked to how effificent the pipeline is

### Little's Law
The law of queuing that predits the through-put of your pipeline
1. L = the average number of units in the system
2. λ = the average arrival rate of units
3. W = the average time a unit spends in the system

The equation ONLY applies to stable systems

A stable system: one in which the rate that work enters the pipeline, or ingress, is equal to the rate in which it exits the system, egress

### How can we decrease W?
The average time a unit spends in the system by a factor of n, we can ONLY decrease the average number of units in the system: 
`L/n = λ * W/n`

### If we increase L (via queues)
We add queues to our stages, which either increases the arival rate of units:
`nL = nλ * W`
Or Increases the average time a unit spends in the system
`nL = λ * nW`

### If we reduce W

`L = λEiWi`

Where Ei is the sum of Wi where i is the amount

To further it, your pipeline will only be as fast as your slowest stage.

## The context Package
This is a <-done channel with extra bells and whistles

### Context type
```Go
type Context interface {
    // Deadline returns the time when work done on behalf of this
    // context should be canceled. Deadline returns ok==false when no
    // deadline is set. Successive calls to Deadline return the same
    // results.
    Deadline() (deadline time.Time, ok bool)
    // Done returns a channel that's closed when work done on behalf
    // of this context should be canceled. Done may return nil if this
    // context can never be canceled. Successive calls to Done return
    // the same value.
    Done() <-chan struct{}
    // Err returns a non-nil error value after Done is closed. Err
    // returns Canceled if the context was canceled or
    // DeadlineExceeded if the context's deadline passed. No other
    // values for Err are defined. After Done is closed, successive
    // calls to Err return the same value.
    Err() error
    // Value returns the value associated with this context for key,
    // or nil if no value is associated with key. Successive calls to
    // Value with the same key returns the same result.
    Value(key interface{}) interface{}
}
```

#### Two Purposes of Context
1. Provide an API for cancelling branches of your call-graph
2. Provide a data-bag for transporting request-scoped data through your call-graph.

### Cancelling Context Types
The context package helps manage all three ways a gouroutine gets terminated

The context type will be the first arg of your function, and there is nothing present in your state. 
And there is nothing that allows the funcion accepting Context to cancel it. Meaning that Context type is able to be protected from children cenceling the context. That and the done method provides a way to match all three requirments to terminating a goroutine.

### Contexts are immutable how can we affect the behavior
Take a look at the functions

```Go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

`WithCancel:` returns a new Context that closes its done channel when the returned cancel function is called
`WithDeadline:` returns a new Context that closes its done channel when the machines clock advances past a given deadline
`WithTimeout:` returns a new context that closes its dones channel after the given timeout duration

Contexts are meant to flow through your programs call-graph. 

### Some contexts you can use (not for prod)
```Go
func Background() Context
func TODO() Context
```
`Background:` returns an empty Context
`TODO:`not meant for production but also returns an empty Context and is intended as a placeholder

### Rules for storing data in a Context
1. The data should transit process or API boundaries
2. The data should be immutable
3. The data should be data, not types with methods
4. The data should trend toward simple types
5. The data should help decorate operations not drive them


