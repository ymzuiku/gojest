# gojest

Press the `f` key to focus on your first error, gojest is like jest(nodejs) in golang.

This gif is focus test first fail on keydown `f`:

![](./gojest.gif)

## feature

- Interactive CLI
- Keep your test behavior, Just use `go test ./...`
- Remove noice log
- Watch when edit code run last action

## install

```sh
$ go install github.com/ymzuiku/gojest@latest
```

## use

In your project run:

```sh
# run all
gojest
```

```sh
# run sub dir
gojest ./dir
```

And when your keydown:

- Test all: `a`
- Test all no cache: `A`
- Focus test first fail: `f`
- Focus test first fail no cache: `F`
- Quit: `q`

## other

[If you need log red path](./README_it.md)
