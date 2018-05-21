package chargen

import (
	"log"
	"os"
)

var charDataLog = log.New(os.Stderr, "CHARDATA: ", 11)

func logger(f func(string) ([]byte, error)) func(string) []byte {
	return func(s string) []byte {
		out, err := f(s)
		if err != nil {
			charDataLog.Println(startErr, err, endErr)
		}
		return out
	}
}
