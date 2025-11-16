package watch

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/ymzuiku/fswatch"
	"github.com/ymzuiku/gojest/keyboard"
	"github.com/ymzuiku/gojest/pwd"
	"github.com/ymzuiku/gojest/stack"
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

var (
	passNum          = 0
	failNum          = 0
	IgnoreTestingLog = true
	IgnoreRuntimeLog = true
)

func filter(line string) string {
	list := strings.Split(line, "\n")
	nextLine := []string{}
	for _, v := range list {
		if IgnoreRuntimeLog && strings.Contains(v, "go/src/runtime/") {
			continue
		}
		if IgnoreTestingLog && strings.Contains(v, "src/testing/testing.go:") {
			continue
		}
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
		if strings.Contains(v, "[no test files]") || strings.Contains(v, "[no tests to run]") || onlyFailReg.MatchString(v) || isonlyPath(v) {
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
		if !strings.Contains(v, "Error Trace:") && v != "" && v != "\n" {
			if strings.Contains(v, ".go:") {

				text := pwd.ReplacePwd(v)
				text = strings.ReplaceAll(text, "\n", "  ")
				text = strings.ReplaceAll(text, "\t", "  ")
				text = strings.Trim(text, " ")
				text = strings.ReplaceAll(text, "  ", "")
				nextLine = append(nextLine, stack.Red(text))
			} else {
				nextLine = append(nextLine, pwd.ReplacePwd(v))
			}
		}

	}
	if len(nextLine) == 0 {
		return ""
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
	url            = ""
	isWatch        = false
	parallel       = ""
	parallelKey    = ""
	timeoutKey     = "-timeout"
	timeout        = "20s"
	runFunctionKey = ""
	runFunction    = ""
	count          = "-count=1"
)

func fixWatchUrl(s string) []string {
	return []string{strings.Replace(s, "...", "", 1)}
}

var input = "a"
var keyboardOpened = false

func Start() {
	app := &cli.App{
		Name:  "gojest",
		Usage: "Go test runner with watch mode",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   "Enable watch mode",
			},
			&cli.StringFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Usage:   "Test timeout",
				Value:   "20s",
			},
			&cli.StringFlag{
				Name:  "count",
				Usage: "Test count",
				Value: "-count=1",
			},
			&cli.StringFlag{
				Name:    "parallel",
				Aliases: []string{"p"},
				Usage:   "Parallel test count",
			},
			&cli.StringFlag{
				Name:  "run",
				Usage: "Run specific test function",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				url = "./..."
			} else {
				url = c.Args().Get(0)
			}

			isWatch = c.Bool("watch")
			timeout = c.String("timeout")
			count = c.String("count")
			if c.String("parallel") != "" {
				parallelKey = "-parallel"
				parallel = c.String("parallel")
			}
			if c.String("run") != "" {
				runFunctionKey = "-run"
				runFunction = c.String("run")
				timeout = ""
				timeoutKey = ""
			}

			// Setup signal handler to restore terminal on SIGINT
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)
			go func() {
				<-sigChan
				cleanupAndExit()
			}()

			if err := keyboard.Open(); err != nil {
				return err
			}
			keyboardOpened = true
			defer func() {
				keyboardOpened = false
				_ = keyboard.Close()
				// Give a moment for terminal to restore
				os.Stdout.WriteString("\r\n")
			}()

			fmt.Println("Press ESC to quit")

			if count != "-count=1" {
				runNoCacheAll()
			} else {
				runAll()
			}

			if isWatch {
				go func() {
					fswatch.Watch(fixWatchUrl(url), []string{}, func(file string) {
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
					return err
				}
				// Handle ESC key
				if key == keyboard.KeyEsc {
					runQuit()
					return nil
				}
				// Handle Ctrl+C
				if key == keyboard.KeyCtrlC {
					runQuit()
					return nil
				}
				// Handle regular character keys (like 'f', 'a', etc.)
				if char != 0 {
					input = string(char)
					if fn, ok := runner[input]; ok {
						fn()
					}
				}
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func callClear() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		return
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func beforeRun() {
	passNum = 0
	failNum = 0
	callClear()
}

func afterRun() {
	if lastFail == "" {
		fmt.Printf("\n--- PASS all: %d,  FAIL: %d", passNum, failNum)
	} else {
		fmt.Printf("\n--- PASS: %d,  FAIL: %d", passNum, failNum)
	}
}

func runCommand(ctx context.Context, filterFn func(string) string, args ...string) error {
	var nextArgs []string
	for _, str := range args {
		if str != "" {
			nextArgs = append(nextArgs, str)
		}
	}
	fmt.Println(strings.Join(nextArgs, " "))
	cmd := exec.CommandContext(ctx, nextArgs[0], nextArgs[1:]...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %s......", err.Error())
		return err
	}

	go asyncLog(stdout, filterFn)
	go asyncLog(stderr, filterFn)
	if err := cmd.Wait(); err != nil {
		if !strings.Contains(err.Error(), "exit status 1") {
			fmt.Println(err.Error())
		}
	}

	return nil
}

func asyncLog(reader io.ReadCloser, filterFn func(string) string) {
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
			line := strings.Join(s[:len(s)-1], "\n")
			if filterFn != nil {
				line = filterFn(line)
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

func runAll() {
	beforeRun()
	lastFail = ""
	lastFailPath = ""
	fmt.Println("Run all:")
	_ = runCommand(context.Background(), filter, "go", "test", runFunctionKey, runFunction, url, parallelKey, parallel, timeoutKey, timeout)
	afterRun()
}

func runNoCacheAll() {
	beforeRun()
	lastFail = ""
	lastFailPath = ""
	fmt.Println("Run all no use cache:")
	_ = runCommand(context.Background(), filter, "go", "clean", "-testcache")
	_ = runCommand(context.Background(), filter, "go", "test", runFunctionKey, runFunction, url, parallelKey, parallel, count, timeoutKey, timeout)
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
		_ = runCommand(context.Background(), filter, "go", "test", runFunctionKey, runFunction, url, "-test.run", lastFail, parallelKey, parallel, timeoutKey, timeout)
	} else {
		fmt.Println("run in file: " + lastFailPath)
		_ = runCommand(context.Background(), filter, "go", "test", runFunctionKey, runFunction, lastFailPath, "-test.run", lastFail, parallelKey, parallel, timeoutKey, timeout)
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
		_ = runCommand(context.Background(), filter, "go", "test", runFunctionKey, runFunction, url, count, "-test.run", lastFail, parallelKey, parallel, timeoutKey, timeout)
	} else {
		fmt.Println("fail in path: " + lastFailPath)
		_ = runCommand(context.Background(), filter, "go", "test", runFunctionKey, runFunction, lastFailPath, count, "-test.run", lastFail, parallelKey, parallel, timeoutKey, timeout)
	}
	afterRun()
}

func cleanupAndExit() {
	if keyboardOpened {
		keyboardOpened = false
		_ = keyboard.Close()
	}
	os.Stdout.WriteString("\r\n")
	os.Exit(0)
}

func runQuit() {
	fmt.Println("Bye~")
	cleanupAndExit()
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
	str := fmt.Sprintf("\nNow action: (%s); Please keydown: (a) All, (shift+a) All no cache, (f) Focus first fail, (h) Helps, (q) Quit", input)

	if isWatch {
		fmt.Println(str + ", Watching...")
	} else {
		fmt.Println(str + "...")
	}
}
