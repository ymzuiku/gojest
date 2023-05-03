package execx

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type filterFn = func(line string) string

func asyncLog(reader io.ReadCloser, filter filterFn) {
	cache := ""
	buf := make([]byte, 1024*8)
	for {
		num, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if num > 0 {
			b := buf[:num]
			s := strings.Split(string(b), "\n")
			line := strings.Join(s[:len(s)-1], "\n") // 取出整行的日志
			if filter != nil {
				line = filter(line)
				if line != "" {
					fmt.Printf("%s%s", cache, line)
				} else {
					fmt.Printf("%s", cache)
				}
			} else {
				fmt.Printf("%s%s\n", cache, line)
			}
			cache = s[len(s)-1]
		}
	}
}

func Run(ctx context.Context, filter filterFn, args ...string) error {
	fmt.Println(strings.Join(args, " "))
	var nextArgs []string
	for _, str := range args {
		if str != "" {
			nextArgs = append(nextArgs, str)
		}
	}
	cmd := exec.CommandContext(ctx, nextArgs[0], nextArgs[1:]...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %s......", err.Error())
		return err
	}

	go asyncLog(stdout, filter)
	go asyncLog(stderr, filter)
	if err := cmd.Wait(); err != nil {
		if !strings.Contains(err.Error(), "exit status 1") {
			fmt.Println(err.Error())
		}
	}

	return nil
}

func RunEmit(ctx context.Context, filter filterFn, args ...string) {
	if err := Run(ctx, filter, args...); err != nil {
		fmt.Println(err)
	}
}
