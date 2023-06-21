# gtboard
A Tensorboard event parser in Go.

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/xxr3376/gtboard/pkg/ingest"
)

func main() {
	r, err := ingest.NewIngester("some name", "your file")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	for {
		err := r.FetchUpdates()
		if err != nil {
			panic(err)
		}

		run := r.GetRun()
        // Do something with newest data
        fmt.Println(len(run.Scalars))

        // Wait for newest data
		time.Sleep(1 * time.Second)
	}
}
```
