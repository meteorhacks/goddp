# Go DDP Server


## Example

```go
package main

import (
  "github.com/meteorhacks/goddp/server"
)

func main() {
  server := server.New()
  server.Method("hello", methodHandler)
  server.Listen(":1337")
}

func methodHandler(p []interface{}) (interface{}, error) {
  return "result", nil
}
```
