package pwd

import (
	"log"
	"os"
)

var pwd = ""

func Load() string {
	if pwd == "" {
		file, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		pwd = file
	}
	return pwd
}
