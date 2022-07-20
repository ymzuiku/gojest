package watch

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/ymzuiku/gojest/execx"
	"github.com/ymzuiku/gojest/pwd"
)

var failReg = regexp.MustCompile(`--- FAIL`)
var onlyFailReg = regexp.MustCompile(`^FAIL`)

var fnReg = regexp.MustCompile(`--- FAIL: (.*?) \(`)

var lastFail = ""

// 是否整行都只是路径
func isonlyPath(v string) bool {
	return strings.Contains(v, pwd.Load()) && !strings.Contains(strings.Replace(v, pwd.Load(), "", 1), " ")
}

func replacePwd(v string) string {
	return strings.ReplaceAll(v, pwd.Load(), ".")
}

func filter(line string) string {
	list := strings.Split(line, "\n")
	nextLine := []string{}
	for _, v := range list {
		if strings.Contains(v, "ok   ") || strings.Contains(v, "(cached)") || strings.Contains(v, "[no test files]") || strings.Contains(v, "[no tests to run]") || onlyFailReg.MatchString(v) || isonlyPath(v) {
			continue
		}
		if failReg.MatchString(v) {
			name := fnReg.FindStringSubmatch(v)[1]
			if lastFail == "" {
				lastFail = name
			} else if strings.Contains(name, lastFail+"/") {
				lastFail = name
			}
		}
		nextLine = append(nextLine, replacePwd(v))
	}

	return strings.Join(nextLine, "\n")
}

var runner = map[string]func(){
	"a": runAll,
	"A": runNoCacheAll,
	"f": runFocus,
	"F": runNoCacheFocus,
	"q": runQuit,
}

func Start() {
	var input string

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")

	runAll()
	for {
		fmt.Println("\nPlease keydown: (a) All, (A) All no cache, (f) Focus first fail, (F) Focus first fail no cache, (q) Quit...")
		char, key, err := keyboard.GetKey()

		if err != nil {
			panic(err)
		}
		input = string(char)
		if fn, ok := runner[input]; ok {
			fn()
		} else if key == keyboard.KeyCtrlC {
			runQuit()
		}
	}
}

func runAll() {
	lastFail = ""
	execx.CallClear()
	fmt.Println("Run all ...")
	execx.Run(context.Background(), filter, "go", "test", "./...")
}

func runNoCacheAll() {
	lastFail = ""
	execx.CallClear()
	fmt.Println("Run all no use cache ...")
	execx.Run(context.Background(), filter, "go", "test", "./...", "-count=1")
}

func runFocus() {
	execx.CallClear()
	if lastFail == "" {
		fmt.Println("Not have last fails")
		return
	}
	fmt.Println("Run last fails: " + lastFail + " ...")
	execx.Run(context.Background(), filter, "go", "test", "./...", "-test.run", lastFail)
}

func runNoCacheFocus() {
	execx.CallClear()
	if lastFail == "" {
		fmt.Println("Not have last fails")
		return
	}
	fmt.Println("Run last fails no cache: " + lastFail + " ...")
	execx.Run(context.Background(), filter, "go", "test", "./...", "-count=1", "-test.run", lastFail)
}

func runQuit() {
	fmt.Println("bye.")
	os.Exit(0)
}
