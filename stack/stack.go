package stack

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/ymzuiku/gojest/pwd"
)

var Debug = false

func parseFile(file string) string {
	return "." + strings.Replace(file, pwd.Load(), "", 1)
}

func New(err error) error {
	if err == nil {
		return err
	}
	if !Debug {
		return err
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return errors.New("[errox]WrapError runtime.Caller Fail")
	}

	return fmt.Errorf("%s:%d %w", parseFile(file), line, err)
}

func NewCaller(err error, caller int) error {
	if !Debug {
		return err
	}
	_, file, line, ok := runtime.Caller(caller)
	if !ok {
		return errors.New("[errox]WrapError runtime.Caller Fail")
	}

	return fmt.Errorf("%s:%d %w", parseFile(file), line, err)
}

func FileLine(caller int) string {
	if !Debug {
		return ""
	}
	_, file, line, ok := runtime.Caller(caller)
	if !ok {
		return "runtime.Caller out"
	}

	return fmt.Sprintf("%s:%d", parseFile(file), line)
}

func Red(str string) string {
	return fmt.Sprintf("\033[0;37;41m%s\033[0m", str)
}

var fails = 0

// func Pin(fail bool) bool {
// 	if !fail {
// 		fails += 1
// 		fmt.Printf("\nFaile:%d %s\n", fails, Red(FileLine(2)))
// 	}
// 	return fail
// }
