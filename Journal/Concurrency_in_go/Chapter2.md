# Modeling Your Code: Communicating Sequential Processes
## The Difference between Concurrency and Parallelism
Concurrency: a property of the code
Parallelism: a property of the running program


Interesting thought: We do not write parallel code, only concurrent code we hope will be run in parallel


Summary:

Go uses two models of concurrency 

CSP --> Communicatiing Sequntial Processes
I.E. linking goroutines with channels as opposed to Mutex

The main point however is to try to figure out when to use channels versus when to use primitives.

A general rule of thumb is that when trying to sync large mutating data points, having a synchronization point with a Channel makes a lot of sense. It is not that mutex cannot be used, but may be more difficult to implement.

But the overarching principle in Go, is keep concurrency simple, so when applicable use channels and do not worry about goroutine costs.

