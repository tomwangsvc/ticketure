## Ticketure take-home coding question

This question asks you to fix, finish implementing and test an `http.Response` cache

### Introduction

The included `go` project implements a very simple http server that listens on port 9000 and responds to all inputs with a random string.  It should compile and run successfully without changes.

Imagine that the tasks this server does are expensive, and that the response to different inputs will be different, but that it is reasonable to re-present a prior response (for the same input) for a short period of time.

This is a situation where a response cache makes sense. A cache of recent (fresh) responses can be kept and for every incoming request the server can check the cache first for a response before doing the expensive work required to compute a new response.

The provided solution will (in a buggy manner) cache responses to GET requests.  

### Cache expectations

 - If a request for an uncached item is received, then all requests from all callers for that item in the subsequent 30 seconds should return the same response and the expensive work should only happen once.
 - Responses older than 30 seconds should never be returned.

### Ground rules

The `doExpensiveWork` handler is naive in its generation of random strings but does not need to change.  The `time.Sleep(4 * time.Second)` represents the cost of the work

The `main` function also should not change.  The inner handler `doExpensiveWork` should not know it is being wrapped by a cache.

The rest of the code is up for grabs. You can change it all. However an important part of the challenge is to identify and be able to talk about issues and concerns with the existing code. So if you choose to replace it wholesale rather than fix it, you must still be able to walk us through the problems you identified and justify your choice to re-write.

Your solution should be your own. Prefer the `go` standard library and if you must import modules from elsewhere they should be small and you should be able to justify choosing them and describe how they work.

### The challenge

1. Identify bugs and problems with the `cache.go` code and fix them
2. Implement the shutdown function for the cache to cleanly terminate the backgound go-routine
3. Extend `CachedHandler` to handle POST requests correctly
4. Demonstrate that your solution works (or doesn't) under load and fix any load related issues you find
