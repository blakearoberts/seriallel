# Seriallel

Seriallel attempts to increase the performance of serializing large data structures to JSON. Namely, when an array or slice is encountered, each item is serialized in its own goroutine.

## Getting Started

```go
package main

import "github.com/blakearoberts/seriallel/json"

type Big struct {
    // ...
}

func main() {
    bigs := []Big{/* ... */}
    bytes, err := json.Marshal(bigs)
    if err != nil {
        panic(err)
    }
    println(string(bytes))
}
```

## Performance

Really bad.

In practice, I was able to reduce serialization time from ~40ms to ~15ms for ~512KB payloads where the data structure resembled the following:

```js
[
    {
        "field": [
            {},
            ..., m // m ~ 10
        ],
        ...,
    },
    ..., n // n ~ 20
]
```
