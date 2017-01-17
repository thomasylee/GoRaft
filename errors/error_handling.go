package errors

import (
	"log"
	"os"
)

func HandleError(prefix string, err error, exit bool) bool {
	if err == nil {
		return false
	}

	log.Fatal(prefix, err.Error())
	if exit {
		os.Exit(1)
	}
	return true
}
