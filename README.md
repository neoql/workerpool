# workerpool

Concurrency limiting goroutine pool.

## Install

```
go get -v github.com/neoql/workerpool
```

## Example

```go
package main

import (
    "fmt"
    "time"
    "sync"

    "github.com/neoql/workerpool"
)

func main() {
    var wg sync.WaitGroup

    fn := func(argv workerpool.Argv) {
        no := argv.(int)
        fmt.Printf("fn-%d enter...\n", no)
        time.Sleep(time.Second*3)
        fmt.Printf("fn-%d exit.\n", no)
        wg.Done()
    }

    pool := workerpool.New(2, time.Second*1)

    wg.Add(3)
    for i := 0; i < 3; i++ {
        if !pool.Spawn(fn, i) {
            fmt.Printf("launch fn-%d failed, WorkerPool is full", i)
            fmt.Println("Waiting for an empty position")
            pool.WaitSpawn(fn, i)
        }
    }

    wg.Wait()
}
```
```
fn-0 enter...
fn-1 enter...
launch fn-2 failed, WorkerPool is fullWaiting for an empty position
fn-1 exit.
fn-0 exit.
fn-2 enter...
fn-2 exit.
```

## License

MIT, read more [here](./LICENSE).