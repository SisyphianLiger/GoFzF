# Concurrency at Scale
Using the patterns of Chp4 to composable systems that scale

## Error Propogation
Easy for something in your system to go wrong and difficult to understand why it happened

### Philosophy of Error Propogation
1. What is an Error: a system has entered a state which it cannot fulfill an operation a user either implicitly or explicitely requested.

Need to Explain:
    - What happened
    - When and where it occured
    - friendly user-facing message
    - How the user can more information

Errors (in the book, are Bugs or Known Edge Cases)

### LETS Study a System

SYSTEM:
    CLI Component --> INTERMEDIARY COMPONENT --> LOW LEVEL COMPONENT

Well if we have a problem in the Low-Level component, 
there is certain errors we want our users to see, as in information, and certain info that would not really help them, so we can create some type of struct to withold this info.

A simple example:
```Go
    type MyError struct {
        Inner error
        Message string
        StackTrace string
        Misc map[string]interface{}
    }
    func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
        return MyError{
            Inner: err, // 1
            Message: fmt.Sprintf(messagef, msgArgs...),
            StackTrace: string(debug.Stack()), // 2
            Misc: make(map[string]interface{}), //3
        }
    }
    func (err MyError) Error() string {
        return err.Message
    }
```
1. We store the error we are wrapping. because we alwasy want to be able to get back to the lowest-level error for investigation
2. this line of code is our stack trace where we can see where it happened
3. Catch-all of information can be stored here that may give more context to the bug

#### Low Level Module
This of course would be applied to our Low-Level Module
```Go
    // "lowlevel" module
    type LowLevelErr struct {
        error
    }
    func isGloballyExec(path string) (bool, error) {
        info, err := os.Stat(path)
        if err != nil {
            return false, LowLevelErr{(wrapError(err, err.Error()))} // 1
        }
        return info.Mode().Perm()&0100 == 0100, nil
    }
```

We wrap the rawerror here from calling os.Stat, and are OK with the message coming out of this error so no need for masking


#### Intermediate Module
```GO
    // "intermediate" module
    type IntermediateErr struct {
        error
    }
    func runJob(id string) error {
        const jobBinPath = "/bad/job/binary"

        isExecutable, err := isGloballyExec(jobBinPath)

        if err != nil {
            return err // 1
        } else if isExecutable == false {
            return wrapError(nil, "job binary is not executable")
        }
        return exec.Command(jobBinPath, "--id="+id).Run() // 1
    }
```

1. Here we passing on erros from the lowlevel module because of are architectural decision to consider erros pased on from other modules withou wrapping them in our own type bugs, this causes issues later...



#### Top Level
```Go
    func handleError(key int, err error, message string) {
        log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
        log.Printf("%#v", err) // 3
        fmt.Printf("[%v] %v", key, message)
    }

    func main() {
        log.SetOutput(os.Stdout)
        log.SetFlags(log.Ltime|log.LUTC)
        err := runJob("1")
        if err != nil {
            msg := "There was an unexpected issue; please report this as a bug."
            if _, ok := err.(IntermediateErr); ok { // 1
                msg = err.Error()
            }
            handleError(1, err, msg) // 2
        }
    }
```
1. Here we check to see if the error is of the expected type and if so it'a a well crafted terror and can go to the user

2. We bind the log and error mesage together with ID 1 

3. log out full error in case someone needs to dige into what happened

If an error message happens here say in the top level we get a nonsensical bug because the error happened at the intermediate level. So if we wrap the error in intermediate as followed

```Go
// "intermediate" module
    type IntermediateErr struct {
        error
    }

    func runJob(id string) error {
        const jobBinPath = "/bad/job/binary"
        isExecutable, err := isGloballyExec(jobBinPath)
        if err != nil {
            return IntermediateErr{wrapError( // 1 
                err,
                Error Propagation | 153
                "cannot run job %q: requisite binaries not available",
                id,
                )}
        } else if isExecutable == false {
            return wrapError(
                nil,
                "cannot run job %q: requisite binaries are not executable",
                id,
                )
        }
    return exec.Command(jobBinPath, "--id="+id).Run()
}
```
1. This is where we are wrapping the error message to make sense to the user that we are in the intermediate level

## Timeouts and Cancellation
So sometimes when we have concurrent requests, we want the ability to upon a certain time, cancel

### Why do we want concurrent processes to support time-outs?
1. System Saturation --> is a systems ability to process requests is at capacity
2. Stale Data --> data sometimes has a contextual window for which it must be processed before more data can come in
3. Attempting to prevent deadlocks --> in large systesms it can be diff to discern the data flow and edge cases, timeouts here can then guarentee your system won't deadlock
        --> Obviously if the system is build correctly its not a problem

### Causes of cancellation 
1. Timeouts: a timeout is an implicit cancellation
2. User intervention: a good user experience sometimes allows for users to cancel the operation they've started
3. parent cancellation: if any kind of parent of a concurrent op stops the child of that parent stops too
4. replicated requests: wish to send data to multiple concurrent processes to get faster response from one

### How can we do this?
```Go
var value interface{}

select {
    case <-done:
        return
    case value = <-valueStream:
}

result := reallyLongCalculation(value)

select {
    case <-done:
        return
    case resultStream <-result:
}
```

This couples the value stream to write to the result stream. But `reallyLongCalculation` is not preemtable and may take awhile:

```Go
reallyLongCalculation := func(done <-chan interface{}, value interface{},) interface{} {
    intermediateResult := longCalculation(value)
    select {
        case <-done:
            return nil
        default:
    }
    return longCalculation(intermediateResult)
}
```
This means now that we have halved the problem, reallyLongCalculation is now preemtable but we can only preempt reallyLongCalculation, so we need to make longCalculation preemtable as well

```Go
reallyLongCalculation := func(done <-chan interface{}, value interface{},) interface{} {
    intermediateResult := longCalculation(done, value)
    return longCalculation(intermediateResult)
}
```

Two make the logic preemtable we mush do two things.

1. define the period within which our concurrent process is preemtbale, 
2. ensure that any functionality that takes more time that this period is also preemtable.


So synopsis is again, we can split tasks in a pipeline via different calls the function which should hopefully reduce load on a particular node.

## Heartbeats
Heartbeats are a way for concurrent process to signal life outside partties

Heartbeat signifies life to an observer
### Two types
1. Heartbeats occur on a time interval
2. Heartbeats occur at the begining of a unit of work

Heartbeats are intense...

```Go
doWork := func( done <-chan interface{}, pulseInterval time.Duration,) 
    (<-chan interface{}, <-chan time.Time) {

        heartbeat := make(chan interface{}) // 1
        results := make(chan time.Time)

        go func() {
            defer close(heartbeat)
            defer close(results)

            pulse := time.Tick(pulseInterval) // 2
            workGen := time.Tick(2*pulseInterval) // 3

            sendPulse := func() {
                select {
                case heartbeat <-struct{}{}:
                default: // 4
                }
            }
            sendResult := func(r time.Time) {
                for {
                    select {
                    case <-done:
                        return
                    case <-pulse: // 5
                        sendPulse()
                    case results <- r:
                        return
                    }
                }
            }
            for {
                select {
                case <-done:
                    return
                case <-pulse: // 5
                    sendPulse()
                case r := <-workGen:
                    sendResult(r)
                }
            }
        }()
        return heartbeat, results
    }

```
1. Here is the heartbeat channel to send heartbeats to 
2. we set heartbeat to pulse at the pulseInterval we werer given.
3. Another ticker to simulate work coming in
4. We include a default here to guard against the fact that no one may be listening to the heartbeat
5. just like done anytime you perform a send or recieve you also need to include a case for pulse


How can we utilize this function to consume events it emits
```Go
    done := make(chan interface{})
    time.AfterFunc(10*time.Second, func() { close(done) }) // 1

    const timeout = 2*time.Second // 2
    heartbeat, results := doWork(done, timeout/2) // 3
    for {
        select {
        case _, ok := <-heartbeat: // 4
            if ok == false {
                return
            }
            fmt.Println("pulse")
        case r, ok := <-results: // 5
            if ok == false {
                return
            }
            fmt.Printf("results %v\n", r.Second())
        case <-time.After(timeout): // 6
            return
        }
    }
```

1. We set up the standard done channel and close it after 10 seconds
2. We set out timeout period to couple our heartbeat interval to our timeout
3. pass timeout/2 here gives heartbeat an extra tick to respond
4. select on heartbeat, when no results guarenteed a mseeage from hb channel every timeout/2
5. select from results channel nothing fancy
6. timeout if we haven't recieved either heartbead or new results


#### Main Use
With Heartbeats you are able to make your tests for concurrent code deterministic, buy making sure that each go routine reachs a point of completion or a timeout occurs. So pretty nifty testing, but caveat is that if there is a pipeline to be tested, with a func that is large, this will maybe not be the best method of testing.

## Replicated Requests
APplication is servicing a users HTTP request or retriving a blob of data, you can either replicate the request to multiple handlers and one of them will return faster than the other ones and can immediately return the result. 

Downside is resources being used so watch out to not replicate processes, servers, data centers

What does this process look like?
### Go Code to replicate
```Go
    doWork := func(
        done <-chan interface{},
        id int,
        wg *sync.WaitGroup,
        result chan<- int,
    ) {
            started := time.Now()
            defer wg.Done()
            // Simulate random load
            simulatedLoadTime := time.Duration(1+rand.Intn(5))*time.Second
            select {
            case <-done:
            case <-time.After(simulatedLoadTime):
            }
            select {
            case <-done:
            case result <- id:
            }
            took := time.Since(started)
            // Display how long handlers would have taken
            if took < simulatedLoadTime {
                took = simulatedLoadTime
            }
            fmt.Printf("%v took %v\n", id, took)
    }

    done := make(chan interface{})
    result := make(chan int)
    var wg sync.WaitGroup
    wg.Add(10)
    for i:=0; i < 10; i++ { // 1
        go doWork(done, i, &wg, result)
    }
    firstReturned := <-result // 2
    close(done) // 3
    wg.Wait()
    fmt.Printf("Received an answer from #%v\n", firstReturned)
```

1. We create 10 Handlers to handle our request
2. Line grabs first value returned from handlers
3. close the remaining handlers to stop work

#### Takeaways
With this technique you can up your speed to deliver a service if its important, however, with speed comes a memory overhead cost, and the requests should only replicate if you are using handlers that are somewhat different from one another. Handlers with uniformity will not gain the same level of speed here.

## Rate Limiting
We we rate limit, cuz bad actors and other reasons but lets make one!

```Go
    func Open() *APIConnection {
        return &APIConnection{}
    }
    
    type APIConnection struct{}

    func (a *APIConnection) ReadFile(ctx context.Context) error {
        // Pretend we read hear
        return nil
    }

    func (a *APIConnection) ResolveAddress(ctx context.Context) error {
        // Pretend we read hear
        return nil
    }
```
In theory this request is going over the wire, we take `context.Context` as the first arg in case it needs to cancel the request or pass values over the server.

Solution: Create simple driver to access API, needs to read 10 files and resolve 10 addresses

```Go
func main() {
    defer log.Printf("Done.")
    log.SetOutput(os.Stdout)
    log.SetFlags(log.Ltime | log.LUTC)
    apiConnection := Open()
    var wg sync.WaitGroup

    wg.Add(20)
    for i := 0; i < 10; i++ {
        go func() {
            defer wg.Done()
            err := apiConnection.ReadFile(context.Background())
            if err != nil {
                log.Printf("cannot ReadFile: %v", err)
            }
            log.Printf("ReadFile")
        }()
    }
    for i := 0; i < 10; i++ {
        go func() {
            defer wg.Done()
            err := apiConnection.ResolveAddress(context.Background())
            if err != nil {
                log.Printf("cannot ResolveAddress: %v", err)
            }
            log.Printf("ResolveAddress")
        }()
    }
    wg.Wait()
}
```

So here is the thing right now if this is our driver, it can be accessed many many times, and that is expensive to say the least so what do we need to fix it?

### The Limit Package (Token BUcket rate limiter) 
From: golang.org/x/time/rate

```Go
    // Limit defines the maximum frequency of some events. Limit is
    // represented as number of events per second. A zero Limit allows no
    // events.
    type Limit float64
    // NewLimiter returns a new Limiter that allows events up to rate r
    // and permits bursts of at most b tokens.
    func NewLimiter(r Limit, b int) *Limiter
```

1. Limit: represented as the numebr of events per second
2. NewLimiter: allows limits up to r

OK OK OK OK Ratelimiter hype!

```Go
    func Open() *APIConnection {
        return &APIConnection{
            rateLimiter: rate.NewLimiter(rate.Limit(1), 1), // 1    
        }
    }
    
    type APIConnection struct{
        rateLimiter: rate.Limiter,  
    }

    func (a *APIConnection) ReadFile(ctx context.Context) error {
        if err := a.rateLimiter.Wait(ctx); err != nil {
            return err
        }
        // Pretend we read hear
        return nil
    }

    func (a *APIConnection) ResolveAddress(ctx context.Context) error {
        if err := a.rateLimiter.Wait(ctx); err != nil {
            return err
        }
        // Pretend we read hear
        return nil
    }
```

Super easy to implement! Just simply att the Rate limit package adn then we check before each api call that a limit from our context has been met or not.

The way the Rate limiter is set up is to restrict the events by 1 second, the output of this new code will log such values. In production we establish tiers of rate limits etc.

## Healing Unhealthy Goroutines
We create logic that monitors a goroutines health called a steward
Stewards are responsible for restarting a ward's goroutine should it become unhealthy

```Go
// Type definition for the goroutine starter function
type startGoroutineFn func(done <-chan interface{}, pulseInterval time.Duration) 
(heartbeat <-chan interface{}) // 1

// Creates a new steward that monitors a goroutine
newSteward := func(timeout time.Duration, startGoroutine startGoroutineFn) startGoroutineFn { // 2
    return func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} {
        heartbeat := make(chan interface{})
        
        go func() {
            defer close(heartbeat)
            
            var wardDone chan interface{}
            var wardHeartbeat <-chan interface{}
            
            startWard := func() { // 3
                wardDone = make(chan interface{})
                wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2)
            }
            
            startWard()
            pulse := time.Tick(pulseInterval)
            
        monitorLoop:
            for {
                timeoutSignal := time.After(timeout)
                for {
                    select {
                    case <-pulse:
                        select {
                        case heartbeat <- struct{}{}:
                        default:
                        }
                    case <-wardHeartbeat:
                        continue monitorLoop
                    case <-timeoutSignal:
                        log.Println("steward: ward unhealthy; restarting")
                        close(wardDone)
                        startWard()
                        continue monitorLoop
                    case <-done:
                        return
                    }
                }
            }
        }()
        
        return heartbeat
    }
}
```

The main thing here is the call to startWard() when a timeout Signal is given. Think of it this way, say I make a request to a third party API. That request takes an incredibly long time, and will not be done anytime soon. Because I have "wrapped" this function with a steward, if it is the case on our end that the function takes too long to call, then the Steward acts as an API capable of restarting this routine and reapplying the request as well as logging the issue.


