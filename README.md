# gojest

Like Nodejs's jest in golang.

Reduced a lot of log noise, only keeping the logs that the author considers necessary for debugging during the testing process.

[中文文档](./README_cn.md)

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

**In your project run**

run all:

```sh
gojest
```

Run sub dir:

```sh
gojest ./dir/...
```

Auto run last action when your edit code:

```sh
gojest -w
```

Run Terst function:

```sh
gojest -run TestMyFunction  -w
```

Timeout 5s:

```sh
gojest -w -t 5s
```

parallel 10:

```sh
gojest -w -p 10
```

And when your keydown:

- Test all: `a`
- Focus test first fail: `f`
- Test all no cache: `shift+a`
- Focus test first fail no cache: `shift+f`
- Helps: `h`
- Quit: `q`
