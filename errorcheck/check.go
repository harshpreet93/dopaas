package errorcheck

import (
	"log"
	"os"
)

func ExitOn(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
		os.Exit(1)
	}
}
