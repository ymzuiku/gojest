# gojest

gojest is like jest(nodejs) in golang.

![](./gojest.gif)

## install

```sh
$ go install github.com/ymzuiku/gojest@latest
```

## use

In your project run:

```sh
gojest
```

It's all.

## other

If you need red path log, you can use `gojest/it`, it package from `testify/assert`

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
