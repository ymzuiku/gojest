package pwd

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var pwd = ""

func Pwd() string {
	if pwd == "" {
		file, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatalln(err)
		}
		pwd = file
	}
	return pwd
}

func ReplacePwd(v string) string {
	return strings.ReplaceAll(v, Pwd(), ".")
}
