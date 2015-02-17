# Go DDP Server


## Example

```go
package main

import (
  "github.com/meteorhacks/goddp"
)

func main() {
  server := goddp.NewServer()
  server.Method("hello", methodHandler)
  server.Listen(":1337")
}

func methodHandler(p []interface{}) (interface{}, error) {
  return "result", nil
}
```
