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
var lastFailPath = ""

// 是否整行都只是路径
func isonlyPath(v string) bool {
	return strings.Contains(v, pwd.Pwd()) && !strings.Contains(strings.Replace(v, pwd.Pwd(), "", 1), " ")
}

var lastLine = ""

func filter(line string) string {
	// defer func() {
	// 	lastFail = line
	// }()
	list := strings.Split(line, "\n")
	nextLine := []string{}
	for _, v := range list {
		if lastLine != "" {
			lastLine = ""
			arr := strings.Split(v, pwd.Pwd())
			if len(arr) == 2 {
				path := strings.Split(arr[1], ".go:")[0]
				lastFailPath = "." + strings.Trim(path, " ") + ".go"
			}
		}
		if strings.Contains(v, "Error Trace:") {
			lastLine = v
		}
		// if lastItFail == "" && itFailReg.MatchString(v) {
		// 	arr := strings.Split(v, pwd.Pwd())
		// 	if len(arr) == 2 {
		// 		path := strings.Split(arr[1], ".go:")[0]
		// 		lastItFail = "." + strings.Trim(path, " ") + ".go"
		// 	}
		// }
		if strings.Contains(v, "ok   ") || strings.Contains(v, "(cached)") || strings.Contains(v, "[no test files]") || strings.Contains(v, "[no tests to run]") || onlyFailReg.MatchString(v) || isonlyPath(v) {
			continue
		}
		if failReg.MatchString(v) {
			name := fnReg.FindStringSubmatch(v)[1]
			if lastFail == "" {
				lastFail = name
				lastFailPath = ""
			} else if strings.Contains(name, lastFail+"/") {
				lastFail = name
				lastFailPath = ""
			}
		}
		nextLine = append(nextLine, pwd.ReplacePwd(v))
	}
	if len(nextLine) == 0 {
		return "-"
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

var url string

func Start() {
	if len(os.Args) < 2 {
		url = "./..."
	} else {
		url = os.Args[1]
	}

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
	fmt.Println("Run all:")
	execx.Run(context.Background(), filter, "go", "test", url)
}

func runNoCacheAll() {
	lastFail = ""
	execx.CallClear()
	fmt.Println("Run all no use cache:")
	execx.Run(context.Background(), filter, "go", "test", url, "-count=1")
}

func runFocus() {
	execx.CallClear()
	if lastFail == "" {
		fmt.Println("Not have last fails")
		return
	}
	fmt.Println("Run last fails: " + lastFail)
	if lastFailPath == "" {
		execx.Run(context.Background(), filter, "go", "test", url, "-test.run", lastFail)
	} else {
		fmt.Println("fail in path: " + lastFailPath)
		execx.Run(context.Background(), filter, "go", "test", lastFailPath, "-test.run", lastFail)
	}

}

func runNoCacheFocus() {
	execx.CallClear()
	if lastFail == "" {
		fmt.Println("Not have last fails")
		return
	}
	fmt.Println("Run last fails no cache: " + lastFail)
	if lastFailPath == "" {
		execx.Run(context.Background(), filter, "go", "test", url, "-count=1", "-test.run", lastFail)
	} else {
		fmt.Println("fail in path: " + lastFailPath)
		execx.Run(context.Background(), filter, "go", "test", lastFailPath, "-count=1", "-test.run", lastFail)
	}
}

func runQuit() {
	fmt.Println("Bye.")
	os.Exit(0)
}
