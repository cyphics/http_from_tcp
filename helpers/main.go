// Package helpers defines common functions
package helpers

import (
	"log"
)

func CheckFatal(err error, msg string) {
	if err != nil {
		log.Fatalf("Fatal Error: %s, %s\n", msg, err.Error())
	}
}
