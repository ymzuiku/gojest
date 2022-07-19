# gojest

gojest is list jest(nodejs) in golang.

[](./gojest.gif)

## install

```sh
$ go install github.com/ymzuiku/gojest@latest
```

## use

In your project run:

```sh
gojest
```

## red path assert (option)

If you need red path log, you can use `gojest/it`, it package from `testify/assert`

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
