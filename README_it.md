# gojest/it

If you need log red path, you can use `gojest/it`, it package from `testify/assert`

1. First get pkg in your mod:

```sh
$ go get github.com/ymzuiku/gojest
```

2. Use `gojest/it` in your test:

```go
package xxx

import (
	"testing"

	"github.com/ymzuiku/gojest/it"
)

func TestGojest(t *testing.T) {

	t.Run("error", func(t *testing.T) {
		it.Equal(t, 20, 1)
	})
}
```

## FailNow in everyone

In one test item, if your need block in first assert fail, you can set `it.UseFailNow` to `false`:

```go
package xxx

import (
	"testing"

	"github.com/ymzuiku/gojest/it"
)

func init(){
	it.UseFailNow = true
}

```
