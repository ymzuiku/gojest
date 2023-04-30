# gojest/expect

If you need fail now in every assert, you can use `gojest/expect`, it package from `testify/assert`

1. First get pkg in your mod:

```sh
$ go get github.com/ymzuiku/gojest
```

2. Use `gojest/expect` in your test:

```go
package xxx

import (
	"testing"

	"github.com/ymzuiku/gojest/expect"
)

func TestGojest(t *testing.T) {

	t.Run("error", func(t *testing.T) {
		expect.Equal(t, 20, 1)
	})
}
```

## FailNow in everyone

In one test item, if your need block in first assert fail, you can set `expect.UseFailNow` to `false`:

```go
package xxx

import (
	"testing"

	"github.com/ymzuiku/gojest/expect"
)

func init(){
	expect.UseFailNow = true
}

```
