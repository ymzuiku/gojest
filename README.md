# gojest

Press the `f` key to focus on your first error, gojest is like jest(nodejs) in golang.

This gif is focus test first fail on keydown `f`:

![](./gojest.gif)

## feature

- Interactive CLI
- Keep your test behavior, Just use `go test ./...`
- Remove noice log

## install

```sh
$ go install github.com/ymzuiku/gojest@latest
```

## use

In your project run:

```sh
gojest
```

And when your keydown:

- Test all: `a`
- Test all no cache: `A`
- Focus test first fail: `f`
- Focus test first fail no cache: `F`
- Quit: `q`

## other

[If you need log red path](./README_it.md)
