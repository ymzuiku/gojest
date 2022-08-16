package watch

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ymzuiku/fswatch"
	"github.com/ymzuiku/gojest/execx"
	"github.com/ymzuiku/gojest/keyboard"
	"github.com/ymzuiku/gojest/pwd"
)

var (
	testGoReg = regexp.MustCompile(`(\.|_)test\.go`)
	// itGoReg      = regexp.MustCompile(`it\.go\:`)
	passReg      = regexp.MustCompile(`^ok `)
	failReg      = regexp.MustCompile(`--- FAIL`)
	onlyFailReg  = regexp.MustCompile(`^FAIL`)
	fnReg        = regexp.MustCompile(`--- FAIL: (.*?) \(`)
	lastFail     = ""
	lastFailPath = ""
)

// 是否整行都只是路径
func isonlyPath(v string) bool {
	return strings.Contains(v, pwd.Pwd()) && !strings.Contains(strings.Replace(v, pwd.Pwd(), "", 1), " ")
}

// var lastLine = ""
var passNum = 0
var failNum = 0

func filter(line string) string {
	list := strings.Split(line, "\n")
	nextLine := []string{}
	for _, v := range list {
		if lastFail != "" && lastFailPath == "" && testGoReg.MatchString(v) {
			arr := strings.Split(v, pwd.Pwd())
			if len(arr) == 2 {
				path := strings.Split(arr[1], ".go:")[0]
				lastFailPath = "." + strings.Trim(path, " ")
				list := strings.Split(lastFailPath, "/")
				lastFailPath = strings.Join(list[:len(list)-1], "/")
			}
		}
		if passReg.MatchString(v) {
			passNum += 1
		}
		if strings.Contains(v, "(cached)") || strings.Contains(v, "[no test files]") || strings.Contains(v, "[no tests to run]") || onlyFailReg.MatchString(v) || isonlyPath(v) {
			continue
		}
		if failReg.MatchString(v) {
			failNum += 1
			name := fnReg.FindStringSubmatch(v)[1]
			if lastFail == "" {
				lastFail = "^" + name + "$"
				lastFailPath = ""
			}
		}
		if !strings.Contains(v, "Error Trace:") {
			nextLine = append(nextLine, pwd.ReplacePwd(v))
		}

	}
	if len(nextLine) == 0 {
		return "-"
	}
	return strings.Join(nextLine, "\n") + "\n"
}

var runner = map[string]func(){
	"a": runAll,
	"A": runNoCacheAll,
	"f": runFocus,
	"F": runNoCacheFocus,
	"q": runQuit,
	"h": runHelp,
}

var (
	url     string
	isWatch bool
)

func fixWatchUrl(s string) []string {
	return []string{strings.Replace(s, "...", "", 1)}
}

var input = "a"

func Start() {
	if len(os.Args) < 2 {
		fmt.Println("\nerror, gojest need input path, like: ./...")
		os.Exit(1)
	}

	url = os.Args[1]

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")

	runAll()
	for _, arg := range os.Args {
		if arg == "-w" {
			isWatch = true
		}
	}

	if isWatch {
		go func() {
			fswatch.Watch(fixWatchUrl(url), func(_ string) {
				if fn, ok := runner[input]; ok {
					fn()
					printTip()
				}
			})
		}()
	}

	for {
		printTip()
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

func beforeRun() {
	passNum = 0
	failNum = 0
	execx.CallClear()
}

func afterRun() {
	if lastFail == "" {
		fmt.Printf("\n--- PASS all: %d,  FAIL: %d", passNum, failNum)
	} else {
		fmt.Printf("\n--- PASS: %d,  FAIL: %d", passNum, failNum)
	}
}

func runAll() {
	beforeRun()
	lastFail = ""
	lastFailPath = ""
	fmt.Println("Run all:")
	execx.Run(context.Background(), filter, "go", "test", url)
	afterRun()
}

func runNoCacheAll() {
	beforeRun()
	lastFail = ""
	lastFailPath = ""
	fmt.Println("Run all no use cache:")
	execx.Run(context.Background(), filter, "go", "test", url, "-count=1")
	afterRun()
}

func runFocus() {
	beforeRun()
	if lastFail == "" {
		fmt.Println("Not have last fails, run all")
		runAll()
		return
	}
	fmt.Println("Run last fails: " + lastFail)
	if lastFailPath == "" {
		execx.Run(context.Background(), filter, "go", "test", url, "-test.run", lastFail, "-timeout", "30s")
	} else {
		fmt.Println("run in file: " + lastFailPath)
		execx.Run(context.Background(), filter, "go", "test", lastFailPath, "-test.run", lastFail, "-timeout", "30s")
	}
	afterRun()
}

func runNoCacheFocus() {
	beforeRun()
	if lastFail == "" {
		fmt.Println("Not have last fails")
		runAll()
		return
	}
	fmt.Println("Run last fails no cache: " + lastFail)
	if lastFailPath == "" {
		execx.Run(context.Background(), filter, "go", "test", url, "-count=1", "-test.run", lastFail, "-timeout", "30s")
	} else {
		fmt.Println("fail in path: " + lastFailPath)
		execx.Run(context.Background(), filter, "go", "test", lastFailPath, "-count=1", "-test.run", lastFail, "-timeout", "30s")
	}
	afterRun()
}

func runQuit() {
	fmt.Println("Bye.")
	os.Exit(0)
}

func runHelp() {
	fmt.Println("\nPlease keydown:")
	fmt.Println("Run all test: (a)")
	fmt.Println("Run all test and no cache: (shift+a)")
	fmt.Println("Run first fail: (f)")
	fmt.Println("Run first fail and no cache: (shift+f)")
	fmt.Println("View helps: (f)")
	fmt.Println("Quit: (q)")
}

func printTip() {
	str := fmt.Sprintf("\nNow action: (%s); Please keydown: (a) All, (f) Focus first fail, (h) Helps, (q) Quit", input)

	if isWatch {
		fmt.Println(str + ", Watching...")
	} else {
		fmt.Println(str + "...")
	}
}
