package execx

import (
	"os"
	"os/exec"
	"runtime"
)

var clear map[string]func() //create a map for storing clear funcs

func clearLinux() {
	{
		cmd := exec.Command(`printf`, `"\033c"`)
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	{
		cmd := exec.Command(`reset`)
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	{
		cmd := exec.Command(`clear`)
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func clearWindows() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func init() {
	clear = map[string]func(){
		"linux":   clearLinux,
		"darwin":  clearLinux,
		"windows": clearWindows,
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic(runtime.GOOS + " Your platform is unsupported! I can't clear terminal screen :(")
	}
}
