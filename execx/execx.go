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
			line := strings.Join(s[:len(s)-1], "\n") //取出整行的日志
			if filter != nil {
				line = filter(line)
			}
			if line != "" {
				fmt.Printf("%s%s\n", cache, line)
			} else {
				fmt.Printf("%s", cache)
			}
			cache = s[len(s)-1]
		}
	}
}

func Run(ctx context.Context, filter filterFn, args ...string) error {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

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
