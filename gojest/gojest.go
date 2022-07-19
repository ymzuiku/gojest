package gojest

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/ymzuiku/gojest/execx"
)

var failReg = regexp.MustCompile(`--- FAIL`)

var fnReg = regexp.MustCompile(`--- FAIL: (.*?) \(`)

var lastFail = ""

var pwd = ""

func loadFileDir() string {
	if pwd == "" {
		file, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		pwd = file
	}
	return pwd
}

// 是否整行都只是路径
func isonlyPath(v string) bool {
	return strings.Contains(v, loadFileDir()) && !strings.Contains(strings.Replace(v, loadFileDir(), "", 1), " ")
}

func filter(line string) string {
	list := strings.Split(line, "\n")
	nextLine := []string{}
	for _, v := range list {
		if strings.Contains(v, "ok   ") || strings.Contains(v, "(cached)") || strings.Contains(v, "[no test files]") || strings.Contains(v, "[no tests to run]") || isonlyPath(v) {
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
		nextLine = append(nextLine, v)
	}

	return strings.Join(nextLine, "\n")
}

var runner = map[string]func(){
	"a": runAll,
	"f": runFocus,
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
		fmt.Println("\nPlease input: all(a), focus last(f), quit(q)")
		char, key, err := keyboard.GetKey()

		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyCtrlC {
			runQuit()
		}
		input = string(char)
		// fmt.Scan(&input)
		if fn, ok := runner[input]; ok {
			fn()
		}
	}
}

func runAll() {
	lastFail = ""
	execx.CallClear()
	fmt.Println("run all ...")
	execx.Run(context.Background(), filter, "go", "test", "./...")
}

func runFocus() {
	execx.CallClear()
	if lastFail == "" {
		fmt.Println("not have last fails")
		return
	}
	fmt.Println("run last fails: " + lastFail + " ...")
	execx.Run(context.Background(), filter, "go", "test", "./...", "-test.run", lastFail)
}

func runQuit() {
	fmt.Println("bye.")
	os.Exit(0)
}
