# Goroutines and the Go Runtime
Several Ideas from academia!

## Work Stealing
Go handles multiplexing goroutines onto OS threads for you, this is known as Work Stealing.

What does WorkStealing mean?
1. At a fork point, add tasks to the tail of the deque associated with the thread
2. if the thread is idel, steal work from the head of deque associated with some other random thread
3. at a join point that cannot be realized yet (i.e. goroutine is synchronized is not completed yet) pop work off the tail of the thread's own deque
4. if the threads deque is empty:
    - Stall at a join
    - Steal work from the head of a random thread's associated deque

## Stealing Tasks or Continuations
1. In Go goroutines are tasks
2. Everything after a goroutine is called the continuation

### What does stealing coninutations versus tasks do for us
Continuations: 
    1. will lead to a smaller queue size
    2. A serial order of execution
    3. and nonstalling join points

Go's Scheduler 3 main concepts
    1. G: A go tourint
    2. M: An OS thread references as a machine in the source code
    3. P: A context referenced as a processor in th source code

Goes runtime is started when M hosts P and scheduled G.

VERY IMPORTANT GUARENTEED DURING RUNTIME:There will always be at least enough OS threads availbale to handle hosting every context

This allows for runtime optimizations

### Thread Optimizations
1. Go dissociates the context from the OS thread so that the context can be handed off to another, unblocked OS thread
2. When the Goroutine becomes unblocked, the host OS thread attempes to steal back a context from another OS Threads to continue on the previously blocked goroutine
3. If 2 is not possible the gouroutine will be placed on a global context, and the thread will go to sleep
4. Periodically, a context will check the global context to see if there are any goroutintes there
