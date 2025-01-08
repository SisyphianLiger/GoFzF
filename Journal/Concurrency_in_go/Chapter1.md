# Intro to Concurrency

1. Concurrency Meaning -->  a process that occurs simultaneoulsy with one
                            or more processes
2. Amdahl's law: gains are bounded by how much of the program must be written sequentially

## Embarissingly Parallel
Example --> Calculating PI
Meaning: Sometimes there is the case where the problem domain really can be solved horizontally. I.E. more compute across different sections of the problem. 

WTF is cloud computing --> a new kind of scale and approach to application deployments and horizontal scaling. No longer machines curated and maintained but pools of resources provisioned into machiens for workloads on demand.

### Race Condition
```Go
    var data int

    go func() {
        data++
    }()

    if data == 0 {
        fmt.Printf("the value is %v.\n", data)
    }
```
Here the problem arrises that two lines of code are trying to access the same point in memory, and there is no deterministic ordering here

Three examples that could happen
1. Nothing is printed, because the func was executed before the if
2. The if block was successful but then the data gets increased printing 1 instead of 0
3. The if block prints 0 given that if block is executed to completion first

When writing concurrent code, you have to meticulously iterate through the possible scenarios

### Atomicity
Within the context of an operation it is indivisible or uninterruptible

When thinking about Atomoicity think of its scope

### Memory Access Synchronization
```Go
    var data int

    go func() { data++}()

    if data == 0 {
        fmt.Println("the value is 0.")
    } else {
        fmt.Printf("the value is %v.\n", data)
    }
```

Here there are areas we need to watch out for

Go routine accessing the data
the if accesssing the data
and the print statements accesing data

Because Data is shared between routine and main block we have 3
"critical sections"

### Deadlocks, Livelocks, Starvation
#### Deadlock
All concurrent processes are waiting on another on

Coffman Conditions: basis for techniques to help detect deadlocks

Mutual Exclusion:A concurrent process holds exclusive rights to a resourceat any one time.

Wait For Condition: A concurrent process must simultaneously hold a resource and be waiting for an additional resource.

No Preemption: A resource held by a concurrent process can only be released by that process, so it fulfills this condition.

Circular Wait: A concurrent process (P1) must be waiting on a chain of other concurrent processes (P2), which are in turn waiting on it (P1), so it fulfills this final condition
too.

#### Livelocks
Live locks, being stuck in a hallway but trying to let the person who you are stuck with go one way and you the other but end up blocked forever...

If for example there is a sync lock between two functions, the result can lead to a function being executed at the same speed. In the case of the example both members would try going left, then right until the for loop is done.

#### Starvation
Any situation where a concurrent process cannot get all the resources it needs to perform work

In the example there is a shared lock between a func that holds the lock once for 3 nano seconds and another that holds the lock for 1 nano second 3 times. A test shows that there is way more "work" don't from the greedy worker. This would indicate that the greedy worker has unnecessarily expanded its hold on the shared lock beyond the critical section and is preventing the polite workers goroutines from working

To figure out if starvation is happening, metrics are a good way to test


