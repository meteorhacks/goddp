# Go DDP Server

DDP server implemented with go.

## Example

```go
package main

import (
  "github.com/meteorhacks/goddp/server"
)

func main() {
  s := server.New()
  s.Method("double", handler)
  s.Listen(":1337")
}

func handler(ctx server.MethodContext) {
  n, ok := ctx.Params[0].(float64)

  if !ok {
    ctx.SendError("invalid parameters")
  } else {
    ctx.SendResult(n * 2)
  }

  ctx.SendUpdated()
}
```
