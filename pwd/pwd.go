package pwd

import (
	"log"
	"os"
	"path/filepath"
)

var pwd = ""

func Load() string {
	if pwd == "" {
		file, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatalln(err)
		}
		pwd = file
	}
	return pwd
}
